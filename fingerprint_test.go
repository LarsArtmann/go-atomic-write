package atomicwrite

import (
	"testing"
)

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
