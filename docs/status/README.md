# Status & Planning Reports

Point-in-time snapshots. **Never rewrite a report** — resolve items inline
(strikethrough + verdict, per the docs-health ANNOTATE canon) and archive when
every item carries a disposition. Living documents live at the repo root
(`README.md`, `FEATURES.md`, `TODO_LIST.md`, `ROADMAP.md`, `CHANGELOG.md`,
`AGENTS.md`); this directory is history.

## Current reports (docs/status/)

| Report                                                                                 | Topic                                                                                                                     |
| -------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `2026-11-03_20-42_no-linter-submodule-and-rename-failure-temp-leak-fix.md`             | Linter sub-module rejected; rename-failure temp-file leak fixed                                                           |
| `2026-11-03_22-20_todo-sweep-ci-hardening-vuln-sweep-ts-self-heal-benchmark-ranges.md` | TODO sweep, CI hardening, npm vuln sweep, TS6 self-heal, benchmarks                                                       |
| `2026-11-03_22-21_buildflow-gate-red-to-green-audit-typecheck-jscpd-fixes.md`          | BuildFlow gate audit: typecheck/jscpd fixes                                                                               |
| `2026-11-03_23-11_brand-identity-atomic-swap-logo.md`                                  | Brand identity build log ("the atomic swap")                                                                              |
| `2026-11-03_23-29_brand-identity-atomic-swap-comprehensive-status.md`                  | Brand verdict + gaps (prime HARVEST source for the current TODO_LIST)                                                     |
| `2026-11-03_23-44_oxlint-cleanup-and-config-tuning.md`                                 | oxlint cleanup + config tuning                                                                                            |
| `2026-11-04_00-14_buildflow-file-size-split-refactor-status.md`                        | BuildFlow file-size budget: the Go source split that shipped in v0.6.1                                                    |
| `2026-11-04_02-06_docs-health-audit-annotation-archive-and-living-doc-rebuild.md`      | docs-health audit: archive sweep + living-doc rebuild (source of the 02-12 Pareto plan)                                   |
| `2026-11-04_05-36_release-fleet-and-website-quality-execution-status.md`               | v0.6.1 release + 17-repo fleet bump + website quality execution (T01–T19); §f is the T20–T34 seed                         |
| `2026-11-04_06-15_fleet-bump-v0.6.1-summary.md`                                        | Fleet bump ledger: 17 repos × hashes, 4 skips + reasons, 4 filed consumer issues                                          |
| `2026-11-04_07-09_lane-4-5-gates-brand-seo-docs-landing-execution-status.md`           | Lane 4–5 execution: gates, brand pipeline, SEO, docs, landing (T20–T31); WAIT handoff report                              |
| `2026-11-04_19-55_writefuncwithperm-issue20-session-status.md`                         | Concurrent session: WriteFuncWithPerm feature (issue #20) + website deps hardening; §f brainstorm                         |
| `2026-11-04_20-45_fleet-bump-v062-summary.md`                                          | Fleet bump ledger v0.6.2: 17 repos × hashes, vendor-ordering bug + fix, skips + reasons (crm#7 filed)                     |
| `2026-06-28_03-32_atomicity-durability-overhaul.html`                                  | LEAVE-ALONE (v0.2.0 release dashboard; referenced by `CHANGELOG.md` [0.2.0] — moving it would break that historical link) |

## Archived (resolved — every item carries an inline verdict)

`archived/` is closed history; do not reopen it for harvesting.

### docs/status/archived/ — 2026-11-03 sweep manifest

| File                                                                  | Classification      | Deciding reason                                                                     |
| --------------------------------------------------------------------- | ------------------- | ----------------------------------------------------------------------------------- |
| `2026-07-09_14-52_marketing-overhaul-and-website-creation.md`         | ANNOTATE → ARCHIVED | All 50 backlog items + findings resolved (shipped, superseded, routed, or declined) |
| `2026-07-11_22-12_landing-page-design-overhaul.md`                    | ANNOTATE → ARCHIVED | Design-overhaul backlog fully resolved by 09-18/11-03 sessions                      |
| `2026-07-26_17-53_docs-health-run-and-brutal-self-review.md`          | ANNOTATE → ARCHIVED | First docs-health run; every §e/§f/§g item closed by the same-day 18-27 session     |
| `2026-07-26_18-27_split-brain-fix-and-website-modernization.md`       | ANNOTATE → ARCHIVED | 50-item list resolved: CI/CSP leftovers routed to TODO_LIST, rest done/moot         |
| `2026-07-26_18-49_todo-list-clearance-and-brutal-self-review.md`      | ANNOTATE → ARCHIVED | 8 cleared TODOs verified; 50-item list resolved (Lighthouse/brand items routed)     |
| `2026-07-26_19-59_v0.4.0-release-session-review.md`                   | ANNOTATE → ARCHIVED | Release completed in-file; all follow-ups resolved or routed                        |
| `2026-09-17_14-22_v0.5.2-release-session.md`                          | ANNOTATE → ARCHIVED | 15-item harvest cleared 11-03; release shipped                                      |
| `2026-09-17_20-57_website-unfreeze-and-ci-deploy-repair.md`           | ANNOTATE → ARCHIVED | 50-item list: CI hardening done 11-03; remainder routed                             |
| `2026-09-18_06-11_landing-page-messaging-and-hero-split-brain-fix.md` | ANNOTATE → ARCHIVED | Messaging shipped in v0.6.0; residuals routed (voice pass, FAQ links)               |
| `2026-09-18_10-37_backlog-clearance-and-build-verification-gate.md`   | ANNOTATE → ARCHIVED | 50-item §f fully dispositioned in the 2026-11-03 docs-health sweep                  |
| `2026-11-04_17-05_consumer-fixes-t34-gates-preview-channels-status.md` | ANNOTATE → ARCHIVED | T32–T34 closed, previews live, 4 consumer fixes landed (same-day)                   |
| `2026-11-04_19-35_v062-release-fleet-bump-inflight-status.md`          | ANNOTATE → ARCHIVED | Superseded by 20-45 ledger + 20-46 closeout; release + bump done                    |

### docs/planning/archived/

| File                                   | Classification      | Deciding reason                                                                                                                                               |
| -------------------------------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `2026-06-03_16-19-full-code-review.md` | ANNOTATE → ARCHIVED | Review of the pre-v0.2.0 codebase; improvement plan fully resolved (two edge tests routed to TODO_LIST), stale `.bak`/3-arg-`Write` premises inline-corrected |

## Annotation conventions

- Inline verdicts: `~~item~~ done at <hash>` / `superseded: …` / `→ TODO_LIST: …` / `declined: …`. Open items stay bare.
- Batch tooling: the docs-health skill's `annotate.go` (`repair`/`check-rows`); `check-rows` must exit clean before archiving.
