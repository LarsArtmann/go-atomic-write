# Todo List

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

> Rebuilt 2026-11-03 from a full docs-health sweep; re-harvested 2026-11-04
> twice: after the release/quality session (v0.6.1, fleet bump, CSP unblock)
> and after the Lane 4–5 execution (T20–T33: gitleaks, build-script regression
> tests, CopyButton, brand pipeline + 2× rasters, brand tokens, SEO batch,
> link checker, CONTRIBUTING/deploy docs, npm-sweep tracker, when-not-to-use +
> FAQ pages, platforms matrix, Playwright smoke suite + PR preview channels).
> Evidence in `CHANGELOG.md` [Unreleased] and
> `docs/status/2026-11-04_07-09_lane-4-5-gates-brand-seo-docs-landing-execution-status.md`.

## Medium Impact

| Task                                                                                                                                                        | Status  | Impact | Effort | Evidence                                                                                         |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ | ------ | ------------------------------------------------------------------------------------------------ |
| Keyboard tab-order + theme first-paint verification — the Playwright smoke suite covers copy payloads and `astro:page-load` animation binding since `bb3c47b`; tab-order asserts and a theme-persistence check remain | 🔴 TODO | Medium | M      | `website/tests/smoke.spec.mjs`; supersedes the old interactive-nav row (plan T32)                |
| Brand pixel-QA: scripted glyph-ink assertion + footer-logo screenshot QA (the rest of the brand batch shipped — `brand:png` rasters, `brand-tokens.json`, SVG well-formedness gate, measured wordmark air) | 🔴 TODO | Medium | M      | `scripts/generate-brand-assets.mjs`; `docs/status/2026-11-03_23-29…` §5.2–5.7                    |

## Low Impact

| Task                                                                                                                                                        | Status     | Impact | Effort | Evidence                                                                      |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- | ------ | ------ | ----------------------------------------------------------------------------- |
| `minimatch@3.1.5` moderate advisory (glob@7 ← rimraf@3 ← chrome-launcher ← @lhci/cli): overriding to 11 breaks glob@7's CJS API at runtime — wait for rimraf/chrome-launcher upstream majors, then re-audit | 🔵 BLOCKED | Low    | S      | `website/pnpm-workspace.yaml` `overrides`; audit 2026-11-03 (1 advisory left) |
| Owner actions bundle: OG card in real previewers (Slack/X/Facebook), real Buttondown subscribe test, iOS apple-touch-icon check, GitHub social-preview upload (`assets/brand/social-preview.png` → repo Settings → Social preview), delete the 7 stale `firebase-adminsdk` keys | 🔴 TODO | Low | S | `assets/brand/README.md` (upload recipe); AGENTS.md SA-cap gotcha             |
| Deploy-secret follow-through: migrate CI deploys to Workload Identity Federation (Option A of the decided posture doc); the decision doc itself shipped (`6142c78`) | 🔴 TODO | Low | M | `docs/deploy-secret-posture.md`                                               |
| Real-Windows test run of the fingerprint/temp-path suite — cross-compile builds every push, but the tests have never executed on Windows hardware (`rename_windows.go` retry loop + `FlushFileBuffers` durability need a windows CI test leg or one manual run) | 🔴 TODO | Low | M | routed from `docs/status/2026-11-03_20-42…` row 9 (docs-gates audit 2026-11-04) |
| Visual-regression screenshot harness for header/docs/OG surfaces (the smoke suite covers the landing page only)                                              | 🔴 TODO | Low | M | routed from `docs/status/2026-11-04_02-06…` row 45 (docs-gates audit 2026-11-04) |
| Deterministic collision test for `tempPathFor` (seed/inject `randomSuffix`); the current 200-writer test is statistical                                      | 🔴 TODO | Low | S | `atomicwrite_test.go`; routed from `docs/status/2026-11-04_05-36…` f.28       |
| File the CV Go-test chronic red (10/10 jobs failing on every master commit) after confirming the concurrent storm session does not own it                    | 🔴 TODO | Low | S | `docs/status/2026-11-04_05-36…` f.24                                          |
| `pnpm run lighthouse:deployed` against the live site after the next deploy settles (config exists since v0.6.0, never observed)                              | 🔴 TODO | Low | S | `lighthouserc.deployed.json`; `docs/status/2026-11-04_05-36…` f.30            |

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
