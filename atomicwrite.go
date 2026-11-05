// Package atomicwrite provides TOCTOU-safe file writes using xxhash64
// fingerprint verification, cross-platform file locking via gofrs/flock,
// atomic rename, and fsync for crash durability.
//
// Missing parent directories of the target path are created on demand.
package atomicwrite

import (
	"io"
	"io/fs"
)

// Write writes data to path with crash durability (fsync + atomic rename).
// It does NOT perform TOCTOU verification — use WriteVerified for that.
// Use Write when concurrent modification is not a concern (e.g., temp files,
// first-time creation, or single-writer scenarios).
func Write(path string, data []byte) error {
	tmpPath, perm, err := prepareTempPath(path)
	if err != nil {
		return err
	}

	stageErr := writeAndSync(tmpPath, data, perm)
	if stageErr != nil {
		return stageErr
	}

	return commitRename(path, tmpPath)
}

// WriteWithPerm writes data to path with an explicit file permission and
// crash durability (fsync + atomic rename). The temp file is created with the
// requested permission, so the final file appears atomically with the correct
// mode — there is no chmod-after-rename window in which the file is briefly
// visible with the wrong permissions.
//
// The explicit permission takes precedence over an existing file's current
// mode. Use Write instead to preserve an existing file's mode (or the 0o644
// default for new files).
//
// It does NOT perform TOCTOU verification — use WriteVerified for that.
func WriteWithPerm(path string, data []byte, perm fs.FileMode) error {
	tmpPath, err := tempPathFor(path)
	if err != nil {
		return err
	}

	stageErr := writeAndSync(tmpPath, data, perm)
	if stageErr != nil {
		return stageErr
	}

	return commitRename(path, tmpPath)
}

// WriteVerified writes data to path with TOCTOU protection and crash durability.
// Data is staged to a unique temp file, fsync'd, then verified against the
// fingerprint before atomic rename.
//
// The fingerprint must be computed via FingerprintFile before reading/modifying
// the file. A zero-value fingerprint indicates the file should not yet exist;
// the write will fail with ErrConcurrentModification if another process creates
// it first. This prevents the silent-skip footgun where a caller forgets to
// compute a fingerprint.
func WriteVerified(path string, data []byte, fingerprint Fingerprint) error {
	tmpPath, perm, err := prepareTempPath(path)
	if err != nil {
		return err
	}

	stageErr := writeAndSync(tmpPath, data, perm)
	if stageErr != nil {
		return stageErr
	}

	return commitVerified(path, tmpPath, fingerprint)
}

// WriteIfChanged writes data to path only if it differs from the current
// on-disk content. Returns changed=true if the file was written, false if
// the content was identical and the write was skipped.
//
// This is the idiomatic primitive for config-file writers and code generators
// that must not produce spurious diffs on re-runs: no content change means
// no file mutation, no mtime bump, no file-watcher trigger.
//
// Race-safe for existing files: if another process modifies the file between
// the fingerprint check and the atomic rename, returns
// ErrConcurrentModification. First-write (file does not exist) uses a plain
// atomic write — there is no prior content to protect.
func WriteIfChanged(path string, data []byte) (bool, error) {
	existing, err := FingerprintFile(path)
	if err != nil {
		return false, err
	}

	// First write: file does not exist yet — no content to compare or protect.
	if existing.IsZero() {
		err = Write(path, data)

		return true, err
	}

	// Content unchanged — skip the write entirely.
	if FingerprintFromBytes(data) == existing {
		return false, nil
	}

	// Content differs — verified write catches concurrent modification.
	err = WriteVerified(path, data, existing)
	if err != nil {
		return false, err
	}

	return true, nil
}

// WriteFunc writes to path via a streaming callback with crash durability.
// The callback receives a buffered writer (64KB buffer) and may stream content
// incrementally without holding the full payload in memory.
// It does NOT perform TOCTOU verification — use WriteFuncVerified for that.
//
// Use WriteFunc instead of Write when the content is large or produced
// incrementally (e.g., JSON encoders, diagram renderers).
func WriteFunc(path string, fn func(w io.Writer) error) error {
	tmpPath, perm, err := prepareTempPath(path)
	if err != nil {
		return err
	}

	stageErr := writeFuncAndSync(tmpPath, fn, perm)
	if stageErr != nil {
		return stageErr
	}

	return commitRename(path, tmpPath)
}

// WriteFuncWithPerm writes to path via a streaming callback with an explicit
// file permission and crash durability (fsync + atomic rename). See
// WriteWithPerm for permission semantics and WriteFunc for streaming behavior.
// The temp file is created with the requested permission, so the final file
// appears atomically with the correct mode — there is no chmod-after-rename
// window in which the file is briefly visible with the wrong permissions.
//
// The explicit permission takes precedence over an existing file's current
// mode. Use WriteFunc instead to preserve an existing file's mode (or the
// 0o644 default for new files).
//
// It does NOT perform TOCTOU verification — use WriteFuncVerified for that.
func WriteFuncWithPerm(path string, fn func(io.Writer) error, perm fs.FileMode) error {
	tmpPath, err := tempPathFor(path)
	if err != nil {
		return err
	}

	stageErr := writeFuncAndSync(tmpPath, fn, perm)
	if stageErr != nil {
		return stageErr
	}

	return commitRename(path, tmpPath)
}

// WriteFuncVerified writes to path via a streaming callback with TOCTOU
// protection and crash durability. See WriteVerified for fingerprint semantics.
func WriteFuncVerified(path string, fn func(w io.Writer) error, fingerprint Fingerprint) error {
	tmpPath, perm, err := prepareTempPath(path)
	if err != nil {
		return err
	}

	stageErr := writeFuncAndSync(tmpPath, fn, perm)
	if stageErr != nil {
		return stageErr
	}

	return commitVerified(path, tmpPath, fingerprint)
}
