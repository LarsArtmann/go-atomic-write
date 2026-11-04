# Status Report: Full Docs-Health Audit — 2026-0* Corpus Annotated & Archived, Living Docs Rebuilt

**When:** 2026-11-04 02:06 CET (Wednesday, night session)
**Session scope:** docs-health AUDIT over all twelve `**/2026-0*` files — VERIFY against code, ANNOTATE every numbered item inline (683 verdicts), ARCHIVE the fully-resolved reports, HARVEST survivors into TODO_LIST/ROADMAP, and repair every living doc the audit indicted (CHANGELOG release bookkeeping, FEATURES line-refs, TODO split brain, AGENTS release checklist, README consistency, plus two in-audit code fixes: the Website deploy path-filter bug and the README CI badge).
**Go source touched:** 0. CI workflow touched: yes (`.github/workflows/website.yml` paths, YAML-validated).
**Companion artifacts:** `docs/status/README.md` (new index + archive manifest), `TODO_LIST.md` (rebuilt), git commit `eec3cd6` (daemon-authored message captured the work well).

---

## 0. Scope — what this report verified vs. did not

**Verified with commands/artifacts this session:** every annotation applied via the docs-health `annotate.go repair` tool and gated with `check-rows` (11/11 files COMPLETE); archive completeness via `grep -L '~~'` (empty = every archived file carries strikethroughs); golangci-lint + dprint-format + full BuildFlow dev pass (54 steps, 0 failed; `go test -race` green including the concurrent session's new split test files); both workflow YAMLs parse (`js-yaml`); every moved CHANGELOG entry dated by `git log v0.5.2..v0.6.0` and `git log -S` queries; line references in FEATURES re-derived by grepping the actual function declarations.

**NOT verified, and why:** the live deployed site (~~7 commits unpushed — `origin/master..HEAD = 7`~~ CORRECTED 02:15: the concurrent session pushed at 23:26Z and the Website deploy already fired successfully — `gh run list` success 23:26:55Z, which also proves the path-filter fix works; only 2 doc commits remain unpushed; live rendering still unverified); verification is the existing TODO row); the rendered OG card in real previewers; GCP key hygiene (owner console action); sibling-repo deploy health (fleet scope, declined); a full 100 % audit of the two sub-agents' per-item verdicts (sampled ~10 lines per file + targeted greps of every verdict that changed a disposition — see d.3 for the one that was wrong); whether any repo-external process consumed deleted artifacts (unverifiable, declined in-annotation).

---

## a) FULLY DONE

