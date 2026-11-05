package atomicwrite

import (
	"errors"
	"os"
	"path/filepath"
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

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if string(data) != want {
		t.Errorf("content = %q, want %q", string(data), want)
	}
}

func TestWriteFirstRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "testfile")

	err := Write(path, []byte("hello"))
	if err != nil {
		t.Fatalf("Write first run: %v", err)
	}

	assertFileContent(t, path, "hello")

	_, statErr := os.Stat(path + ".bak")
	if !os.IsNotExist(statErr) {
		t.Error("expected no .bak file on first run")
	}
}

func TestWriteVerifiedFirstRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "newfile")

	err := WriteVerified(path, []byte("first"), Fingerprint{})
	if err != nil {
		t.Fatalf("WriteVerified first run: %v", err)
	}

	assertFileContent(t, path, "first")
}

func TestWriteVerifiedZeroFingerprintRejectsExistingFile(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "already here")

	err := WriteVerified(path, []byte("overwrite"), Fingerprint{})
	if !errors.Is(err, ErrConcurrentModification) {
		t.Fatalf("expected ErrConcurrentModification for zero fingerprint on existing file, got %v", err)
	}

	assertFileContent(t, path, "already here")
}

func TestWriteWithFingerprint(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	fingerprint := FingerprintFromBytes([]byte("original"))

	err := WriteVerified(path, []byte("updated"), fingerprint)
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

	err := os.WriteFile(path, []byte(modified), 0o644)
	if err != nil {
		t.Fatalf("modify file: %v", err)
	}

	writeErr := WriteVerified(path, []byte("updated"), fingerprint)
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

	writeErr := WriteVerified(path, []byte("updated"), fingerprint)
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

	err := WriteVerified(path, []byte("updated"), fingerprint)
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

	writeErr := WriteVerified(path, []byte("updated"), fingerprint)
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

func TestWriteWithPerm_SetsPermOnCreate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "secret")

	if err := WriteWithPerm(path, []byte("sensitive"), 0o600); err != nil {
		t.Fatalf("WriteWithPerm: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("perm = %o, want 600", got)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "sensitive" {
		t.Errorf("content = %q, want %q", data, "sensitive")
	}
}

func TestWriteWithPerm_OverridesExistingPerm(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "old")

	if err := WriteWithPerm(path, []byte("new"), 0o600); err != nil {
		t.Fatalf("WriteWithPerm: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("perm = %o, want 600 (explicit perm takes precedence)", got)
	}
}

func TestWriteWithPerm_LeavesNoTempFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "file")

	if err := WriteWithPerm(path, []byte("data"), 0o644); err != nil {
		t.Fatalf("WriteWithPerm: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 || entries[0].Name() != "file" {
		names := make([]string, 0, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}

		t.Errorf("dir contains leftover temp files: %v", names)
	}
}

func TestWriteCreatesMissingParentDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "cache", "nested", "file.json")

	err := Write(path, []byte("hello"))
	if err != nil {
		t.Fatalf("Write into a missing parent dir: %v", err)
	}

	assertFileContent(t, path, "hello")
}

func TestWriteVerifiedCreatesMissingParentDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "cache", "file.json")

	err := WriteVerified(path, []byte("first"), Fingerprint{})
	if err != nil {
		t.Fatalf("WriteVerified into a missing parent dir: %v", err)
	}

	assertFileContent(t, path, "first")
}

func TestWriteWithPermCreatesMissingParentDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "conf", "file.json")

	err := WriteWithPerm(path, []byte("hello"), 0o600)
	if err != nil {
		t.Fatalf("WriteWithPerm into a missing parent dir: %v", err)
	}

	assertFileContent(t, path, "hello")

	info, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatalf("stat written file: %v", statErr)
	}

	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("perm = %o, want 600", got)
	}
}
