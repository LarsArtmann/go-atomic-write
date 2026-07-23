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

func TestWriteFunc_FirstRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "output.json")

	err := WriteFunc(path, func(w io.Writer) error {
		_, err := fmt.Fprintf(w, `{"hello":"world"}`)

		return err
	}, Fingerprint{})
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
		_, err := w.Write([]byte("first"))

		return err
	}, Fingerprint{})
	if err != nil {
		t.Fatalf("first WriteFunc: %v", err)
	}

	err = WriteFunc(path, func(w io.Writer) error {
		_, err := w.Write([]byte("second"))

		return err
	}, Fingerprint{})
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
			_, err := w.Write([]byte(chunk))
			if err != nil {
				return err
			}
		}

		return nil
	}, Fingerprint{})
	if err != nil {
		t.Fatalf("WriteFunc failed: %v", err)
	}

	data, err := os.ReadFile(path) //nolint:gosec // test
	if err != nil {
		t.Fatal(err)
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

	wantErr := errors.New("boom")

	err := WriteFunc(path, func(w io.Writer) error {
		return wantErr
	}, Fingerprint{})
	if err == nil {
		t.Fatal("expected error from callback")
	}

	if _, statErr := os.Stat(path); statErr == nil {
		t.Error("temp file should be cleaned up on callback error")
	}

	if _, statErr := os.Stat(path + ".*.tmp"); statErr == nil {
		t.Error("temp file should not linger")
	}
}

func TestWriteFunc_PreservesPermissions(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "perm.txt")

	err := WriteFunc(path, func(w io.Writer) error {
		_, err := w.Write([]byte("first"))

		return err
	}, Fingerprint{})
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chmod(path, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	err = WriteFunc(path, func(w io.Writer) error {
		_, err := w.Write([]byte("second"))

		return err
	}, Fingerprint{})
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
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
		_, err := w.Write([]byte("data"))

		return err
	}, Fingerprint{})
	if err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			t.Errorf("leftover temp file: %s", entry.Name())
		}
	}
}