| #  | Item                                                                                                                                                                                                                   | Evidence                                                                                                                              |
| -- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| 1  | **All twelve `**/2026-0*` files viewed** (11 `.md` + 1 `.html`); 10 `.md` reports + 1 planning review fully ANNOTATED — **683 numbered items carry inline verdicts** (`done at <hash>` / `superseded` / `→ TODO_LIST` / `→ ROADMAP` / `declined` / observation) | `annotate.go repair` pass reports (606 + 77 action lines); `check-rows` → `11 file(s) complete`                                       |
| 2  | **All 11 annotated files ARCHIVED** via `git mv` → `docs/status/archived/` (10) + `docs/planning/archived/` (1); completeness gate re-run post-move (green)                                                                 | `docs/status/archived/` listing; `grep -L '~~'` over archived files prints nothing                                                    |
| 3  | **`[0.6.0] - 2026-09-22` CHANGELOG section reconstructed** — the tag was cut with NO section; 12 commit-verified entries moved out of `[Unreleased]` (go-directive raise, landing messaging, module-path constant, FAQ, verify-build gate, html-validate, i18n, hero/copy/README fixes, deploy unblock, workflow_dispatch) | `git log v0.5.2..v0.6.0`; `git show 2234e1d` (release commit added the entry to `[Unreleased]` and stopped there)                     |
| 4  | **Keep-a-Changelog compare links appended** (10 links, `[Unreleased]` → `[0.1.0]`) — open since the v0.4.0 review's item 34                                                                                                | `CHANGELOG.md` tail                                                                                                                   |
| 5  | **TODO_LIST rebuilt** — stale "Cut release v0.6.0" split brain removed (tag exists since 09-22); replaced by `v0.6.1`; 20 evidence-cited rows total (1 High, 8 Medium, 11 Low incl. 1 BLOCKED)                              | `TODO_LIST.md`; `git rev-list --count origin/master..HEAD` = 7 cited in the deploy row                                                |
| 6  | **FEATURES.md rebuilt** — line references re-derived against the source split (`atomicwrite.go`/`fingerprint.go`/`stage.go`/`commit.go`); missing `WriteWithPerm` row added (shipped v0.5.0, absent since); CI/CD flipped PARTIALLY → FULLY_FUNCTIONAL (deploys proven through 11-03); `~27 GB/s` → `~23–27 GB/s`; brand-system and build-gate rows added | grep of `func ` declarations; `gh run list` (last success 2026-11-03 21:34Z)                                                          |
| 7  | **README repaired** — CI badge added (CI existed since 07-26, badge never did); benchmark figure in Design decisions aligned with the re-measured ranges                                                                    | `README.md:15-17`, `:198`                                                                                                             |
| 8  | **Website deploy path-filter bug FIXED** — `website.yml` paths now include root `CHANGELOG.md`, `README.md`, `go.mod` (push AND pull_request); this was High-impact item f.45/f.46 of the 09-18 report — the exact freeze class that went unnoticed for 7 weeks | `.github/workflows/website.yml:7-19`; YAML parse via js-yaml OK                                                                       |
| 9  | **AGENTS.md hardened** — Release Checklist gained the v0.6.0 lesson ("changelog section BEFORE tagging — v0.6.0 shipped without one") and the raster-regen step (brand rasters are not byte-gated); `commitRename` cleanup line added to architecture; oxlint fleet-sweep gotcha appended to the npm-sweep gotcha | `AGENTS.md` Release Checklist steps 4–9; `docs/status/2026-11-03_23-29…` §5                                                           |
| 10 | **`docs/status/README.md` created** — report index (6 current + 1 LEAVE-ALONE) + archive manifest (11 files, classification + deciding reason) — closes the "status live-index rot" and "archive manifest" requirements at once                                              | `docs/status/README.md`                                                                                                               |
| 11 | **DOMAIN_LANGUAGE.md + CONTRIBUTING.md** — `WriteWithPerm` domain row; "brand artifacts are generated, never hand-edit" section                                                                                            | both files; grep-verified absence before adding                                                                                       |
| 12 | **ROADMAP extended** — ~15 harvested raw ideas routed into themes 2/3/4 (API-compat golden test, Playwright smoke suite, when-not-to-use page, PR preview channel, RSS/print/a11y-statement, personas, visual-polish cluster, Search Console, demo material, siteConfig JSON projection) | `ROADMAP.md` themes 2–4                                                                                                               |
| 13 | **Gates all green post-work** — BuildFlow dev pass 54 steps / 0 failed (go test-race, golangci-lint, art-dupl, gomod guards), explicit `golangci-lint` and `dprint-format` passes, both workflow YAMLs parse                                                               | BuildFlow summaries; `js-yaml` load OK for ci.yml + website.yml                                                                       |
| 14 | **June 2026 code review inline-corrected** — stale premises struck in place (127 LOC claim, "`.bak` good safety net"), status banner added, improvement-plan rows 1–10 individually verdicted (two edge tests survived → TODO_LIST) | `docs/planning/archived/2026-06-03_16-19-full-code-review.md:11-15,29`                                                                |
| 15 | **Sub-agent inventories cross-verified** — one false verdict caught and corrected by grep before it reached an archived file (see d.3)                                                                                     | `grep -c "Release Checklist" AGENTS.md` = 1 vs agent claim "absent"                                                                   |

## b) PARTIALLY DONE

