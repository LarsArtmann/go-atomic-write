# TODO List

> Short-term, actionable, bounded work items, verified against the actual code.
> For long-term vision and unrefined ideas, use ROADMAP.md.
> Items are ranked by impact. Status is verified, not assumed.

## Status legend

| Status           | Meaning                                                     |
| ---------------- | ----------------------------------------------------------- |
| 🔴 `TODO`        | Not started. Needs doing.                                   |
| 🟡 `IN_PROGRESS` | Actively being worked on.                                   |
| 🔵 `BLOCKED`     | Cannot proceed, external dependency or decision needed.     |
| 🟢 `DONE`        | Completed. Remove from this list and log in `CHANGELOG.md`. |

> All items from the 2026-07-26 docs-health run were cleared then (see
> `CHANGELOG.md` `[0.4.0]`). The list below was harvested on 2026-09-17 from
> `docs/status/2026-09-17_14-22_v0.5.2-release-session.md` section (f);
> items already resolved in the 2026-09-17 session (Firebase secret, live-site
> verification, AGENTS.md gotcha refresh, CHANGELOG `[Unreleased]` backfill,
> LSP restart) are NOT listed — see `CHANGELOG.md` `[Unreleased]`.

## High Impact

| Task                                                                                                                  | Status  | Impact | Effort | Evidence                                                                         |
| --------------------------------------------------------------------------------------------------------------------- | ------- | ------ | ------ | -------------------------------------------------------------------------------- |
| Add `GOOS=windows` + `GOOS=darwin` build matrix to CI (v0.5.1 was exactly such a break; Linux-only CI can't catch it) | 🔴 TODO | High   | S      | `.github/workflows/ci.yml`; report §f #3                                         |
| npm vulnerability sweep: clear the 10 high Dependabot alerts (astro ecosystem bumps)                                  | 🔴 TODO | High   | M      | https://github.com/LarsArtmann/go-atomic-write/security/dependabot; report §f #6 |

## Medium Impact

| Task                                                                                                                                                | Status  | Impact | Effort | Evidence                                         |
| --------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ | ------ | ------------------------------------------------ |
| Investigate why the Dependabot `npm_and_yarn` run fails                                                                                             | 🔴 TODO | Medium | S      | Dependabot run #1579905290; report §f #7         |
| Add CI guard job: `pnpm install --frozen-lockfile` + `go mod tidy -diff` (auto-commit daemon drift tripwire — 4 CI breaks: 09-09/09-13/09-15/09-17) | 🔴 TODO | Medium | S      | `AGENTS.md` daemon gotcha; report §f #9          |
| Add `govulncheck` job for the Go module in CI                                                                                                       | 🔴 TODO | Medium | S      | `.github/workflows/ci.yml`; report §f #10        |
| Find what wrote `go 1.27.1` into `go.mod` at 2026-09-17 ~09:29 (auto-commit 77b79cc) and prevent recurrence                                         | 🔴 TODO | Medium | S      | commit `77b79cc`; report §f #11                  |
| Consumer propagation: enumerate LarsArtmann repos importing go-atomic-write; bump them to v0.5.2                                                    | 🔴 TODO | Medium | M      | pkg.go.dev shows "Imported by: 0"; report §f #12 |

## Low Impact

| Task                                                                                                            | Status  | Impact | Effort | Evidence                                                           |
| --------------------------------------------------------------------------------------------------------------- | ------- | ------ | ------ | ------------------------------------------------------------------ |
| Migrate `exhaustruct` → `exhaustruct_v5` in `.golangci.yml` (deprecated since golangci-lint v2.13.0)            | 🔴 TODO | Low    | S      | `.golangci.yml:42,151`; lint warning; report §f #13                |
| Consistency gate for the version trio (`go.mod` / `.golangci.yml` run.go / `AGENTS.md`)                         | 🔴 TODO | Low    | S      | `AGENTS.md:122` drift gotcha; report §f #14                        |
| Update README benchmark header "Go 1.26.3" → current toolchain (1.26.7)                                         | 🔴 TODO | Low    | S      | `README.md:217`; report §f #15                                     |
| Remove `minimumReleaseAgeExclude: astro@7.3.3` from `pnpm-workspace.yaml` once its release-age threshold passes | 🔴 TODO | Low    | S      | `website/pnpm-workspace.yaml`; report §f #16                       |
| Pin `GOTOOLCHAIN` in CI workflows (setup-go manifest lag risk)                                                  | 🔴 TODO | Low    | S      | `.github/workflows/ci.yml`; report §f #18                          |
| Re-check Dependabot config caps (5 PRs/ecosystem) after the npm vuln sweep                                      | 🔴 TODO | Low    | S      | `.github/dependabot.yml`; report §f #20                            |
| Add release-checklist section to AGENTS.md (CI-green-first ordering, cross-compile gate)                        | 🔴 TODO | Low    | S      | report §d self-critique; report §f #21                             |
| Run Lighthouse CI against the deployed site post-deploy                                                         | 🔴 TODO | Low    | S      | `website/lighthouserc.json`; report §f #22                         |
| Track + bump Node-20-deprecated action SHAs when compatible versions ship                                       | 🔴 TODO | Low    | S      | Website workflow Node.js 20 deprecation annotations; report §f #24 |

---

<!-- Guidance for the builder filling this in:
  - Source of truth is the CODE. Verify each item before adding, many
    documented TODOs are already done.
  - One task per row. If it takes more than ~2 hours, split it into smaller
    tasks.
  - Cite evidence (file:line) so the next person can verify without re-deriving.
  - DONE items should be REMOVED, not kept. Use CHANGELOG.md for history.
  - If a task is vague ("improve X"), refine it into concrete steps or move
    it to ROADMAP.md.
  - Deduplicate by semantic intent, not by text match.
  - For 80/20 impact prioritization, use the pareto-planning skill AFTER
    building the list here.
-->
