# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- **`WriteIfChanged(path, data) (changed bool, err error)`** — writes only if content differs from disk. Returns `changed=true` when the file was written, `false` when skipped (identical content). This is the idiomatic primitive for config-file writers and code generators that must not produce spurious diffs on re-runs: no content change means no file mutation, no mtime bump, no file-watcher trigger. Composes `FingerprintFile` + `WriteVerified`; first-write uses plain `Write` (no prior content to protect).

### Changed

- **BREAKING: `Write` / `WriteFunc` API split** (`atomicwrite.go`) — the previous `Write(path, data, fingerprint)` mixed two distinct intents behind a zero-value `Fingerprint{}`, which silently skipped TOCTOU verification (a footgun). The API is now split by intent:
  - `Write(path, data)` / `WriteFunc(path, fn)` — plain atomic write, no TOCTOU check (the common case).
  - `WriteVerified(path, data, fingerprint)` / `WriteFuncVerified(path, fn, fingerprint)` — atomic write WITH fingerprint verification. A zero-value fingerprint here means "the file must NOT already exist" — a concurrent creation is treated as a conflict and fails, rather than being silently ignored.
- This is a source-incompatible change: every call site passing `Fingerprint{}` to skip verification migrates to `Write`/`WriteFunc`; every call site passing a real fingerprint migrates to the `*Verified` variant. All in-repo consumers were updated.

### Migration

| Before                                       | After                             |
| -------------------------------------------- | --------------------------------- |
| `Write(path, data, Fingerprint{})`           | `Write(path, data)`               |
| `Write(path, data, fp)` (real fingerprint)   | `WriteVerified(path, data, fp)`   |
| `WriteFunc(path, fn, Fingerprint{})`         | `WriteFunc(path, fn)`             |
| `WriteFunc(path, fn, fp)` (real fingerprint) | `WriteFuncVerified(path, fn, fp)` |

## [0.3.0] - 2026-07-23

### Added

- **`WriteFunc(path, fn, fingerprint)`** — streaming atomic writes via a callback. The callback receives a `bufio.Writer` (64KB buffer) so callers can stream large or incrementally-produced content (JSON encoders, diagram renderers, etc.) without holding the full payload in memory. The temp file is fsync'd before atomic rename, and TOCTOU fingerprint verification works identically to `Write`.
- **Comprehensive `WriteFunc` tests** — first run, overwrite, large stream (400KB across 100 chunks), callback error cleanup, permission preservation, and no-leftover-files verification

### Changed