| # | Item                                                                                                              | What works now                                                                                                   | What is missing                                                                                                                                                                                             | Effort |
| - | ----------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| 1 | **Routing granularity of annotations**                                                                            | Every routed item says `→ TODO_LIST: <slug>` / `→ ROADMAP: <theme>`                                              | Slugs are consolidated batch names ("browser verification sweep"), not 1:1 with TODO row titles — mapping lives in this report + `docs/status/README.md`, not in a greppable 1:1 table                                                                 | S      |
| 2 | **Brand-report (23-29) §f harvest**                                                                               | Its High/Medium cluster consolidated into the TODO "Brand system batch" row + the a11y decision row               | ~15 Low/brand-cosmetic items (concepts contact sheet, visual-regression harness, per-page OG cards, generator manifest, dprint-excludes decision, animated favicon, gallery page, OG wordmark A/B, docs-logo hover, jscpd check, AUTHORS mark, README PNG fallback, fleet brand roll-out) were neither routed nor individually ledgered — declined-by-omission, not decided | M      |
| 3 | **Done-ledger tables in the archived §a sections**                                                                | Left untouched (uniform UNTOUCHED passes the gate today)                                                          | No `<!-- check-rows:ledger -->` exemption comments inserted — a future partial strike in those tables would silently flip them to PARTIAL; exemption would make the intent explicit                                                                          | S      |
| 4 | **The 06-28 HTML dashboard**                                                                                      | Classified LEAVE-ALONE in the manifest (referenced by CHANGELOG [0.2.0])                                          | No superseded-banner rendered in the HTML itself; a reader landing there directly gets no staleness signal                                                                                                   | S      |
| 5 | **Sub-agent verdict verification**                                                                                | Every verdict that CHANGED a disposition (done vs open vs routed) was re-verified against code                    | Verdicts that merely restated in-file evidence were trusted on a ~10-line-per-file sample; the line numbers were verified only where sampled (one out-of-range citation existed in agent output and was caught by `emit-keys`)                              | M      |
| 6 | **Website CI verification of the paths fix**                                                                      | YAML-validated locally; BuildFlow green (it does not exercise GitHub workflows)                                    | The fix has not fired on a real push yet — first push will prove it (and deploy the backlog of content simultaneously)                                                                                       | S (rides next push) |
| 7 | **Fresh 11-03 reports**                                                                                           | Indexed in `docs/status/README.md`; their open items routed where High/Medium                                     | They are unannotated by design (fresh snapshots) — the 20-42/22-20/22-21/23-11 reports were only skimmed for next-tasks, not fully read                                                                     | (later ANNOTATE pass) |

## c) NOT STARTED

Everything below is routed and evidenced in `TODO_LIST.md`; listed here for the snapshot. (Impact/Effort as rated there.)

1. Cut release `v0.6.1` from `[Unreleased]` (High/S) — depends on nothing; the fleet bump depends on it.
2. Fleet bump: 22 enumerated LarsArtmann importers via go-ecosystem-upgrade (Medium/L).
3. ~~Push `master` + deploy verification~~ Deploy already fired 23:26Z: brand, benchmarks, changelog live (rendering unverified); OG card in real previewers; post-deploy fetch assertion in the workflow (Low/S but blocks all live verification).
4. Browser verification sweep: 320px screenshots, axe/contrast over the FULL palette, keyboard tab order, theme first-paint (Medium/L).
5. Lighthouse: baseline-measured budgets + docs-page scores + LHCI gate wired into `website.yml` (Medium/M).
6. Newsletter CSP: `form-action 'self'` vs Buttondown POST — never browser-tested (Medium/S).
7. Brand system batch: `brand:png`, SVG parse in `brand:check`, `brand-tokens.json`, wordmark real-font measure, footer screenshot QA, pixel-QA script, logo a11y decision (Medium/M).
8. Build-gate hardening: source-level hero parity, docs-page meta assertion, FAQ count, CSP unhashed-inline, exact copy-button count, stale-API grep, `pnpm run check` (Medium/M).
9. Landing voice pass completion — titles/eyebrows/badge/Footer/Newsletter/CTA (Medium/S).
10. FAQ ↔ `guides/benchmarks` + `guides/error-handling` links (Medium/S).
11. Fleet npm sweep tracker — identify the sweep, teach it the constraints (Medium/M).
12. Fingerprint edge tests: `FingerprintFile(directory)`, `Matches(nil)` — open since the June review (Low/S).
13. `animations.js` after ClientRouter navigation — verified binding-once bug (Low/S).
14. OG tagline/CTA voice in the generator — generated `og-image.svg` still says "TOCTOU-safe" + github URL (Low/S).
15. PWA/app icons: apple-touch-icon 180, manifest 192/512 (Low/S).
16. Docs/SEO hygiene batch: sitemap lastmod, `go get` in docs entry points, `theme-color` light, `og:image:type`/`site_name`/`alt`, `engines.node`, flake.lock update, internal link checker (Low/M).
17. `CopyButton.astro` extraction (Low/S).
18. Unique-temp-file invariant test (Low/S).
19. `gitleaks` step in CI (Low/S).
20. Deploy-secret posture: WIF federation vs rotation cadence; delete 7 stale `firebase-adminsdk` keys (owner GCP action) (Low/S).
21. CONTRIBUTING: release section + deploy/secret recipe (Low/S).
22. Everything in ROADMAP themes 1–4 that has no TODO row (raw ideas by design).

