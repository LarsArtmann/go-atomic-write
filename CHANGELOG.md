# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

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
