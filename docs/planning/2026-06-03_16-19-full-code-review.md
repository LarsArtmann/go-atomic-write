# Full Code Review — go-atomic-write

**Date:** 2026-06-03  
**Reviewer:** Senior Staff Architect  
**Scope:** Every source file, test file, and supporting document

---

## Executive Summary

This is a **well-crafted, focused library** that does exactly one thing and does it correctly. 127 LOC of production code, 295 LOC of tests, clean API, minimal dependencies. The codebase demonstrates good taste.

**Overall Grade: A-** — Excellent foundation with a few meaningful improvements possible.

All tests pass (including race detector). `go vet` clean. Build clean. Benchmarks confirm the xxhash64 choice (~10× faster than SHA-256).

---

## File-by-File Review

### `atomicwrite.go` (127 LOC) — Grade: A-

**Strengths:**

- Clean, linear control flow — easy to read top-to-bottom
- `Fingerprint` as a named type (`[8]byte`) with methods — good domain modeling
- `ErrConcurrentModification` as a sentinel error with `errors.Is` support — idiomatic
- Permission preservation from existing file — thoughtful
- `.bak` creation for crash recovery — good safety net
- Proper error wrapping with `%w` throughout

**Issues Found:**

| #   | Severity   | Issue                                                                                                                                                                                                                                                                                                               | Location | Verdict    |
| --- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---------- |
| 1   | **Medium** | Lock is acquired but `defer fileLock.Close()` comes after — the `defer` only runs when the function returns. If `atomicRename` panics, the lock is released but the `.tmp` file is orphaned. Not a real-world concern in Go (no exceptions), but worth noting for defensive programming.                            | L96      | Acceptable |
| 2   | **Medium** | `.bak` from a previous write is never cleaned up. On repeated writes, the old `.bak` is silently overwritten — this is correct behavior but not documented in the README's error contract section.                                                                                                                  | L115     | Doc gap    |
| 3   | **Low**    | `Fingerprint` is `[8]byte` — the `IsZero()` method compares against zero value. This means an empty file produces a non-zero fingerprint (xxhash of empty bytes ≠ 0), which is correct. But a caller who computes `FingerprintFromBytes(nil)` gets a non-zero result. This is semantically fine but could surprise. | L36-41   | By design  |
| 4   | **Low**    | `os.Stat` error is silently ignored (only `err == nil` path extracts perm). If Stat fails for a reason other than NotExist (e.g., permission denied on parent dir), we silently fall through to `0644` and may fail later at `WriteFile`. This is arguably fine — the later error is more actionable.               | L68-71   | Acceptable |
| 5   | **Low**    | No `sync.Once` or similar protection for the `defaultFilePerm` constant — not needed since it's a compile-time constant. Clean.                                                                                                                                                                                     | L64      | Clean      |

**Architecture observations:**

- Data flow is linear and correct: stage → lock → verify → rename
- `commitWithVerification` properly cleans up `.tmp` on all error paths
- `atomicRename` creates `.bak` before rename — if `.bak` rename fails (old file doesn't exist), it's ignored. If `.tmp` → target rename fails, error is returned. This ordering is correct.
- The `flock` is on the _target_ file, not the `.tmp` file — this is the right choice for coordinating multiple writers

### `atomicwrite_test.go` (295 LOC) — Grade: A

**Strengths:**

- 11 test functions covering all major paths
- `t.Parallel()` on most tests — good practice
- `t.TempDir()` for isolation — correct
- `tempFile` helper — clean, DRY
- Concurrency test (`TestConcurrentWriteRACE`) actually exercises flock contention, not just goroutine races
- Race detector passing (`go test -race`)

**Issues Found:**

| #   | Severity | Issue                                                                                                                                                                                                                         | Location | Verdict       |
| --- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------------- |
| 1   | **Low**  | `TestConcurrentWriteRACE` uses `wg.Go` (Go 1.25+ feature) — good, but the test's success criteria (`successes >= 1`) is weak. It doesn't verify that the final file content is valid, only that at least one write succeeded. | L234-277 | Could improve |
| 2   | **Low**  | `TestConcurrentWriteRACE` is missing `t.Parallel()` — probably intentional since it's a concurrency test, but inconsistent with other tests.                                                                                  | L234     | Acceptable    |
| 3   | **Low**  | No test for writing empty content (`[]byte{}`) — edge case worth covering.                                                                                                                                                    | —        | Gap           |
| 4   | **Low**  | No test for very large files (e.g., 100MB) — not critical for a library, but the README benchmarks go up to 1MB.                                                                                                              | —        | Nice-to-have  |
| 5   | **Info** | `TestAtomicRenameReportsErrorOnFailure` tests internal function directly — fine for a single-package layout, and the test validates error propagation.                                                                        | L279-295 | Clean         |

**Missing test coverage:**

- Empty file write
- Write to read-only directory (permission error path)
- `.bak` overwrite on second write to same file
- `FingerprintFile` with a directory path instead of a file
- `Fingerprint.Matches` with nil input

### `hash_bench_test.go` (104 LOC) — Grade: A

**Strengths:**

- Clean benchmark structure with both batch and streaming modes
- Uses `b.SetBytes` for throughput calculation
- Uses `min()` builtin (Go 1.21+)
- Data initialization is correct and done before the timer

**Issues Found:**

| #   | Severity | Issue                                                                                                                                                      | Location             | Verdict    |
| --- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------- | ---------- |
| 1   | **Low**  | Data initialization pattern (`byte(i % 256)`) is repeated in 4 benchmark functions. Could extract to a helper, but at 104 LOC this is acceptable.          | L27-29, L41-43, etc. | Acceptable |
| 2   | **Info** | `byte(i % 256)` — `i % 256` is always `i` for `byte` range, but this works because `i` exceeds 255 and wraps via `byte(i % 256)`. Could just be `byte(i)`. | L28                  | Trivial    |

### `go.mod` — Grade: A

- Go 1.26.3 — current
- 2 direct dependencies, 1 indirect — minimal
- Both dependencies are well-maintained, purpose-fit, and not candidates for replacement

### `README.md` (162 LOC) — Grade: A

**Strengths:**

- Clear "Why" section explaining the problem
- Complete usage examples including error handling
- API table — excellent reference
- Benchmark table with hardware specs
- Error contract section — sets expectations

**Issues:**

- The benchmark run command says `cd lib/go-atomic-write` — this path doesn't match the project structure (it's the repo root, not `lib/`)
- No mention of `.bak` cleanup behavior in the error contract section

