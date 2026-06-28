package atomicwrite

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
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

func TestWriteLeavesNoLeftoverFiles(t *testing.T) {
	t.Parallel()

	originalContent := "original"
	path := tempFile(t, originalContent)

	fingerprint := FingerprintFromBytes([]byte(originalContent))

	err := Write(path, []byte("updated"), fingerprint)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	assertFileContent(t, path, "updated")

	_, bakErr := os.Stat(path + ".bak")
	if !os.IsNotExist(bakErr) {
		t.Error("expected no .bak file after write")
	}

	tmpMatches, globErr := filepath.Glob(filepath.Join(filepath.Dir(path), "*.tmp"))
	if globErr != nil {
		t.Fatalf("glob: %v", globErr)
	}

	if len(tmpMatches) > 0 {
		t.Errorf("expected no temp files after write, found: %v", tmpMatches)
	}
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

	tmpMatches, globErr := filepath.Glob(filepath.Join(filepath.Dir(path), "testfile.*.tmp"))
	if globErr != nil {
		t.Fatalf("glob: %v", globErr)
	}

	if len(tmpMatches) > 0 {
		t.Errorf("temp files should be cleaned up on error, found: %v", tmpMatches)
	}
}

func TestConcurrentWriteRACE(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	var successes, conflicts atomic.Int32

	var waitGroup sync.WaitGroup

	const writers = 10

	for writerIndex := range writers {
		waitGroup.Go(func() {
			fingerprint, fpErr := FingerprintFile(path)
			if fpErr != nil {
				t.Logf("FingerprintFile: %v", fpErr)

				return
			}

			payload := "writer-" + strconv.Itoa(writerIndex)

			writeErr := Write(path, []byte(payload), fingerprint)
			if writeErr == nil {
				successes.Add(1)
			} else if errors.Is(writeErr, ErrConcurrentModification) {
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
		t.Fatal("expected at least one successful write")
	}

	if conflicts.Load() < 1 {
		t.Log("no conflicts detected — race window may be too narrow")
	}

	data, readErr := os.ReadFile(path) //nolint:gosec // test reads from t.TempDir
	if readErr != nil {
		t.Fatalf("reading final file: %v", readErr)
	}

	content := string(data)
	validPayload := false

	for writerIndex := range writers {
		if content == "writer-"+strconv.Itoa(writerIndex) {
			validPayload = true

			break
		}
	}

	if !validPayload {
		t.Errorf("final content %q is not a valid writer payload — possible corruption", content)
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
