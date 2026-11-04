# Status: Landing redesign, docs-gates CI, fleet repairs, and the pareto-plan closure (2026-11-05 00:40)

Execution session for `docs/planning/2026-11-04_21-05_landing-first-pareto-plan.md` — the owner said go. This report covers the landing lane (T01–T04 + deploy), the trust lane (T05–T13 prep), and the completeness batch that was locally verifiable. **Live remainder: none** — every unshipped item below is either owner-gated, billing-walled, or tracked in TODO_LIST.

## a) FULLY DONE

| Item | Evidence |
| --- | --- |
| **T01–T04 landing redesign shipped and deployed** — dark-first "commit ledger" system: always-dark editor artifact as the hero centerpiece, selection-highlighted `byte` in the headline, glow deleted, warm-paper light theme, green-tinted ink dark theme, stats ledger, ruled FAQ/use-case rows, mismatch callout as the bold dark block, eyebrows dropped, CTA symbol list extended with `WriteFuncWithPerm` | Live at atomicwrite.web.app serving `editor-surface` + `bg-accent/25` markup; deploy from `3a50ac9` (Website run success 22:11Z); commit `3a50ac9` |
| **Redesign gates all green** — contrast AA both themes (0 failing pairs), hero source parity (11 tokens), copy payloads exact, CSP hashes, smoke 4/4→6/6, Lighthouse budgets pass, brand mirror, link check, typecheck 0 errors | `audit:contrast` 0 failures; `test:ui` 6 passed; `lhci autorun` "Done running autorun"; per-gate outputs in session log |
| **Install-pill clipping fixed** — the command it sells was cut mid-word at 1440px; pill widened + type stepped down | Full command visible in post-fix screenshots |
| **T05 changelog entries** — vendor-ordering fix + lockfile-drift recovery documented in `[Unreleased]` | `288ee5f` |
| **T11 archive sweep** — 17-05 + 19-35 reports archived with remainder-none banners; cadence 16 → 14 (gate green again) | `dab37f6`; docs-gates PASS |
| **T06 docs-gates CI** — new `docs-gates` job (pinned flake rev `6717741`, nix-installer v23, `fetch-depth: 0` after the shallow-clone hash-resolution failure); PR #23 verified green then merged | `d777c64` merge; job pass 34s on `5826927` |
| **T08 consumer CI sweep (14 repos)** — short-SHA trap found (first pass was a false all-clear); 2 bump-caused reds repaired (golangci-lint-auto-configure vendorHash → green `28340f6`; dependabot-auto-configure vendorHash + dprint → green `b8cc53e`); 2 fleet lint-pin fixes (v2.13.2 → v2.14.0: dependabot-auto-configure, go-workflow-auditlog `623fd38`); 4 baseline-reds confirmed same-failure-set pre/post (not bump-caused); go-workflow-auditlog's unmasked lint debt filed | sweep ledger in `2026-11-04_20-45…` summary addendum; `86abd26`; gowa#44 |
| **T07 crm#7 repaired and closed** — tidy aligned the graph (templ-components v1.27.0), `regen-ui.sh` regenerated templ + tailwind artifacts (bare `templ generate` misses the derived artifacts and fails the freshness E2E), full local suite green, v0.6.2 bump pushed (`fc7a56d`, `259132b`); closed on local evidence with the billing-wall caveat | crm#7 comment + CLOSED; `12F` CI never started (billing) |
| **T09 lighthouse collection gap root-caused** — `outputDir` is cwd-relative; 25 fresh reports written from `website/`; TODO row re-closed honestly + AGENTS.md lesson | `4ccb1f9` |
| **T10 pkg.go.dev banner** — fresh fetch shows NO banner; it was CDN lag after the release (proxy `@latest` = v0.6.2, page shows v0.6.2 + Latest badge) | fetch + agentic_fetch outputs |
| **T12 a11y asserts** — tab-order (skip link first → header nav) + theme first-paint (class on `<html>` at domcontentloaded both themes; theme-color pair mirrors tokens) added; smoke suite now 6 tests | `b28fdbd`; TODO row removed |
| **T15 deterministic collision test** — `randomSuffix` became an injectable seam; 512 concurrent writers assert distinct shaped paths; fixed-suffix + wrapped-error cases pinned; sequential-by-design to avoid racing parallel readers; lint `0 issues` | `deba219`, `8b5a8b5`; `stage_collision_test.go` |
| **T14 brand pixel-QA** — brand-tokens gains `ogBackground` + `accentLight`; generator derives those inks from tokens (two stale literals found and killed: OG bg `#0a0908` was deliberate-only-by-accident, muted label `#8b827a` was the pre-v0.6.2 value); `brand:check` now asserts every stroke/fill hex in all 8 SVGs is a token value; OG raster regenerated | `4f0d555`; "8 SVG ink sets on-token" |
| **T19 owner-actions checklist** — six manual steps with per-item verification (social preview, OG previewers, Buttondown, apple-touch-icon, SA-key cleanup, the Actions billing wall); linked from brand README + CONTRIBUTING | `c858c10`; `docs/owner-actions-checklist.md` |
| **T20 lessons batch** — annotate.go value grammar, gh short-SHA trap, bump-unmasks-lint-debt pattern, btrfs D-state signature (lighthouse cwd lesson landed with T09) | `bb7088a`; AGENTS.md gotchas |
| **T21 19:55 harvest** — sitemap postbuild gate (`check-sitemap.mjs`: index resolves, ≥10 urls, lastmod dates) + 4 negative tests wired into `postbuild`; CTA symbol-list drift fixed; CRM-export + umask rows harvested to TODO_LIST; 07:09 e.4 + 19:55 f.3 annotated done | `4ccb1f9`, `06a2f11` |
| **T22** — oxlint posture vs AGENTS.md: zero drift (19 disables, categories, 6 warn re-enables all match); gitleaks v3 observed green on a real push (PR #23, 7s) | config dump vs AGENTS.md list |
| **T13 prep** — v0.6.3 release-lane gates pre-run green: uncached lint `0 issues`, race suite, tidy-diff clean, trio 1.27 aligned, build OK. Execution stays gated per plan | session log |

## b) PARTIALLY DONE

| Item | Done | Missing |
| --- | --- | --- |
| T16 Windows leg | Local `GOOS=windows` build + vet green (plus darwin/arm64) | The CI windows-test job cannot be observed: **GitHub Actions billing wall** ("recent account payments have failed or your spending limit needs to be increased", hit ~22:45Z, affects the whole fleet — crm runs `41929139694`/`41929174862` never started). Owner action #6 in the checklist; row stays TODO |
| T24 final sweep | TODO rows closed incrementally; docs-gates 5/5 at handoff | This report's own freshness window (next session harvests nothing — all rows struck here) |

## c) NOT STARTED (gated, by design)

- **T17 visual-regression harness** — deferred: the redesign just changed every visual baseline; a harness built tonight would freeze tomorrow's iteration. Build it after the design settles.
- **T18 WIF migration** — needs GCP owner auth (checklist item 5 is adjacent).
- **T23 dispositions** — owner decisions (BDSM-World, media-atlas mirror), restated in the checklist §7.
- **v0.6.3 execution** — prep green; ships when the release is cut (all gates pre-verified).

## d) TOTALLY FUCKED UP

1. **`gh run list --commit` with short SHAs produced a false all-clear** across all 14 repos on the first sweep pass. Lesson recorded (AGENTS.md); the second pass with full SHAs found 7 real failures.
2. **Two blind find-replace edits on the generator** (stripped quoting inside a template literal, wrong nolint placement after golines rewrap) — both caught by the gates the same minute (XML parse + nolintlint), but the pattern is the 07:09 session's "composition failure" class again. The ink-assertion itself then had an inverted-comparison bug (`hex` vs `#hex`) that the new gate caught on its own first run — the gate working as designed.
3. **The theme first-paint test's first draft was wrong** — an `addInitScript` re-runs on reload and overwrote the dark-theme case back to light; restructured into two contexts. The site was right; the test wasn't.
4. **nc/kill lesson did not apply but `gh issue close --comment --body-file` silently no-op'd** (no comment, no close) — the two-step `gh issue comment` + `gh issue close` form is the reliable one.

## e) WHAT WE SHOULD IMPROVE

1. **Gates keep catching agent composition errors** (d.2/d.3) — keep the negative-first test habit; the ink gate, XML validity, and smoke suite each fired on agent-authored breakage tonight.
2. **The billing wall is a fleet-wide single point of failure** — every "verify CI green" step across 21 repos depends on one billing account. The checklist now names it; WIF (T18) is unrelated but the spend-limit trip suggests reviewing per-repo concurrency limits.
3. **Short-form GH queries need full refs** — a sweep script should resolve SHAs internally (the sweep script now does).

## f) TOP NEXT

1. `[T]` After billing is restored: re-run crm CI on `259132b` and confirm the crm#7 closure's CI caveat (owner, 5 min).
2. `[T]` Windows CI leg job draft + observe (T16 remainder; blocked on billing).
3. `[T]` Visual-regression harness (T17) once the redesign settles.
4. `[R]` WIF migration + dispositions (T18/T23) stay owner-lane.

## g) QUESTIONS

- None new; the three open questions from 20-46 §g are restated in `docs/owner-actions-checklist.md` §7 for one-stop owner review.
