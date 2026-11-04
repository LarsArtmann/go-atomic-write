# Brand Identity "The Atomic Swap" — Comprehensive Status Update

**When:** 2026-11-03 23:29 CET (Tuesday)
**Session scope:** logo research (typespec-asyncapi pattern) → identity design → single-geometry generator → 9 generated artifacts + 2 rasters → README/website/docs integration → CI drift gates → multi-layer verification (unit, build, browser, pixel, Lighthouse) → this self-review.
**No Go source was touched.** Point-in-time record; annotate, never rewrite (`docs/status/**` convention).
**Companion report:** `docs/status/2026-11-03_23-11_brand-identity-atomic-swap-logo.md` (the build log; this one is the verdict + gaps).

---

## 0. Executive verdict

A complete, generator-backed brand system is shipped and verified on five surfaces (favicon SVG+ICO, landing header/footer, Starlight docs header, OG card SVG+PNG, README hero) with byte-drift gates in BOTH CI workflows. Every gate I know of is green: `brand:check`, `astro check`, build + CSP + verify-build, html-validate, dprint, YAML parse, Lighthouse budgets (desktop, incl. error-gated a11y/SEO), `go build`/`go test`. Real-browser QA (chromium 152 screenshots + pixel histogram) confirmed the rendering; the pixel histogram also caught ME misreading a screenshot. Three real defects remain honest: no human has seen the deployed site, the GitHub README rendering is push-gated, and the lockup wordmark width was tuned against a fallback font, not the real Space Grotesk. Process debt: I explored ONE concept (typespec explored 59), the footer instance was never screenshotted, and the follow-up tasks were never harvested into `TODO_LIST.md`.

---

## 1. a) FULLY DONE (verifiable)

