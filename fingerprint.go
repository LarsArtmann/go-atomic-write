package atomicwrite

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/cespare/xxhash/v2"
)

// Fingerprint is an xxhash64 digest of file content at read time.
// A zero-value Fingerprint indicates no prior file existed.
type Fingerprint [8]byte

// IsZero returns true if the fingerprint represents no prior file.
func (fp Fingerprint) IsZero() bool {
	return fp == Fingerprint{}
}

// Matches returns true if the given content produces the same fingerprint.
func (fp Fingerprint) Matches(content []byte) bool {
	return FingerprintFromBytes(content) == fp
}

// FingerprintFromBytes computes an xxhash64 Fingerprint from raw content.
func FingerprintFromBytes(data []byte) Fingerprint {
	var fp Fingerprint
	binary.BigEndian.PutUint64(fp[:], xxhash.Sum64(data))

	return fp
}

// FingerprintFile computes an xxhash64 Fingerprint from a file's current content.
// Returns a zero-value Fingerprint if the file does not exist.
func FingerprintFile(path string) (Fingerprint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Fingerprint{}, nil
		}

		return Fingerprint{}, fmt.Errorf("reading %s for fingerprint: %w", path, err)
	}

	return FingerprintFromBytes(data), nil
}
