/**
 * Brand asset generator: ONE geometry definition -> every brand artifact.
 *
 * The identity is "the atomic swap": an emerald squircle holding two opposing
 * arrows - new bytes flow in (top, pointing right), old bytes flow out
 * (bottom, pointing left) - one indivisible replace operation, which is
 * exactly what the library guarantees.
 *
 * Before this script the same squircle + glyph was copy-pasted into the
 * favicon, the Logo component and the OG image, and they had already drifted
 * (the OG mark was a hand-scaled copy of the favicon path). Now every artifact
 * is assembled from the constants below by exactly one function.
 *
 *   node scripts/generate-brand-assets.mjs           rewrite the artifacts
 *   node scripts/generate-brand-assets.mjs --check   exit 1 if any artifact drifted
 *   node scripts/generate-brand-assets.mjs --png     rasterize favicon.ico, og-image.png, app icons (needs ImageMagick 7)
 *
 * Never hand-edit the generated files; edit this script and regenerate.
 * Run from anywhere: paths resolve relative to this file, not the cwd.
 */
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { execFileSync } from "node:child_process";
import path from "node:path";

const repoRoot = path.join(import.meta.dirname, "..");

// ---------------------------------------------------------------------------
// Geometry (64x64 design space)
// ---------------------------------------------------------------------------

const CANVAS = 64;
const SQUIRCLE_RADIUS = 14;

// Two swap arrows: top shaft flows right (new bytes in), bottom shaft flows
// left (old bytes out). Heads are filled triangles; shafts are round-capped
// strokes. The bbox is x 15.5-48.5, y 16.5-47.5 - dead center of the canvas.
const GLYPH_STROKE_WIDTH = 6.5;
const GLYPH_INK = "#0c0a09";
const SWAP_ARROWS = [
	{ shaft: "M17 24H39", head: "37 16.5 48.5 24 37 31.5" },
	{ shaft: "M47 40H25", head: "27 32.5 15.5 40 27 47.5" },
];

// ---------------------------------------------------------------------------
// Palette (mirrors website/src/styles/global.css tokens)
// ---------------------------------------------------------------------------

const ACCENT = "#10b981";
const ACCENT_DEEP = "#059669";
const INK_ON_LIGHT = "#0c0a09";
const INK_ON_DARK = "#f5f5f4";
const MUTED_ON_LIGHT = "#78716c";
const MUTED_ON_DARK = "#a8a29e";

const PALETTES = {
	light: { ink: INK_ON_LIGHT, muted: MUTED_ON_LIGHT },
	dark: { ink: INK_ON_DARK, muted: MUTED_ON_DARK },
};

const STATIC_SANS = "'Space Grotesk', 'Helvetica Neue', Arial, sans-serif";
const GENERATED_NOTICE = `<!--
  GENERATED FILE - do not edit by hand.
  Source: scripts/generate-brand-assets.mjs (run: node scripts/generate-brand-assets.mjs)
-->
`;

// ---------------------------------------------------------------------------
// Assembly
// ---------------------------------------------------------------------------

const swapGlyph = (ink) =>
	SWAP_ARROWS.map(
		(arrow) =>
			`<path d="${arrow.shaft}" stroke="${ink}" stroke-width="${GLYPH_STROKE_WIDTH}" stroke-linecap="round" fill="none" />\n  <polygon points="${arrow.head}" fill="${ink}" />`,
	).join("\n  ");

const squircleGradient = (id) =>
	`<linearGradient id="${id}" x1="0" y1="0" x2="1" y2="1">\n    <stop offset="0" stop-color="${ACCENT}" />\n    <stop offset="1" stop-color="${ACCENT_DEEP}" />\n  </linearGradient>`;

/** The standalone 64x64 emblem: emerald squircle + swap glyph. */
const markSvg = () => `\
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${CANVAS} ${CANVAS}" width="${CANVAS}" height="${CANVAS}" role="img" aria-labelledby="gawMarkTitle">
<title id="gawMarkTitle">go-atomic-write: the atomic swap</title>
<defs>
  ${squircleGradient("gawMarkEm")}
</defs>
<rect width="${CANVAS}" height="${CANVAS}" rx="${SQUIRCLE_RADIUS}" fill="url(#gawMarkEm)" />
  ${swapGlyph(GLYPH_INK)}
</svg>
`;

/** Favicon flavor of the mark: no a11y attributes, plain id namespace. */
const faviconSvg = () => `\
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${CANVAS} ${CANVAS}" width="${CANVAS}" height="${CANVAS}" fill="none">
  <defs>
    ${squircleGradient("gawFaviconEm")}
  </defs>
  <rect width="${CANVAS}" height="${CANVAS}" rx="${SQUIRCLE_RADIUS}" fill="url(#gawFaviconEm)" />
  ${swapGlyph(GLYPH_INK)}
</svg>
`;