- Benchmark data generation extracted into a shared `benchData` helper — eliminates duplicated `make([]byte, size)` + fill-by-index across 4 benchmark functions
- `.gitattributes` added for cross-platform line-ending normalization
- Full marketing website launched with Astro + Starlight + Tailwind v4 at [atomicwrite.lars.software](https://atomicwrite.lars.software)

## [0.2.0] - 2026-06-28

### Fixed

- **Removed `.bak` two-rename pattern** — the old `atomicRename` did `path→.bak` then `tmp→path`, opening a window where the target didn't exist. A concurrent reader could see `ENOENT`, and a crash between the two renames left the target missing. The library named "atomic-write" was not atomic. Now a single `os.Rename(tmp, path)` is one atomic syscall on POSIX and atomically replaces on Windows via `MoveFileEx`.
- **Concurrent writer corruption** — `tmpPath` was computed before the lock (`path + ".tmp"`), so two concurrent `Write()` calls to the same target both truncated and wrote the same staging file, interleaving bytes. Now each writer gets a unique temp file via `crypto/rand` suffix (`path + "." + randomHex + ".tmp"`).
- **Concurrency test couldn't catch the bug** — `TestConcurrentWriteRACE` used identical content from all writers, so byte interleaving was undetectable. Rewritten with 10 writers using divergent content (`writer-0` through `writer-9`) plus an integrity check verifying the final file is exactly one writer's payload.
- **"Crash safety" doc claim was false** — no `fsync` was called anywhere. Now `writeAndSync()` opens the temp file, writes data, calls `file.Sync()`, and closes (with cleanup on error). On POSIX, the target directory is also `fsync`'d after rename via `syncDir()`. The package doc now says "atomic rename, and fsync for crash durability" — and means it.

### Added

- **`fsync` for crash durability** — temp file is fsync'd before rename; target directory is fsync'd after rename (POSIX only; Windows has no equivalent directory sync)
- **Unique temp file names** — `crypto/rand` 4-byte hex suffix prevents concurrent writer staging-file corruption
- **Platform-specific rename** — split into `rename_unix.go` (single `rename(2)` + directory `fsync`) and `rename_windows.go` (retry loop on `ERROR_ACCESS_DENIED`/`ERROR_SHARING_VIOLATION`, 5 retries with exponential backoff 1–16ms, matching Go's own `cmd/go` fix for issue 36568)
- **`writeAndSync()`** — opens temp file, writes data, fsyncs, closes; cleans up temp file on any error
- **`randomSuffix()`** — generates random hex suffix for unique temp file names
- **`syncDir()`** — opens and fsyncs the target directory after rename (POSIX only)
- **`TestWriteLeavesNoLeftoverFiles`** — verifies no `.bak` and no `.tmp` files remain after a successful write
- **Status report** — full HTML dashboard at `docs/status/2026-06-28_03-32_atomicity-durability-overhaul.html`

### Changed

- Package doc: "atomic rename for crash safety" → "atomic rename, and fsync for crash durability"
- `Write()` now stages to a unique temp file instead of a shared `path + ".tmp"`
- `atomicRename()` no longer creates `.bak` files — single `os.Rename` replaces atomically
- `TestConcurrentWriteRACE` rewritten: 10 writers, divergent content, integrity check
- `TestTempFileCleanedUpOnError` uses `filepath.Glob` instead of hardcoded temp path
- AGENTS.md, README.md updated: new architecture, data flow, file structure, platform-specific behavior, removed all `.bak` references
- `go.sum` tidied — removed stale testify/go-spew/go-difflib/yaml.v3 checksums

### Removed

- `TestWriteCreatesBackup` — tested the `.bak` pattern as a feature; the pattern was a bug
- `.bak` file creation — no longer created alongside the target

### Caveats

- Windows rename retry loop and `FlushFileBuffers` compile via build tags but are **untested on real Windows hardware**
- No POSIX-equivalent directory fsync on Windows — the directory entry after rename may not be durable after a crash
- No stale temp file cleanup — crashed writes leave `.tmp` files (unique names prevent interference but don't self-clean)

## [0.1.1] - 2026-06-03

### Changed

- **LICENSE changed from PROPRIETARY to MIT** — the README already stated MIT; the LICENSE file now matches
- `.golangci.yml` now permits `github.com/cespare/xxhash/v2` and `github.com/gofrs/flock` in the `depguard` allow list (both are intentional runtime dependencies, not candidates for replacement)
- Tests use `//nolint:gosec` annotations with rationale for `os.ReadFile`/`os.WriteFile` calls that operate exclusively on `t.TempDir()` paths
- `TestConcurrentWriteRACE` now calls `t.Parallel()` to satisfy `paralleltest` and to actually run in parallel with other tests
- Benchmark helpers call `b.Helper()` so failures are attributed to the caller
- Variable names lengthened for `varnamelen` compliance: `fp` → `fingerprint`, `wg` → `waitGroup`, `h` → `hasher`, `mb` → `megabyte`
- Inline `if err := …; err != nil { … }` blocks in tests refactored to plain assignment plus `if` for `noinlineerr` compliance
- Documentation formatting and accuracy fixes across README, CHANGELOG, CONTRIBUTING, DOMAIN_LANGUAGE, and AGENTS.md

### Fixed

- `errcheck`: `h.Write` return value is now explicitly discarded (`_, _ = hasher.Write(…)`) in the streaming benchmark
- `*.bak` added to `.gitignore` — backup files from atomic writes were showing as untracked

## [0.1.0] - 2026-06-02

### Added

- TOCTOU-safe file writes via xxhash64 fingerprint verification
- Cross-platform file locking (`flock` on Unix, `LockFileEx` on Windows) via `gofrs/flock`
- Atomic rename with `.bak` backup of previous content
- `Fingerprint` type (`[8]byte`) with `IsZero()` and `Matches()` methods
- `FingerprintFile` — computes fingerprint from file path (zero-value for nonexistent files)
- `FingerprintFromBytes` — computes fingerprint from raw bytes
- `Write(path, data, fingerprint)` — main API with TOCTOU protection
- `ErrConcurrentModification` sentinel error for fingerprint mismatch detection
- File permission preservation from existing file (defaults to `0644` for new files)
- xxhash64 vs SHA-256 benchmarks demonstrating ~11× speed advantage
- Unit tests, permission preservation tests, and concurrent writer contention test
