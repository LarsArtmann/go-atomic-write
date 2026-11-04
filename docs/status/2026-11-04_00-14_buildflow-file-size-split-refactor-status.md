# Status Report — BuildFlow file-size debt: source split, verification, doc sync

**Date:** 2026-11-04 00:14
**Session scope:** One task — resolve the two `buildflow -s file-size-check` warning findings (`atomicwrite.go` 391/350, `atomicwrite_test.go` 647/350) end-to-end: split, verify, keep docs honest.

---

## Scope line — what this report VERIFIED vs. what it did NOT

**Verified this session (command + artifact):**

| Claim                                    | Evidence                                                                                                                                                        |
| ---------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 8 Go files, all ≤ 282 lines (budget 350) | `wc -l` output: 151/46/76/136 (src), 282/70/161/153 (tests), 200 (`writefunc_test.go`, untouched)                                                               |
| Zero tests lost or gained                | `diff` of sorted `^func Test` names, full set (31 incl. `writefunc_test.go`), baseline = `785632f~1` (the true pre-split commit): `TEST_NAMES_IDENTICAL`        |
| Build + static analysis clean            | `go build ./... && go vet ./...` → exit 0                                                                                                                       |
| Tests pass incl. race detector           | `go test -race ./...` → `ok ... 1.023s` (fresh run, then cached on identical bytes)                                                                             |
| golangci-lint: 0 findings                | `BUILDFLOW_NO_RESULT_CACHE=1 buildflow -s golangci-lint --format finding` → `"total": 0`                                                                        |
| file-size-check: 0 findings              | `BUILDFLOW_NO_RESULT_CACHE=1 buildflow -s file-size-check --format finding` → `"findings": []`                                                                  |
| FEATURES.md citations fresh              | 11 edits applied; new line numbers from `grep -n '^func \|^var Err\|^type Fingerprint\|^const writeFunc'` on the actual files                                   |
| Website `.mdx` docs carry no stale refs  | grep over `website/src/content/docs/`: only the GENERATED `changelog.mdx` mentions `atomicwrite.go` (historical prose from CHANGELOG sync, not a line citation) |
| TODO_LIST.md release state               | Row 27: "Cut release `v0.6.1` from `[Unreleased]` (temp-file leak fix, source split, …)" — 🔴 TODO, High, S; my split rides this existing task                  |

**NOT verified, and why:**

- **`buildflow timings --regressions` never run.** My cold-cache explanation of the "+15301%" regression verdict is _asserted from the env-guard warning + baseline values_ (74 ms baselines = cached no-ops), not mechanically confirmed by the command BuildFlow itself suggests.
- **Full `buildflow` pipeline not run post-split.** I ran single steps (`golangci-lint`, `file-size-check`) with cache bypass; the whole-DAG green run (incl. cross-compile, govulncheck, drift guard) is CI's job and was not triggered.
- **`buildflow format` not run.** Formatting cleanliness is _inferred_ from golangci-lint's 0 findings (which includes the formatter checks), not measured by the formatter pass itself.
- **Website build not run.** The new CHANGELOG entry flows into the generated changelog page on next `pnpm run build`; I checked the entry has no raw `<` (the known MDX breaker) but did not build.
- **Windows/macOS cross-compile not run locally.** Split touched no platform-specific files, but the v0.5.1 lesson (Linux-green, Windows-broken at tag time) means only CI proves it.
- **Machine-level BuildFlow env issues** (dead `/mnt/buildcache` mount, cache DB VACUUM, stale binary 74a9212 vs HEAD 21a0bb6, `GOEXPERIMENT=jsonv2`) — observed via preflight warnings, out of repo scope, not acted on.
- **Benchmarks not run.** `go test -bench=. -benchmem` compiles (build passed) but was not executed; README's perf numbers were already stale before this session (old report item 47).

---

## Session narrative (what happened, compressed)

