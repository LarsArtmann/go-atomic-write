# go-atomic-write

TOCTOU-safe file writes for Go — fingerprint verification, cross-platform file locking, and atomic rename.

## Why

Writing a file safely is harder than it looks:

- **TOCTOU races** — between reading and writing, another process can modify the file
- **Partial writes** — a crash mid-write leaves a corrupt file
- **Concurrent writers** — multiple processes writing the same file

This library solves all three in ~130 lines with a minimal dependency footprint.

## How it works

1. **Fingerprint** — compute an xxhash64 digest when you read the file
2. **Write to `.tmp`** — stage new content in a temp file
3. **Lock + verify** — acquire an exclusive file lock (`flock` on Unix, `LockFileEx` on Windows), re-read the file, and verify the fingerprint still matches
4. **Atomic rename** — rename the temp file to the target (creates `.bak` of the old content)

If the fingerprint doesn't match, `Write` returns `ErrConcurrentModification` — the caller should re-read and retry.

## Install

```bash
go get github.com/larsartmann/go-atomic-write
```

## Usage

```go
package main

import (
    "fmt"
    "os"

    atomicwrite "github.com/larsartmann/go-atomic-write"
)

func main() {
    path := "/path/to/config.json"

    // Read + fingerprint
    data, err := os.ReadFile(path)
    if err != nil && !os.IsNotExist(err) {
        panic(err)
    }
    fp := atomicwrite.FingerprintFromBytes(data)

    // Modify
    newData := []byte(`{"updated": true}`)

    // Write with TOCTOU protection
    err = atomicwrite.Write(path, newData, fp)
    if err != nil {
        panic(err)
    }
}
```

### First write (no existing file)

Pass a zero-value `Fingerprint` to skip verification:

```go
err := atomicwrite.Write(path, data, atomicwrite.Fingerprint{})
```

### Detecting concurrent modification

```go
err := atomicwrite.Write(path, newData, fp)
if errors.Is(err, atomicwrite.ErrConcurrentModification) {
    // Re-read the file, merge changes, and retry
}
```

### Fingerprinting a file

```go
// From raw bytes
fp := atomicwrite.FingerprintFromBytes([]byte("content"))

// From a file path (returns zero Fingerprint if file doesn't exist)
fp, err := atomicwrite.FingerprintFile(path)

// Check if it represents no prior file
if fp.IsZero() {
    // First run — no file existed
}

// Check if content matches
if fp.Matches(currentData) {
    // File hasn't changed
}
```

## API

| Symbol | Description |
|--------|-------------|
| `Fingerprint` | `[8]byte` — xxhash64 digest of file content at read time |
| `Fingerprint.IsZero()` | Returns `true` for zero-value (no prior file) |
| `Fingerprint.Matches(data)` | Returns `true` if data produces the same fingerprint |
| `FingerprintFromBytes(data)` | Computes fingerprint from raw bytes |
| `FingerprintFile(path)` | Computes fingerprint from a file (zero if nonexistent) |
| `Write(path, data, fp)` | Writes data with TOCTOU protection |
| `ErrConcurrentModification` | Sentinel error: file changed between read and write |

## Error contract

- All errors are wrapped with `fmt.Errorf` and `%w` — use `errors.Is` / `errors.As` to inspect
- `ErrConcurrentModification` is a sentinel error you can check with `errors.Is`
- Temp files (`.tmp`) are cleaned up on any failure
- File permissions are preserved from the existing file (defaults to `0644` for new files)
- A `.bak` of the previous content is created on successful writes (not on first write)

## Dependencies

| Dependency | Purpose |
|-----------|---------|
| [`cespare/xxhash`](https://github.com/cespare/xxhash) | Fast non-cryptographic hash for fingerprinting |
| [`gofrs/flock`](https://github.com/gofrs/flock) | Cross-platform file locking (`flock` on Unix, `LockFileEx` on Windows) |

## Benchmarks

xxhash64 vs `crypto/sha256` — 100K iterations per benchmark, single core.

**Hardware:** AMD Ryzen AI MAX+ 395, Go 1.26.3, linux/amd64

| Hash | Size | ns/op | Throughput | Allocations |
|------|------|------:|----------:|:------------|
| xxhash64 | 1 KB | 42 | 24,443 MB/s | 0 |
| SHA-256 | 1 KB | 486 | 2,106 MB/s | 1 × 32 B |
| xxhash64 | 10 KB | 383 | 26,708 MB/s | 0 |
| SHA-256 | 10 KB | 4,253 | 2,408 MB/s | 1 × 32 B |
| xxhash64 | 100 KB | 3,744 | 27,349 MB/s | 0 |
| SHA-256 | 100 KB | 41,760 | 2,452 MB/s | 1 × 32 B |
| xxhash64 | 1 MB | 37,954 | 27,628 MB/s | 0 |
| SHA-256 | 1 MB | 429,268 | 2,443 MB/s | 1 × 32 B |

**xxhash64 is ~11× faster** than SHA-256, zero allocations, and hits ~27 GB/s — effectively RAM bandwidth. No hardware accelerator exists or is needed; the hash is memory-bound, not compute-bound. SHA-256 results already include SHA-NI hardware acceleration.

Run yourself:

```bash
cd lib/go-atomic-write && go test -bench=. -benchmem -benchtime=100000x
```

## Platform support

Works everywhere Go compiles. File locking uses:
- `flock(2)` on Linux, macOS, BSD
- `LockFileEx` on Windows

## License

MIT
