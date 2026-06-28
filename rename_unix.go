//go:build !windows

package atomicwrite

import (
	"fmt"
	"os"
	"path/filepath"
)

func atomicRename(path, tmpPath string) error {
	err := os.Rename(tmpPath, path)
	if err != nil {
		return fmt.Errorf("renaming %s to %s: %w", tmpPath, path, err)
	}

	return syncDir(filepath.Dir(path))
}

func syncDir(dir string) error {
	dirFile, err := os.Open(dir) //nolint:gosec // dir derived from caller-controlled path
	if err != nil {
		return fmt.Errorf("opening directory %s for sync: %w", dir, err)
	}

	defer func() { _ = dirFile.Close() }()

	syncErr := dirFile.Sync()
	if syncErr != nil {
		return fmt.Errorf("syncing directory %s: %w", dir, syncErr)
	}

	return nil
}