1. User ran `buildflow -s file-size-check --format finding` → 2 warnings, `fixStrategy: none` (manual fix required). Same findings were already inventoried as debt item 22 in `docs/status/2026-11-03_22-21_*` (at that time 376/591 lines — still growing, which justified splitting over suppressing).
2. Read all 4 source + 2 test files fully; identified concern seams.
3. Split `atomicwrite.go` (391) → `atomicwrite.go` (151, package doc + public API), `fingerprint.go` (46), `commit.go` (76), `stage.go` (136). Same flat single package — the "single-package layout is deliberate" decision is untouched; only file count changed.
4. Split `atomicwrite_test.go` (647) → `atomicwrite_test.go` (282, write tests + shared helpers `tempFile`/`assertFileContent`), `fingerprint_test.go` (70), `writeifchanged_test.go` (161), `commit_test.go` (153, RACE + rename-failure). Code moved verbatim; no test bodies changed.
5. Verified: build, vet, race tests, test-name identity vs pre-split commit, golangci-lint 0, file-size-check 0.
6. Doc sync: FEATURES.md 11 line citations recomputed; AGENTS.md structure table rewritten (+ added the previously missing `writefunc_test.go` row; fixed stale `commitWithVerification()` → `commitVerified`); CHANGELOG `[Unreleased]` entry; new AGENTS.md gotcha for the stale result-cache replay.
7. Discovered + documented: BuildFlow's result cache replayed the stale findings after the cause was deleted; bypass is `BUILDFLOW_NO_RESULT_CACHE=1` (documented triage in the skill, now also in project AGENTS.md Gotchas).

---

## a) FULLY DONE

| #  | Item                                                                                                                                                            | Evidence                                                                      |
| -- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------- |
| 1  | `atomicwrite.go` split into 4 concern files (391 → 151/46/76/136)                                                                                               | `wc -l`; `go build`                                                           |
| 2  | `atomicwrite_test.go` split into 4 files (647 → 282/70/161/153)                                                                                                 | `wc -l`                                                                       |
| 3  | Test integrity: 31/31 test function names identical to pre-split baseline                                                                                       | `diff <(git show 785632f~1:…)` → `TEST_NAMES_IDENTICAL`                       |
| 4  | `go build`, `go vet`, `go test -race` all green                                                                                                                 | command outputs in scope table                                                |
| 5  | golangci-lint gate: 0 issues on the new file layout                                                                                                             | BuildFlow findings JSON, cache bypassed                                       |
| 6  | file-size-check: 0 findings (the original complaint)                                                                                                            | BuildFlow findings JSON, cache bypassed                                       |
| 7  | FEATURES.md: all 11 line citations recomputed against new files (incl. 4 bare `symbol:line` forms like `commitVerified:249` → `commit.go:20`)                   | multiedit applied; `grep -n '\.go:[0-9]' FEATURES.md` now shows fresh numbers |
| 8  | AGENTS.md structure table rewritten for 12 files; missing `writefunc_test.go` row added; stale `commitWithVerification()` fixed to `commitVerified (commit.go)` | AGENTS.md diff                                                                |
| 9  | CHANGELOG `[Unreleased]` refactor entry (no API change, 31 tests verbatim)                                                                                      | CHANGELOG.md diff                                                             |
| 10 | New AGENTS.md gotcha: stale result-cache replay + `BUILDFLOW_NO_RESULT_CACHE=1` bypass                                                                          | AGENTS.md Gotchas                                                             |
| 11 | `.mdx` docs verified clean of stale internal references                                                                                                         | grep evidence in scope table                                                  |
| 12 | TODO_LIST.md checked: v0.6.1 release already tracked (row 27, High); no duplicate release item added by this session                                            | grep evidence                                                                 |

## b) PARTIALLY DONE

| # | Item                        | Done                                                                                                                | Missing                                                                                                                                                                                                                                              |
| - | --------------------------- | ------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | Regression-verdict triage   | Root cause identified (cold caches from dead `/mnt/buildcache` mount; baselines of 1–419 ms were cached no-op runs) | `buildflow timings --regressions` never executed to mechanically confirm; machine-level fix (restore mount) out of repo scope and not done                                                                                                           |
| 2 | Doc-sync sweep              | FEATURES/AGENTS/CHANGELOG/TODO_LIST/ROADMAP/README + `.mdx` all checked                                             | Old status/planning reports (historical archives) not annotated: `2026-11-03_22-21` still lists file-size as open debt item 22; `2026-07-26`'s "re-verify FEATURES citations" next-action is now done but unannotated (docs-health ANNOTATE pending) |
| 3 | The standard BuildFlow loop | Steps 3, 5, 7 executed (run, drill in, close the loop via AGENTS.md gotcha)                                         | Step 4 (`--fix` detect→repair→verify) not run — nothing was auto-fixable (`fixStrategy: none`), but the formatter pass wasn't exercised either                                                                                                       |

