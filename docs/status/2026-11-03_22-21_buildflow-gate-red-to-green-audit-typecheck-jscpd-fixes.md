# Status Report — BuildFlow Gate Red→Green: Audit, Type-Check, jscpd Fixes + Upstream Filings

**Timestamp:** 2026-11-03 22:21 CET
**Repo:** `github.com/larsartmann/go-atomic-write` @ master (HEAD at write time: c0767f4, auto-commit daemon active)
**Session window:** 2026-11-03 ~21:15 → 22:15 CET
**Trigger:** `buildflow --fix --build-mode=full` exited **69** — 34 error-severity findings (pnpm-audit 29, type-check 4, jscpd 1).

## Scope: what this report VERIFIED vs did NOT verify

**Verified with commands/artifacts this session:**

- 3 full `buildflow --fix --build-mode=full` runs (21:47 → EXIT 69 stale-cache; 22:12 → **EXIT 0**) + targeted single-step runs (`-s <tool> --format finding`)
- `astro check` via `pnpm run typecheck`: **0 errors / 0 warnings, 33 files** (twice: before and after all dependency changes)
- `pnpm audit`: 62 findings → **1 moderate** (accepted, documented); lockfile contains `basic-ftp@6.2.5`, `glob@14.0.3`, `minimatch@11.0.4`
- `pnpm install --frozen-lockfile` ✓, `go mod tidy -diff` ✓, `go build ./...` ✓ (CI drift-guard equivalents)
- Require-chain smoke tests via `createRequire` into `.pnpm` store contexts: basic-ftp 6.2.5, get-uri, proxy-agent, glob 14.0.3 (globSync exercised), chrome-launcher 0.13.4 + 1.2.2, minimatch, lighthouse — all load and function
- jscpd glob experiments: 4 pattern variants tested directly against jscpd 4.0.9 — only `**/dist/**` matches
- BuildFlow source reads: `tools/providers/js_providers.go:21` (generic `tsc --noEmit`), `js_extra_checks.go:12` (vue-tsc special case exists, no Astro), `jscpd_provider*.go` (BuildFlow#48 merge, `--gitignore` only passed for major <4), delta `74a9212..21a0bb6` inspected (no type-check/astro fixes)
- GitHub issues **#84, #85, #86 filed and verified landed** (`gh issue view`, all OPEN); drafts passed `check-draft.go` 0 FAIL / 0 WARN

**NOT verified (and why):**

- `astro build` / `pnpm run validate` / lighthouse end-to-end — not run this session; dependency graph changed under the website, so the postbuild verify-gates are unproven against the new lockfile
- Effectiveness of `exclude: website/dist` in `.buildflow.yml` per tool — jscpd provably ignored it (needed `.jscpd.json`); go-structure-linter still reported its 2 dist-driven warnings in the final green run
- The `typescript@>=7` override guard (added by the concurrent session) — never simulated a resolution to prove it folds `>=7` back to 6
- oxlint's **323 findings** — never drilled into file distribution (dist? node_modules? sources?)
- The 9–10 "tools unavailable" health-check failures — `--verbose` never run
- govulncheck's 7 warnings (Go 1.27 toolchain-internal paths) — untriaged
- Concurrent session's parallel edits (README benchmark rewrite, CHANGELOG.md, TODO_LIST.md, `lighthouserc.deployed.json`) — present in `git status`, diffs not reviewed by me
- CI runs — push-dependent; nothing pushed (nothing was pushed per harness rules)

**Context constraint:** a second agent session worked this repo concurrently the whole time (its 20:42 status report is `2026-11-03_20-42_no-linter-submodule-and-rename-failure-temp-leak-fix.md`). My edits collided with its writes twice (both caught by the read-before-edit gate). All its changes were read and judged on merits before being accepted, never reverted.

---

## a) FULLY DONE

| #  | Item                                                                                                                   | Evidence                                                                                                                                                                                                                          |
| -- | ---------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1  | **BuildFlow findings gate red → green**: exit 69 → **EXIT 0**                                                          | Final full run: "BuildFlow passed with warnings 56/63, 7m44s, 363 findings"                                                                                                                                                       |
| 2  | **jscpd error fixed** (dist-vs-public duplication)                                                                     | `.jscpd.json` `ignore: ["**/dist/**"]` (unions over BuildFlow defaults per BuildFlow#48); targeted run: `"findings": []`                                                                                                          |
| 3  | **type-check 4 errors resolved**                                                                                       | `skip_steps: [type-check]` in new `.buildflow.yml` with rationale comment; justified by astro check 0-errors + `website.yml:46` already gating on it + BuildFlow source showing generic `tsc --noEmit` (no Astro provider exists) |
| 4  | **pnpm-audit 29 errors → 0**                                                                                           | Overrides in `website/pnpm-workspace.yaml`: `basic-ftp ^6.2.1` (→6.2.5), `glob@>=13 <14 → ^14.0.2` (→14.0.3), `minimatch@>=10 <11 → ^11.0.1` (→11.0.4); fresh audit shows only the accepted moderate                              |
| 5  | **Discovered pnpm 11 config split-brain** and documented it                                                            | `overrides` in package.json is silently ignored (even `--force` no-op'd, 71ms "Already up to date"); workspace.yaml is authoritative; dead package.json block deleted                                                             |
| 6  | **CI drift guards verified locally**                                                                                   | `pnpm install --frozen-lockfile` ✓, `go mod tidy -diff` ✓, `go build ./...` ✓                                                                                                                                                     |
| 7  | **TypeScript ping-pong reverted** (3 occurrences today)                                                                | bd2f100 pinned ^7.0.2→^6.0.3; buildflow's pnpm-update re-bumped (3cf0176); re-reverted; concurrent session added `"typescript@>=7": ^6.0.3` override as install-graph guard                                                       |
| 8  | **Root `npm init -y` stub + lockfile-less Dependabot entry removed** (2nd recurrence of the documented 09-22 incident) | Stub trashed, `npm directory: /` entry removed per AGENTS.md doctrine; gotcha updated                                                                                                                                             |
| 9  | **Stale `global.out.css` gone + regression-proofed**                                                                   | Concurrent session deleted it (1625 lines) and added `src/styles/*.out.css` to `website/.gitignore` — verified unreferenced before accepting                                                                                      |
| 10 | **`.gitignore` comment fixed on sight**                                                                                | "CI uses npm" → CI is pnpm-only (verified in both workflows: corepack pnpm@11.20.0)                                                                                                                                               |
| 11 | **AGENTS.md: 4 gotchas written**                                                                                       | pnpm-overrides ownership, audit baseline + accepted-risk rationale, BuildFlow scanning policy (`.buildflow.yml`/`.jscpd.json`/type-check skip), root-manifest recurrence                                                          |
| 12 | **Upstream closed-loop: #84 #85 #86 filed, verified**                                                                  | pnpm-update range-crossing; tsc-vs-astro-check; audit result-cache staleness — all with source-level verification, drafts in `BuildFlow/docs/drafts/`                                                                             |
| 13 | **Overridden package chain smoke-tested**                                                                              | All load; glob 14 named-export shape confirmed functional (`globSync` exercised)                                                                                                                                                  |

## b) PARTIALLY DONE

1. **pnpm-audit at "1 accepted moderate"** — `minimatch@3.1.5` (globparts DoS) via `@lhci/cli@0.16.0 → chrome-launcher@0.13.4 → rimraf@3 → glob@7`. Every fix path verified broken at source level (rimraf v6 removes the callback API chrome-launcher calls: `this.rimraf(dir, cb)`; minimatch 11 breaks glob 7; latest lhci still pins chrome-launcher ^0.13.4). Documented non-fix in AGENTS.md; revisit when lhci ships chrome-launcher ^1.
2. **`exclude: website/dist` in `.buildflow.yml`** — kept on principle (gitignored build output shouldn't be scanned), but only PROVEN inert for jscpd (which needed `.jscpd.json`) and ineffective for go-structure-linter's 2 warnings (still reported in the green run). Its actual reach is unmeasured.
3. **TS ping-pong protection** — the `typescript@>=7` override guards the _installed graph_, not the _manifest string_; pnpm-update will likely keep rewriting `package.json` until BuildFlow #84 lands. Manual reverts remain the interim reality.
4. **BuildFlow binary freshness** — preflight warned binary (74a9212) vs HEAD (21a0bb6); I inspected the delta for relevance (none for my errors) but never rebuilt/reinstalled. Concurrent session keeps committing, so drift grows.
5. **BuildFlow result-cache staleness** — root-caused and filed (#86), workaround identified (`BUILDFLOW_NO_RESULT_CACHE=1`), but CI may hit the same phantom-red until #86 lands.
6. **Concurrent-session coordination** — handled reactively (re-read before every edit); no protocol, no lock, two collided edits (no damage).

## c) NOT STARTED

1. oxlint's **323 findings** — zero triage (file distribution unknown; possibly dist/node_modules noise that `exclude` should kill)
2. art-dupl 6 findings (10-token test-helper clones), file-size-check 2 (atomicwrite.go 376 > 350; atomicwrite_test.go 591 > 350), doc-files-age 4 (README/AGENTS/FEATURES/TODO_LIST freshness) — all warning-severity, untouched
3. `buildflow doctor` full pass / the 9–10 unavailable-tools health failures
4. govulncheck 7 warnings (all toolchain-internal paths under `/mnt/buildcache/...go1.27.0`) — untriaged
5. Website production build + `validate` + lighthouse budgets after the dependency changes
6. HARVEST of section (f) into `TODO_LIST.md` (skill mandates it once the session continues; I am waiting per instructions)
7. CI green confirmation for the new lockfile/overrides (needs the next workflow run)
8. `pnpm peers check` warning (`@puppeteer/browsers@3.2.6` wants `proxy-agent >=8.0.1`, has 6.5.0) — pre-existing, noted, not triaged
9. jscpd v4+ `--gitignore` gap (BuildFlow passes it only for major <4; jscpd 4.0.9 scanned gitignored dist) — candidate upstream issue, not filed (needs one more isolation test)

## d) TOTALLY FUCKED UP (my mistakes this session — honest list)

1. **Analytical error:** initially judged `minimatch@3.1.5` "below the vulnerable range, not vulnerable" — wrong: the range is `<11.0.1`, which includes 3.1.5. Caught two steps later when the audit still listed it. Reading advisory ranges backwards is exactly the verify-before-filing failure class.
2. **Unverified config guess:** wrote `.jscpd.json` with `website/dist/**` without testing pattern semantics first — provably not matching; cost a full extra detect cycle. Should have run the 4-variant experiment before writing anything.
3. **Fell into the pnpm split-brain trap I then documented:** added overrides to package.json (the file pnpm 11 ignores) and only discovered it via a suspiciously fast no-op `--force`. The lockfile's `overrides:` block had been telling me the truth the whole time.
4. **Two edit-before-read gate trips** — both caused by concurrent-session writes landing between my read and edit; zero damage, but I should have re-checked mtimes proactively given the live-session signals (daemon commits at 20:50/20:54/21:03).
5. **Almost reverted concurrent work:** my first instinct on `global.out.css` and the TS pin was to fix them myself; in both cases the concurrent session had already done it correctly. The safety rules (read → judge → don't touch what you didn't author) prevented sabotage — but my reflex was wrong twice.
6. **Nothing catastrophic:** no data loss, no reverted third-party work, no pushed commits, no broken build introduced. The tree at session end is gate-green and drift-guard-clean.

## e) WHAT WE SHOULD IMPROVE

1. **Test tool semantics in isolation BEFORE writing config** — the jscpd glob lesson generalizes: every config file format (globs, override syntax, ignore semantics) gets a 30-second direct experiment before it enters the repo.
2. **Never trust a manifest when the tool has a different config home** — pnpm 11 moved `overrides` to pnpm-workspace.yaml; the lockfile's recorded `overrides:` section is the ground truth of what the resolver saw. Check the effective state, not the file you expect.
3. **Verify "done" with the same detector that failed** — the stale result-cache episode means `buildflow -s X` after a fix must bypass the cache (`BUILDFLOW_NO_RESULT_CACHE=1`) until #86 lands; otherwise you re-read a verdict, not reality.
4. **Stop treating warnings as invisible** — "passed with warnings 56/63 (363 findings)" hides real debt (323 oxlint findings!). The gate only sees errors; a weekly warnings sweep or `warnings_budget` baselines would keep the warning count honest instead of ambient.
5. **Act on preflight advisories in background** — binary-freshness and DB-size warnings were noted and dropped; both are one-command fixes that should have run concurrently with the main work.
6. **Concurrent-session protocol** — two agents raced on the same manifests today with two near-misses. A cheap convention (check `git log --since=5 minutes` + mtimes before any manifest edit) should be an AGENTS.md bullet, not tribal memory.
7. **File issues with the mechanical checker every time** — the check-draft gate (0 FAIL/0 WARN on all three) is fast and catches register mistakes; it should never be skipped, even for "obvious" filings.
8. **When a fix path is broken at source level, write it down immediately** — the chrome-launcher/rimraf/minimatch incompatibility chain took four source reads; if it isn't in AGENTS.md, the next session re-derives it (now it is).

## f) TOP 50 THINGS TO GET DONE NEXT

_Sorted roughly by impact; first block = this report's direct consequences. Items 1–2 are skill-mandated loop-closers._

| #  | Task                                                                                                                                                                  | Why / source                                                                                |
| -- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| 1  | HARVEST sections (b)(c)(f) into `TODO_LIST.md` (docs-health HARVEST mode)                                                                                             | Skill mandate: report findings must not be entombed in timestamped files                    |
| 2  | Verify `git status` → daemon committed today's manifest+lockfile changes TOGETHER (split-commit = CI drift-guard break)                                               | Documented daemon hazard; dependabot/AGENTS/package-deletion were uncommitted at last check |
| 3  | Run `astro build` + `pnpm run validate` + postbuild verify-gates against the NEW lockfile                                                                             | Dependency graph changed; build gates unproven (verify-build.mjs fails loudly on drift)     |
| 4  | Make CI immune to the stale-audit-cache until #86 lands: set `BUILDFLOW_NO_RESULT_CACHE=1` for audit steps in CI, or pin buildflow + document                         | #86: phantom-red gate on a clean repo                                                       |
| 5  | Drill oxlint's 323 findings (`buildflow -s oxlint --format finding`) — classify dist/node_modules noise vs real                                                       | Largest single warning block; possibly killed by `exclude`                                  |
| 6  | Measure `exclude: website/dist` per tool; delete the key if provably inert (config-honesty)                                                                           | b.2; config lint reports stale mirrors — don't keep a dead key                              |
| 7  | Fix go-structure-linter's 2 remaining warnings (assets//internal/ from dist chunks) — verify whether `exclude` reaches it, else route via tool config                 | Green run still reports them                                                                |
| 8  | Constrain or skip `pnpm-update` in `.buildflow.yml` (`dep_update_mode`) until #84 lands                                                                               | Manifest keeps getting rewritten out-of-range                                               |
| 9  | Simulate the `typescript@>=7 → ^6.0.3` override guard (dry resolution) to prove it works                                                                              | b.4; unverified protection                                                                  |
| 10 | Find what creates the root `package.json` stub (2× today) — fix the source, not the symptom                                                                           | g.1; twice in one day means a tool does it                                                  |
| 11 | Rebuild + reinstall BuildFlow binary (`nix build . && nix run .#reinstall`) — 74a9212 vs 21a0bb6+ and growing                                                         | Preflight advisory, ignored                                                                 |
| 12 | `buildflow doctor` full pass; resolve the 9–10 unavailable tools (`--verbose` for list)                                                                               | Health-check noise every run                                                                |
| 13 | Triage govulncheck's 7 warnings (Go 1.27 toolchain internals — likely false positives; confirm and document or fix env)                                               | Warning debt from both red runs                                                             |
| 14 | Decide `lighthouserc.deployed.json` (untracked, concurrent session's): commit, gitignore, or trash                                                                    | Ambiguous artifact in tree                                                                  |
| 15 | Review concurrent session's README benchmark-table rewrite (unreviewed diff in tree)                                                                                  | Unreviewed content changes shipped by daemon                                                |
| 16 | Review CHANGELOG.md / TODO_LIST.md / changelog.mdx concurrent edits (all modified, unreviewed)                                                                        | Same class as 15                                                                            |
| 17 | Run lighthouse budget check (chromium via nix devShell) against the new dependency tree                                                                               | LHCI never exercised post-changes                                                           |
| 18 | Evaluate chrome-launcher@1.2.2 override WITH a real lighthouse run → kill the accepted minimatch moderate                                                             | g.3; turns documented non-fix into a fix                                                    |
| 19 | `pnpm peers check` follow-up: proxy-agent 6.5.0 vs @puppeteer/browsers wanting >=8.0.1 — upgrade path or document                                                     | c.8; pre-existing peer mismatch in lhci chain                                               |
| 20 | Review `minimumReleaseAgeExclude` list (tailwind 4.3.4 era entries) — still accurate after updates?                                                                   | Stale excludes weaken supply-chain policy                                                   |
| 21 | art-dupl 6 findings (10-token test clones) — deduplicate helpers or accept + budget                                                                                   | Warning debt                                                                                |
| 22 | file-size-check: atomicwrite.go 376/350, atomicwrite_test.go 591/350 — split files or set explicit budget with rationale                                              | Warning debt; single-package layout is deliberate                                           |
| 23 | doc-files-age: refresh FEATURES.md / TODO_LIST.md (AGENTS.md/README clear via today's writes; the others ~1100–1488h stale)                                           | Freshness drift                                                                             |
| 24 | VACUUM `~/.cache/buildflow/cache.db` (0.22 GB, 56% free pages); assess 0.71 GB state DB purge                                                                         | Preflight warning, one command                                                              |
| 25 | Unset `GOEXPERIMENT=jsonv2` globally if no project needs it                                                                                                           | Preflight warning every run                                                                 |
| 26 | Add AGENTS.md bullet: use `BUILDFLOW_NO_RESULT_CACHE=1 buildflow -s pnpm-audit` after audit fixes until #86                                                           | So the next session doesn't re-derive it                                                    |
| 27 | Isolate-and-file the jscpd v4 gitignore gap upstream (does jscpd 4 respect --gitignore at all? BuildFlow only passes it <4)                                           | Candidate issue; needs one isolation test                                                   |
| 28 | Runtime-test `.jscpd.json` union merge with a second ignore entry (merges verified in source, not runtime)                                                            | Trust-before-reliance                                                                       |
| 29 | Confirm next Dependabot run passes for `/website` (new overrides + lockfile)                                                                                          | Config changed under Dependabot's feet                                                      |
| 30 | Confirm next Website workflow (typecheck/build/deploy) goes green post-changes                                                                                        | Firebase deploy rides it; silent-freeze history                                             |
| 31 | Controlled dependency-update pass (astro 7.4.5, starlight 0.42.8, tailwind 4.3.4) done deliberately instead of via racing auto-updaters                               | Updates happened adversarially today                                                        |
| 32 | Add concurrent-session protocol bullet to AGENTS.md (mtime + recent-commits check before manifest edits)                                                              | Two near-misses today                                                                       |
| 33 | Warnings-baseline policy: set `warnings_budget` for accepted noise (go-auto-upgrade style) or sweep weekly                                                            | 363 ambient findings erode signal                                                           |
| 34 | Independent confirmation run: `go test -race ./...` + `golangci-lint run ./...` directly (buildflow ran them; direct run = belt and braces for the release checklist) | Release checklist step 1/3                                                                  |
| 35 | Sweep stale `//nolint` comments (AGENTS.md says only hash_bench `makezero` should remain)                                                                             | Doc/code drift check                                                                        |
| 36 | Check whether CI's drift-guard trio check needs the AGENTS.md Go-version line updated (untouched today — verify still in sync)                                        | Trio drift has bitten 4×                                                                    |
| 37 | File/follow the pnpm peers issue upstream (lhci→puppeteer proxy-agent) if #19 confirms unfixable locally                                                              | Chain hygiene                                                                               |
| 38 | Fleet-share the pnpm-workspace.yaml overrides pattern (other LarsArtmann Astro sites likely carry the same 29 audit errors)                                           | The fixes are copy-pasteable; #86/#84 make them necessary                                   |
| 39 | Follow up #84/#85/#86: tick MANUALLY REVIEWED after reading; link any duplicates found later                                                                          | Filing ledger hygiene                                                                       |
| 40 | Verify `data-copy`/verify-build gates still pass after any future hero/website edits (they gate the build)                                                            | Guard-rail documentation check                                                              |
| 41 | Assess whether `website/dist` should be excluded from the link scan too (39 files scanned; 12 under excluded paths)                                                   | Link-scan scope hygiene                                                                     |
| 42 | Update `docs/status/2026-11-03_20-42_*` (concurrent session's report) if today's late findings invalidate its claims                                                  | Cross-session consistency                                                                   |
| 43 | Consider pinning buildflow version in CI vs fleet drift (74a9212 binary vs 21a0bb6 HEAD vs CI's version)                                                              | Reproducibility                                                                             |
| 44 | Decide whether `.buildflow.yml` should also carry `env:` documentation for GOEXPERIMENT (per-project env policy)                                                      | Config hygiene                                                                              |
| 45 | Run `pnpm run sync:changelog` + inspect generated changelog.mdx after CHANGELOG edits settle                                                                          | Generated-file drift                                                                        |
| 46 | Audit remaining `overrides` entries (brace-expansion ^5.0.12, fast-uri ^3.1.9) — still needed? When can they be dropped?                                              | Override debt decays badly                                                                  |
| 47 | Check `git town` state / sync status (master, git-town.toml — untouched all session)                                                                                  | VCS hygiene                                                                                 |
| 48 | Sweep for other nested `.gitignore` rules invisible to fleet tools (jscpd proved nested-ignores don't propagate)                                                      | Same failure class, other tools                                                             |
| 49 | Evaluate `code-quality-scan` HTML dashboard for the 363-finding warning backlog (visual triage)                                                                       | Debt visualization                                                                          |
| 50 | Re-run the full buildflow gate tomorrow to confirm cache/warm-state stability of EXIT 0                                                                               | Guard against one-off green                                                                 |

## g) QUESTIONS I CAN NOT FIGURE OUT MYSELF

1. **The root `package.json` stub regenerated TWICE today** (09-22 incident pattern, this time committed alongside a Dependabot entry). Do you know which tool, agent, or habit runs `npm init -y` (or equivalent) at this repo's root — so we fix the source instead of deleting symptoms every time?
2. **Until BuildFlow #84 ships, how should pnpm-update be handled here?** Options I see: skip the step repo-wide (loses legitimate updates), set `dep_update_mode` to constrain it (I haven't verified which modes exist), or keep manual reverts. Which tradeoff do you want?
3. **Is the accepted `minimatch@3.1.5` moderate a standing risk you accept**, or should I spend a chromium lighthouse run verifying chrome-launcher@1.2.2 so the override can kill it? (The run needs the nix devShell chromium and a few minutes; result would let me replace the documented non-fix with a real fix.)

---

_Report written 2026-11-03 22:21 CET by the 21:15–22:15 session. Point-in-time snapshot; concurrent session may have advanced the tree since the last command ran. Waiting for instructions._
