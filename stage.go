package atomicwrite

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// prepareTempPath computes the file mode (preserving the existing file's
// mode or using the default) and returns a unique temp file path alongside it.
func prepareTempPath(path string) (string, fs.FileMode, error) {
	const defaultFilePerm = fs.FileMode(0o644)

	perm := defaultFilePerm

	info, err := os.Stat(path)
	if err == nil {
		perm = info.Mode().Perm()
	}

	tmpPath, tmpErr := tempPathFor(path)
	if tmpErr != nil {
		return "", 0, tmpErr
	}

	return tmpPath, perm, nil
}

// tempPathFor returns a unique temp file path alongside the target file.
func tempPathFor(path string) (string, error) {
	suffix, suffixErr := randomSuffix()
	if suffixErr != nil {
		return "", fmt.Errorf("generating temp file suffix: %w", suffixErr)
	}

	return path + "." + suffix + ".tmp", nil
}

// randomSuffix is a var so tests can inject a fixed suffix and pin the
// temp-path format (and the collision precondition) deterministically.
var randomSuffix = func() (string, error) { //nolint:gochecknoglobals // deliberate test seam, not mutable state in production paths
	var buf [4]byte

	_, err := rand.Read(buf[:])
	if err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}

	return hex.EncodeToString(buf[:]), nil
}

// writeFuncBufferSize is the default buffer size for streaming writes.
const writeFuncBufferSize = 65536

func writeAndSync(tmpPath string, data []byte, perm fs.FileMode) error {
	file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("creating temp file %s: %w", tmpPath, err)
	}

	_, writeErr := file.Write(data)
	if writeErr != nil {
		_ = file.Close() //nolint:erraudit // best-effort close on a failed write; the write error is returned

		cleanupTmp(tmpPath)

		return fmt.Errorf("writing temp file %s: %w", tmpPath, writeErr)
	}

	syncErr := file.Sync()
	if syncErr != nil {
		_ = file.Close() //nolint:erraudit // best-effort close on a failed sync; the sync error is returned

		cleanupTmp(tmpPath)

		return fmt.Errorf("syncing temp file %s: %w", tmpPath, syncErr)
	}

	closeErr := file.Close()
	if closeErr != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("closing temp file %s: %w", tmpPath, closeErr)
	}

	return nil
}

func writeFuncAndSync(tmpPath string, fn func(io.Writer) error, perm fs.FileMode) error {
	file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("creating temp file %s: %w", tmpPath, err)
	}

	bw := bufio.NewWriterSize(file, writeFuncBufferSize)

	writeErr := fn(bw)

	flushErr := bw.Flush()

	syncErr := file.Sync()

	closeErr := file.Close()

	if writeErr != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("writing temp file %s: %w", tmpPath, writeErr)
	}

	if flushErr != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("flushing temp file %s: %w", tmpPath, flushErr)
	}

	if syncErr != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("syncing temp file %s: %w", tmpPath, syncErr)
	}

	if closeErr != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("closing temp file %s: %w", tmpPath, closeErr)
	}

	return nil
}

func cleanupTmp(tmpPath string) {
	_ = os.Remove(tmpPath) //nolint:erraudit // best-effort cleanup; the primary error is returned by callers
}
