package atomicwrite

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
)

func tempFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")

	err := os.WriteFile(path, []byte(content), 0o644) //nolint:gosec // test fixture uses 0644 inside t.TempDir
	if err != nil {
		t.Fatal(err)
	}

	return path
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	data, err := os.ReadFile(path) //nolint:gosec // test reads from t.TempDir
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if string(data) != want {
		t.Errorf("content = %q, want %q", string(data), want)
	}
}

func TestFingerprintFromBytes(t *testing.T) {
	t.Parallel()

	data := []byte("hello world")
	fingerprint := FingerprintFromBytes(data)

	if fingerprint.IsZero() {
		t.Error("fingerprint should not be zero for non-empty data")
	}

	if !fingerprint.Matches(data) {
		t.Error("fingerprint should match the same data")
	}

	if fingerprint.Matches([]byte("different")) {
		t.Error("fingerprint should not match different data")
	}
}

func TestFingerprintIsZero(t *testing.T) {
	t.Parallel()

	var zero Fingerprint
	if !zero.IsZero() {
		t.Error("zero-value Fingerprint should be zero")
	}

	nonZero := FingerprintFromBytes([]byte("data"))
	if nonZero.IsZero() {
		t.Error("fingerprint from data should not be zero")
	}
}

func TestFingerprintFile(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "content")

	fingerprint, err := FingerprintFile(path)
	if err != nil {
		t.Fatalf("FingerprintFile: %v", err)
	}

	if fingerprint.IsZero() {
		t.Error("fingerprint should not be zero for existing file")
	}

	if !fingerprint.Matches([]byte("content")) {
		t.Error("fingerprint should match file content")
	}
}

func TestFingerprintFileNonexistent(t *testing.T) {
	t.Parallel()

	fingerprint, err := FingerprintFile("/nonexistent/path/file")
	if err != nil {
		t.Fatalf("FingerprintFile nonexistent: %v", err)
	}

	if !fingerprint.IsZero() {
		t.Error("nonexistent file should return zero fingerprint")
	}
}

func TestWriteFirstRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")

	err := Write(path, []byte("hello"), Fingerprint{})
	if err != nil {
		t.Fatalf("Write first run: %v", err)
	}

	assertFileContent(t, path, "hello")

	_, statErr := os.Stat(path + ".bak")
	if !os.IsNotExist(statErr) {
		t.Error("expected no .bak file on first run")
	}
}

func TestWriteWithFingerprint(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	fingerprint := FingerprintFromBytes([]byte("original"))

	err := Write(path, []byte("updated"), fingerprint)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	assertFileContent(t, path, "updated")
}

func TestWriteRejectsConcurrentModification(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	fingerprint := FingerprintFromBytes([]byte("original"))

	modified := "modified"

	err := os.WriteFile(path, []byte(modified), 0o644) //nolint:gosec // test fixture in t.TempDir
	if err != nil {
		t.Fatalf("modify file: %v", err)
	}

	writeErr := Write(path, []byte("updated"), fingerprint)
	if writeErr == nil {
		t.Fatal("expected error for concurrent modification, got nil")
	}

	if !errors.Is(writeErr, ErrConcurrentModification) {
		t.Errorf("expected ErrConcurrentModification, got: %v", writeErr)
	}

	_, statErr := os.Stat(path + ".tmp")
	if statErr == nil {
		t.Error("temp file should be cleaned up on error")
	}
}

func TestWritePreservesPermissions(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	err := os.Chmod(path, 0o600)
	if err != nil {
		t.Fatalf("chmod: %v", err)
	}

	fingerprint := FingerprintFromBytes([]byte("original"))

	writeErr := Write(path, []byte("updated"), fingerprint)
	if writeErr != nil {
		t.Fatalf("Write: %v", writeErr)
	}

	info, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatalf("stat: %v", statErr)
	}

	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("permissions = %o, want 0o600", perm)
	}
}

func TestWriteCreatesBackup(t *testing.T) {
	t.Parallel()

	originalContent := "original"
	path := tempFile(t, originalContent)

	fingerprint := FingerprintFromBytes([]byte(originalContent))

	err := Write(path, []byte("updated"), fingerprint)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	assertFileContent(t, path+".bak", originalContent)
}

func TestTempFileCleanedUpOnError(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	fingerprint := FingerprintFromBytes([]byte("original"))

	err := os.Remove(path)
	if err != nil {
		t.Fatalf("remove original: %v", err)
	}

	writeErr := Write(path, []byte("updated"), fingerprint)
	if writeErr == nil {
		t.Fatal("expected error when original file deleted during verification, got nil")
	}

	_, statErr := os.Stat(filepath.Dir(path) + "/testfile.tmp")
	if statErr == nil {
		t.Error("temp file should be cleaned up on error")
	}
}

func TestConcurrentWriteRACE(t *testing.T) {
	t.Parallel()

	content := "original"
	path := tempFile(t, content)

	var successes, conflicts atomic.Int32

	var waitGroup sync.WaitGroup

	const writers = 5

	for range writers {
		waitGroup.Go(func() {
			fingerprint, fpErr := FingerprintFile(path)
			if fpErr != nil {
				t.Logf("FingerprintFile: %v", fpErr)

				return
			}

			writeErr := Write(path, []byte("updated"), fingerprint)
			if writeErr == nil {
				successes.Add(1)
			} else {
				conflicts.Add(1)
			}
		})
	}

	waitGroup.Wait()

	total := successes.Load() + conflicts.Load()
	if int(total) != writers {
		t.Errorf("expected %d total outcomes, got %d (successes=%d, conflicts=%d)",
			writers, total, successes.Load(), conflicts.Load())
	}

	if successes.Load() < 1 {
		t.Error("expected at least one successful write")
	}

	if conflicts.Load() < 1 {
		t.Log("no conflicts detected — race window may be too narrow with only 5 writers")
	}
}

func TestAtomicRenameReportsErrorOnFailure(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	ghostPath := filepath.Join(dir, "nonexistent", "dir", "testfile")
	tmpPath := filepath.Join(dir, "temp-file")

	err := os.WriteFile(tmpPath, []byte("content"), 0o644) //nolint:gosec // test fixture in t.TempDir
	if err != nil {
		t.Fatalf("write temp: %v", err)
	}

	renameErr := atomicRename(ghostPath, tmpPath)
	if renameErr == nil {
		t.Fatal("expected error for rename into nonexistent directory, got nil")
	}
}