## c) NOT STARTED

| # | Item                                                                                                                                | Why it's listed                                                                                                   |
| - | ----------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| 1 | docs-health HARVEST of this report's section (f) into `TODO_LIST.md` / `ROADMAP.md`                                                 | Status-report skill mandates it after the report; deferred per "THEN WAIT FOR INSTRUCTIONS"                       |
| 2 | docs-health ANNOTATE of `docs/status/2026-11-03_22-21_*` (file-size item resolved) and `2026-07-26` (citation re-verification done) | Archive convention: annotate, never rewrite; not yet done                                                         |
| 3 | Machine-level BuildFlow env fixes (cache mount, VACUUM, binary rebuild, GOEXPERIMENT)                                               | Out of repo scope; awaiting user decision (see g3)                                                                |
| 4 | v0.6.1 release                                                                                                                      | Already tracked in TODO_LIST row 27 (High) — this session's CHANGELOG entry joins it; release process not started |
| 5 | Benchmark freshness (re-run `hash_bench_test.go`, refresh README numbers)                                                           | Pre-existing debt (old report item 47); untouched this session                                                    |

## d) TOTALLY FUCKED UP

**Nothing in the repository.** All gates green, no data loss, no API change, no test loss — verified, not claimed.

Honest near-misses (self-caught, fixed within one step each, listed because they were sloppy):

1. **Wrong-baseline integrity check.** First `git show HEAD:atomicwrite_test.go` comparison read the _new_ file — the auto-commit daemon had already committed mid-task, shifting HEAD under me. Wasted a step and briefly produced a meaningless "11 vs 31" comparison. Correct fix: walk history to `785632f~1` and diff against that. **Lesson: pin the pre-change commit BEFORE the first write when a commit daemon is active.**
2. **Asymmetric diff construction.** First test-name diff omitted `writefunc_test.go` from the BEFORE side → false "6 tests added", exit 1. Self-corrected next command. The conclusion ("pure move") was only drawn after the symmetric diff passed — no false claim was made in between.
3. **Added-then-removed comment.** I wrote a doc comment on `randomSuffix()` during the move (violating the pure-move discipline and the no-new-comments rule), caught it myself and removed before building.
4. **Stale-cache re-run confusion.** First `file-size-check` re-run (no bypass) printed the old 2 findings; I presented that output before diagnosing it as the documented cache-replay failure mode. One wasted step; the skill literally listed this symptom.

Machine-level, genuinely broken but outside this repo: the dead `/mnt/buildcache` mount making every regression verdict noise (env-guard warning), and the stale BuildFlow binary. See (c3)/(g3).

## e) WHAT WE SHOULD IMPROVE (brutal self-review)

**What did I forget?**

- `.mdx` in the doc-drift grep pattern (caught by post-hoc check; turned out clean, but I got lucky, not thorough).
- TODO_LIST.md — never opened during the work; the release row and the absence of a file-size row were only verified when writing this report.
- Mechanical confirmation of the timing-regression explanation (`buildflow timings --regressions`).

**What could I have done better?**

- Pin the pre-change commit before editing (daemon-aware verification discipline).
- Run the symmetric diff right the first time — verification commands deserve the same exactness as edits.
- Run `buildflow format` and a full `buildflow` pass to close the standard loop completely instead of stopping at single-step checks that were sufficient-for-claims but not exhaustive.
- Byte-level proof of the verbatim move (`git diff --color-moved=dimmed-zebra`) instead of name-identity + tests-green inference. Strong, but not airtight.

**What is stupid that we do anyway?**

- Status reports inventory debt (item 22 existed since 22:21 with the numbers already growing: 376→391, 591→647 in hours) but nothing harvests them into TODO_LIST until someone asks. The split happened because the user pasted tool output, not because the tracked process caught it.
- Machine-level BuildFlow env debt (dead cache mount, stale binary) makes every run's regression verdict noise that humans must manually discount — a standing tax on every session in every repo.

**Split brains created or found?**

