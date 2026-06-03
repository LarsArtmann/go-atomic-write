# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

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
