//go:build windows

package atomicwrite

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

func atomicRename(path, tmpPath string) error {
	const maxRetries = 5
	const baseBackoff = time.Millisecond

	var lastErr error

	for attempt := range maxRetries {
		lastErr = os.Rename(tmpPath, path)
		if lastErr == nil {
			return nil
		}

		if !isRetryableRename(lastErr) {
			break
		}

		time.Sleep(baseBackoff << uint(attempt))
	}

	return fmt.Errorf("renaming %s to %s: %w", tmpPath, path, lastErr)
}

func isRetryableRename(err error) bool {
	var linkErr *os.LinkError
	if !errors.As(err, &linkErr) {
		return false
	}

	var errno syscall.Errno
	if !errors.As(linkErr.Err, &errno) {
		return false
	}

	return errno == syscall.ERROR_ACCESS_DENIED || errno == errSharingViolation
}

// errSharingViolation is ERROR_SHARING_VIOLATION (WinAPI error code 32): a
// Windows rename fails with it while another process still holds the target
// open. Go's stdlib syscall package for Windows does not define this constant
// (only a small errno subset), so referencing it by name broke compilation for
// GOOS=windows. Defined here by value instead.
const errSharingViolation = syscall.Errno(32)
