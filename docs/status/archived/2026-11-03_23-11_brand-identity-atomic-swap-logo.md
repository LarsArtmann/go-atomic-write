# Brand Identity: "The Atomic Swap" Logo — Status Report

> **Archived 2026-11-05 (cadence).** Live remainder: none — the brand system shipped in v0.6.2 (generator, tokens, rasters, gates); the ink-palette pixel-QA extension landed 2026-11-05 (commit `4f0d555`).

**When:** 2026-11-03 23:11 CET
**Session scope:** logo concept → single-geometry generator → 9 generated artifacts → README + website + docs-header integration → CI drift gates → verification.
**No Go source was touched.** Point-in-time record; annotate, never rewrite (`docs/status/**` convention).

---

## 0. Executive verdict

A complete brand system shipped, modeled on typespec-asyncapi's generator pattern: ONE geometry ("the atomic swap" — emerald squircle, new bytes in →, old bytes out ←) produces every surface (favicon SVG+ICO, header/footer logo, Starlight docs header, OG card SVG+PNG, README lockup dark/light, standalone marks), and `brand:check` byte-gates all text artifacts in BOTH CI workflows. All website gates green. Honest gaps: no real-browser screenshot of the deployed pages (raster QA via ImageMagick only), the README `<picture>` is unverified on github.com until the push lands, and the two raster artifacts are not byte-gated.

## 1. What was wrong before

The old identity was an emerald squircle with the letter "A", hand-copied into three places that had already drifted: `favicon.svg` (dark `#0c0a09` glyph), `Logo.astro` (theme-flipping glyph via `var(--color-bg-primary)`), and `og-image.svg` (hand-scaled 72px copy). No lockups, no README logo, no docs-header logo, no palette record, no provenance.

## 2. DONE (verified)

| #  | Item                                                                                                                                                                                             | Evidence                                                                                                                                                                                      |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ~~1~~  | ~~Generator `scripts/generate-brand-assets.mjs` — one geometry → 9 artifacts, `--check` mode, cwd-independent paths, parent-dir creation~~ done (docs-health pass 2026-11-05) | ~~`node scripts/generate-brand-assets.mjs --check` → `OK (9 artifacts match the generator)`~~ |
| ~~2~~  | ~~`assets/brand/`: `mark.svg` (64×64), horizontal + stacked lockups in light/dark, `README.md` (compositions, files, palette, embedding, provenance, favicon-scale test)~~ done (docs-health pass 2026-11-05) | ~~files committed; `brand:check` covers all 5 SVGs~~ |
| ~~3~~  | ~~`Logo.astro` regenerated — flat token-driven fills (`var(--color-accent)` / `var(--color-on-accent)`), no `defs`, keeps Props compat so header/footer need no edits~~ done (docs-health pass 2026-11-05) | ~~both tokens exist in BOTH theme blocks (global.css:23,68); 2 instances in `dist/index.html` carry the new glyph; `{stroke:var(--color-on-accent)}` + both fill utilities present in built CSS~~ |
| ~~4~~  | ~~`favicon.svg` + `favicon.ico` (16/32/48) + `website/src/assets/brand.svg` regenerated~~ done (docs-health pass 2026-11-05) | ~~`magick icon:auto-resize` output identical structure to previous ICO~~ |
| ~~5~~  | ~~OG card: `og-image.svg` now generated (mark slot can never drift) → `og-image.png` 1200×630 re-rasterized~~ done (docs-health pass 2026-11-05) | ~~`identify`: 1200x630; visual QA render checked~~ |
| ~~6~~  | ~~README hero `<picture>` with `prefers-color-scheme`, absolute `raw.githubusercontent` URLs, 300×64~~ done (docs-health pass 2026-11-05) | ~~markup mirrors typespec-asyncapi's GitHub-proven pattern~~ |
| ~~7~~  | ~~Starlight docs-header logo — `logo: { src: "./src/assets/brand.svg" }`~~ done (docs-health pass 2026-11-05) | ~~built docs pages render `<img src="/_astro/brand.CqrkVa9N.svg">`~~ |
| ~~8~~  | ~~CI gates: `brand:check` as first drift-guard step (ci.yml, runs on every push) and first Website-workflow step~~ done (docs-health pass 2026-11-05) | ~~both YAMLs parse; step needs only system node~~ |
| ~~9~~  | ~~Website package scripts `brand:generate` / `brand:check` (no lockfile impact — scripts are not dependencies)~~ done (docs-health pass 2026-11-05) | ~~`pnpm run brand:check` green from website/~~ |
| ~~10~~ | ~~Gates: `astro check` 0 errors / 0 warnings; build 11 pages + CSP patch 11/11 + verify-build OK; `html-validate` silent; `brand:check` OK~~ done (docs-health pass 2026-11-05) | ~~full gate run 2026-11-03 23:41–42~~ |
| ~~11~~ | ~~Visual QA: mark at 64/28/22/16px, header variant (flat token) on dark `#0a0908` and light `#fafaf9`, both lockups, OG card — all eyeballed as PNG renders~~ done (docs-health pass 2026-11-05) | ~~swap arrows legible at 16px (stroke ≈1.6px)~~ |
| ~~12~~ | ~~Docs: CHANGELOG `[Unreleased]` entry; AGENTS.md (commands ×2, structure rows ×3, gotchas ×2 incl. corrected "OG image is editable" staleness, Starlight logo-src trap); `assets/brand/README.md`~~ done (docs-health pass 2026-11-05) | ~~in tree~~ |
| ~~13~~ | ~~dprint clean (BuildFlow `dprint-format` step made zero changes)~~ done (docs-health pass 2026-11-05) | ~~git diff after run: only intended AGENTS/CHANGELOG edits~~ |

