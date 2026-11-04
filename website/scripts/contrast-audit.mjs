import { readFileSync } from "node:fs";

const css = readFileSync("src/styles/global.css", "utf8");

const srgbToLinear = (c) => {
	c /= 255;
	return c <= 0.04045 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4;
};
const hexToRgb = (hex) => {
	const h = hex.replace("#", "");
	const full =
		h.length === 3
			? h
					.split("")
					.map((c) => c + c)
					.join("")
			: h;
	return [0, 2, 4].map((i) => Number.parseInt(full.slice(i, i + 2), 16));
};
const parseColor = (str) => {
	const s = str.trim();
	if (s.startsWith("#")) {
		return hexToRgb(s);
	}
	const m = s.match(/rgba?\((?<vals>[^)]+)\)/u);
	if (!m) {
		return null;
	}
	const [r, g, b, a = 1] = m.groups.vals.split(",").map((x) => Number.parseFloat(x.trim()));
	return [r, g, b, a];
};
const blend = (fg, bg) => {
	const [fr, fgc, fb, fa = 1] = fg;
	const [br, bgc, bb] = bg;
	return [
		Math.round(fr * fa + br * (1 - fa)),
		Math.round(fgc * fa + bgc * (1 - fa)),
		Math.round(fb * fa + bb * (1 - fa)),
	];
};
const lum = (rgb) => {
	const [r, g, b] = rgb.map(srgbToLinear);
	return 0.2126 * r + 0.7152 * g + 0.0722 * b;
};
const ratio = (c1, c2) => {
	const [l1, l2] = [lum(c1), lum(c2)].toSorted((a, b) => b - a);
	return (l1 + 0.05) / (l2 + 0.05);
};

const parseBlock = (name, block) => {
	const tokens = {};
	for (const m of block.matchAll(/--color-([\w-]+):\s*(?<val>[^;]+);/gu)) {
		tokens[m[1]] = parseColor(m.groups.val);
	}
	return { name, tokens };
};

const themeStart = css.indexOf("@theme {");
const themeEnd = css.indexOf("}", css.indexOf("--animate", themeStart));
const lightStart = css.indexOf(":root.light");
const bodyStart = css.indexOf("body {");
const dark = parseBlock("dark", css.slice(themeStart, themeEnd));
const light = parseBlock("light", css.slice(lightStart, bodyStart));

const effective = (theme, name) => {
	const c = theme.tokens[name];
	if (!c) {
		return null;
	}
	return (c[3] ?? 1) < 1 ? blend(c, theme.tokens["bg-primary"]) : c.slice(0, 3);
};

const texts = [
	"text-primary",
	"text-secondary",
	"text-muted",
	"accent",
	"accent-light",
	"accent-hover",
	"amber",
	"danger",
	"success",
	"code-comment",
];
// success/danger render check/X icons (WCAG 1.4.11 non-text: 3:1), not body text.
const graphics = new Set(["danger", "success"]);
const threshold = (t) => (graphics.has(t) ? 3 : 4.5);
const bgs = ["bg-primary", "bg-card-solid", "bg-elevated"];

let failures = 0;
for (const theme of [dark, light]) {
	console.log(`\n=== ${theme.name} ===`);
	for (const t of texts) {
		const tc = theme.tokens[t];
		if (!tc) {
			continue;
		}
		const targets = t === "code-comment" ? ["bg-code", "code-inline-bg"] : bgs;
		for (const b of targets) {
			const bgc = effective(theme, b);
			if (!bgc) {
				continue;
			}
			const fg = (tc[3] ?? 1) < 1 ? blend(tc, bgc) : tc.slice(0, 3);
			const r = ratio(fg, bgc);
			const need = threshold(t);
			if (r >= need) {
				continue;
			}
			failures++;
			console.log(`FAIL ${t} on ${b}: ${r.toFixed(2)}:1 (needs ${need}:1)`);
		}
	}
	const r = ratio(theme.tokens["on-accent"].slice(0, 3), theme.tokens["accent"].slice(0, 3));
	if (r < 4.5) {
		failures++;
		console.log(`FAIL on-accent on accent(button): ${r.toFixed(2)}:1`);
	} else {
		console.log(`ok   on-accent on accent(button): ${r.toFixed(2)}:1`);
	}
}
console.log(`\n${failures} failing pair(s)`);
