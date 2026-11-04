import { readFileSync } from "node:fs";

const css = readFileSync("src/styles/global.css", "utf8");

const hexToRgb = (hex) => {
  const h = hex.replace("#", "");
  const full = h.length === 3 ? h.split("").map((c) => c + c).join("") : h;
  return [0, 2, 4].map((i) => parseInt(full.slice(i, i + 2), 16));
};
const parseColor = (str) => {
  str = str.trim();
  if (str.startsWith("#")) return hexToRgb(str);
  const m = str.match(/rgba?\(([^)]+)\)/);
  if (!m) return null;
  const [r, g, b, a = 1] = m[1].split(",").map((s) => parseFloat(s.trim()));
  return [r, g, b, a];
};
const blend = (fg, bg) => {
  const [fr, fg2, fb, fa = 1] = fg;
  const [br, bg2, bb] = bg;
  return [
    Math.round(fr * fa + br * (1 - fa)),
    Math.round(fg2 * fa + bg2 * (1 - fa)),
    Math.round(fb * fa + bb * (1 - fa)),
  ];
};
const lum = ([r, g, b]) => {
  const ch = (c) => {
    c /= 255;
    return c <= 0.04045 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
  };
  return 0.2126 * ch(r) + 0.7152 * ch(g) + 0.0722 * ch(b);
};
const ratio = (c1, c2) => {
  const [l1, l2] = [lum(c1), lum(c2)].sort((a, b) => b - a);
  return (l1 + 0.05) / (l2 + 0.05);
};

const parseBlock = (name, block) => {
  const tokens = {};
  for (const m of block.matchAll(/--color-([\w-]+):\s*([^;]+);/g)) {
    tokens[m[1]] = parseColor(m[2]);
  }
  return { name, tokens };
};

const rootBlock = css.slice(css.indexOf(":root {"), css.indexOf(":root.light"));
const lightBlock = css.slice(css.indexOf(":root.light"), css.indexOf("body {"));
const dark = parseBlock("dark", rootBlock);
const light = parseBlock("light", lightBlock);

// effective solid backgrounds (blend alpha bgs onto bg-primary)
const eff = (theme, name) => {
  const c = theme.tokens[name];
  if (!c) return null;
  return (c[3] ?? 1) < 1 ? blend(c, theme.tokens["bg-primary"]) : c.slice(0, 3);
};

const texts = [
  "text-primary", "text-secondary", "text-muted",
  "accent", "accent-light", "accent-hover",
  "amber", "danger", "success", "code-comment", "on-accent",
];
const bgs = ["bg-primary", "bg-card-solid", "bg-elevated"];

let failures = 0;
for (const theme of [dark, light]) {
  console.log(`\n=== ${theme.name} ===`);
  for (const t of texts) {
    const tc = theme.tokens[t];
    if (!tc) continue;
    const targets = t === "code-comment" ? ["bg-code", "code-inline-bg"] : bgs;
    for (const b of targets) {
      let bgc = eff(theme, b);
      if (!bgc) continue;
      const fg = (tc[3] ?? 1) < 1 ? blend(tc, bgc) : tc.slice(0, 3);
      const r = ratio(fg, bgc);
      const pass = r >= 4.5;
      const largePass = r >= 3;
      if (!pass) {
        failures++;
        console.log(
          `FAIL ${t} on ${b}: ${r.toFixed(2)}:1 ${largePass ? "(AA-large ok)" : "(fails even AA-large)"}`,
        );
      }
    }
  }
  // on-accent: text on accent button
  const oa = theme.tokens["on-accent"];
  const ac = theme.tokens["accent"];
  const r = ratio(oa.slice(0, 3), ac.slice(0, 3));
  if (r < 4.5) { failures++; console.log(`FAIL on-accent on accent(button): ${r.toFixed(2)}:1`); }
  else console.log(`ok   on-accent on accent(button): ${r.toFixed(2)}:1`);
}
console.log(`\n${failures} failing pair(s)`);