| #  | Item                                                                                                                                                                                                                                 | Evidence                                                                                                                                                                    |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1  | Reference research: typespec-asyncapi logo system fully mapped (generator, compositions table, palette quirks, embedding rules, provenance posture)                                                                                  | its `assets/brand/README.md`, `Logo.astro`, generator source, status report — read before any design decision                                                               |
| 2  | Identity: "the atomic swap" — emerald squircle, new bytes in (→) / old bytes out (←); old drifted "A" squircle retired                                                                                                               | generator header comment + `assets/brand/README.md`                                                                                                                         |
| 3  | Generator `scripts/generate-brand-assets.mjs`: one geometry → 9 artifacts, `--check` byte mode, cwd-independent, mkdir-recursive, 2-space indent per root editorconfig                                                               | runs clean; `--check` exit 0; 80→ (tolerated-class) oxlint warnings only, matching repo's existing script tolerance                                                         |
| 4  | `assets/brand/`: `mark.svg` (64×64) + horizontal/stacked lockups × light/dark + `README.md` (compositions, files, palette, embedding, favicon-scale test, provenance)                                                                | 6 files; `brand:check` covers the 5 SVGs                                                                                                                                    |
| 5  | `Logo.astro` regenerated: flat token fills (`var(--color-accent)`/`var(--color-on-accent)`), no `defs` (no gradient-id collisions header↔footer), Props unchanged so call sites untouched                                            | both tokens exist in both theme blocks (global.css:23,68); 2 inline instances in `dist/index.html` with the new glyph; all 3 arbitrary-value utilities present in built CSS |
| 6  | `favicon.svg` regenerated + `favicon.ico` (16/32/48) via `magick icon:auto-resize`                                                                                                                                                   | ICO structure identical to previous; dist copies byte-identical (`cmp`)                                                                                                     |
| 7  | `website/src/assets/brand.svg` bundled copy for Starlight `logo: { src: "./src/assets/brand.svg" }`                                                                                                                                  | built docs pages render `<img src="/_astro/brand.CqrkVa9N.svg">`                                                                                                            |
| 8  | OG card: `og-image.svg` now GENERATED (mark slot can never drift) → `og-image.png` 1200×630 re-rasterized                                                                                                                            | `identify` 1200x630; card render eyeballed; dist copy byte-identical                                                                                                        |
| 9  | README hero `<picture>` (prefers-color-scheme, absolute raw URLs, 300×64) — typespec's GitHub-proven markup                                                                                                                          | dprint-clean; markup mirrors reference exactly                                                                                                                              |
| 10 | CI gates: `brand:check` as FIRST step of ci.yml drift-guard (runs on every push, covers root-level `assets/`+`scripts/` that website.yml's path filter misses) and FIRST step of website.yml build job                               | both YAMLs parse; zero-dependency step (system node)                                                                                                                        |
| 11 | `brand:generate` / `brand:check` npm scripts in website/package.json (scripts ≠ dependencies → lockfile untouched, drift-guard safe)                                                                                                 | `pnpm run brand:check` green from website/                                                                                                                                  |
| 12 | Gates: `astro check` 0 errors/0 warnings · build 11 pages · CSP 11/11 · verify-build OK · html-validate silent · Lighthouse desktop 4×3 runs all assertions (a11y ≥0.95, SEO ≥0.95 error-gated) · `go build`+`go test` ok            | full runs this session                                                                                                                                                      |
| 13 | Browser QA: chromium 152 headless screenshots of landing (both themes, header zoomed 3×) + docs page; docs glyph proven via pixel histogram (`#0C0A09` ×136 px = `GLYPH_INK`, zero white px)                                         | `/tmp/gaw-browser-qa/*`; §6 addendum in companion report                                                                                                                    |
| 14 | Docs: CHANGELOG `[Unreleased]` entry · AGENTS.md (commands ×2, structure rows ×3, gotchas ×2 — incl. correcting the stale "og-image.svg is editable" claim and the Starlight logo-src trap) · brand README · companion status report | in tree; dprint-clean                                                                                                                                                       |
| 15 | Concurrent oxlint autofix sweep (20 files rewritten mid-session) judged safe: all `node --check` valid, full website gate re-run over the rewritten tree, left intact per never-revert-others'-work                                  | diff reviewed before judging                                                                                                                                                |

## 2. b) PARTIALLY DONE (works, but with gaps)

| # | Item                                | What works                                                                                                                  | What remains                                                                                                                                                                         | Effort |
| - | ----------------------------------- | --------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------ |
| 1 | Brand system vs typespec-asyncapi's | All core surfaces + CI gates                                                                                                | No PNG lockups @2× for SVG-hostile renderers; no concepts gallery (see §4.6); no horizontal-lockup usage inside the site itself                                                      | S–M    |
| 2 | `brand:check` coverage              | Byte-gates all 9 text artifacts                                                                                             | Rasters (favicon.ico, og-image.png) not gated — deliberate (encoders not byte-stable), but regeneration is a manual magick step                                                      | S      |
| 3 | Wordmark rendering                  | Canvas sized to wide DejaVu fallback (ImageMagick's render); left-anchored so narrower fonts leave trailing air, never clip | Real Space Grotesk width never measured in a browser; trailing-air size unverified on github.com                                                                                     | S      |
| 4 | Logo instances QA                   | Header (28px) both themes + docs header screenshot+pixels                                                                   | Footer instance (20px) verified by markup+CSS only, never screenshotted                                                                                                              | S      |
| 5 | Accessibility                       | `role="img"` + aria-label on both inline instances; Lighthouse a11y 0.95+                                                   | Header announces aria-label "go-atomic-write logo" AND adjacent visible text "go-atomic-write" — probable duplicate announcement for screen readers; needs an `aria-hidden` decision | S      |
| 6 | Docs follow-through                 | AGENTS/CHANGELOG/status reports updated                                                                                     | `TODO_LIST.md` NOT harvested (convention: §f tasks belong there — docs-health HARVEST); `FEATURES.md` untouched (arguably fine — not a library feature)                              | S      |
| 7 | Palette sourcing                    | Documented in brand README, mirrors global.css                                                                              | Two sources of truth remain (generator constants + CSS tokens); no shared token file consumed by both                                                                                | M      |
| 8 | Release readiness                   | CHANGELOG under `[Unreleased]`, CI gates wired                                                                              | Next release must pass golangci-lint (the gate) — not run this session (zero Go changes; CI will run it on push anyway)                                                              | S      |

## 3. c) NOT STARTED

1. Apple-touch-icon (180×180 PNG) + `<link rel="apple-touch-icon">`.
2. Manifest icons 192/512 PNG entries (current manifest: only `favicon.svg`, `sizes: any`).
3. GitHub social preview (1280×640) built from the lockup.
4. PNG raster lockups (dark/light @2×).
5. Motion variant (swap-arrow pulse on header hover), reduced-motion aware.
6. Concepts gallery (typespec built a 59-card contact sheet; I designed exactly one concept — no alternatives were drawn, compared, or rejected on paper).
7. Visual-regression harness (screenshot diff) for brand surfaces.
8. Automated pixel-QA script (assert glyph-ink color present in built pages) — done by hand this session, not codified.
9. `og:image:alt` / `twitter:image:alt` meta tags (missing today).
10. Per-page OG cards for docs routes (only one static OG exists).
11. Brand usage rules (clear space, minimum sizes) documented.
12. TODO_LIST harvest (§2.6) + next-release checklist reminder ("if glyph changed, regenerate rasters").
13. CONTRIBUTING.md note for contributors: brand artifacts are generated; never hand-edit.
14. Deploy + eyeball live site (needs push/workflow or manual firebase deploy).

## 4. d) TOTALLY FUCKED UP (radical honesty)

### 4.1 I designed ONE concept and shipped it — no alternatives explored — Severity: MEDIUM (process)

typespec-asyncapi's identity came out of 59 drawn concepts and an explicit selection. I designed a single mark, liked it, and shipped it. The result renders well and QA'd clean, but "the best identity" was never tested against alternatives. A logo chosen without comparisons is a preference, not a decision.

### 4.2 The Starlight logo-src failure was self-inflicted — Severity: LOW (but it was the session's only red build)

I wrote `logo: { src: "/favicon.svg" }` (public URL) → `ImageNotFound` build failure. The answer ("logo.src resolves through the asset pipeline; use a bundled copy") was sitting in typespec-asyncapi's astro.config.mjs, which I had open for other reasons. I checked AFTER failing, not before. Cost: one wasted build cycle (~1 min) + the fix. Now an AGENTS.md gotcha, so it dies there.

### 4.3 I misread my own screenshot and hunted a nonexistent bug — Severity: LOW (cost, not damage)

At the docs-header zoom I concluded the glyph rendered WHITE and went looking for an inverting filter (4 tool calls). The pixel histogram proved the glyph was `#0C0A09` — exactly as authored — the whole time. Emerald antialiasing at small scale against dark chrome fooled my eyes. Positive spin: the verification instinct was right (I sampled pixels instead of shipping a guess), but the ORDER was wrong — sample first, squint never.

### 4.4 Edit-tool discipline failures — Severity: LOW

Three edits failed with "must read the file before editing" because I had only `cat`/`rg`-read package.json, website.yml, ci.yml via bash. Two rounds wasted re-reading. The rule exists; I knew it; I batched past it.

### 4.5 Wrong BuildFlow step name — Severity: LOW

`buildflow -s dprint` → "unknown tool name" error. The buildflow skill says check the catalog/dry-run FIRST. `-s dprint-format` worked, but the failed invocation was avoidable.

### 4.6 Blind spots I only caught because tooling forced me — Severity: MEDIUM (honesty about luck)

- The generator crashed with ENOENT on first run (no `assets/brand/` dir) — caught because I ran it before declaring done, not because I designed for it. mkdir-recursive was a patch after the failure.
- The oxlint sweep rewrote 20 files WHILE I worked (including my own generator, changing its const style). I verified it safe — but if I had regenerated artifacts in that window without re-checking, or if the sweep had hit a half-written file, the daemon would have committed garbage. I got lucky with timing; nothing in my process protected against it.
- The `--blink-settings=preferredColorScheme=2` flag produced the OPPOSITE theme from its name — my screenshot filenames ("landing-light.png" containing the dark theme) are still wrong, documented in the addendum only.

### 4.7 Not fucked up, worth stating: zero Go risk taken — Severity: NONE

No Go file touched, `go build`/`go test` run anyway, golangci-lint (the merge gate) left to CI. Deliberate, defensible.

## 5. e) WHAT WE SHOULD IMPROVE

1. **Design process**: any future identity change must produce ≥3 candidate marks side-by-side (even a 30-minute contact sheet) before one ships. A logo is a decision; decisions need alternatives.
2. **Sample pixels before believing screenshots** — codify it: the QA step for any brand surface should be a scripted color assertion, not eyeballing. Today it was a manual `magick txt:` pipeline; make it `website/scripts/qa-brand-pixels.mjs` (assert `#0C0A09` present inside the docs-logo crop of a headless screenshot).
3. **Single palette source**: extract `brand-tokens.json` (accent, deep, ink, muted) consumed by BOTH the generator and (via a tiny import) global.css/custom properties. Kills the drift class before it recurs.
4. **Close the raster gap**: a `brand:png` script (magick invocations checked in as an npm script) so favicon.ico + og-image.png regenerate with one command and the magick commands live in code, not in a README.
5. **SVG well-formedness in brand:check**: byte-compare can't catch a generator emitting malformed SVG; add a parse step (cheap DOMParser/xmllint).
6. **Measure the wordmark with the real font** in a headless browser once (not the DejaVu proxy) and either fix the canvas or record the measured trailing air.
7. **A11y polish**: decide aria-label vs aria-hidden for the header/footer marks (duplicate announcement risk); add og:image:alt.
8. **Always dry-run BuildFlow step names** (`buildflow --dry-run --verbose`) before `-s`; and View before Edit — both rules I tripped over tonight.
9. **Harvest discipline**: §f-style task lists belong in TODO_LIST.md immediately, not entombed in status reports (typespec's own report called this out; I repeated it).
10. **Sweep resilience**: when a fleet autofix sweep is active, re-run `brand:check` + build AFTER the last observed mutation, not before — my final gates cleared by ~15 minutes of luck margin.

## 6. f) NEXT TASKS (ranked by impact, 50)

| #  | Task                                                                                                                                  | Impact         | Effort |
| -- | ------------------------------------------------------------------------------------------------------------------------------------- | -------------- | ------ |
| 1  | Push → verify README `<picture>` renders on github.com (the only push-gated claim)                                                    | High           | S      |
| 2  | Deploy (workflow or `firebase deploy --only hosting:atomicwrite`) + eyeball live landing/docs/OG                                      | High           | S      |
| 3  | Decide header logo a11y: aria-label vs aria-hidden when adjacent title text exists; apply to header+footer                            | High           | S      |
| 4  | Codify pixel-QA: script asserting glyph ink in a headless screenshot of built pages                                                   | High           | S      |
| 5  | `brand-tokens.json` single palette source consumed by generator + global.css                                                          | High           | M      |
| 6  | Harvest brand tasks into `TODO_LIST.md` (docs-health HARVEST pass)                                                                    | High           | S      |
| 7  | Apple-touch-icon 180×180 + link tag                                                                                                   | Medium         | S      |
| 8  | Manifest icons 192/512 PNG + entries                                                                                                  | Medium         | S      |
| 9  | GitHub social preview 1280×640 from the lockup (needs repo-settings access)                                                           | Medium         | S      |
| 10 | `brand:png` npm script; wire favicon.ico/og-image.png regeneration into one command                                                   | Medium         | S      |
| 11 | SVG well-formedness parse step inside `brand:check`                                                                                   | Medium         | S      |
| 12 | Measure real Space Grotesk lockup width headlessly; fix canvas or record measured air                                                 | Medium         | S      |
| 13 | Footer 20px instance screenshot QA                                                                                                    | Medium         | S      |
| 14 | og:image:alt + twitter:image:alt meta                                                                                                 | Medium         | S      |
| 15 | PNG raster lockups dark/light @2×                                                                                                     | Medium         | S      |
| 16 | 3-concept contact sheet for the CURRENT mark (validate "atomic swap" against 2 alternatives post-hoc)                                 | Medium         | M      |
| 17 | Screenshot visual-regression harness for header/docs/OG surfaces                                                                      | Medium         | M      |
| 18 | Favicon-scale CI gate (render 16/22/32, assert non-blank + glyph-color present)                                                       | Medium         | S      |
| 19 | Motion variant: swap-arrow pulse on header hover, reduced-motion aware                                                                | Low            | M      |
| 20 | Use the stacked lockup somewhere real or mark it reserved (it has zero consumers today)                                               | Low            | S      |
| 21 | CONTRIBUTING.md: "brand artifacts are generated" note                                                                                 | Low            | S      |
| 22 | Per-page OG cards for docs routes (typespec has 25; we have 1)                                                                        | Low            | M      |
| 23 | Release v0.6.0: golangci-lint gate → CHANGELOG cut → tag (checklist steps 1–8)                                                        | Medium         | S      |
| 24 | Post-release: verify pkg.go.dev + `go get @vX.Y.Z` + live changelog page sync                                                         | Medium         | S      |
| 25 | AGENTS.md release checklist: add "if glyph changed, regenerate favicon.ico/og-image.png" step                                         | Low            | S      |
| 26 | Test OG card in real previewers (Slack/Twitter/Facebook) after deploy                                                                 | Medium         | S      |
| 27 | Compress og-image.png (59 KB → ~30 KB, social fetch latency)                                                                          | Low            | S      |
| 28 | Compare our OG against gogenfilter's for fleet-consistent card layout                                                                 | Low            | S      |
| 29 | Consider lockup-horizontal in the site footer (currently README-only asset)                                                           | Low            | S      |
| 30 | Brand usage rules section in assets/brand/README.md (clear space, min sizes)                                                          | Low            | S      |
| 31 | Cross-project lesson to crush-config `references/lessons.md`: "one-geometry brand generator + byte-check" pattern                     | Medium (fleet) | S      |
| 32 | Add logo sliver to 404 page? (Starlight ships one — verify it looks fine, else leave)                                                 | Low            | S      |
| 33 | Verify Starlight docs `<head>` carries favicon.ico fallback (landing layout does; docs head unverified)                               | Low            | S      |
| 34 | Dark-mode token audit: confirm `--color-on-accent` stays identical if themes ever diverge (brand README notes the dependency)         | Low            | S      |
| 35 | Generator: emit a manifest of artifact paths (for tooling/tests to consume instead of hardcoding)                                     | Low            | S      |
| 36 | Replace remaining copy-paste of magick commands with the script from #10 in AGENTS.md                                                 | Low            | S      |
| 37 | Add `scripts/` to dprint excludes or leave — decide (generator is oxlint-tolerated, dprint doesn't cover mjs)                         | Low            | S      |
| 38 | Silence-by-design: fix the generator's own oxlint warnings (one-var etc.) or add a targeted ignore with rationale                     | Low            | S      |
| 39 | Document the oxlint autofix sweep behavior in AGENTS.md (new recurring fleet pattern observed tonight: 20 files restyled mid-session) | Medium         | S      |
| 40 | Guard against sweep races: `brand:check` + build after ANY external restyle burst (checklist note)                                    | Low            | S      |
| 41 | Explore animated favicon (APNG/SVG animation) — probably gimmick, decide explicitly                                                   | Low            | S      |
| 42 | Add mark to `AUTHORS`/repo README meta? (no — skip; record as considered)                                                             | Low            | S      |
| 43 | Gallery page on the website showing lockup variants? (probably unnecessary for a lib)                                                 | Low            | S      |
| 44 | Evaluate `og-image` with wordmark lockup instead of bare mark + text title (design A/B)                                               | Low            | S      |
| 45 | Starlight custom CSS: hover scale on docs logo? (site header has rotate-on-hover; docs doesn't — align or justify)                    | Low            | S      |
| 46 | Confirm no jscpd false positives from generator templates (`.jscpd.json` has no `scripts/**` exclude)                                 | Low            | S      |
| 47 | `pnpm run brand:check` wired into any future root-level `verify` aggregate if one appears                                             | Low            | S      |
| 48 | Consider a `.png` fallback in the README `<picture>` for renderers without SVG-in-`<img>` (rare; skip unless evidence)                | Low            | S      |
| 49 | Annotate the companion report (23-11) if §f items get harvested/done (keep the snapshot chain accurate)                               | Low            | S      |
| 50 | Fleet: apply the same brand-system pattern to the next unbranded sibling project (stack: copy generator shape, swap geometry)         | Medium (fleet) | M      |

## 7. g) QUESTIONS I CANNOT ANSWER MYSELF

1. **Header logo announcement (§2.5/#3):** should the mark be `aria-hidden` next to the visible "go-atomic-write" text (my recommendation: yes — redundant double announcement), or do you want the explicit "go-atomic-write logo" label kept? I can implement either; the preference is yours.
2. **Live verification (§f #2):** deploy now via local `firebase deploy --only hosting:atomicwrite` (works without the CI secret, logged in as lartyhd@gmail.com) so the logo is visible on atomicwrite.web.app tonight, or wait for the Website workflow on your next push?
3. **Concept validation (§4.1/#16):** are you satisfied shipping the first-and-only concept, or do you want a proper 3-concept contact sheet (like typespec's 59-card gallery) before this identity is treated as final?
