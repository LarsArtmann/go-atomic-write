package atomicwrite

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

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

			writeErr := WriteVerified(path, []byte(payload), fingerprint)
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

	data, readErr := os.ReadFile(path)
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

	err := os.WriteFile(tmpPath, []byte("content"), 0o644)
	if err != nil {
		t.Fatalf("write temp: %v", err)
	}

	renameErr := atomicRename(ghostPath, tmpPath)
	if renameErr == nil {
		t.Fatal("expected error for rename into nonexistent directory, got nil")
	}
}

func TestRenameFailureCleansUpTempFile(t *testing.T) {
	t.Parallel()

	writers := map[string]func(path string) error{
		"Write": func(path string) error {
			return Write(path, []byte("updated"))
		},
		"WriteWithPerm": func(path string) error {
			return WriteWithPerm(path, []byte("updated"), 0o600)
		},
		"WriteFunc": func(path string) error {
			return WriteFunc(path, func(writer io.Writer) error {
				_, writeErr := writer.Write([]byte("updated"))

				return writeErr
			})
		},
	}

	for name, write := range writers {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			target := filepath.Join(dir, "target")

			// A non-empty directory at the target path makes the final
			// rename fail on every supported platform.
			mkdirErr := os.Mkdir(target, 0o700)
			if mkdirErr != nil {
				t.Fatalf("mkdir target: %v", mkdirErr)
			}

			innerErr := os.WriteFile(filepath.Join(target, "keep"), []byte("x"), 0o600)
			if innerErr != nil {
				t.Fatalf("write inner file: %v", innerErr)
			}

			writeErr := write(target)
			if writeErr == nil {
				t.Fatal("expected error when renaming onto a directory, got nil")
			}

			tmpMatches, globErr := filepath.Glob(filepath.Join(dir, "*.tmp"))
			if globErr != nil {
				t.Fatalf("glob: %v", globErr)
			}

			if len(tmpMatches) > 0 {
				t.Errorf("rename failure should remove the staged temp file, found: %v", tmpMatches)
			}
		})
	}
}
