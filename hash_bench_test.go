package atomicwrite

import (
	"crypto/sha256"
	"testing"

	"github.com/cespare/xxhash/v2"
)

func BenchmarkXxhash64_1KB(b *testing.B)   { benchXxhash(b, kb) }
func BenchmarkXxhash64_10KB(b *testing.B)  { benchXxhash(b, 10*kb) }
func BenchmarkXxhash64_100KB(b *testing.B) { benchXxhash(b, 100*kb) }
func BenchmarkXxhash64_1MB(b *testing.B)   { benchXxhash(b, megabyte) }

func BenchmarkSHA256_1KB(b *testing.B)   { benchSHA256(b, kb) }
func BenchmarkSHA256_10KB(b *testing.B)  { benchSHA256(b, 10*kb) }
func BenchmarkSHA256_100KB(b *testing.B) { benchSHA256(b, 100*kb) }
func BenchmarkSHA256_1MB(b *testing.B)   { benchSHA256(b, megabyte) }

const (
	kb       = 1024
	megabyte = 1024 * 1024
)

// benchData generates a deterministic byte pattern of the given size for hashing benchmarks.
// Fills by index, not append, so the makezero lint about non-zero length is intentional.
func benchData(size int) []byte {
	data := make([]byte, size) //nolint:makezero // pre-allocated buffer filled by index, not append
	for index := range data {
		data[index] = byte(index % 256)
	}

	return data
}

func benchXxhash(b *testing.B, size int) {
	b.Helper()

	data := benchData(size)

	b.SetBytes(int64(size))
	b.ResetTimer()

	for range b.N {
		_ = xxhash.Sum64(data)
	}
}

func benchSHA256(b *testing.B, size int) {
	b.Helper()

	data := benchData(size)

	b.SetBytes(int64(size))
	b.ResetTimer()

	for range b.N {
		hasher := sha256.New()
		hasher.Write(data)
		_ = hasher.Sum(nil)
	}
}

func BenchmarkXxhash64_Streaming_1MB(b *testing.B) { benchXxhashStreaming(b, megabyte) }
func BenchmarkSHA256_Streaming_1MB(b *testing.B)   { benchSHA256Streaming(b, megabyte) }

func benchXxhashStreaming(b *testing.B, size int) {
	b.Helper()

	data := benchData(size)

	const chunkSize = 4096

	b.SetBytes(int64(size))
	b.ResetTimer()

	for range b.N {
		hasher := xxhash.New()

		for off := 0; off < len(data); off += chunkSize {
			end := min(off+chunkSize, len(data))

			_, _ = hasher.Write(data[off:end])
		}

		_ = hasher.Sum(nil)
	}
}

func benchSHA256Streaming(b *testing.B, size int) {
	b.Helper()

	data := benchData(size)

	const chunkSize = 4096

	b.SetBytes(int64(size))
	b.ResetTimer()

	for range b.N {
		hasher := sha256.New()

		for off := 0; off < len(data); off += chunkSize {
			end := min(off+chunkSize, len(data))

			hasher.Write(data[off:end])
		}

		_ = hasher.Sum(nil)
	}
}
