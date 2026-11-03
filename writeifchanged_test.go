package atomicwrite

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteIfChanged_NewFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	changed, err := WriteIfChanged(path, []byte(`{"key":"value"}`))
	if err != nil {
		t.Fatalf("WriteIfChanged new file: %v", err)
	}

	if !changed {
		t.Error("expected changed=true for new file")
	}

	assertFileContent(t, path, `{"key":"value"}`)
}

func TestWriteIfChanged_SameContent(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "unchanged")

	info, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatalf("stat before: %v", statErr)
	}

	originalTime := info.ModTime()

	changed, err := WriteIfChanged(path, []byte("unchanged"))
	if err != nil {
		t.Fatalf("WriteIfChanged same content: %v", err)
	}

	if changed {
		t.Error("expected changed=false for identical content")
	}

	assertFileContent(t, path, "unchanged")

	updatedInfo, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatalf("stat after: %v", statErr)
	}

	if !originalTime.Equal(updatedInfo.ModTime()) {
		t.Error("mtime should not change when content is skipped")
	}
}

func TestWriteIfChanged_DifferentContent(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "old")

	changed, err := WriteIfChanged(path, []byte("new"))
	if err != nil {
		t.Fatalf("WriteIfChanged different content: %v", err)
	}

	if !changed {
		t.Error("expected changed=true for different content")
	}

	assertFileContent(t, path, "new")
}

func TestWriteIfChanged_EmptyFile_EmptyData(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "")

	changed, err := WriteIfChanged(path, []byte{})
	if err != nil {
		t.Fatalf("WriteIfChanged empty: %v", err)
	}

	if changed {
		t.Error("expected changed=false for identical empty content")
	}
}

func TestWriteIfChanged_NewFile_EmptyData(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")

	changed, err := WriteIfChanged(path, []byte{})
	if err != nil {
		t.Fatalf("WriteIfChanged new empty file: %v", err)
	}

	if !changed {
		t.Error("expected changed=true for new file even with empty data")
	}

	assertFileContent(t, path, "")
}

func TestWriteIfChanged_PreservesPermissions(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	err := os.Chmod(path, 0o600)
	if err != nil {
		t.Fatalf("chmod: %v", err)
	}

	changed, writeErr := WriteIfChanged(path, []byte("updated"))
	if writeErr != nil {
		t.Fatalf("WriteIfChanged: %v", writeErr)
	}

	if !changed {
		t.Error("expected changed=true for different content")
	}

	info, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatalf("stat: %v", statErr)
	}

	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("permissions = %o, want 0o600", perm)
	}
}

func TestWriteIfChanged_LeavesNoLeftoverFiles(t *testing.T) {
	t.Parallel()

	path := tempFile(t, "original")

	changed, err := WriteIfChanged(path, []byte("updated"))
	if err != nil {
		t.Fatalf("WriteIfChanged: %v", err)
	}

	if !changed {
		t.Error("expected changed=true for different content")
	}

	tmpMatches, globErr := filepath.Glob(filepath.Join(filepath.Dir(path), "*.tmp"))
	if globErr != nil {
		t.Fatalf("glob: %v", globErr)
	}

	if len(tmpMatches) > 0 {
		t.Errorf("expected no temp files after write, found: %v", tmpMatches)
	}
}
