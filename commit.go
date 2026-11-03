package atomicwrite

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofrs/flock"
)

// ErrConcurrentModification indicates the file was modified by another
// process between the fingerprint read and the write attempt.
var ErrConcurrentModification = errors.New("file was modified concurrently since read")

// commitVerified performs TOCTOU verification. For a zero fingerprint
// (first-write), the target file must not exist; the check runs BEFORE
// acquiring the flock because gofrs/flock creates the file via O_CREATE,
// which would cause a false "created concurrently" conflict. First-write
// uses a plain atomic rename — there is no prior content to protect.
func commitVerified(path, tmpPath string, fingerprint Fingerprint) error {
	if fingerprint.IsZero() {
		if _, err := os.Stat(path); err == nil {
			cleanupTmp(tmpPath)

			return fmt.Errorf("%w: %s was created concurrently; removed tmp file %s",
				ErrConcurrentModification, path, tmpPath)
		} else if !os.IsNotExist(err) {
			cleanupTmp(tmpPath)

			return fmt.Errorf("stating %s for verification (tmp file %s): %w", path, tmpPath, err)
		}

		return commitRename(path, tmpPath)
	}

	fileLock := flock.New(path)

	lockErr := fileLock.Lock()
	if lockErr != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("acquiring exclusive lock on %s (tmp file %s): %w", path, tmpPath, lockErr)
	}
	defer func() { _ = fileLock.Close() }()

	current, err := os.ReadFile(path)
	if err != nil {
		cleanupTmp(tmpPath)

		return fmt.Errorf("re-reading %s for verification (tmp file %s): %w", path, tmpPath, err)
	}

	if !fingerprint.Matches(current) {
		cleanupTmp(tmpPath)

		return fmt.Errorf("%w: %s was modified since read; removed tmp file %s",
			ErrConcurrentModification, path, tmpPath)
	}

	return commitRename(path, tmpPath)
}

// commitRename renames the staged temp file onto the target path. When the
// rename fails, the staged file is removed so a failed commit leaves no .tmp
// residue next to the target. The rename error is returned unchanged because
// the platform implementations already wrap it with both paths.
func commitRename(path, tmpPath string) error {
	renameErr := atomicRename(path, tmpPath)
	if renameErr != nil {
		cleanupTmp(tmpPath)

		return renameErr
	}

	return nil
}