## d) TOTALLY FUCKED UP

Honest list. Nothing was destroyed; these are real defects of THIS session's process.

1. **I missed one of the twelve files until mid-verification — and my gate did not catch it.** `2026-09-17_20-57_website-unfreeze-and-ci-deploy-repair.md` was in neither sub-agent's assignment nor my own tracking; I caught it by recounting coverage against the glob, NOT via `check-rows` (which had passed the file with zero specs because uniformly-UNTOUCHED tables are not offenders). Root cause: I built the agent assignment list from memory of the glob instead of mechanically (`fd 2026-0 docs/`). The fix was cheap (66+11 lines annotated, file archived) but the miss class is exactly what docs-health exists to prevent. Mitigation now: `docs/status/README.md` indexes every report; lesson below in (e).
2. **Three failed dry-runs before the annotation batch applied** — I composed ~600 spec lines before reading the `repair` spec grammar in the tool source; `kind:value` joined-field vs separate-tab, a sed character class missing `s:`, and one duplicate key. The atomic dry-run contained the blast radius (zero partial writes), but the iterations were pure waste.
3. **One sub-agent verdict was flat wrong and nearly shipped:** "release checklist in AGENTS.md — grep: absent" (09-17_14-22 §f.21). The checklist exists (`grep -c` = 1). I caught it by cross-checking the two agents against each other, then grep. If I had trusted the agent wholesale, an archived report would carry a false STILL-OPEN verdict inviting future sessions to re-create existing work.
4. **An off-by-two line range in my own 10-37 spec** (§b/§c tables) — my verdicts were written against section-relative rows, not the viewed absolute lines. Caught by re-deriving from the view output before apply; had I applied blind, 24 rows would carry verdicts for their neighbors.
5. **I tripped a trap this repo's own fresh report documented.** `buildflow -s dprint` → "unknown tool name". The 2026-11-03 23-29 report literally recorded the same failure and the correct name (`-s dprint-format`). I read that report hours earlier. Reading-not-heeded is its own failure mode; the fix cost one catalog lookup.
6. **Route-slug drift between annotations and TODO rows.** Annotation values like "Lighthouse CI gate (browser batch)" don't textually match the TODO row "Browser verification sweep". Semantically sound, textually ungreppable. Root cause: I annotated while still designing the TODO consolidation, instead of fixing row names first.
7. **The two in-audit "fixes" (website.yml paths, README badge) were verified by YAML parse and BuildFlow, not by a triggered workflow run.** Confidence is high (list-extended edit), but "proven on GitHub" remains open — correctly reflected in b.6 rather than claimed.

## e) WHAT WE SHOULD IMPROVE

1. **Mechanize the annotate-target list.** `fd '2026-' docs/ > targets` → assign every file to an agent or self, before any dispatch. The d.1 miss is unexcusable given the glob was already in hand.
2. **Read the tool's spec grammar before composing specs, not after the first failure.** `annotate.go`'s usage output named the modes but not the TSV shape; one `grep -A20 loadSpecs` upfront would have saved three dry-run cycles.
3. **Bound and declare agent-trust explicitly.** Either verify 100 % of verdicts (cheap for verdict-relevant greps) or state the sampling rate in the report. "Sampled" must be written down, not implied.
4. **Fix route slugs before annotating.** Build the TODO rows first, then generate annotation values from the row titles — makes the archived verdicts greppable directly into the backlog.
5. **Insert `check-rows:ledger` exemptions while archiving.** One comment per done-ledger table makes future intent explicit and future partial strikes loud.
6. **Heed in-repo documentation as a gate, not trivia.** The dprint step-name trap was pre-documented; a `grep -i "unknown tool" docs/status/*.md` before invoking a BuildFlow step would be a cheap reflex.
7. **Banner stale HTML dashboards even when LEAVE-ALONE** — a one-line injected banner is cheaper than the next reader's confusion.
8. **Keep the two-score audit discipline** — the Accuracy/Fitness split from the health-report format is what surfaced the CHANGELOG bookkeeping Critical (a purely factual check would have passed it as "entries exist").