### `docs/DOMAIN_LANGUAGE.md` — Grade: D (Scaffold)

This is entirely a placeholder template. None of the domain terms are filled in. For this project's scope, it's arguably not needed — the domain is small enough that the README and code comments suffice.

### `AGENTS.md` — Grade: A

Excellent project context for AI sessions. Accurate, concise, covers the important gotchas.

---

## Architecture Analysis

### Data Flow (correct)

```
Read file → Fingerprint → Modify data → Write(path, data, fp)
  ├─ Zero fp:  stage .tmp → atomicRename (no lock)
  └─ Non-zero: stage .tmp → flock → re-read → verify → atomicRename → unlock
       └─ Mismatch: cleanup .tmp → ErrConcurrentModification
```

### Type Safety — Grade: A

- `Fingerprint` is a named type, not `uint64` or `[]byte` — good
- `IsZero()` and `Matches()` methods prevent direct byte comparison mistakes
- `ErrConcurrentModification` is a typed sentinel — `errors.Is` works correctly
- The only "stringly-typed" thing is the path parameter, which is unavoidable in Go's `os` API

### Split Brains — None detected

No duplicate type definitions, no conflicting state representations. Single source of truth for all types.

### Duplications — Minimal

- `byte(i % 256)` data init pattern in benchmarks — acceptable at this scale
- `.tmp` path computation (`path + ".tmp"`) appears once — not duplicated
- `cleanupTmp` is called in 2 error paths — properly extracted to a function

### Composability — Grade: B+

The library is correctly composable:

- `FingerprintFromBytes` can be used without `Write` (e.g., for comparison)
- `FingerprintFile` composes `FingerprintFromBytes` — clean
- `Write` accepts a `Fingerprint` rather than computing it — the caller controls when fingerprinting happens

**What could be better:** The `Write` function takes `(path, data, fingerprint)` — a flat parameter list. For future extensibility (e.g., custom temp dir, custom permissions, options pattern), this would need an options struct. But for a 127 LOC library, YAGNI applies.

### Ghost Systems — None

Every exported symbol is used. No dead code. No unused dependencies. Every test exercises real code paths.

---

## Brutal Self-Review Answers

1. **What did we forget?**  
   Empty file edge case test. `.bak` overwrite behavior test. The README has a wrong path in the benchmark command.

2. **What is something stupid?**  
   The README benchmark command path (`lib/go-atomic-write`) is copy-pasted from a monorepo context and doesn't work for this standalone project.

3. **What could be better?**  
   `docs/DOMAIN_LANGUAGE.md` is an empty scaffold — either fill it in or delete it. Currently it's noise.

4. **What could still be improved?**  
   Missing test for empty content writes. The concurrent test could validate final file state.

5. **Did we lie?**  
   The README implies `.bak` is always created on successful writes, but it silently ignores `.bak` rename failure. The error contract says "A `.bak` of the previous content is created on successful writes (not on first write)" — this is technically true for the _attempt_, but the `.bak` may not exist if the OS rename fails.

6. **How can we be less stupid?**  
   Delete or fill in `DOMAIN_LANGUAGE.md`. Fix the README benchmark path.

7. **Ghost systems?**  
   None. Every piece of code is integrated and used.

8. **Scope creep?**  
   None — the library is admirably focused.

9. **Removed something useful?**  
   Not applicable — clean history.

10. **Split brains?**  
    None detected.

11. **Testing?**  
    Good but could be stronger: empty content, very large files, `.bak` overwrite, directory-as-path edge cases are missing.

---

## Improvement Plan (Pareto-Ranked)

### The 1% that delivers 51% of the value

