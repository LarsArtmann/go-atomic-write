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

	err := os.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	return path
}

func TestFingerprintFromBytes(t *testing.T) {
	t.Parallel()

	data := []byte("hello world")
	fp := FingerprintFromBytes(data)

	if fp.IsZero() {
		t.Error("fingerprint should not be zero for non-empty data")
	}

	if !fp.Matches(data) {
		t.Error("fingerprint should match the same data")
	}

	if fp.Matches([]byte("different")) {
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

	fp, err := FingerprintFile(path)
	if err != nil {
		t.Fatalf("FingerprintFile: %v", err)
	}

	if fp.IsZero() {
		t.Error("fingerprint should not be zero for existing file")
	}

	if !fp.Matches([]byte("content")) {
		t.Error("fingerprint should match file content")
	}
}

func TestFingerprintFileNonexistent(t *testing.T) {
	t.Parallel()

	fp, err := FingerprintFile("/nonexistent/path/file")
	if err != nil {
		t.Fatalf("FingerprintFile nonexistent: %v", err)
	}

	if !fp.IsZero() {
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

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("reading file: %v", readErr)
	}

	if string(data) != "hello" {
		t.Errorf("content = %q, want %q", string(data), "hello")
	}

	if _, statErr := os.Stat(path + ".bak"); !os.IsNotExist(statErr) {
		t.Error("expected no .bak file on first run")
	}
}

func TestWriteWithFingerprint(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	fp := FingerprintFromBytes([]byte("original"))

	err := Write(path, []byte("updated"), fp)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("reading file: %v", readErr)
	}

	if string(data) != "updated" {
		t.Errorf("content = %q, want %q", string(data), "updated")
	}
}

func TestWriteRejectsConcurrentModification(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	fp := FingerprintFromBytes([]byte("original"))

	modified := "modified"
	if err := os.WriteFile(path, []byte(modified), 0o644); err != nil {
		t.Fatalf("modify file: %v", err)
	}

	err := Write(path, []byte("updated"), fp)
	if err == nil {
		t.Fatal("expected error for concurrent modification, got nil")
	}

	if !errors.Is(err, ErrConcurrentModification) {
		t.Errorf("expected ErrConcurrentModification, got: %v", err)
	}

	if _, statErr := os.Stat(path + ".tmp"); statErr == nil {
		t.Error("temp file should be cleaned up on error")
	}
}

func TestWritePreservesPermissions(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatalf("chmod: %v", err)
	}

	fp := FingerprintFromBytes([]byte("original"))

	err := Write(path, []byte("updated"), fp)
	if err != nil {
		t.Fatalf("Write: %v", err)
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

	fp := FingerprintFromBytes([]byte(originalContent))

	err := Write(path, []byte("updated"), fp)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	backup, readErr := os.ReadFile(path + ".bak")
	if readErr != nil {
		t.Fatalf("read backup: %v", readErr)
	}

	if string(backup) != originalContent {
		t.Errorf("backup = %q, want %q", string(backup), originalContent)
	}
}

func TestTempFileCleanedUpOnError(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	fp := FingerprintFromBytes([]byte("original"))

	if err := os.Remove(path); err != nil {
		t.Fatalf("remove original: %v", err)
	}

	err := Write(path, []byte("updated"), fp)
	if err == nil {
		t.Fatal("expected error when original file deleted during verification, got nil")
	}

	if _, statErr := os.Stat(filepath.Dir(path) + "/testfile.tmp"); statErr == nil {
		t.Error("temp file should be cleaned up on error")
	}
}

func TestConcurrentWriteRACE(t *testing.T) {
	content := "original"
	path := tempFile(t, content)

	var successes, conflicts atomic.Int32

	var wg sync.WaitGroup

	const writers = 5

	for range writers {
		wg.Go(func() {
			fp, fpErr := FingerprintFile(path)
			if fpErr != nil {
				t.Logf("FingerprintFile: %v", fpErr)

				return
			}

			writeErr := Write(path, []byte("updated"), fp)
			if writeErr == nil {
				successes.Add(1)
			} else {
				conflicts.Add(1)
			}
		})
	}

	wg.Wait()

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

	err := os.WriteFile(tmpPath, []byte("content"), 0o644)
	if err != nil {
		t.Fatalf("write temp: %v", err)
	}

	renameErr := atomicRename(ghostPath, tmpPath)
	if renameErr == nil {
		t.Fatal("expected error for rename into nonexistent directory, got nil")
	}
}
