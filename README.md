<h1 align="center">go-atomic-write</h1>

<p align="center"><strong>Crash-safe, race-free file writes for Go.</strong></p>

<p align="center">
<a href="https://pkg.go.dev/github.com/larsartmann/go-atomic-write"><img src="https://pkg.go.dev/badge/github.com/larsartmann/go-atomic-write.svg" alt="Go Reference"></a>

<a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License: MIT"></a>
</p>

<p align="center">
<a href="https://atomicwrite.lars.software">Documentation</a> · <a href="https://pkg.go.dev/github.com/larsartmann/go-atomic-write">API Reference</a>
</p>

---

Every `os.WriteFile` call has three failure modes that silently corrupt data. This library eliminates all three with a minimal dependency footprint — fingerprint verification, cross-platform file locking, atomic rename, and fsync for crash durability.

## Why?

Writing a file safely is harder than it looks:

- **TOCTOU races** — between reading and writing, another process can modify the file. You overwrite their changes without knowing.
- **Partial writes** — a crash mid-write leaves a corrupt, half-written file. The old content is gone.
- **Concurrent writers** — two processes write the same file. Bytes interleave. Data is lost.

`os.WriteFile` handles none of these. A naive write-to-temp-then-rename handles one. **go-atomic-write handles all three.**

## Comparison

| Approach            | TOCTOU-safe | Crash-durable | Concurrent-safe |   Dependencies    |
| ------------------- | :---------: | :-----------: | :-------------: | :---------------: |
| `os.WriteFile`      |             |               |                 |       None        |
| DIY write + rename  |             |    Partial    |                 |       None        |
| **go-atomic-write** |      ✓      |       ✓       |        ✓        | 2 (xxhash, flock) |

## How it works

1. **Fingerprint** — compute an xxhash64 digest when you read the file
2. **Write to unique `.tmp`** — stage new content in a uniquely-named temp file (prevents concurrent writers from corrupting each other)
3. **fsync** — flush the temp file to disk for crash durability
4. **Lock + verify** — acquire an exclusive file lock (`flock` on Unix, `LockFileEx` on Windows), re-read the file, and verify the fingerprint still matches
5. **Atomic rename** — rename the temp file to the target (single `rename(2)` on POSIX, `MoveFileEx` on Windows)
6. **fsync directory** — sync the directory entry to make the rename durable (POSIX)

If the fingerprint doesn't match, `Write` returns `ErrConcurrentModification` — the caller should re-read and retry.

## Use cases