/**
 * Emblem at `size` px, top-left corner at (offsetX, offsetY). Used by the
 * lockups and the OG card; `gradientId` must be unique per embedding document.
 */
const markGroup = ({ offsetX, offsetY, size, gradientId }) => {
	const scale = size / CANVAS;
	return `<g transform="translate(${offsetX} ${offsetY}) scale(${scale})">
  <rect width="${CANVAS}" height="${CANVAS}" rx="${SQUIRCLE_RADIUS}" fill="url(#${gradientId})" />
  ${swapGlyph(GLYPH_INK)}
</g>`;
};

const wordmark = ({ x, y, size, palette, anchor = "start" }) =>
	`<text x="${x}" y="${y}" font-family="${STATIC_SANS}" font-size="${size}" font-weight="700" letter-spacing="-0.5" text-anchor="${anchor}"><tspan fill="${palette.muted}">go-</tspan><tspan fill="${palette.ink}">atomic-write</tspan></text>`;

/** Emblem left, wordmark right. One text line - for READMEs and banners. */
const lockupHorizontal = (width, variant) => {
	const palette = PALETTES[variant];
	const id = `gawLh${variant}`;
	return `\
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${width} 64" width="${width}" height="64" role="img" aria-labelledby="${id}Title">
<title id="${id}Title">go-atomic-write: the atomic swap</title>
<defs>
  ${squircleGradient(`${id}Em`)}
</defs>
${markGroup({ offsetX: 4, offsetY: 4, size: 56, gradientId: `${id}Em` })}
${wordmark({ x: 74, y: 42.5, size: 30, palette })}
</svg>
`;
};

/** Emblem centered above the centered wordmark. For square-ish slots. */
const lockupStacked = ({ width, height, baselineY, variant }) => {
	const palette = PALETTES[variant];
	const id = `gawLs${variant}`;
	const centerX = width / 2;
	return `\
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${width} ${height}" width="${width}" height="${height}" role="img" aria-labelledby="${id}Title">
<title id="${id}Title">go-atomic-write: the atomic swap</title>
<defs>
  ${squircleGradient(`${id}Em`)}
</defs>
${markGroup({ offsetX: centerX - 32, offsetY: 4, size: 64, gradientId: `${id}Em` })}
${wordmark({ x: centerX, y: baselineY, size: 30, palette, anchor: "middle" })}
</svg>
`;
};

/**
 * The Astro component keeps the small-size simplification: a flat accent fill
 * (no gradient) driven by site tokens, so the mark stays theme-aware in the
 * header/footer and avoids gradient-id collisions between page instances.
 */
const logoAstro = `\
---
// GENERATED FILE - do not edit by hand.
// Source: scripts/generate-brand-assets.mjs (run: node scripts/generate-brand-assets.mjs)
export interface Props {
  width?: number;
  height?: number;
  class?: string;
}

const { width = 28, height = 28, class: className } = Astro.props;
---

<svg class={className} width={width} height={height} viewBox="0 0 ${CANVAS} ${CANVAS}" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
  <rect width="${CANVAS}" height="${CANVAS}" rx="${SQUIRCLE_RADIUS}" class="fill-[var(--color-accent)]" />
  <path d="${SWAP_ARROWS[0].shaft}" stroke-width="${GLYPH_STROKE_WIDTH}" stroke-linecap="round" class="stroke-[var(--color-on-accent)]" fill="none" />
  <polygon points="${SWAP_ARROWS[0].head}" class="fill-[var(--color-on-accent)]" />
  <path d="${SWAP_ARROWS[1].shaft}" stroke-width="${GLYPH_STROKE_WIDTH}" stroke-linecap="round" class="stroke-[var(--color-on-accent)]" fill="none" />
  <polygon points="${SWAP_ARROWS[1].head}" class="fill-[var(--color-on-accent)]" />
</svg>
`;

/**
 * The OG card is static (rasterized by ImageMagick, no runtime rendering), so
 * the whole composition is templated here and the mark slot can never drift
 * from the emblem again. Ink colors and layout are carried over verbatim from
 * the previous hand-authored card.
 */
