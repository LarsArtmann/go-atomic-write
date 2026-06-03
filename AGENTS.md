# AGENTS.md — go-atomic-write

## Project

Single-package Go library providing TOCTOU-safe file writes via xxhash64 fingerprint verification, cross-platform file locking (`flock`/`LockFileEx`), and atomic rename.

- **Module:** `github.com/larsartmann/go-atomic-write`
- **Go version:** 1.26.3
- **Main branch:** `master` (configured in `git-town.toml`)

## Commands

| Command                            | Purpose                                                                |
| ---------------------------------- | ---------------------------------------------------------------------- |
| `go test ./...`                    | Run all tests                                                          |
| `go test -race ./...`              | Run all tests with the race detector                                   |
| `go test -bench=. -benchmem`       | Run benchmarks (in `hash_bench_test.go`)                               |
| `go vet ./...`                     | Static analysis                                                        |
| `go build ./...`                   | Verify compilation                                                     |
| `golangci-lint run ./...`          | Run all configured linters — MUST exit with `0 issues` before merging |
| `golangci-lint run ./... --fix`    | Auto-fix formatting/imports (gci, gofumpt, goimports, golines)         |

`golangci-lint` config is the **gating quality check** of this project. The configuration lives in `.golangci.yml` and enables ~100 linters at their strictest defaults. Adding a new third-party import requires updating the `depguard.rules.main.allow` list.

No Makefile, no flake.nix, no CI config. All commands are plain `go` toolchain.

## Structure

Flat single-package layout — all source in the repository root:

| File                  | Purpose                                                                                         |
| --------------------- | ----------------------------------------------------------------------------------------------- |
| `atomicwrite.go`      | Entire public API (~130 LOC): `Fingerprint`, `Write`, `FingerprintFile`, `FingerprintFromBytes` |
| `atomicwrite_test.go` | Unit + concurrency tests                                                                        |
| `hash_bench_test.go`  | xxhash64 vs SHA-256 benchmarks                                                                  |

## Architecture & Data Flow

```
Caller reads file → computes Fingerprint → modifies data → calls Write()
  └─ Write() stages to .tmp → locks + verifies fingerprint → atomic rename
       └─ On mismatch: returns ErrConcurrentModification, cleans up .tmp
```

Key internal functions:

- `Write()` — entry point; writes `.tmp`, branches on fingerprint presence
- `commitWithVerification()` — acquires exclusive `flock`, re-reads target, verifies match, renames
- `atomicRename()` — creates `.bak` of old file, renames `.tmp` to target (ignores `.bak` rename failure)

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
- Concurrency test (`TestConcurrentWriteRACE`) uses `sync.WaitGroup.Go` and `atomic.Int32` — not a standard race detector pattern, it exercises the actual `flock` contention path

## Gotchas

- `.tmp` and `.bak` files are created alongside the target file (same directory) — callers need write permissions on the directory, not just the file
- `.tmp` and `.bak` files are in `.gitignore`
- `.bak` creation failure is silently ignored (`_ = os.Rename(path, path+".bak")`) — this is intentional; the old file may not exist on first write
- `Fingerprint` is `[8]byte` (xxhash64), stored big-endian — do not compare with `string` or `[]byte` representations of hashes
- `FingerprintFile` returns zero-value (not an error) for nonexistent files — this is the "first write" sentinel
- File permissions are preserved from the existing file, defaulting to `0644` for new files
- `ErrConcurrentModification` is a sentinel `errors.New` value — always check with `errors.Is`, not string matching
- The `//nolint:gosec` comments on `os.ReadFile` calls in `atomicwrite.go` are intentional — `path` is caller-controlled, not user input
- The `//nolint:gosec` comments on `os.ReadFile`/`os.WriteFile` calls in `atomicwrite_test.go` are intentional — all paths are `t.TempDir()`-rooted and never user input
- Lint is gating: a single `golangci-lint` issue blocks merging. Common pitfalls:
  - Single-letter variables (`fp`, `wg`, `h`, `mb`) trigger `varnamelen` — use full words
  - Inline `if err := …; err != nil` triggers `noinlineerr` — assign first, then check
  - `b.Helper()` is required on every benchmark helper that takes `*testing.B` (thelper)
  - 0o644 in tests triggers `gosec G306` — annotate with `//nolint:gosec` and rationale
  - `os.ReadFile(path)` with a variable path triggers `gosec G304` — annotate with `//nolint:gosec` and rationale
  - New third-party imports require updating `.golangci.yml` `depguard.rules.main.allow`
