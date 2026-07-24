// Package atomicwrite provides TOCTOU-safe file writes using xxhash64
// fingerprint verification, cross-platform file locking via gofrs/flock,
// atomic rename, and fsync for crash durability.
package atomicwrite

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/cespare/xxhash/v2"
	"github.com/gofrs/flock"
)

// ErrConcurrentModification indicates the file was modified by another
// process between the fingerprint read and the write attempt.
var ErrConcurrentModification = errors.New("file was modified concurrently since read")

// Fingerprint is an xxhash64 digest of file content at read time.
// A zero-value Fingerprint indicates no prior file existed.
type Fingerprint [8]byte

// IsZero returns true if the fingerprint represents no prior file.
func (fp Fingerprint) IsZero() bool {
	return fp == Fingerprint{}
}

// Matches returns true if the given content produces the same fingerprint.
func (fp Fingerprint) Matches(content []byte) bool {
	return FingerprintFromBytes(content) == fp
}

// FingerprintFromBytes computes an xxhash64 Fingerprint from raw content.
func FingerprintFromBytes(data []byte) Fingerprint {
	var fp Fingerprint
	binary.BigEndian.PutUint64(fp[:], xxhash.Sum64(data))

	return fp
}

// FingerprintFile computes an xxhash64 Fingerprint from a file's current content.
// Returns a zero-value Fingerprint if the file does not exist.
func FingerprintFile(path string) (Fingerprint, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is caller-controlled
	if err != nil {
		if os.IsNotExist(err) {
			return Fingerprint{}, nil
		}

		return Fingerprint{}, fmt.Errorf("reading %s for fingerprint: %w", path, err)
	}

	return FingerprintFromBytes(data), nil
}

// Write writes data to path with TOCTOU protection and crash durability.
// Data is staged to a unique temp file, fsync'd, then atomically renamed
// over the target. The target directory is fsync'd after rename (POSIX).
// If fingerprint is non-zero, it verifies the file hasn't changed since the
// fingerprint was computed, using cross-platform file locking (flock on Unix,
// LockFileEx on Windows) and atomic rename.
// A zero-value fingerprint skips verification (first run).
func Write(path string, data []byte, fingerprint Fingerprint) error {
	tmpPath, perm, err := prepareTempPath(path)
	if err != nil {
		return err
	}

	stageErr := writeAndSync(tmpPath, data, perm)
	if stageErr != nil {
		return stageErr
	}

	return commitOrRename(path, tmpPath, fingerprint)
}

// writeFuncBufferSize is the default buffer size for streaming writes.
const writeFuncBufferSize = 65536

// WriteFunc writes to path via a streaming callback with TOCTOU protection
// and crash durability. The callback receives a buffered writer (64KB buffer)
// and may stream content incrementally without holding the full payload in
// memory. The temp file is fsync'd before atomic rename.
//
// If fingerprint is non-zero, it verifies the file hasn't changed since the
// fingerprint was computed. A zero-value fingerprint skips verification.
//
// Use WriteFunc instead of Write when the content is large or produced
// incrementally (e.g., JSON encoders, diagram renderers).
func WriteFunc(path string, fn func(w io.Writer) error, fingerprint Fingerprint) error {
	tmpPath, perm, err := prepareTempPath(path)
	if err != nil {
		return err
	}

	stageErr := writeFuncAndSync(tmpPath, fn, perm)
	if stageErr != nil {
		return stageErr
	}

	return commitOrRename(path, tmpPath, fingerprint)
}

// prepareTempPath computes the file mode (preserving the existing file's
// mode or using the default) and returns a unique temp file path alongside it.
func prepareTempPath(path string) (string, fs.FileMode, error) {
	const defaultFilePerm = fs.FileMode(0o644)

	perm := defaultFilePerm

	info, err := os.Stat(path)
	if err == nil {
		perm = info.Mode().Perm()
	}

	suffix, suffixErr := randomSuffix()
	if suffixErr != nil {
		return "", 0, fmt.Errorf("generating temp file suffix: %w", suffixErr)
	}

	return path + "." + suffix + ".tmp", perm, nil
}

// commitOrRename verifies fingerprint and renames, or skips verification when
// the fingerprint is the zero sentinel for a first write.
func commitOrRename(path, tmpPath string, fingerprint Fingerprint) error {
	if !fingerprint.IsZero() {
		return commitWithVerification(path, tmpPath, fingerprint)
	}

	return atomicRename(path, tmpPath)
}

func writeFuncAndSync(tmpPath string, fn func(io.Writer) error, perm fs.FileMode) error {
	file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm) //nolint:gosec // caller-controlled
	if err != nil {
		return fmt.Errorf("creating temp file %s: %w", tmpPath, err)
	}

	bw := bufio.NewWriterSize(file, writeFuncBufferSize)

	writeErr := fn(bw)

	flushErr := bw.Flush()

	syncErr := file.Sync()

	closeErr := file.Close()

	if writeErr != nil {
		_ = os.Remove(tmpPath)

		return fmt.Errorf("writing temp file %s: %w", tmpPath, writeErr)
	}

	if flushErr != nil {
		_ = os.Remove(tmpPath)

		return fmt.Errorf("flushing temp file %s: %w", tmpPath, flushErr)
	}

	if syncErr != nil {
		_ = os.Remove(tmpPath)

		return fmt.Errorf("syncing temp file %s: %w", tmpPath, syncErr)
	}

	if closeErr != nil {
		_ = os.Remove(tmpPath)

		return fmt.Errorf("closing temp file %s: %w", tmpPath, closeErr)
	}

	return nil
}

func commitWithVerification(path, tmpPath string, fingerprint Fingerprint) error {
	fileLock := flock.New(path)

	lockErr := fileLock.Lock()
	if lockErr != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("acquiring exclusive lock on %s: %w", path, lockErr)
	}
	defer func() { _ = fileLock.Close() }()

	current, err := os.ReadFile(path) //nolint:gosec // path is caller-controlled
	if err != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("re-reading %s for verification: %w", path, err)
	}

	if !fingerprint.Matches(current) {
		cleanupTmp(tmpPath)

		return fmt.Errorf("%w: %s was modified since read", ErrConcurrentModification, path)
	}

	return atomicRename(path, tmpPath)
}

func randomSuffix() (string, error) {
	var buf [4]byte

	_, err := rand.Read(buf[:])
	if err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}

	return hex.EncodeToString(buf[:]), nil
}

func writeAndSync(tmpPath string, data []byte, perm fs.FileMode) error {
	file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm) //nolint:gosec // caller-controlled
	if err != nil {
		return fmt.Errorf("creating temp file %s: %w", tmpPath, err)
	}

	_, writeErr := file.Write(data)
	if writeErr != nil {
		_ = file.Close()
		_ = os.Remove(tmpPath)

		return fmt.Errorf("writing temp file %s: %w", tmpPath, writeErr)
	}

	syncErr := file.Sync()
	if syncErr != nil {
		_ = file.Close()
		_ = os.Remove(tmpPath)

		return fmt.Errorf("syncing temp file %s: %w", tmpPath, syncErr)
	}

	closeErr := file.Close()
	if closeErr != nil {
		_ = os.Remove(tmpPath)

		return fmt.Errorf("closing temp file %s: %w", tmpPath, closeErr)
	}

	return nil
}

func cleanupTmp(tmpPath string) {
	_ = os.Remove(tmpPath)
}