const ogImageSvg = `\
<svg width="1200" height="630" viewBox="0 0 1200 630" fill="none" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <pattern id="dots" width="28" height="28" patternUnits="userSpaceOnUse">
      <circle cx="1" cy="1" r="1" fill="rgba(231,229,228,0.03)"/>
    </pattern>
    <linearGradient id="emerald" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0%" stop-color="${ACCENT}"/>
      <stop offset="100%" stop-color="${ACCENT_DEEP}"/>
    </linearGradient>
    <linearGradient id="glow" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="rgba(16,185,129,0.12)"/>
      <stop offset="100%" stop-color="rgba(16,185,129,0)"/>
    </linearGradient>
  </defs>

  <!-- Background -->
  <rect width="1200" height="630" fill="#0a0908"/>
  <rect width="1200" height="630" fill="url(#dots)"/>

  <!-- Top glow -->
  <rect x="0" y="0" width="1200" height="315" fill="url(#glow)"/>

  <!-- Logo mark -->
  <g transform="translate(100, 90) scale(1.125)">
    <rect width="64" height="64" rx="14" fill="url(#emerald)"/>
    ${swapGlyph(GLYPH_INK)}
  </g>

  <!-- Title -->
  <text x="100" y="260" font-family="'Noto Sans','DejaVu Sans',sans-serif" font-size="72" font-weight="700" fill="#f5f5f4" letter-spacing="-2">go-atomic-write</text>

  <!-- Tagline -->
  <text x="100" y="320" font-family="'Noto Sans','DejaVu Sans',sans-serif" font-size="32" font-weight="400" fill="#a8a29e">Crash-safe, race-free file writes for Go</text>

  <!-- Feature badges -->
  <g transform="translate(100, 400)" font-family="'Noto Sans','DejaVu Sans',sans-serif" font-size="24" font-weight="500">
    <g>
      <rect width="200" height="48" rx="24" fill="rgba(16,185,129,0.1)" stroke="rgba(16,185,129,0.3)" stroke-width="1"/>
      <text x="100" y="32" text-anchor="middle" fill="#6ee7b7">No lost updates</text>
    </g>
    <g transform="translate(220, 0)">
      <rect width="200" height="48" rx="24" fill="rgba(16,185,129,0.1)" stroke="rgba(16,185,129,0.3)" stroke-width="1"/>
      <text x="100" y="32" text-anchor="middle" fill="#6ee7b7">Crash-durable</text>
    </g>
    <g transform="translate(440, 0)">
      <rect width="240" height="48" rx="24" fill="rgba(16,185,129,0.1)" stroke="rgba(16,185,129,0.3)" stroke-width="1"/>
      <text x="120" y="32" text-anchor="middle" fill="#6ee7b7">Concurrent-write safe</text>
    </g>
  </g>

  <!-- Bottom URL -->
  <text x="100" y="540" font-family="'JetBrains Mono','DejaVu Sans Mono',monospace" font-size="22" font-weight="400" fill="#8b827a">pkg.go.dev/github.com/larsartmann/go-atomic-write</text>

  <!-- Accent line -->
  <rect x="0" y="624" width="1200" height="6" fill="url(#emerald)"/>
</svg>
`;

// ---------------------------------------------------------------------------
// Artifact table + write/check
// ---------------------------------------------------------------------------

const LOCKUP_HORIZONTAL_WIDTH = 300;
const LOCKUP_STACKED_WIDTH = 340;
const LOCKUP_STACKED_HEIGHT = 117;
const LOCKUP_STACKED_BASELINE = 106;

const ARTIFACTS = [
	{ path: "assets/brand/mark.svg", content: () => GENERATED_NOTICE + markSvg() },
	{
		path: "assets/brand/lockup-horizontal-light.svg",
		content: () => GENERATED_NOTICE + lockupHorizontal(LOCKUP_HORIZONTAL_WIDTH, "light"),
	},
	{
		path: "assets/brand/lockup-horizontal-dark.svg",
		content: () => GENERATED_NOTICE + lockupHorizontal(LOCKUP_HORIZONTAL_WIDTH, "dark"),
	},
	{
		path: "assets/brand/lockup-stacked-light.svg",
		content: () =>
			GENERATED_NOTICE +
			lockupStacked({
				width: LOCKUP_STACKED_WIDTH,
				height: LOCKUP_STACKED_HEIGHT,
				baselineY: LOCKUP_STACKED_BASELINE,
				variant: "light",
			}),
	},
	{
		path: "assets/brand/lockup-stacked-dark.svg",
		content: () =>
			GENERATED_NOTICE +
			lockupStacked({
				width: LOCKUP_STACKED_WIDTH,
				height: LOCKUP_STACKED_HEIGHT,
				baselineY: LOCKUP_STACKED_BASELINE,
				variant: "dark",
			}),
	},
	{ path: "website/src/components/Logo.astro", content: () => logoAstro },
	{ path: "website/public/favicon.svg", content: () => GENERATED_NOTICE + faviconSvg() },
	{ path: "website/src/assets/brand.svg", content: () => faviconSvg() },
	{ path: "website/public/og-image.svg", content: () => GENERATED_NOTICE + ogImageSvg },
];

