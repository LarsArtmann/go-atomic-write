package atomicwrite

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var errTestCallback = errors.New("test callback failure")

func TestWriteFunc_FirstRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "output.json")

	err := WriteFunc(path, func(w io.Writer) error {
		_, writeErr := fmt.Fprintf(w, `{"hello":"world"}`)
		if writeErr != nil {
			return fmt.Errorf("writing content: %w", writeErr)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("WriteFunc failed: %v", err)
	}

	assertFileContent(t, path, `{"hello":"world"}`)
}

func TestWriteFunc_Overwrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "output.json")

	err := WriteFunc(path, func(w io.Writer) error {
		_, writeErr := w.Write([]byte("first"))
		if writeErr != nil {
			return fmt.Errorf("writing content: %w", writeErr)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("first WriteFunc: %v", err)
	}

	err = WriteFunc(path, func(w io.Writer) error {
		_, writeErr := w.Write([]byte("second"))
		if writeErr != nil {
			return fmt.Errorf("writing content: %w", writeErr)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("second WriteFunc: %v", err)
	}

	assertFileContent(t, path, "second")
}

func TestWriteFunc_LargeStream(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "large.txt")

	chunk := strings.Repeat("A", 4096)
	wantChunks := 100

	err := WriteFunc(path, func(w io.Writer) error {
		for range wantChunks {
			_, writeErr := w.Write([]byte(chunk))
			if writeErr != nil {
				return fmt.Errorf("writing chunk: %w", writeErr)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatalf("WriteFunc failed: %v", err)
	}

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}

	want := strings.Repeat(chunk, wantChunks)
	if string(data) != want {
		t.Errorf("expected %d bytes, got %d", len(want), len(data))
	}
}

func TestWriteFunc_CallbackError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "failed.json")

	err := WriteFunc(path, func(_ io.Writer) error {
		return errTestCallback
	})
	if err == nil {
		t.Fatal("expected error from callback")
	}

	_, statErr := os.Stat(path)
	if statErr == nil {
		t.Error("temp file should be cleaned up on callback error")
	}

	matches, globErr := filepath.Glob(path + ".*.tmp")
	if globErr == nil && len(matches) > 0 {
		t.Error("temp file should not linger")
	}
}

func TestWriteFunc_PreservesPermissions(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "perm.txt")

	err := WriteFunc(path, func(w io.Writer) error {
		_, writeErr := w.Write([]byte("first"))
		if writeErr != nil {
			return fmt.Errorf("writing content: %w", writeErr)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	chmodErr := os.Chmod(path, 0o600)
	if chmodErr != nil {
		t.Fatal(chmodErr)
	}

	err = WriteFunc(path, func(w io.Writer) error {
		_, writeErr := w.Write([]byte("second"))
		if writeErr != nil {
			return fmt.Errorf("writing content: %w", writeErr)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	info, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatal(statErr)
	}

	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected 0600, got %o", info.Mode().Perm())
	}
}

func TestWriteFunc_LeavesNoLeftoverFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "clean.txt")

	err := WriteFunc(path, func(w io.Writer) error {
		_, writeErr := w.Write([]byte("data"))
		if writeErr != nil {
			return fmt.Errorf("writing content: %w", writeErr)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		t.Fatal(readErr)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			t.Errorf("leftover temp file: %s", entry.Name())
		}
	}
}
