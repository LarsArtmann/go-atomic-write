package atomicwrite

import (
	"crypto/sha256"
	"hash"
	"testing"

	"github.com/cespare/xxhash/v2"
)

func BenchmarkXxhash64_1KB(b *testing.B)   { benchXxhash(b, kb) }
func BenchmarkXxhash64_10KB(b *testing.B)  { benchXxhash(b, 10*kb) }
func BenchmarkXxhash64_100KB(b *testing.B) { benchXxhash(b, 100*kb) }
func BenchmarkXxhash64_1MB(b *testing.B)   { benchXxhash(b, mb) }

func BenchmarkSHA256_1KB(b *testing.B)   { benchSHA256(b, kb) }
func BenchmarkSHA256_10KB(b *testing.B)  { benchSHA256(b, 10*kb) }
func BenchmarkSHA256_100KB(b *testing.B) { benchSHA256(b, 100*kb) }
func BenchmarkSHA256_1MB(b *testing.B)   { benchSHA256(b, mb) }

const (
	kb = 1024
	mb = 1024 * 1024
)

func benchXxhash(b *testing.B, size int) {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	b.SetBytes(int64(size))
	b.ResetTimer()
	for range b.N {
		_ = xxhash.Sum64(data)
	}
}

func benchSHA256(b *testing.B, size int) {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	b.SetBytes(int64(size))
	b.ResetTimer()
	for range b.N {
		h := sha256.New()
		h.Write(data)
		_ = h.Sum(nil)
	}
}

func BenchmarkXxhash64_Streaming_1MB(b *testing.B) { benchXxhashStreaming(b, mb) }
func BenchmarkSHA256_Streaming_1MB(b *testing.B)   { benchSHA256Streaming(b, mb) }

func benchXxhashStreaming(b *testing.B, size int) {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	const chunkSize = 4096
	b.SetBytes(int64(size))
	b.ResetTimer()
	for range b.N {
		h := xxhash.New()
		for off := 0; off < len(data); off += chunkSize {
			end := off + chunkSize
			if end > len(data) {
				end = len(data)
			}
			h.Write(data[off:end])
		}
		_ = h.Sum(nil)
	}
}

func benchSHA256Streaming(b *testing.B, size int) {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}
	const chunkSize = 4096
	b.SetBytes(int64(size))
	b.ResetTimer()
	for range b.N {
		var h hash.Hash = sha256.New()
		for off := 0; off < len(data); off += chunkSize {
			end := off + chunkSize
			if end > len(data) {
				end = len(data)
			}
			h.Write(data[off:end])
		}
		_ = h.Sum(nil)
	}
}