## 3. Deliberate decisions

1. **Flat fill in `Logo.astro`, gradient everywhere else** — at 20–28px the gradient doesn't read; token-driven fill keeps the mark theme-consistent and avoids gradient-id collisions between the header and footer instances on one page. Documented in `assets/brand/README.md`.
2. **Muted `go-` + full-ink `atomic-write`** wordmark — ecosystem is context, promise is the product (echoes typespec-asyncapi's ink+accent split).
3. **Left-anchored horizontal lockup** with canvas sized for the wide DejaVu fallback (ImageMagick's render); Space Grotesk is narrower so browsers show a few px of trailing air — same tradeoff typespec-asyncapi accepted.
4. **Starlight logo via bundled `src/assets/brand.svg`**, not a public URL — `logo: { src: "/favicon.svg" }` fails the build with `ImageNotFound` because `logo.src` resolves through the Astro asset pipeline (hit live this session, now an AGENTS.md gotcha).
5. **`brand:check` in two workflows** — ci.yml's drift-guard runs on every push (covers root-level `assets/` + `scripts/` edits that the website path-filter misses); website.yml covers website-side drift.

## 4. NOT DONE / open

1. **Browser-rendered screenshot of `/` and a docs page** in both themes — QA was ImageMagick rasterization of the same geometry, not a live DOM render; do a 10-second eyeball after the next deploy.
2. **README `<picture>` on github.com** — unverified until pushed (npm is not applicable; this is a Go module).
3. **`favicon.ico` / `og-image.png` not byte-gated** — raster encoders aren't cross-platform stable; regenerate manually after SVG changes (commands in `assets/brand/README.md`).
4. **Horizontal lockup trailing air in browsers** — expected few px; revisit only if it looks unbalanced on the live README.
5. Possible future: motion variant (swap-arrow pulse) for the hero; `apple-touch-icon` PNG (180×180); GitHub social preview (1280×640) from the lockup.

## 5. What I would do differently next time

1. Check the consuming framework's asset-resolution rules BEFORE writing the config (the Starlight logo-src trap cost one build cycle; typespec-asyncapi already knew the answer).
2. The generator's `ENOENT` on first run (missing `assets/brand/`) was caught by actually running it before claiming done — keep that order, never assume mkdir.
3. `buildflow -s dprint` was a wrong step name (`dprint-format` is canonical; `buildflow --dry-run --verbose` lists names) — check the catalog first, per the buildflow skill.

## 6. Addendum (same session, ~23:20) — §4 gaps closed

1. **Browser screenshots DONE** (chromium 152 headless against `astro preview`, not just ImageMagick): landing header zoomed 3× in BOTH themes — swap mark crisp, glyph dark on emerald, correct token rendering; docs header verified two ways: zoomed screenshot + **pixel histogram** (`#0C0A09` ×136 px = the glyph ink, exactly `GLYPH_INK`; zero white pixels). Lesson: at 24–32 px against dark chrome, emerald antialiasing can _look_ like a white glyph — sample pixels before believing eyes.
2. **Lighthouse budgets PASS** (desktop preset, 4 URLs × 3 runs, assertions green incl. the error-gated accessibility ≥0.95 and SEO ≥0.95) with the new inline SVG logo, bundled docs logo, and regenerated favicon/OG in place.
3. **Go gate sanity**: `go build ./...` + `go test ./...` ok (brand work touched no Go source, confirmed empirically).
4. README `<picture>`-on-github.com remains the only §4 item left open (requires the push).