const checkMode = process.argv.includes("--check");
const pngMode = process.argv.includes("--png");

const svgTagPattern =
	/<(?<close>\/?)(?<name>[a-zA-Z][a-zA-Z0-9-]*)(?<attrs>(?:"[^"]*"|'[^']*'|[^>])*?)(?<self>\/?)>/gu;

function scanSvgTags(svg) {
	const stack = [];
	for (const tag of svg.matchAll(svgTagPattern)) {
		if (tag.groups.close) {
			const open = stack.pop();
			if (open !== tag.groups.name) {
				return `mismatched closing tag </${tag.groups.name}> (open: ${open ?? "none"})`;
			}
		} else if (!tag.groups.self) {
			stack.push(tag.groups.name);
		}
	}
	return stack.length > 0 ? `unclosed tag(s): ${stack.join(", ")}` : null;
}

// Well-formedness gate for the generated SVGs: tags balance, no stray markup
// outside tags. Not a full XML parser - the generator owns the input style -
// but it catches hand-edits and template bugs before browsers or ImageMagick
// render garbage (magick fails silently on some malformed SVGs).
function svgWellFormed(svg) {
	const problem = scanSvgTags(svg);
	if (problem) {
		return problem;
	}
	const commentFree = svg.replaceAll(/<!--[\s\S]*?-->/gu, "");
	const tagVolume = [...commentFree.matchAll(svgTagPattern)].reduce((sum, tag) => sum + tag[0].length, 0);
	return commentFree.trim().length > tagVolume * 2 ? "unaccounted markup outside tags" : null;
}

function magick(args) {
	try {
		execFileSync("magick", args, { stdio: "pipe" });
	} catch (error) {
		console.error(`brand:png: magick failed (${args.join(" ")}): ${error.message}`);
		console.error("ImageMagick 7 is required (nix shell nixpkgs#imagemagick provides it).");
		process.exit(1);
	}
}

const drifted = [];
for (const artifact of ARTIFACTS) {
	const target = path.join(repoRoot, artifact.path);
	const expected = artifact.content();
	let current;
	try {
		current = readFileSync(target, "utf8");
	} catch {
		drifted.push(`${artifact.path} (missing)`);
		continue;
	}
	if (current !== expected) {
		drifted.push(artifact.path);
		continue;
	}
	if (artifact.path.endsWith(".svg")) {
		const problem = svgWellFormed(current);
		if (problem) {
			drifted.push(`${artifact.path} (${problem})`);
		}
	}
}

if (pngMode) {
	const site = (rel) => path.join(repoRoot, "website", rel);
	const quantize = (file) => magick([file, "-colors", "256", "-define", "png:compression-level=9", file]);

	magick(["-background", "none", site("public/favicon.svg"), "-define", "icon:auto-resize=16,32,48", site("public/favicon.ico")]);
	console.log("wrote website/public/favicon.ico");

	magick(["-background", "none", site("public/og-image.svg"), "-resize", "1200x630!", site("public/og-image.png")]);
	quantize(site("public/og-image.png"));
	console.log("wrote website/public/og-image.png (256-color quantized)");

	for (const size of [180, 192, 512]) {
		const icon = site(`public/icon-${size}.png`);
		magick(["-background", "#0c0a09", path.join(repoRoot, "assets/brand/mark.svg"), "-resize", `${size}x${size}!`, icon]);
		quantize(icon);
		console.log(`wrote website/public/icon-${size}.png`);
	}
	execFileSync("cp", [site("public/icon-180.png"), site("public/apple-touch-icon.png")]);
	console.log("wrote website/public/apple-touch-icon.png");
	console.log("brand:png: rasters regenerated from the generated SVGs.");
} else if (checkMode) {
	if (drifted.length > 0) {
		console.error(
			`brand:check failed - ${drifted.length} artifact(s) drifted from scripts/generate-brand-assets.mjs:`,
		);
		for (const artifactPath of drifted) {
			console.error(`  - ${artifactPath}`);
		}
		console.error(
			"Run `node scripts/generate-brand-assets.mjs` to regenerate, or `pnpm run brand:generate` in website/.",
		);
		process.exit(1);
	}
	console.log(`brand:check: OK (${ARTIFACTS.length} artifacts match the generator)`);
} else {
	for (const artifact of ARTIFACTS) {
		const target = path.join(repoRoot, artifact.path);
		mkdirSync(path.dirname(target), { recursive: true });
		writeFileSync(target, artifact.content());
		console.log(`wrote ${artifact.path}`);
	}
	console.log(`brand:generate: ${ARTIFACTS.length} artifacts written from one geometry.`);
}
