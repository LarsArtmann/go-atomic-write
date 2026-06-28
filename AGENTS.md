# AGENTS.md — go-atomic-write

## Project

Single-package Go library providing TOCTOU-safe file writes via xxhash64 fingerprint verification, cross-platform file locking (`flock`/`LockFileEx`), atomic rename, and fsync for crash durability.

- **Module:** `github.com/larsartmann/go-atomic-write`
- **Go version:** 1.26.3
- **Main branch:** `master` (configured in `git-town.toml`)

## Commands

| Command                         | Purpose                                                               |
| ------------------------------- | --------------------------------------------------------------------- |
| `go test ./...`                 | Run all tests                                                         |
| `go test -race ./...`           | Run all tests with the race detector                                  |
| `go test -bench=. -benchmem`    | Run benchmarks (in `hash_bench_test.go`)                              |
| `go vet ./...`                  | Static analysis                                                       |
| `go build ./...`                | Verify compilation                                                    |
| `golangci-lint run ./...`       | Run all configured linters — MUST exit with `0 issues` before merging |
| `golangci-lint run ./... --fix` | Auto-fix formatting/imports (gci, gofumpt, goimports, golines)        |

`golangci-lint` config is the **gating quality check** of this project. The configuration lives in `.golangci.yml` and enables ~100 linters at their strictest defaults. Adding a new third-party import requires updating the `depguard.rules.main.allow` list.

No Makefile, no flake.nix, no CI config. All commands are plain `go` toolchain.

## Structure

Flat single-package layout — all source in the repository root:

| File                  | Purpose                                                                                         |
| --------------------- | ----------------------------------------------------------------------------------------------- |
| `atomicwrite.go`      | Public API + staging: `Fingerprint`, `Write`, `FingerprintFile`, `writeAndSync`, `randomSuffix` |
| `rename_unix.go`      | POSIX `atomicRename` (single `rename(2)` + directory `fsync`)                                   |
| `rename_windows.go`   | Windows `atomicRename` (retry on `ERROR_ACCESS_DENIED`/`ERROR_SHARING_VIOLATION`)               |
| `atomicwrite_test.go` | Unit + concurrency + integrity tests                                                            |
| `hash_bench_test.go`  | xxhash64 vs SHA-256 benchmarks                                                                  |

## Architecture & Data Flow

```
Caller reads file → computes Fingerprint → modifies data → calls Write()
  └─ Write() stages to unique .tmp → fsync .tmp → locks + verifies fingerprint → atomic rename → fsync dir
       └─ On mismatch: returns ErrConcurrentModification, cleans up .tmp
```

Key internal functions:

- `Write()` — entry point; generates unique temp path, stages + fsyncs data, branches on fingerprint
- `writeAndSync()` — creates temp file, writes data, fsyncs, closes (with cleanup on error)
- `randomSuffix()` — generates random hex suffix for unique temp file names (crypto/rand)
- `commitWithVerification()` — acquires exclusive `flock`, re-reads target, verifies match, renames
- `atomicRename()` — single `os.Rename` + directory fsync (POSIX); retry loop (Windows)
- `syncDir()` — opens and fsyncs the target directory after rename (POSIX only)

## Dependencies

| Dependency          | Purpose                                                                        |
| ------------------- | ------------------------------------------------------------------------------ |
| `cespare/xxhash/v2` | Non-cryptographic hash for fingerprinting (chosen over SHA-256 for ~11× speed) |
| `gofrs/flock`       | Cross-platform file locking                                                    |

Both are intentional, minimal, and not candidates for replacement.

## Testing Conventions

- Standard `testing` package — no testify or other test utilities
- `t.Parallel()` on most tests
- `t.TempDir()` for filesystem isolation
- `tempFile(t, content)` helper creates a file in a temp dir
- Concurrency test (`TestConcurrentWriteRACE`) uses `sync.WaitGroup.Go` and `atomic.Int32` with divergent content per writer — verifies both `flock` contention and data integrity (no byte interleaving)

## Gotchas

- `.tmp` files use unique names (`path + "." + randomHex + ".tmp"`) to prevent concurrent writers from corrupting a shared staging file
- `.tmp` files are created alongside the target file (same directory) — callers need write permissions on the directory, not just the file
- `.tmp` files are in `.gitignore`; `.bak` files are no longer created (the `.bak` pattern was removed — it broke POSIX atomicity)
- `fsync` is called on both the temp file (before rename) and the target directory (after rename, POSIX only) for crash durability
- On Windows, directory fsync is a no-op (no POSIX-equivalent directory sync concept); file data durability is handled by `FlushFileBuffers` via `file.Sync()`
- On Windows, `os.Rename` is retried up to 5 times with exponential backoff on `ERROR_ACCESS_DENIED`/`ERROR_SHARING_VIOLATION` (antivirus, open handles)
- `Fingerprint` is `[8]byte` (xxhash64), stored big-endian — do not compare with `string` or `[]byte` representations of hashes
- `FingerprintFile` returns zero-value (not an error) for nonexistent files — this is the "first write" sentinel
- File permissions are preserved from the existing file, defaulting to `0644` for new files
- `ErrConcurrentModification` is a sentinel `errors.New` value — always check with `errors.Is`, not string matching
- The `//nolint:gosec` comments on `os.ReadFile`/`os.OpenFile`/`os.Open` calls are intentional — `path` is caller-controlled, not user input
- The `//nolint:gosec` comments on `os.ReadFile`/`os.WriteFile` calls in `atomicwrite_test.go` are intentional — all paths are `t.TempDir()`-rooted and never user input
- Lint is gating: a single `golangci-lint` issue blocks merging. Common pitfalls:
  - Single-letter variables (`fp`, `wg`, `h`, `mb`) trigger `varnamelen` — use full words
  - Inline `if err := …; err != nil` triggers `noinlineerr` — assign first, then check
  - `b.Helper()` is required on every benchmark helper that takes `*testing.B` (thelper)
  - 0o644 in tests triggers `gosec G306` — annotate with `//nolint:gosec` and rationale
  - `os.ReadFile(path)`/`os.OpenFile` with a variable path triggers `gosec G304` — annotate with `//nolint:gosec` and rationale
  - Blank line between error assignment and `if err != nil` triggers `wsl_v5` — put the blank line above the assignment instead
  - New third-party imports require updating `.golangci.yml` `depguard.rules.main.allow`
