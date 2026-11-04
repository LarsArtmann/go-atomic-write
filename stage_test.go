package atomicwrite

import (
	"os"
	"strings"
	"sync"
	"testing"
)

// TestTempPathsNeverCollide pins the invariant the unique-suffix design exists
// for: concurrent writers staging to the SAME target never compute the same
// temp path, or one writer's staging file could be truncated (O_TRUNC) by
// the other before either rename lands.
func TestTempPathsNeverCollide(t *testing.T) {
	t.Parallel()

	const writers = 200

	path := "/mnt/data/config.json"

	var mutex sync.Mutex
	paths := make([]string, writers)
	seen := make(map[string]bool, writers)
	collisions := 0

	var waitGroup sync.WaitGroup

	for writer := range writers {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			tmpPath, err := tempPathFor(path)
			if err != nil {
				t.Errorf("tempPathFor: %v", err)

				return
			}

			mutex.Lock()
			defer mutex.Unlock()

			if seen[tmpPath] {
				collisions++
			}
			seen[tmpPath] = true
			paths[writer] = tmpPath
		}()
	}
	waitGroup.Wait()

	if collisions != 0 {
		t.Errorf("found %d temp path collisions across %d concurrent writers", collisions, writers)
	}

	for _, tmpPath := range paths {
		if !strings.HasPrefix(tmpPath, path+".") || !strings.HasSuffix(tmpPath, ".tmp") {
			t.Errorf("temp path %s does not follow the <path>.<suffix>.tmp layout", tmpPath)
		}
	}
}

// TestPrepareTempPathPreservesMode pins the permission-preservation half of
// prepareTempPath: an existing file's mode is carried into the staging file,
// a fresh target gets the 0644 default.
func TestPrepareTempPathPreservesMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	existing := dir + "/existing"

	//nolint:gosec // fixture in t.TempDir, fixed content
	writeErr := os.WriteFile(existing, []byte("x"), 0o600)
	if writeErr != nil {
		t.Fatalf("seed file: %v", writeErr)
	}

	_, perm, err := prepareTempPath(existing)
	if err != nil {
		t.Fatalf("prepareTempPath: %v", err)
	}

	if perm != 0o600 {
		t.Errorf("existing file's mode 0600 should be preserved, got %o", perm)
	}

	fresh := dir + "/fresh"
	_, perm, err = prepareTempPath(fresh)
	if err != nil {
		t.Fatalf("prepareTempPath: %v", err)
	}

	if perm != 0o644 {
		t.Errorf("fresh file should default to 0644, got %o", perm)
	}
}
