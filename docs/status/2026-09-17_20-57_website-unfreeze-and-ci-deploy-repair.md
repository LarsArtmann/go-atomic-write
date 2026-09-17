# Status Report — Post-Release Recovery: Website Unfreeze & CI Deploy Repair

**Date:** 2026-09-17 20:57 CEST
**Repo:** `github.com/LarsArtmann/go-atomic-write` (`master`, working tree clean)
**Session scope:** Follow-up to the v0.5.2 release session (`docs/status/2026-09-17_14-22_v0.5.2-release-session.md`). The user's directive: *"just use the local CLI and fix it yourself"* — i.e. stop waiting for the 3 open questions and resolve what was resolvable locally.

---

## TL;DR

The website's 2-month deployment outage is **over**. The missing
`FIREBASE_SERVICE_ACCOUNT_LARS_SOFTWARE` secret — the root cause of the freeze
since 2026-07-26 — was minted, registered, and verified end-to-end in CI. The
live site now serves v0.5.2 content. All CI is green, the harvest backlog is in
`TODO_LIST.md`, and 3 of the previous session's open questions are closed.

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| 1 | **Website deployed via local Firebase CLI** — `nix shell nixpkgs#nodejs nixpkgs#firebase-tools -c firebase deploy --only hosting:atomicwrite --project lars-software`; release complete, 64 files | Firebase CLI output; https://atomicwrite.web.app |
| 2 | **Live site verified to show v0.5.2 changelog** (was frozen at 0.4.0 content since 2026-07-26) — confirmed twice by fetching `/changelog` | `fetch https://atomicwrite.web.app/changelog` |
| 3 | **GitHub secret registered** — new JSON key minted for `github-website-deploy@lars-software.iam.gserviceaccount.com` (key `8020dd42…`), set via `gh secret set FIREBASE_SERVICE_ACCOUNT_LARS_SOFTWARE`, local keyfile then `shred -u`'d | `gh secret list` shows the secret, created 2026-09-17T18:31:30Z |
| 4 | **CI deploy job fixed and green** — first run failed with `firebase.json file not found`; added `entryPoint: website` to the `FirebaseExtended/action-hosting-deploy` step; deploy job now passes | run 35259719876, commit `c30d40b` |
| 5 | **`workflow_dispatch` trigger added to the Website workflow** — manual deploys now possible; deploy condition extended so dispatches actually deploy (`(push \|\| workflow_dispatch) && master`) | `.github/workflows/website.yml`, commit `25ae518` + `c30d40b` |
| 6 | **AGENTS.md Phase 8 gotcha refresh completed** (left pending by the release session) — removed the two stale `//nolint:gosec` entries, added: golangci-lint-action must be v7+ with pinned version; `astro check` needs TS 6.x; auto-commit daemon commits manifests without lockfiles; deploy secret requirement + local-deploy recipe; `firebase-adminsdk` SA key-cap; `exhaustruct_v5` migration pending | `AGENTS.md` Gotchas, commit `25ae518` |
| 7 | **CHANGELOG `[Unreleased]` backfilled** — TS6 pin, lint-action v6→v7, and the secret-freeze fix are now documented (the v0.5.2 tag itself is immutable; its gap is mitigated) | `CHANGELOG.md`, commit `25ae518` |
| 8 | **Cross-compile gate run retroactively** — `GOOS=windows`, `GOOS=darwin GOARCH=arm64`, `GOOS=linux` all build; `golangci-lint run ./...` → `0 issues` | local run this session |
| 9 | **docs-health HARVEST executed** — 15 bounded tasks routed into `TODO_LIST.md`, 2 long-term ideas into `ROADMAP.md`; resolved items (secret, live-site verify, gotcha refresh, changelog backfill, LSP restart) verified against code and excluded | `TODO_LIST.md`, `ROADMAP.md`, commit `b102ee6` |
| 10 | **Stale LSP diagnostics cleared** — `lsp_restart` on `gopls` + `golangci_lint_ls` (previous session's Q pending item #19) | this session |
| 11 | **Root-cause intelligence gathered** — `firebase-adminsdk` SA has 12 keys (over GCP's 10-key cap, which is exactly why `keys.create` returned `FAILED_PRECONDITION`); IAM on `lars-software` is only visible to `lartyhd@gmail.com`, not `lars@helpless.ai`; the old `FIREBASE_SERVICE_ACCOUNT` key recipe was found in gogenfilter's docs | `gcloud iam service-accounts keys list` output; `gogenfilter/docs/status/2026-07-20_14-22_ci-readme-recovery.html` |
| 12 | **All CI green on every pushed commit this session** (35259592878/2718, 35259719830/876, 35260101584) and a Dependabot grouped actions-bump PR appeared and is green | `gh run list` |

## b) PARTIALLY DONE

| # | Item | Gap |
|---|------|-----|
| 1 | **v0.5.2 CHANGELOG completeness** — the `[Unreleased]` section now documents the two fixes missing from the 0.5.2 section, but the 0.5.2 section itself can never be amended (tag immutable). Full closure only when `[Unreleased]` ships as 0.5.3 (or folds into a later tag) | Release the `[Unreleased]` block |
| 2 | **Website deploy hardening** — CI deploy works, but the secret is a long-lived SA JSON key (never expiring, `9999-12-31`) stored in GitHub; key rotation and workload-identity federation are not addressed | See (e) #4 |
| 3 | **The previous session's open question #2** (what wrote `go 1.27.1` into go.mod at ~09:29, auto-commit `77b79cc`) — still unanswered, but now converted into a tracked TODO with evidence pointers instead of a dangling question | `TODO_LIST.md` Medium row |
| 4 | **Node.js 20 deprecation warnings** in the Website workflow — actions `checkout`/`setup-node`/`upload-artifact`/`download-artifact` are being force-run on Node 24; they work today but the pinned SHAs need bumping when compatible versions exist | CI annotations; TODO_LIST Low row |

## c) NOT STARTED

| # | Item | Note |
|---|------|------|
| 1 | **npm vulnerability sweep** — 13 Dependabot alerts remain (10 high, 2 moderate, 1 low; astro/extract-zip/js-yaml/sharp/svgo/tmp/uuid). A Dependabot grouped bump PR for actions is open and green, but the npm alerts are untouched | Highest-impact open work |
| 2 | **Cross-platform CI build matrix** (`GOOS=windows`/`darwin`) — verified manually this session, but CI still only builds linux/amd64. v0.5.1's entire content was a Windows-only break Linux CI cannot see | |
| 3 | **`govulncheck` CI job** — not present | |
| 4 | **Lockfile/tidy drift guard** (`pnpm install --frozen-lockfile` + `go mod tidy -diff` gate) — the daemon-drift tripwire doesn't exist yet | |
| 5 | **Dependabot `npm_and_yarn` run failure** (#1579905290) — not investigated | |
| 6 | **Consumer propagation** — pkg.go.dev still shows "Imported by: 0"; no LarsArtmann repos bumped to v0.5.2 | |
| 7 | **`exhaustruct` → `exhaustruct_v5` migration** — deprecated since golangci-lint v2.13.0; the warning prints on every lint run (`.golangci.yml:42,151`) | |
| 8 | Everything else in the Low tier: README benchmark header still says "Go 1.26.3" (`README.md:217`), `minimumReleaseAgeExclude: astro@7.3.3` still present, `GOTOOLCHAIN` unpinned in CI, version-trio consistency gate, release checklist in AGENTS.md, post-deploy Lighthouse, Node-20 SHA tracking | All harvested into `TODO_LIST.md` |

## d) TOTALLY FUCKED UP!

**No new disasters this session.** What follows is honest accounting of mistakes made *during* this session:

1. **Blind-trialed the wrong service account first.** I minted a key for `firebase-adminsdk-dwv0a` (copying the 2026-07 gogenfilter recipe) before listing keys — it failed with `FAILED_PRECONDITION` because that SA already has 12 keys (over GCP's 10 cap). A 5-second `keys list` *before* `keys create` would have routed me straight to `github-website-deploy` (5 keys). Wasted one round trip; no damage.
2. **The first CI deploy run was a predictable failure.** I pushed the secret *before* reading the deploy step closely enough to notice it never told `action-hosting-deploy` where `firebase.json` lives — the job failed with `firebase.json file not found`, and fixing it cost a second push. Should have read the full deploy step when I added `workflow_dispatch` in the same file.
3. **The auto-commit daemon raced me twice.** My CHANGELOG/AGENTS/workflow edits were committed as `chore: auto-commit 3 changed file(s) (heuristic)` before I could write a real commit message, and my first `git commit` attempt hit "nothing to commit, working tree clean". History now carries two meaningless heuristic messages for meaningful changes. (Known daemon behavior — I should have committed within seconds of editing.)
4. **Leftover from the release session, still unresolved:** the v0.5.2 tag's changelog section is permanently incomplete (TS6 pin and action-v7 fixes shipped inside that tag but are only documented under `[Unreleased]`). Unfixable; mitigated only.
5. **Leftover from the release session:** the "investigate who wrote `go 1.27.1`" task remains open — I deferred it again because it's forensic work outside this session's deploy-repair scope, and now it sits in TODO_LIST instead of being answered.

## e) WHAT WE SHOULD IMPROVE!

1. **CI should catch platform breaks automatically.** The manual cross-compile check I ran proves the fix is trivial; it belongs in `ci.yml` as a matrix, not in an agent's memory.
2. **CI should catch lockfile drift automatically.** Four CI breaks (09-09, 09-13, 09-15, 09-17) all share one root cause: manifests change without lockfiles. A `--frozen-lockfile` + `go mod tidy -diff` guard turns each into a 20-second red X with an obvious fix.
3. **Security posture of the deploy secret.** A never-expiring SA JSON key in a GitHub secret is the weakest link in the pipeline. Workload identity federation (GitHub OIDC → GCP SA) would eliminate the secret entirely; short of that, a key-rotation cadence and deleting the 7 stale `firebase-adminsdk` keys from July would shrink the attack surface.
4. **Zero-consumer blindness.** "Imported by: 0" means every breaking-change risk this library worries about is hypothetical. Either propagate it into real consumers or promote it — the library's risk model changes completely once it has one real dependent.
5. **Commit discipline vs. the daemon.** When the daemon is active, write the real commit message immediately after the edit, or `git commit --amend` the heuristic commit before it pushes. Two heuristic messages this session obscure two meaningful changes.
6. **Read the whole deploy/config step before triggering it.** Both of this session's failures (wrong SA, missing `entryPoint`) were "verified by running" instead of "verified by reading". Reading is cheaper than a failed CI round trip.
7. **Keep the version trio in lockstep mechanically.** `README.md:217` still says "Go 1.26.3" while go.mod says 1.26.7 and AGENTS.md says 1.26.7 — the same drift class as the go.mod/AGENTS.md gotcha, one file over. A consistency gate (or simply fixing it in the next doc pass) retires the class.

## f) Up to 50 things we should get done next

_Ranked by impact. Items 1–15 are harvested into `TODO_LIST.md`; 16+ are ROADMAP fuel / smaller ideas surfaced by this session._

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 1 | Add `GOOS=windows` + `GOOS=darwin` build matrix to CI | Critical | S |
| 2 | npm vulnerability sweep — clear the 10 high Dependabot alerts (astro ecosystem bumps) | Critical | M |
| 3 | Investigate + fix the failing Dependabot `npm_and_yarn` run (#1579905290) | High | S |
| 4 | Add CI guard: `pnpm install --frozen-lockfile` + `go mod tidy -diff` (daemon-drift tripwire) | High | S |
| 5 | Add `govulncheck` job for the Go module | High | S |
| 6 | Forensics: what wrote `go 1.27.1` into go.mod (commit `77b79cc`) + prevention | High | S |
| 7 | Merge the open green Dependabot actions-bump PR | High | S |
| 8 | Consumer propagation: enumerate LarsArtmann repos importing go-atomic-write; bump to v0.5.2 | High | M |
| 9 | Migrate `exhaustruct` → `exhaustruct_v5` in `.golangci.yml` (kills a warning on every lint run) | Medium | S |
| 10 | Delete the 7 stale `firebase-adminsdk` keys from Jul 2026 (below the 10-key cap again) | Medium | S |
| 11 | Version-trio consistency gate (go.mod / .golangci.yml / AGENTS.md / README.md:217) | Medium | S |
| 12 | Fix README benchmark header "Go 1.26.3" → 1.26.7 | Medium | S |
| 13 | Add release-checklist section to AGENTS.md (CI-green-first, cross-compile gate, changelog-last ordering) | Medium | S |
| 14 | Release `[Unreleased]` as 0.5.3 (or fold into next feature release) to fully close the v0.5.2 changelog gap | Medium | S |
| 15 | Evaluate workload-identity federation to eliminate the SA JSON secret; if kept, set a key-rotation cadence | Medium | M |
| 16 | Re-check Dependabot caps (5 PRs/ecosystem) after the vuln sweep lands | Low | S |
| 17 | Remove `minimumReleaseAgeExclude: astro@7.3.3` once its release-age threshold passes | Low | S |
| 18 | Pin `GOTOOLCHAIN` in CI workflows (setup-go manifest lag risk) | Low | S |
| 19 | Run Lighthouse CI against the deployed site (post-deploy verification is currently manual/absent) | Low | S |
| 20 | Bump Node-20-deprecated action SHAs when compatible versions ship | Low | S |
| 21 | Custom domains: apply the Namecheap Terraform for `atomicwrite.lars.software` / `go-atomic-write.lars.software` (blocked on API key) | Low | S |
| 22 | Consider v1.0.0: API stable since 0.5.0, zero consumers — cheap to stabilize now | Medium | M |
| 23 | Promotion: the library has zero importers — profile README, show-case post, or flagship adoption | Medium | L |
| 24 | Decide 0.x `--prerelease` tagging convention vs stable releases at v1.0 | Low | S |
| 25 | Real-Windows hardware test of `rename_windows.go` retry loop (build-tag-compiled, never run) | Low | M |
| 26 | Stale `.tmp` sweeper — crashed writes leave uniquely-named temp files forever | Low | M |
| 27 | Fuzz the fingerprint + commit path | Low | M |
| 28 | NFS / network-filesystem durability caveats doc | Low | S |
| 29 | "Migration from os.WriteFile" guide | Low | M |
| 30 | "Understanding TOCTOU races" deep-dive essay | Low | M |
| 31 | "Why xxhash64?" design rationale page | Low | S |
| 32 | Real-world examples: config updater, state manager, code generator | Low | M |
| 33 | FAQ: permissions, large files, network FS, retry strategy | Low | S |
| 34 | Visual polish: alternating section backgrounds, stat bar, Shiki highlighting | Low | M |
| 35 | Dependents page once adoption exists | Low | S |
| 36 | Benchmark dashboard tracked over time | Low | M |
| 37 | Release automation (`release.yml` tag-based GitHub Releases) before v1.0 | Medium | M |
| 38 | Add `--prerelease`/`--latest` tagging decision into the release checklist | Low | S |
| 39 | Add a post-deploy smoke test (fetch `/changelog`, assert latest version string) to the Website workflow — this exact freeze went unnoticed for 2 months | High | S |
| 40 | Add deploy-failure alerting (the deploy job failed silently in CI for 2 months — nobody watched) | High | S |
| 41 | Document the local-deploy recipe + secret minting recipe in CONTRIBUTING.md (currently only in AGENTS.md) | Low | S |
| 42 | Consolidate Firebase deploy knowledge across sibling repos (gogenfilter/templ-components/go-atomic-write each re-derive it) into the website-launch skill | Medium | M |
| 43 | Audit other sibling-repo websites for the same missing-secret freeze pattern | High | S |
| 44 | Upstream the `entryPoint: website` lesson into the website-launch skill template | Low | S |
| 45 | Verify the custom-domain CNAMEs once DNS is applied (fetch both lars.software hostnames) | Low | S |
| 46 | Add renovate/dependabot grouping so vulnerability bumps land as one reviewable PR | Low | S |
| 47 | Benchmark run freshness: re-run `hash_bench_test.go` on the current toolchain; refresh README numbers | Low | S |
| 48 | `docs/DOMAIN_LANGUAGE.md` freshness pass against the current API surface | Low | S |
| 49 | Check whether the gogenfilter-style `FIREBASE_SERVICE_ACCOUNT` key (2026-07-20) is still needed or rotatable there too | Low | S |
| 50 | Prune `docs/status/` — the 2026-07-26 reports are fully resolved; archive them per docs-health ANNOTATE/ARCHIVE flow | Low | S |

## g) Questions I cannot answer myself

1. **Do you want workload-identity federation (GitHub OIDC → GCP) instead of the SA JSON key?** It would delete the secret entirely and remove the key-rotation burden, but requires one-time GCP project config (provider pool + binding) that I can do via gcloud as `lartyhd@gmail.com` if you approve the IAM changes. Otherwise: should I at least delete the 7 stale July keys on `firebase-adminsdk` and set a rotation cadence for the new key?
2. **Which sibling websites should I audit/repair next for the same missing-secret freeze?** templ-components has its secret (deploys work); gogenfilter has one too; the other ~15 sibling `website/` dirs are unverified. A silent 2-month freeze happened once already — I can check every sibling repo's secret list + last successful Website run in one sweep, but I don't know which sites you consider live vs abandoned.
3. **Is there a consumer repo (or a flagship project you want one to be) for go-atomic-write?** "Imported by: 0" makes every compatibility worry theoretical and v1.0.0 cheap. If you name a target consumer, I can do the adoption + propagation pass; if promotion is the answer instead, tell me which channel (HN/show-case/profile README) you want drafted.

---

**Post-report note:** section (f) items 1–15 are ALREADY harvested into `TODO_LIST.md` and items 22–24 into `ROADMAP.md` (commit `b102ee6`) — this report was written after the harvest, so no follow-up HARVEST run is needed. New items 16–50 (including the two new High items, #39/#40 post-deploy smoke test and failure alerting, and #43 sibling-site audit) are report-only and should be routed on the next docs-health pass if adopted.
