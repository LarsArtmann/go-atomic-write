// Package atomicwrite provides TOCTOU-safe file writes using xxhash64
// fingerprint verification, cross-platform file locking via gofrs/flock,
// and atomic rename for crash safety.
package atomicwrite

import (
	"encoding/binary"
	"errors"
	"fmt"
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

// Write writes data to path with TOCTOU protection.
// If fingerprint is non-zero, it verifies the file hasn't changed since the
// fingerprint was computed, using cross-platform file locking (flock on Unix,
// LockFileEx on Windows) and atomic rename for crash safety.
// A zero-value fingerprint skips verification (first run).
func Write(path string, data []byte, fingerprint Fingerprint) error {
	const defaultFilePerm = fs.FileMode(0o644)

	perm := defaultFilePerm

	info, err := os.Stat(path)
	if err == nil {
		perm = info.Mode().Perm()
	}

	tmpPath := path + ".tmp"

	writeErr := os.WriteFile(tmpPath, data, perm)
	if writeErr != nil {
		return fmt.Errorf("writing temp file %s: %w", tmpPath, writeErr)
	}

	if !fingerprint.IsZero() {
		return commitWithVerification(path, tmpPath, fingerprint)
	}

	return atomicRename(path, tmpPath)
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

func atomicRename(path, tmpPath string) error {
	_ = os.Rename(path, path+".bak")

	err := os.Rename(tmpPath, path)
	if err != nil {
		return fmt.Errorf("renaming %s to %s: %w", tmpPath, path, err)
	}

	return nil
}

func cleanupTmp(tmpPath string) {
	_ = os.Remove(tmpPath)
}
