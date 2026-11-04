package atomicwrite

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

// tempPathFor builds temp names from a crypto/rand suffix, so two concurrent
// writers derive the same path only in the astronomically unlikely collision
// case. These tests pin the contract from both sides: the suffix source makes
// real invocations unique, and the (injected) collision case is deterministic
// and visible rather than silently interleaved.

var errEntropyExhausted = errors.New("entropy exhausted")

func TestTempPathForUniqueAcrossWriters(t *testing.T) {
	t.Parallel()

	const writers = 512

	path := tempFile(t, "seed")

	paths := make([]string, writers) //nolint:makezero // indexed by writer id, never appended to

	var waitGroup sync.WaitGroup

	for writerIndex := range writers {
		waitGroup.Go(func() {
			tmpPath, tmpErr := tempPathFor(path)
			if tmpErr != nil {
				t.Errorf("tempPathFor: %v", tmpErr)

				return
			}

			paths[writerIndex] = tmpPath
		})
	}

	waitGroup.Wait()

	seen := make(map[string]bool, writers)

	for _, tmpPath := range paths {
		if seen[tmpPath] {
			t.Fatalf("duplicate temp path across %d writers: %s", writers, tmpPath)
		}

		seen[tmpPath] = true

		if !strings.HasPrefix(tmpPath, path+".") || !strings.HasSuffix(tmpPath, ".tmp") {
			t.Fatalf("temp path %s does not follow the path.<suffix>.tmp format", tmpPath)
		}
	}
}

// No t.Parallel: this test swaps the package-level randomSuffix seam, and
// sequential top-level tests finish before any parallel batch starts, so the
// mutation cannot race the parallel readers of tempPathFor.
func TestTempPathForFixedSuffixIsDeterministic(t *testing.T) { //nolint:paralleltest // mutates shared seam
	originalSuffix := randomSuffix

	t.Cleanup(func() { randomSuffix = originalSuffix })

	randomSuffix = func() (string, error) {
		return "cafe0000", nil
	}

	path := "/data/config.json"

	first, firstErr := tempPathFor(path)
	if firstErr != nil {
		t.Fatalf("tempPathFor: %v", firstErr)
	}

	second, secondErr := tempPathFor(path)
	if secondErr != nil {
		t.Fatalf("tempPathFor: %v", secondErr)
	}

	want := path + ".cafe0000.tmp"

	if first != want || second != want {
		t.Fatalf("fixed suffix must give the deterministic path %s, got %q and %q", want, first, second)
	}

	randomSuffix = func() (string, error) {
		return "", errEntropyExhausted
	}

	_, errorCase := tempPathFor(path)

	if !strings.Contains(errorCase.Error(), "generating temp file suffix") ||
		!strings.Contains(errorCase.Error(), errEntropyExhausted.Error()) {
		t.Fatalf("suffix failure must be wrapped with context, got: %v", errorCase)
	}
}