| #   | Task                                                    | Impact | Effort | Why                                |
| --- | ------------------------------------------------------- | ------ | ------ | ---------------------------------- |
| 1   | Fix README benchmark path (`lib/go-atomic-write` → `.`) | High   | 1 min  | Broken command in user-facing docs |

### The 4% that delivers 64% of the value

| #   | Task                                          | Impact | Effort | Why                                                |
| --- | --------------------------------------------- | ------ | ------ | -------------------------------------------------- |
| 2   | Delete or fill in `docs/DOMAIN_LANGUAGE.md`   | Medium | 5 min  | Empty scaffold is noise, confuses new contributors |
| 3   | Add test for empty content write              | Medium | 5 min  | Edge case not covered                              |
| 4   | Add test for `.bak` overwrite on second write | Medium | 5 min  | Verifies backup rotation works                     |

### The 20% that delivers 80% of the value

| #   | Task                                                                | Impact | Effort | Why                                  |
| --- | ------------------------------------------------------------------- | ------ | ------ | ------------------------------------ |
| 5   | Add `FingerprintFile` test with directory path                      | Low    | 3 min  | Edge case robustness                 |
| 6   | Add `Fingerprint.Matches(nil)` test                                 | Low    | 2 min  | Documents nil-safety behavior        |
| 7   | Validate final file content in `TestConcurrentWriteRACE`            | Low    | 5 min  | Stronger concurrency assertion       |
| 8   | Document `.bak` cleanup/overwrite behavior in README error contract | Low    | 3 min  | Completes the contract documentation |

### Nice-to-have (not in the 80%)

| #   | Task                                                    | Impact  | Effort | Why                                    |
| --- | ------------------------------------------------------- | ------- | ------ | -------------------------------------- |
| 9   | Extract benchmark data init to helper                   | Trivial | 2 min  | DRY                                    |
| 10  | Add `//go:build` constraints or platform-specific tests | Low     | 30 min | Verify Windows `LockFileEx` path works |

---

## Verdict

This is a **high-quality, production-ready library** that solves a real problem with minimal complexity. The code is clean, the tests are solid, the dependencies are appropriate, and the documentation is excellent.

The only actionable items are:

1. Fix the README benchmark path (1 min)
2. Delete or fill in `DOMAIN_LANGUAGE.md` (5 min)
3. Add 2-3 missing edge case tests (15 min)

Total improvement effort: ~20 minutes for meaningful gains.

---

_Review complete. Every file visited. Every question answered honestly._

---

## Resolution (2026-07-26)

> This is a **point-in-time snapshot** from 2026-06-03. It reviews the pre-v0.2.0
> codebase and is now **largely superseded** by three later releases. Kept for
> history; do not act on its findings without cross-checking the current code.

### What changed since this review

| This review says…                                             | Current reality (2026-07-26)                                                                                                        |
| ------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `.bak` two-rename pattern is a "good safety net" / "strength" | **Removed in v0.2.0** — it was a bug (non-atomic window). Single `rename(2)` + directory `fsync` now. See `CHANGELOG.md` [0.2.0]. |
| `Write(path, data, fingerprint)` is the main API              | **Split in [Unreleased]** into `Write`, `WriteVerified`, `WriteIfChanged`, `WriteFunc`, `WriteFuncVerified`. See `atomicwrite.go`.  |
| `127 LOC`, no `fsync`                                         | Library grew with streaming + idempotent APIs; `fsync` is now core (`writeAndSync`, `syncDir`).                                     |

### Improvement-plan items — current status

| #  | Task                                                        | Status      | Note                                                                   |
| -- | ----------------------------------------------------------- | ----------- | ---------------------------------------------------------------------- |
| 1  | Fix README benchmark path                                   | ✅ Done     | Fixed.                                                                 |
| 2  | Delete or fill in `DOMAIN_LANGUAGE.md`                     | ✅ Done     | Filled in (and corrected again on 2026-07-26 — removed stale `.bak` refs). |
| 3  | Test for empty content write                                | ✅ Done     | Covered by `WriteIfChanged` tests.                                     |
| 4  | Test `.bak` overwrite on second write                       | ⛔ Moot     | `.bak` no longer exists.                                               |
| 5  | `FingerprintFile` with directory path                       | 🔲 Open     | Still uncovered — minor.                                               |
| 6  | `Fingerprint.Matches(nil)` test                             | 🔲 Open     | Still uncovered — minor.                                               |
| 7  | Validate final file content in `TestConcurrentWriteRACE`    | ✅ Done     | Rewritten in v0.2.0 with an integrity check.                           |
| 8  | Document `.bak` cleanup in README                           | ⛔ Moot     | `.bak` no longer exists.                                               |
| 9  | Extract benchmark data init to helper                       | ✅ Done     | `benchData` helper added in v0.3.0.                                    |
| 10 | Platform-specific build constraints / tests                 | 🟡 Partial  | `rename_windows.go` has build tags but remains **untested on Windows**. |

**Verdict:** This review is superseded for any `.bak`- or 3-arg-`Write`-related
finding. Open items (5, 6, 10) are tracked in `TODO_LIST.md` if still relevant.
