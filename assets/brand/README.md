# Brand assets — the atomic swap

The project's visual identity: an emerald squircle holding two opposing
arrows — **new bytes flow in** (top, pointing right), **old bytes flow out**
(bottom, pointing left). One indivisible replace operation, which is exactly
the guarantee the library makes.

```
old file  →  [ ⟲ atomic swap ]  →  new file
```

Everything here (plus the favicon, the Logo component and the OG card) is
generated from a single geometry definition in
[`scripts/generate-brand-assets.mjs`](../../scripts/generate-brand-assets.mjs).
Do not hand-edit the SVGs or `website/src/components/Logo.astro`; edit the
generator and regenerate.

## Compositions

One geometry, three compositions — pick by slot shape, not by taste:

| Composition         | Size      | Use for                                                          |
| ------------------- | --------- | ---------------------------------------------------------------- |
| `mark.svg`          | 64 × 64   | Standalone emblem: social avatars, OG art source, favicon source |
| `lockup-horizontal` | 300 × 64  | Wide slots: README hero, banners, social headers                 |
| `lockup-stacked`    | 340 × 117 | Roughly square slots: slides, gallery cards                      |

The two lockups are the same atoms in the same order: the horizontal one runs
left-to-right, the stacked one top-to-bottom. Both share mark size, gap and
wordmark, so an edit to one is an edit to both.

## Files

| File                                      | What it is                                  | Consumed by                                   |
| ----------------------------------------- | ------------------------------------------- | --------------------------------------------- |
| `mark.svg`                                | 64 × 64 emblem (gradient squircle + glyph)  | favicon work, avatar/OG source                |
| `lockup-horizontal-light.svg`             | 300 × 64, dark ink, for light backgrounds   | `README.md` (fallback)                        |
| `lockup-horizontal-dark.svg`              | 300 × 64, light ink, for dark backgrounds   | `README.md` (`prefers-color-scheme: dark`)    |
| `lockup-stacked-light.svg`                | 340 × 117, dark ink                         | — (kept parallel to the dark variant)         |
| `lockup-stacked-dark.svg`                 | 340 × 117, light ink                        | — (kept parallel to the light variant)        |
| `../../website/src/components/Logo.astro` | Compact emblem as an Astro component        | Site header (28 px), site footer (20 px)      |
| `../../website/public/favicon.svg`        | Favicon flavor of the mark                  | Browsers, Starlight docs header, PWA manifest |
| `../../website/public/favicon.ico`        | 16/32/48 raster fallback                    | Legacy favicon `<link rel="alternate icon">`  |
| `../../website/public/og-image.svg`       | 1200 × 630 social card embedding the emblem | Rasterized to `og-image.png`                  |
| `../../website/public/og-image.png`       | 1200 × 630 PNG social card                  | Open Graph / Twitter meta tags                |

## Regenerate

```bash
cd website
pnpm run brand:generate   # rewrite all eight text artifacts
pnpm run brand:check      # exit 1 if any text artifact drifted
```

`brand:check` compares exact bytes, so it is deterministic across machines. It
runs as the first step of the CI drift-guard job and in the Website workflow,
so hand-edits to any generated artifact fail CI. The binary artifacts
(`favicon.ico`, `og-image.png`, `icon-*.png`, `apple-touch-icon.png`) are NOT
gated by `brand:check` (raster encoders are not byte-stable across platforms);
regenerate them from the generated SVGs with ImageMagick 7 in one command:

```bash
cd website && pnpm run brand:png
```

That codifies (and replaces) the previous manual pipeline: favicon.ico at
16/32/48, the OG card at 1200×630, the 180/192/512 app icons on the site
background color, and the apple-touch copy — each PNG 256-color quantized
(the OG card lands at ~28 KB instead of 62 KB raw, with no visible banding
on the flat-design gradients and exact text areas).

## Palette

The palette lives in [`brand-tokens.json`](./brand-tokens.json) — the single
source the generator reads. Four of its values intentionally mirror dark-theme
tokens in `website/src/styles/global.css` (`--color-accent`, `--color-on-accent`,
`--color-text-primary`, `--color-text-secondary`); `brand:check` verifies that
mirror on every run, so the brand and the site cannot drift apart. The
light-theme site tokens are deliberately contrast-tuned and are NOT mirrored
here.

| Token            | Hex                   | Used for                                                           |
| ---------------- | --------------------- | ------------------------------------------------------------------ |
| Accent gradient  | `#10b981` → `#059669` | Squircle (135°), matches the site's emerald accent and the OG card |
| Glyph ink        | `#0c0a09`             | Swap arrows — the site's `--color-on-accent`                       |
| Ink (light bg)   | `#0c0a09`             | `atomic-write` wordmark in every light variant                     |
| Ink (dark bg)    | `#f5f5f4`             | `atomic-write` wordmark in every dark variant                      |
| Muted (light bg) | `#78716c`             | `go-` prefix, light variants                                       |
| Muted (dark bg)  | `#a8a29e`             | `go-` prefix, dark variants                                        |

Three deliberate decisions:

- The `go-` prefix is muted and `atomic-write` is full ink: the ecosystem is
  context, the promise is the product.
- `Logo.astro` uses a **flat accent fill instead of the gradient**, driven by
  site tokens (`var(--color-accent)` / `var(--color-on-accent)`). At 20–28 px
  the gradient does not read, the token-driven fill keeps the mark
  theme-consistent with the rest of the site, and no `defs` means no gradient-id
  collisions between the header and footer instances on one page.
- The horizontal lockup is laid out left-anchored with the canvas sized for the
  wide DejaVu fallback (the render target of ImageMagick on Linux). Space
  Grotesk — the site's actual font — is narrower, so browsers show a measured
  10.1 px of trailing air (wordmark 215.9 px at 30 px/700, canvas 300 − 74 px
  mark slot = 226 px available; measured via headless Chromium against the
  bundled woff2, 2026-11-04). That is ~3% of the canvas: invisible as a layout
  flaw, and resizing the canvas would clip the ImageMagick render, so the air
  is the correct tradeoff. SVG text rendered through `<img>` cannot fetch
  web fonts, so the SVGs declare the `'Space Grotesk', 'Helvetica Neue', Arial,
  sans-serif` fallback stack.

## Favicon scale

The swap glyph keeps its legibility down to 16 px (stroke 6.5/64 ≈ 1.6 px at
16 px, head 2.9 px tall) — verified visually against the dark browser-chrome
background at 16 px and 22 px. Do not redesign the glyph without re-running
that check:

```bash
magick -background none assets/brand/mark.svg -resize 16x16! /tmp/mark-16.png
```

## Embedding

GitHub renders `<picture>` with a `prefers-color-scheme` source, so the README
is theme-aware. There is no npm package here (this is a Go module), but the
README still uses absolute `raw.githubusercontent.com` URLs so the markup
survives being mirrored outside the repo:

```html
<p align="center">
  <picture>
    <source
      media="(prefers-color-scheme: dark)"
      srcset="https://raw.githubusercontent.com/LarsArtmann/go-atomic-write/master/assets/brand/lockup-horizontal-dark.svg"
    />
    <img
      src="https://raw.githubusercontent.com/LarsArtmann/go-atomic-write/master/assets/brand/lockup-horizontal-light.svg"
      alt="go-atomic-write: the atomic swap"
      width="300"
      height="64"
    />
  </picture>
</p>
```

## Provenance

Original artwork, created for this project (no third-party marks involved).
Same license as the repository (MIT).