- **Configuration files** — read-modify-write cycles where another process might edit the same config
- **State and checkpoint files** — database state, job progress, session stores — corruption means data loss
- **Log and data rotation** — append or replace files without readers seeing partial content
- **Cache invalidation** — write-through caches that must never serve a half-written file
- **CI/CD tooling** — generated manifests, lock files, build artifacts written under concurrent pipelines

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
    err = atomicwrite.WriteVerified(path, newData, fp)
    if err != nil {
        panic(err)
    }
}
```

### First write (no existing file)

Use `Write` when there is no prior content to protect (first write, temp files):

```go
err := atomicwrite.Write(path, data)
```

Or pass a zero-value `Fingerprint` to `WriteVerified` if you need concurrent-creation
detection (the write fails if another process creates the file first):

```go
err := atomicwrite.WriteVerified(path, data, atomicwrite.Fingerprint{})
```

### Write only if content changed

Use `WriteIfChanged` when re-running a tool should not produce spurious diffs
(config writers, code generators, lock-file updaters):

```go
changed, err := atomicwrite.WriteIfChanged(path, newData)
if err != nil {
    panic(err)
}
if changed {
    fmt.Println("config updated")
} else {
    fmt.Println("config already up to date")
}
```

No content change means no file mutation, no mtime bump, no file-watcher trigger.

### Detecting concurrent modification

```go
err := atomicwrite.WriteVerified(path, newData, fp)
if errors.Is(err, atomicwrite.ErrConcurrentModification) {
    // Re-read the file, merge changes, and retry
}
```

### Streaming large content with WriteFunc

When content is large or produced incrementally (JSON encoders, diagram renderers),
use `WriteFunc` to stream via a callback instead of holding the full payload in memory:

```go
err := atomicwrite.WriteFuncVerified(path, func(w io.Writer) error {
    enc := json.NewEncoder(w)
    return enc.Encode(largeObject)
}, fp)
```

For streaming without TOCTOU verification, use `WriteFunc(path, fn)`.

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

| Symbol                            | Description                                                |
| --------------------------------- | ---------------------------------------------------------- |
| `Fingerprint`                     | `[8]byte` — xxhash64 digest of file content at read time   |
| `Fingerprint.IsZero()`            | Returns `true` for zero-value (no prior file)              |
| `Fingerprint.Matches(data)`       | Returns `true` if data produces the same fingerprint       |
| `FingerprintFromBytes(data)`      | Computes fingerprint from raw bytes                        |
| `FingerprintFile(path)`           | Computes fingerprint from a file (zero if nonexistent)     |
| `Write(path, data)`               | Atomic write, no TOCTOU check (first write, single-writer) |
| `WriteVerified(path, data, fp)`   | Atomic write with fingerprint race-check                   |
| `WriteIfChanged(path, data)`      | Writes only if content differs; returns `(changed, err)`   |
| `WriteFunc(path, fn)`             | Streams content via callback, no TOCTOU check              |
| `WriteFuncVerified(path, fn, fp)` | Streams content via callback with TOCTOU check             |
| `ErrConcurrentModification`       | Sentinel error: file changed between read and write        |

## Design decisions

- **xxhash64 over SHA-256** — the fingerprint detects changes, not attackers. xxhash64 is ~11× faster, zero allocations, and hits ~27 GB/s. SHA-NI hardware acceleration is already included in the SHA-256 numbers.
- **Unique temp files via `crypto/rand`** — each writer gets a uniquely-named staging file (`path + "." + randomHex + ".tmp"`), preventing concurrent writers from corrupting a shared temp file.
- **`flock`, not mutexes** — file locks protect across processes, not just goroutines. `flock(2)` on Unix, `LockFileEx` on Windows.
- **Single `rename(2)` on POSIX** — one atomic syscall replaces the target. No two-rename window where the file could be missing.
- **fsync before and after** — temp file is synced before rename (data durability); target directory is synced after rename (metadata durability, POSIX only).
- **Minimal dependencies** — only [xxhash](https://github.com/cespare/xxhash) for hashing and [flock](https://github.com/gofrs/flock) for locking. Both are intentional and not candidates for replacement.

## Error contract

- All errors are wrapped with `fmt.Errorf` and `%w` — use `errors.Is` / `errors.As` to inspect
- `ErrConcurrentModification` is a sentinel error you can check with `errors.Is`
- Temp files (`.tmp`) are cleaned up on any failure
- File permissions are preserved from the existing file (defaults to `0644` for new files)
- Data is fsync'd before rename; the directory is fsync'd after (POSIX) for crash durability

## Dependencies

| Dependency                                            | Purpose                                                                |
| ----------------------------------------------------- | ---------------------------------------------------------------------- |
| [`cespare/xxhash`](https://github.com/cespare/xxhash) | Fast non-cryptographic hash for fingerprinting                         |
| [`gofrs/flock`](https://github.com/gofrs/flock)       | Cross-platform file locking (`flock` on Unix, `LockFileEx` on Windows) |

## Benchmarks

xxhash64 vs `crypto/sha256` — 100K iterations per benchmark, single core.

**Hardware:** AMD Ryzen AI MAX+ 395, Go 1.26.3, linux/amd64

| Hash     | Size   |   ns/op |  Throughput | Allocations |
| -------- | ------ | ------: | ----------: | :---------- |
| xxhash64 | 1 KB   |      42 | 24,443 MB/s | 0           |
| SHA-256  | 1 KB   |     486 |  2,106 MB/s | 1 × 32 B    |
| xxhash64 | 10 KB  |     383 | 26,708 MB/s | 0           |
| SHA-256  | 10 KB  |   4,253 |  2,408 MB/s | 1 × 32 B    |
| xxhash64 | 100 KB |   3,744 | 27,349 MB/s | 0           |
| SHA-256  | 100 KB |  41,760 |  2,452 MB/s | 1 × 32 B    |
| xxhash64 | 1 MB   |  37,954 | 27,628 MB/s | 0           |
| SHA-256  | 1 MB   | 429,268 |  2,443 MB/s | 1 × 32 B    |

**xxhash64 is ~11× faster** than SHA-256, zero allocations, and hits ~27 GB/s — effectively RAM bandwidth. No hardware accelerator exists or is needed; the hash is memory-bound, not compute-bound. SHA-256 results already include SHA-NI hardware acceleration.

Run yourself:

```bash
go test -bench=. -benchmem -benchtime=100000x
```

## Platform support

Works everywhere Go compiles. File locking uses:

- `flock(2)` on Linux, macOS, BSD
- `LockFileEx` on Windows

Atomic rename and durability:

- POSIX: single `rename(2)` (atomic replace) + `fsync` on the directory
- Windows: `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING`, retried on `ERROR_ACCESS_DENIED`/`ERROR_SHARING_VIOLATION` (antivirus contention)

## API stability

This library follows [Go module versioning](https://go.dev/doc/modules/version-numbers):

- **Pre-v1.0.0** (`v0.x.y`): The API may change between minor versions. We minimize breaking changes, but reserve the right to rename or adjust public symbols based on user feedback.
- **Post-v1.0.0**: Standard Go compatibility guarantees apply. No breaking changes without a major version bump.

The core `Write` / `Fingerprint` API is stable and unlikely to change.

## License

MIT
