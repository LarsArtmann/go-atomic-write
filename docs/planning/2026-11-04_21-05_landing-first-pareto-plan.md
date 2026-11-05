# Pareto Execution Plan — Landing-First Quality Closure

**Created:** 2026-11-04 21:05 CET
**Author:** Crush session (post-closeout)
**Trigger:** owner verdict on the live landing page ("looks a bit lame") + the standing backlog after the v0.6.2 fleet-bump closeout (`docs/status/2026-11-04_20-46…`).
**Supersedes:** `2026-11-04_02-12_release-fleet-and-site-quality-pareto-plan.md` — that plan's tier 1 (v0.6.1 release) and tier 2 (22-repo fleet bump, OG voice, newsletter CSP) SHIPPED today (v0.6.1 + v0.6.2 tagged, 17 consumers bumped, site verified). Nothing from it is re-dispatched here; its unfinished residue (WIF, owner actions, collision test, Windows leg) is absorbed into tier 4 below with today's evidence.

---

## 0. Ground truth at planning time (verified, not assumed)

- Site: `atomicwrite.lars.software` + `.web.app` both live, deployed from `6180dc7` (Website run success 20:16 CET); **zero site-relevant content changed since** (`git diff 6180dc7..HEAD -- CHANGELOG.md website/` empty).
- Release surfaces agree as of 20:55: v0.6.2 = Latest on GitHub (prerelease flag cleared, `d8ab296`), `[0.6.2]` in CHANGELOG.md, `[0.6.2]` live on the changelog page.
- v0.6.2 GitHub releases: no `--prerelease` flags remain anywhere in history.
- Gates at HEAD `d8ab296`: CI success, Website success, docs-gates 5/5, `go test -race` ok, uncached lint 0 issues.
- Design diagnosis from tonight's screenshots (the "lame" evidence): light theme washes the brand out — near-white on near-white, one giant decorative radial glow eating the hero's lower half, a pale low-contrast code window (the product's core artifact!), emerald reduced to micro-accents, stats row as bare text, and the overall read is the SaaS-card-kit default (identical rounded cards, soft borders, gradient wash). The dark theme exists but is not the default presentation.
- Concurrent session: `WriteFuncWithPerm` sits in `[Unreleased]` (their lane, in flight). The v0.6.3 release lane in this plan is PREP-only until that lands.
- Blocked/gated: crm baseline (crm#7 — repair shape verified, ownership question open); BDSM-World + media-atlas mirror (owner-decision questions open, see 20-46 §g).

**Verschlimmbesser guard:** the landing page is gated six ways (`verify-build`, `check:hero`, `audit:contrast`, `brand:check`, `html-validate`, Lighthouse budgets, Playwright smoke). Every design change below ships THROUGH those gates, not around them. No gate is weakened to make a design easier. The `Logo.astro` and all `assets/brand/*.svg` are GENERATED — brand changes go through `scripts/generate-brand-assets.mjs` only. `hero-code.ts` stays the single source of the hero snippet.

---

## 1. Pareto analysis

| Tier                 | Deliverable                                                                                                                                                                                                                                    | Why it carries the value                                                                                                                                                                                                                                         |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **1% → 51%**         | **The hero, rebuilt** (T01 tokens + T02 hero): dark-first presentation, the radial glow deleted, the code window promoted to first-class artifact, one bold emerald moment                                                                     | The hero IS the verdict. A visitor decides "serious tool vs. template demo" in the first viewport. Tonight's verdict was formed there; the fix lands there. Nothing else on the page can compensate for a washed-out first screen.                               |
| **4% → 64%**         | **+ The full landing system** (T03 sections, T04 verification battery): de-duplicate the card kit, the swap-motif pipeline visual, real contrast discipline                                                                                    | The page below the fold must sustain the hero's quality or the first screen reads as luck. The verification battery is what makes the redesign _shippable_ — contrast AA both themes, budgets, smoke — so the fix cannot regress the gates.                      |
| **20% → 80%**        | **+ Release-trust and doc-trust closure** (T05–T13): the two missing changelog entries, docs-gates in CI, crm#7 repair, the 14-repo consumer CI sweep, lighthouse root-cause, pkg.go.dev banner, archive sweep, tab-order asserts, v0.6.3 prep | These are the surfaces a skeptical senior engineer checks _after_ the landing page convinces them: the changelog tells the truth, the docs gates can't silently rot, every fleet consumer actually built green, the release lane is ready for WriteFuncWithPerm. |
| **final 20% → 100%** | **+ Completeness batches** (T14–T24): brand pixel-QA, collision test, Windows leg, visual-regression, WIF, owner actions, lessons, dispositions                                                                                                | Real but deferrable: no known user harm. Bounded so they cost hours, not sprawl. Two items are owner-decision-gated (T23).                                                                                                                                       |

**Anti-Pareto warning:** do NOT let tier-4 breadth delay the landing fix. The graph is two lanes: the landing lane (T01→T05) ships alone if it has to; the trust lane (T06–T13) is independent and parallel.

---

## 2. Comprehensive plan — 30–100 min tasks, ALL open TODOs, tier → value → effort

24 tasks ≈ 19 h. Covers: tonight's redesign brief, all 9 TODO_LIST rows, the 20-46 status report §f survivors, and the open items from the 19:35/19:55 sessions. minimatch advisory stays BLOCKED (upstream) — not planned.

| #   | Tier | Task                                                                                                                        | Value | Eff  | Covers                                   |
| --- | ---- | --------------------------------------------------------------------------------------------------------------------------- | ----- | ---- | ---------------------------------------- |
| T01 | 1    | Landing design system: token encoding (dark-first emerald scale, type roles, backdrop discipline) + brand-mirror check      | High  | 90m  | "lame" brief; TODO pixel-QA context      |
| T02 | 1    | Hero rebuild: glow removal, headline treatment, code window as artifact, stats ledger, CTA row                              | High  | 90m  | "lame" brief; hero-code.ts single-source |
| T03 | 2    | Landing sections pass: features, how-it-works (swap motif), mismatch callout, comparison, use cases, FAQ, footer            | High  | 100m | "lame" brief; card-kit de-duplication    |
| T04 | 2    | Redesign verification battery: contrast AA both themes, smoke, budgets, screenshots, brand:check                            | High  | 60m  | TODO contrast/pixel rows; gates          |
| T05 | 2    | Add the 2 missing `[Unreleased]` changelog entries (fleet vendor-ordering fix + lockfile-drift recovery)                    | Med   | 30m  | 20-46 §f.10 — FORGOT tonight             |
| T06 | 2    | Wire docs-gates into CI (freshness/annotate/routing/cadence job)                                                            | High  | 60m  | TODO docs-gates row; 19:35 f.33          |
| T07 | 3    | Repair crm baseline per crm#7 (tidy commit → templ generate → build → bump → close)                                         | Med   | 90m  | crm#7; 20-46 §f.2                        |
| T08 | 3    | Consumer CI sweep: bulk-check the 14 unobserved bump commits, triage failures                                               | Med   | 45m  | 20-46 b3/f5 residue                      |
| T09 | 3    | Root-cause lighthouse:deployed collection gap (per-run report dir), re-close TODO row honestly                              | Med   | 45m  | 20-46 b1                                 |
| T10 | 3    | Resolve pkg.go.dev latest-banner anomaly (cold-cache @latest check)                                                         | Low   | 30m  | 20-46 b2                                 |
| T11 | 3    | Archive sweep: fully-annotated reports → archived/ (cadence is at the 15/15 threshold)                                      | Med   | 60m  | docs-cadence gate; 20-46 §f.13           |
| T12 | 3    | Keyboard tab-order + theme first-paint Playwright asserts                                                                   | Med   | 60m  | TODO_LIST Medium row                     |
| T13 | 3    | v0.6.3 release-lane prep (gates pre-run, changelog section draft, trio check; execution gated on WriteFuncWithPerm landing) | Med   | 45m  | 19:55 session; release checklist         |
| T14 | 4    | Brand pixel-QA: scripted glyph-ink assertion + footer-logo screenshot QA                                                    | Low   | 60m  | TODO_LIST Medium row                     |
| T15 | 4    | Deterministic collision test for `tempPathFor` (seed/inject `randomSuffix`)                                                 | Low   | 45m  | TODO_LIST Low row                        |
| T16 | 4    | Real-Windows test leg (CI matrix job; triage rename/AV findings)                                                            | Low   | 90m  | TODO_LIST Low row                        |
| T17 | 4    | Visual-regression screenshot harness (header/docs/OG surfaces)                                                              | Low   | 90m  | TODO_LIST Low row                        |
| T18 | 4    | WIF migration for CI deploys (posture doc Option A; needs GCP owner auth)                                                   | Low   | 90m  | TODO_LIST Low row                        |
| T19 | 4    | Owner-actions prep: step-by-step checklist doc (OG previewers, social upload, key deletion)                                 | Low   | 30m  | TODO_LIST Low row                        |
| T20 | 4    | AGENTS.md lessons batch (annotate `,` grammar; lighthouse invocation; btrfs D-state hangs)                                  | Low   | 30m  | 20-46 §f.25/26/40                        |
| T21 | 4    | 19:55 report harvest remainder + sitemap assertion into verify-build (`07:09 f.17`)                                         | Low   | 45m  | 19:55 §f; 07:09 f.17                     |
| T22 | 4    | gitleaks v3 runtime observation + oxlint config drift check                                                                 | Low   | 30m  | 20-46 §f.12/45                           |
| T23 | 4    | BDSM-World + media-atlas mirror dispositions (**owner-decision gated** — 20-46 §g.2/3)                                      | Low   | 30m  | 20-46 §f.7/9                             |
| T24 | 4    | Final sweep: TODO_LIST harvest of survivors, docs-gates run, status-index rows                                              | Low   | 30m  | docs-health closure                      |

---

## 3. Detailed breakdown — ≤12 min micro-tasks, ALL TODOs

| ID    | Task                                                                                                          | Min | Tier | Depends on |
| ----- | ------------------------------------------------------------------------------------------------------------- | --- | ---- | ---------- |
| T01.1 | Screenshot inventory: full page light+dark, 375px mobile, list concrete failure modes                         | 12  | 1    | —          |
| T01.2 | Read the current system: `global.css` tokens, `LandingLayout`, section components, `astro.config` fonts       | 10  | 1    | —          |
| T01.3 | Write the design plan: palette (4–6 hex), type roles, layout concept, principles (frontend-design two-pass)   | 12  | 1    | T01.1-2    |
| T01.4 | Self-review the plan against the AI-tell calibration list; revise what reads generic                          | 12  | 1    | T01.3      |
| T01.5 | Encode tokens in `global.css` (`@theme`, light+dark; keep the 4 brand-mirrored tokens coherent)               | 12  | 1    | T01.4      |
| T01.6 | Brand-mirror contract: update `brand-tokens.json` + `brand:generate` + `brand:check` if mirrored tokens moved | 12  | 1    | T01.5      |
| T01.7 | Fonts audit: display + mono roles confirmed or adjusted in `astro.config`                                     | 10  | 1    | T01.4      |
| T01.8 | Token verification: build + typecheck green                                                                   | 10  | 1    | T01.5-7    |
| T02.1 | Delete the radial glow; define the hero backdrop (brand-geometry-derived texture or none)                     | 12  | 1    | T01.8      |
| T02.2 | Headline/subheadline type scale + spacing; replace the single-word-accent tell with structural emphasis       | 12  | 1    | T02.1      |
| T02.3 | Code window rebuild: dark editor surface, emerald syntax roles, disciplined shadow                            | 12  | 1    | T02.1      |
| T02.4 | Stats row as a bordered ledger (tabular numbers, dividers, not bare text)                                     | 12  | 1    | T02.2      |
| T02.5 | CTA row: solid-emerald primary, ghost secondary, visible focus states                                         | 12  | 1    | T02.2      |
| T02.6 | `hero-code.ts` token colors adapt (single source preserved; `data-copy` payloads untouched)                   | 12  | 1    | T02.3      |
| T02.7 | Gates: `check:hero` + `verify-build` parity green                                                             | 10  | 1    | T02.6      |
| T03.1 | Features grid: break the identical-card kit (claim hierarchy, border discipline)                              | 12  | 2    | T02.7      |
| T03.2 | HowItWorks: the swap motif — brand arrows as phase connectors                                                 | 12  | 2    | T03.1      |
| T03.3 | Mismatch callout (`fingerprint != disk`) as the page's one bold moment                                        | 12  | 2    | T03.2      |
| T03.4 | Comparison table: emerald check column, zebra discipline                                                      | 12  | 2    | T03.1      |
| T03.5 | Use cases + platforms: compact two-column treatment                                                           | 12  | 2    | T03.1      |
| T03.6 | FAQ: spacing + `details` styling                                                                              | 12  | 2    | T03.1      |
| T03.7 | Footer/newsletter: emerald treatment, live region + CSP form-action intact                                    | 12  | 2    | T03.1      |
| T03.8 | Scroll-reveal audit: ONE orchestrated moment; reduced-motion respected; `data-animate-bound` contract intact  | 12  | 2    | T03.1-7    |
| T04.1 | `audit:contrast` both themes; fix any new AA violations                                                       | 12  | 2    | T03.8      |
| T04.2 | Playwright smoke 4/4 + 320px overflow check                                                                   | 12  | 2    | T03.8      |
| T04.3 | Lighthouse budgets local run                                                                                  | 12  | 2    | T04.1      |
| T04.4 | Screenshot review (fold/mobile/dark) against the "lame" criteria; iterate once                                | 12  | 2    | T04.1-3    |
| T04.5 | `html-validate` + `brand:check` + `test:scripts`                                                              | 12  | 2    | T04.4      |
| T05.1 | Add vendor-ordering fix + lockfile-drift recovery to `[Unreleased]` Fixed                                     | 10  | 2    | —          |
| T05.2 | `sync:changelog` + MDX escape check                                                                           | 10  | 2    | T05.1      |
| T05.3 | Commit (daemon carrier-amend pattern)                                                                         | 10  | 2    | T05.2      |
| T06.1 | Design the job: which gates, `nix run` invocation, caching posture                                            | 12  | 2    | —          |
| T06.2 | Add the job to `ci.yml` (pinned input; fail on findings)                                                      | 12  | 2    | T06.1      |
| T06.3 | Verify the job green on a branch push                                                                         | 12  | 2    | T06.2      |
| T06.4 | Document in CONTRIBUTING quality-gates section                                                                | 12  | 2    | T06.3      |
| T06.5 | Merge + observe master CI                                                                                     | 12  | 2    | T06.4      |
| T07.1 | crm: confirm crm#7 unclaimed; branch state clean                                                              | 12  | 3    | —          |
| T07.2 | Commit the tidy (go.mod ↔ graph alignment, templ-components 1.27.0)                                           | 12  | 3    | T07.1      |
| T07.3 | `templ generate` with the resolved runtime; review the generated diff                                         | 12  | 3    | T07.2      |
| T07.4 | build + test green                                                                                            | 12  | 3    | T07.3      |
| T07.5 | Push repair; observe crm CI                                                                                   | 12  | 3    | T07.4      |
| T07.6 | `go get go-atomic-write@v0.6.2` + tidy; bump commit + push                                                    | 12  | 3    | T07.5      |
| T07.7 | Close crm#7 with hashes + evidence comment (github-voice)                                                     | 10  | 3    | T07.6      |
| T08.1 | Bulk `gh run list` over the 14 unobserved consumers; collect conclusions                                      | 12  | 3    | —          |
| T08.2 | Triage any failures (root-cause or file)                                                                      | 12  | 3    | T08.1      |
| T08.3 | Record results in the 20-45 fleet summary (annotate)                                                          | 10  | 3    | T08.2      |
| T08.4 | Update TODO/report cross-references                                                                           | 9   | 3    | T08.3      |
| T09.1 | Read `lighthouserc.deployed.json` + lhci autorun report-dir semantics                                         | 12  | 3    | —          |
| T09.2 | Reproduce with a per-run report dir; confirm fresh collection lands                                           | 12  | 3    | T09.1      |
| T09.3 | Fix the config (unique dir per run or assert/collection split)                                                | 12  | 3    | T09.2      |
| T09.4 | Re-run; re-close the TODO row honestly + AGENTS.md lesson                                                     | 9   | 3    | T09.3      |
| T10.1 | Cold module cache: `go get @latest` + proxy `.info` check                                                     | 12  | 3    | —          |
| T10.2 | Identify the banner trigger (graph state vs page render)                                                      | 10  | 3    | T10.1      |
| T10.3 | Record finding (doc note or pkg.go.dev upstream report)                                                       | 8   | 3    | T10.2      |
| T11.1 | Candidate list: fully-annotated 11-03/11-04 reports                                                           | 12  | 3    | —          |
| T11.2 | `check-rows --strict` + `bare-rows` per candidate                                                             | 12  | 3    | T11.1      |
| T11.3 | `git mv` to `archived/` + manifest rows in the status README                                                  | 12  | 3    | T11.2      |
| T11.4 | docs-gates re-run (cadence back under threshold)                                                              | 12  | 3    | T11.3      |
| T11.5 | Commit                                                                                                        | 12  | 3    | T11.4      |
| T12.1 | Tab-order test: logical focus path hero → nav → sections                                                      | 12  | 3    | T04.5      |
| T12.2 | Theme first-paint test: no wrong-theme flash; `theme-color` pair correct                                      | 12  | 3    | T12.1      |
| T12.3 | Run suite; fix findings                                                                                       | 12  | 3    | T12.2      |
| T12.4 | Wire into smoke spec; CI green                                                                                | 12  | 3    | T12.3      |
| T12.5 | TODO row closed (removed + changelog)                                                                         | 12  | 3    | T12.4      |
| T13.1 | Pre-run gates: lint uncached, race, tidy-diff                                                                 | 12  | 3    | —          |
| T13.2 | Changelog `[Unreleased]` → `[0.6.3]` section draft (coordinate with the concurrent session's entry)           | 12  | 3    | T13.1      |
| T13.3 | Version trio check (`go.mod` / `.golangci.yml` / AGENTS.md)                                                   | 10  | 3    | T13.2      |
| T13.4 | Execution-gate note: ships only when WriteFuncWithPerm lands + the redesign is deployed                       | 11  | 3    | T13.3      |
| T14.1 | Glyph-ink assertion design (SVG ink coverage vs brand tokens)                                                 | 12  | 4    | —          |
| T14.2 | Implement in `generate-brand-assets.mjs --check`                                                              | 12  | 4    | T14.1      |
| T14.3 | Footer-logo screenshot QA harness                                                                             | 12  | 4    | T14.2      |
| T14.4 | Wire into `brand:check` or a sibling script                                                                   | 12  | 4    | T14.3      |
| T14.5 | Run + record results                                                                                          | 12  | 4    | T14.4      |
| T15.1 | Design the seam: seed/inject `randomSuffix`                                                                   | 12  | 4    | —          |
| T15.2 | Implement the deterministic collision test                                                                    | 12  | 4    | T15.1      |
| T15.3 | Run + `-race`                                                                                                 | 12  | 4    | T15.2      |
| T15.4 | TODO row closed                                                                                               | 9   | 4    | T15.3      |
| T16.1 | Decide: CI windows matrix vs manual-run doc                                                                   | 12  | 4    | —          |
| T16.2 | CI job draft (`windows-latest`, `go test ./...`)                                                              | 12  | 4    | T16.1      |
| T16.3 | Local `GOOS=windows` sanity build/test                                                                        | 12  | 4    | T16.2      |
| T16.4 | Push + observe the windows job                                                                                | 12  | 4    | T16.3      |
| T16.5 | Triage failures (rename retry, AV, path cases)                                                                | 12  | 4    | T16.4      |
| T16.6 | Record + row update                                                                                           | 10  | 4    | T16.5      |
| T17.1 | Choose approach: Playwright screenshots + threshold compare                                                   | 12  | 4    | —          |
| T17.2 | Baseline capture: header, docs, OG surfaces                                                                   | 12  | 4    | T17.1      |
| T17.3 | Compare workflow + CI-wiring decision                                                                         | 12  | 4    | T17.2      |
| T17.4 | Implement the harness                                                                                         | 12  | 4    | T17.3      |
| T17.5 | Run; commit baselines                                                                                         | 12  | 4    | T17.4      |
| T17.6 | Docs + TODO row update                                                                                        | 10  | 4    | T17.5      |
| T18.1 | Read `docs/deploy-secret-posture.md` Option A (OIDC commands)                                                 | 12  | 4    | —          |
| T18.2 | Create the GCP pool/provider (**needs gcloud owner auth**)                                                    | 12  | 4    | T18.1      |
| T18.3 | Update `website.yml` deploy to WIF auth                                                                       | 12  | 4    | T18.2      |
| T18.4 | Test deploy on a branch                                                                                       | 12  | 4    | T18.3      |
| T18.5 | Rotate/remove the SA key                                                                                      | 12  | 4    | T18.4      |
| T18.6 | Record + row update                                                                                           | 10  | 4    | T18.5      |
| T19.1 | Draft the owner checklist (OG previewers, social upload, key deletion)                                        | 12  | 4    | —          |
| T19.2 | Link it from `assets/brand/README.md` + CONTRIBUTING                                                          | 10  | 4    | T19.1      |
| T19.3 | TODO row annotated (owner execution pending)                                                                  | 8   | 4    | T19.2      |
| T20.1 | AGENTS.md: annotate.go `,`-split grammar lesson                                                               | 10  | 4    | —          |
| T20.2 | AGENTS.md: lighthouse invocation lesson (post T09)                                                            | 10  | 4    | T09.4      |
| T20.3 | AGENTS.md: btrfs D-state hang lesson                                                                          | 10  | 4    | —          |
| T21.1 | Harvest remaining 19:55 §f items → TODO_LIST/ROADMAP                                                          | 12  | 4    | —          |
| T21.2 | Implement the sitemap assertion in `verify-build.mjs` (exists, ≥N urls, lastmod present)                      | 12  | 4    | —          |
| T21.3 | `test:scripts` coverage for the new assertion                                                                 | 12  | 4    | T21.2      |
| T21.4 | Annotate 07:09 f.17 done; run the battery                                                                     | 9   | 4    | T21.1-3    |
| T22.1 | Observe the gitleaks v3 job on a real push (runtime migration)                                                | 12  | 4    | —          |
| T22.2 | `.oxlintrc.json` vs AGENTS.md disables drift check                                                            | 10  | 4    | —          |
| T22.3 | Record findings                                                                                               | 8   | 4    | T22.1-2    |
| T23.1 | BDSM-World: remote-state check + disposition record (**needs owner answer**)                                  | 12  | 4    | —          |
| T23.2 | media-atlas: mirror intent record (**needs owner answer**)                                                    | 10  | 4    | —          |
| T23.3 | Update fleet docs + AGENTS.md with the decisions                                                              | 8   | 4    | T23.1-2    |
| T24.1 | TODO_LIST harvest of this plan's survivors                                                                    | 12  | 4    | all        |
| T24.2 | docs-gates run                                                                                                | 10  | 4    | T24.1      |
| T24.3 | Status-index rows for new reports                                                                             | 8   | 4    | T24.1      |

**Micro-task count:** 99 — every one ≤12 min. minimatch advisory (BLOCKED, upstream) intentionally unplanned.

---

## 4. Execution graph

```mermaid
flowchart TD
    subgraph T1["1% → 51% — the hero"]
        T01[T01 design system + tokens] --> T02[T02 hero rebuild]
    end
    subgraph T2["4% → 64% — full landing + trust gates"]
        T02 --> T03[T03 sections pass]
        T03 --> T04[T04 verification battery]
        T05[T05 changelog entries]:::quick
        T06[T06 docs-gates CI]:::par
    end
    subgraph T3["20% → 80% — release/doc trust"]
        T07[T07 crm#7 repair]:::par
        T08[T08 consumer CI sweep]:::par
        T09[T09 lighthouse root-cause]:::par
        T10[T10 pkg.go.dev banner]:::par
        T11[T11 archive sweep]:::par
        T12[T12 tab-order asserts]
        T13[T13 v0.6.3 prep]:::gated
    end
    subgraph T4["final 20% — completeness"]
        T14[T14 brand pixel-QA]:::par
        T15[T15 collision test]:::par
        T16[T16 windows leg]:::par
        T17[T17 visual-regression]:::par
        T18[T18 WIF]:::owner
        T19[T19 owner-actions prep]:::par
        T20[T20 AGENTS.md lessons]:::par
        T21[T21 19:55 harvest + sitemap]:::par
        T22[T22 gitleaks + oxlint]:::par
        T23[T23 dispositions]:::owner
        T24[T24 final sweep]
    end
    T01 -.tokens feed.-> T02
    T04 --> SHIP[[deploy + verify live]]
    T05 --> T13
    T06 --> T24
    T02 -.blocks v0.6.3 site story.-> T13
    classDef quick fill:#047857,color:#fff
    classDef par fill:#e2e8f0,color:#111
    classDef gated fill:#fef3c7,color:#111
    classDef owner fill:#fee2e2,color:#111
```

Reading order: **T01 → T02 → T03 → T04 → SHIP** is the landing lane (ships alone if it has to). **T05, T06** are quick parallel wins. Everything marked `par` has no cross-lane dependency and can run any time. `gated` (T13) waits for the concurrent session's WriteFuncWithPerm + the redesign deploy. `owner` (T18, T23) waits on GCP auth / your answers to the three open questions in `docs/status/2026-11-04_20-46…` §g.

---

## 5. Errata vs. the prior plan and backlog

- The 02-12 plan's T01/T02 (sync + v0.6.1 cut) and its fleet lane: DONE today (v0.6.1 + v0.6.2, 17 consumers) — not re-planned.
- "22 importer repos" (02-12) → the fleet universe is 21 entries of which BDSM-World no longer exists on disk (17 bumped, 3 skipped with reasons, 1 absent) — corrected in the 20-45 ledger.
- TODO_LIST's lighthouse row was CLOSED tonight on assertion-only evidence; T09 re-opens the question honestly rather than trusting that closure.
- The 20-46 report's §f.10 ("two forgotten changelog entries") is confirmed still open (`Fixed: Nothing yet.`) and is T05 here.
- New tonight and unplanned anywhere before: the landing redesign (T01–T04) — this plan's reason to exist.