## f) 50 things we should get done next

Ranked: impact first. (HARVEST already landed the High/Medium/Low survivors in `TODO_LIST.md`; the tail below is ROADMAP-fuel per the docs-health anti-pattern — treat 30–50 as a brainstorm, not a commitment.)

| #  | Task                                                                                              | Impact | Effort |
| -- | ------------------------------------------------------------------------------------------------- | ------ | ------ |
| 1  | Cut `v0.6.1` (checklist steps 1–9; changelog section BEFORE tag — the new rule)                   | High   | S      |
| 2  | ~~Push `master` (7 commits)~~ Verify the 23:26Z deploy rendered: brand + benchmarks + changelog live; OG card in real previewers | High   | S      |
| 3  | Post-release: `pkg.go.dev` shows v0.6.1; `go get @v0.6.1` resolves; changelog page synced         | High   | S      |
| 4  | Fleet bump: all 22 enumerated importer repos via go-ecosystem-upgrade                             | High   | L      |
| 5  | Browser verification sweep: 320px + desktop screenshots, theme first-paint, ClientRouter nav      | High   | M      |
| 6  | Full-palette WCAG contrast audit (amber/danger/success/code-comment, both themes) + fixes         | High   | M      |
| 7  | Lighthouse: baseline-measured budgets, docs-page scores, LHCI gate in `website.yml`               | High   | M      |
| 8  | Newsletter CSP: browser-test the Buttondown POST, fix `form-action` (header CSP or allowlist)     | High   | S      |
| 9  | Header/footer logo a11y: `aria-hidden` vs `aria-label` decision, apply to both instances           | High   | S      |
| 10 | Codify pixel-QA: `qa-brand-pixels.mjs` asserting glyph ink in headless screenshots                 | High   | S      |
| 11 | `brand-tokens.json` — one palette source for generator + `global.css`                             | High   | M      |
| 12 | Build-gate hardening batch (source-level parity, meta-description, FAQ count, CSP assert, button count, stale-API grep, `pnpm run check`) | High | M |
| 13 | Landing voice pass completion (titles/eyebrows/badge/Footer/Newsletter/CTA)                        | Medium | S      |
| 14 | FAQ ↔ `guides/benchmarks` / `guides/error-handling` links                                          | Medium | S      |
| 15 | `brand:png` script — rasters regenerated from code, not README                                     | Medium | S      |
| 16 | SVG well-formedness parse inside `brand:check`                                                     | Medium | S      |
| 17 | Wordmark: measure real Space Grotesk width headlessly, fix canvas or record the air                | Medium | S      |
| 18 | Footer 20px logo screenshot QA                                                                     | Medium | S      |
| 19 | Favicon-scale CI gate (render 16/22/32, assert non-blank + glyph color)                            | Medium | S      |
| 20 | OG tagline/CTA voice in the generator (kill "TOCTOU-safe", pkg.go.dev CTA) + re-raster PNG         | Medium | S      |
| 21 | Test OG card in Slack/Twitter/Facebook previewers                                                  | Medium | S      |
| 22 | Compress `og-image.png` (59 KB → ~30 KB)                                                           | Medium | S      |
| 23 | apple-touch-icon 180×180 + link tag; manifest icons 192/512 PNG                                    | Medium | S      |
| 24 | GitHub social preview 1280×640 built from the lockup                                               | Medium | S      |
| 25 | Fingerprint edge tests: `FingerprintFile(directory)` + `Matches(nil)`                              | Medium | S      |
| 26 | `animations.js`: hook `astro:page-load` (fixes post-navigation dead animations)                    | Medium | S      |
| 27 | Internal link checker for docs + README (CI-able)                                                  | Medium | S      |
| 28 | Build-script regression tests: `fix-csp.mjs` idempotency, `sync-changelog.mjs` escaping            | Medium | S      |
| 29 | `gitleaks` step in CI                                                                              | Medium | S      |
| 30 | Deploy-secret posture: WIF federation vs rotation cadence; delete 7 stale SA keys (GCP console)    | Medium | S      |
| 31 | CONTRIBUTING: release-cutting section + deploy/secret recipe                                       | Medium | S      |
| 32 | Fleet npm sweep: identify, teach it the TS6/pnpm constraints or exclude this repo                  | Medium | M      |
| 33 | Sitemap `lastmod` + `go get` in every docs entry point                                             | Low    | S      |
| 34 | `theme-color` light variant + `og:image:type`/`og:site_name`/`og:image:alt`/`twitter:image:alt`    | Low    | S      |
| 35 | `engines.node` field + `website/flake.lock` update                                                 | Low    | S      |
| 36 | `CopyButton.astro` extraction (repeated ~300-char class string)                                    | Low    | S      |
| 37 | Unique-temp-file invariant test (`stage.go`)                                                       | Low    | S      |
| 38 | "When not to use this" docs page (promote FAQ item)                                                | Low    | M      |
| 39 | Port the landing FAQ into docs (site-search coverage)                                              | Low    | M      |
| 40 | Surface `WriteIfChanged`/`WriteFunc`/`WriteWithPerm` on the landing page                           | Low    | M      |
| 41 | Playwright smoke suite (copy payload, FAQ opens, install bar at 320px)                             | Low    | M      |
| 42 | PR preview channel (Firebase preview per PR)                                                       | Low    | M      |
| 43 | PNG raster lockups (dark/light @2×) + decide stacked-lockup usage                                  | Low    | S      |
| 44 | 3-concept contact sheet validating the atomic-swap mark post-hoc                                   | Low    | M      |
| 45 | Screenshot visual-regression harness for header/docs/OG surfaces                                   | Low    | M      |
| 46 | Per-page OG cards for docs routes                                                                  | Low    | M      |
| 47 | Brand usage rules (clear space, min sizes) in `assets/brand/README.md`                             | Low    | S      |
| 48 | Generator oxlint warnings fixed or rationale-ignored; `.jscpd` scripts/** exclusion check          | Low    | S      |
| 49 | Annotate the 11-03 reports once their items land (keep the snapshot chain accurate)                | Low    | S      |
| 50 | Fleet: roll the brand-system pattern (one geometry + byte-check) to the next unbranded sibling     | Low    | M      |

## g) Questions I cannot answer myself

**Q1. ~~Push now, or do you want to review locally first?~~ ANSWERED BY REALITY (02:15):** the concurrent session pushed at 23:26Z and production deployed successfully — the brand is live (rendering unverified). Revised question: ~~Seven commits are unpushed~~ do you want a live-URL review of the deployed brand now, or proceed straight to the v0.6.1 cut? My recommendation: skim the live site, then cut. the website.yml path-filter fix and the full brand identity. The next push triggers BOTH workflows: CI (with the new gates) and a production deploy of the brand + all backlogged content. Everything is gate-green locally; the only thing I cannot decide is whether the atomic-swap logo goes live before you have eyeballed it (`pnpm run preview` from `website/` renders it). My recommendation: push, review the live URL.

**Q2. Is the Buttondown newsletter form supposed to be live?** The meta CSP's `form-action 'self'` almost certainly blocks its cross-origin POST, and this has been true since the CSP shipped (v0.4.0) with zero verification either way. If the form is meant to work, the CSP fix jumps ~20 rows up the queue; if it is decorative-for-now (or being removed), I would drop the concern and note the form as non-functional in the FAQ instead.

**Q3. `v0.6.1` or `v0.7.0` for the next release?** `[Unreleased]` mixes a behavior fix (temp-file leak), internal restructuring (source split, no API change), website features (brand system, gates), and CI hardening. Patch fits the fix-only reading; minor fits the additive reading (v0.6.0 itself was cut for a one-line go.mod change, so this repo's precedent leans minor). The fleet bump consumes whatever we cut, so the number sets the fleet's floor — your call.

---

**Waiting for instructions.**