- Created: none detected. FEATURES/AGENTS/CHANGELOG were updated in the same session as the split, so no window existed for drift.
- Found (pre-existing, now fixed): stale `commitWithVerification()` name in AGENTS.md Architecture section.
- Residual risk (tracked, not fixed): historical status reports still contain stale line numbers — by archive convention that is correct, but they are only annotated, never rewritten, so readers must know the convention.

**Ghost systems?** None in scope: the package is small, every function is wired and tested (31 tests name-identical, all green). No capability exists in the repo that isn't reachable from the public API.

**Did I lie?** No. Every green claim in this report maps to a command in the scope table; asserted-but-unverified items (timings, formatter pass, website build, cross-compile) are explicitly listed as such.

**Testing state?** Unchanged and green: 31 tests (incl. RACE concurrency test), standard library only, `t.Parallel()` everywhere. The refactor correctly added zero tests. Standing gap (pre-existing): no test exercises `syncDir` failure paths (POSIX dir-fsync), and Windows paths are only compile-checked by CI, never run.

## f) NEXT (up to 50, impact-sorted; brainstorm — most are ROADMAP fuel pending HARVEST)

| #  | Pri | Eff | Task                                                                                                                                                                                                                                   |
| -- | --- | --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1  | P1  | S   | Run `buildflow timings --regressions` — turn the cold-cache regression explanation from assertion into measurement                                                                                                                     |
| 2  | P1  | S   | HARVEST this report's (f) list into `TODO_LIST.md` / `ROADMAP.md` (docs-health HARVEST mode)                                                                                                                                           |
| 3  | P1  | S   | ANNOTATE `docs/status/2026-11-03_22-21_*`: debt item 22 (file-size 376/591) resolved by this split (0 findings, evidence in scope table)                                                                                               |
| 4  | P1  | S   | Execute the existing TODO_LIST row 27: cut release `v0.6.1` from `[Unreleased]` (now also carries the split entry); CI green on master incl. Windows cross-compile BEFORE tagging (v0.5.1 lesson)                                      |
| 5  | P2  | S   | Full `buildflow` fast-mode run post-split — whole-DAG green with current files, warms caches, resets the noise                                                                                                                         |
| 6  | P2  | S   | `buildflow format` once — measure formatting cleanliness instead of inferring it from golangci-lint                                                                                                                                    |
| 7  | P2  | S   | Byte-level verbatim-move proof: review `git diff --color-moved=dimmed-zebra 785632f~1..HEAD -- '*.go'`; archive verdict in a session note                                                                                              |
| 8  | P2  | S   | Website build check: `cd website && pnpm run sync:changelog && pnpm run build` — new CHANGELOG entry must render on the generated changelog page                                                                                       |
| 9  | P2  | S   | Decide file-size policy (needs user, see g1): keep default 350 + split-on-breach as standing practice, or override budget in `.buildflow.yml` with rationale                                                                           |
| 10 | P2  | S   | Decide whether file-size warnings should gate in CI (`--strict` promotion) or stay advisory (policy)                                                                                                                                   |
| 11 | P2  | S   | CONTRIBUTING.md: document the 350-line budget + "split along concern seams" so contributors don't re-accumulate the debt                                                                                                               |
| 12 | P2  | S   | ANNOTATE `docs/status/2026-07-26_*`: its next-action "re-verify FEATURES.md citations" is done (11 citations recomputed this session)                                                                                                  |
| 13 | P2  | M   | Test gap (pre-existing): no test covers `syncDir` failure path (POSIX dir-fsync failure mid-commit) — decide if injectable without SE_COMBINE hacks or accept as platform-untestable and document                                      |
| 14 | P2  | M   | Test gap (pre-existing): Windows rename-retry loop (`ERROR_ACCESS_DENIED`/`ERROR_SHARING_VIOLATION` backoff) has no test; CI only compile-checks it — consider a table-driven unit test for the backoff policy function if extractable |
| 15 | P3  | S   | Machine: restore `/mnt/buildcache` mount or fix cache env (env-guard) — kills the standing regression noise                                                                                                                            |
| 16 | P3  | S   | Machine: `sqlite3 ~/.cache/buildflow/cache.db VACUUM` (0.22 GB, 54% free pages); review the 0.71 GB state DB                                                                                                                           |
| 17 | P3  | S   | Machine: rebuild BuildFlow (`nix build . && nix run .#reinstall`) — binary at 74a9212 vs HEAD 21a0bb6; current results may lag BuildFlow fixes                                                                                         |
| 18 | P3  | S   | Machine: decide `GOEXPERIMENT=jsonv2` — unset globally if no project needs it                                                                                                                                                          |
| 19 | P3  | S   | Machine: `buildflow format` to re-verify cached verdicts after the templ-fmt formatter-fingerprint drift                                                                                                                               |
| 20 | P3  | S   | After cache mount restored: confirm regression verdicts reset to clean baselines                                                                                                                                                       |
| 21 | P3  | S   | `go test -bench=. -benchmem` sanity run: benchmarks compile (build-verified) but haven't executed post-split                                                                                                                           |
| 22 | P3  | S   | Refresh benchmark numbers in README after 21 (old report item 47 — still open, pre-existing)                                                                                                                                           |
| 23 | P3  | S   | Verify whether BuildFlow `pre-commit` mode includes file-size-check (don't assume); if yes, confirm hook behaves post-split; if no, note the gap                                                                                       |
| 24 | P3  | S   | Grep `.github/` for any hard-coded test/source file names that the split could have broken (workflows covered by earlier grep, but verify explicitly)                                                                                  |
| 25 | P3  | S   | Confirm `go mod tidy -diff` clean and the version-trio drift guard still passes (untouched this session, but cheap to verify together with item 5)                                                                                     |
| 26 | P3  | S   | AGENTS.md Gotchas: evaluate whether the stale-cache gotcha and the skill's documented bypass create a duplicate-source-of-truth; keep one canonical wording (docs-health VERIFY pass)                                                  |
| 27 | P3  | S   | Sweep remaining docs (`docs/guides`, `getting-started`) for internal `symbol:line` forms (the session grep covered `\.go:` and bare forms in FEATURES only — full-repo VERIFY pass belongs to docs-health)                             |
| 28 | P3  | S   | Consider a one-line ADR-style note (in CHANGELOG entry it's covered; only add if release notes need the rationale repeated — avoid duplication)                                                                                        |
| 29 | P3  | S   | Check `api-reference.mdx` renders the same public API post-split (doc source is handwritten; symbol set unchanged — verify by diff of exported names)                                                                                  |
| 30 | P3  | S   | Exported-symbol set proof: `go doc -all` diff pre/post split to mechanically confirm zero API change (CHANGELOG claims it; measure it)                                                                                                 |
| 31 | P3  | M   | Standing process improvement: auto-harvest BuildFlow findings into TODO_LIST (the 22:21 report's item 22 aged ~2h before action) — BuildFlow-dev question, fleet-wide (see skill: extension belongs upstream)                          |
| 32 | P3  | M   | If file-size breaches recur: add the file-size baseline to the periodic status report template so debt is re-inventoried automatically                                                                                                 |

(Items considered and rejected: renaming `hash_bench_test.go` (gratuitous, referenced in AGENTS/README); renaming `writefunc_test.go` (already mirror-named); suppressing the 350 budget instead of splitting (rejected — files were still growing and the seams were real).)

## g) Questions I cannot figure out myself

1. **File-size policy:** Is "keep BuildFlow's default 350-line budget and split along concern seams on breach" now the standing practice for this repo (and fleet-wide), or do you want an explicit budget override with rationale in `.buildflow.yml`? The 22:21 report called the single-package layout deliberate — I preserved it (same package, more files) but you may prefer the budget route for future cases.
2. **Release:** TODO_LIST row 27 says cut `v0.6.1` from `[Unreleased]` and that "the fleet bump depends on it". `[Unreleased]` now additionally carries this refactor. Do you want me to start the release checklist now (CI green incl. Windows → changelog section → annotated tag), or wait for more `[Unreleased]` accumulation?
3. **Machine env debt:** The dead `/mnt/buildcache` mount, cache-DB VACUUM, stale BuildFlow binary, and `GOEXPERIMENT=jsonv2` are machine-level and affect every repo. Should I fix them in this session (the commands are known), or do you handle BuildFlow machine maintenance separately in the BuildFlow repo?

---

_Report per status-report skill; session-scoped per user instruction. Auto-commit daemon will pick this file up (manual commit skipped per harness contract)._
