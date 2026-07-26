// Post-build CSP patcher: injects a per-file hash-based Content-Security-Policy
// <meta> tag into every built HTML page in dist/.
//
// Hash-based = NO 'unsafe-inline' for scripts. Each inline <script> (except
// JSON-LD data blocks and external src scripts) gets a 'sha256-<b64>' entry.
// 'unsafe-inline' is kept for style-src (Tailwind/Astro inline critical CSS),
// which is the accepted tradeoff since style-based exfiltration is limited.
//
// Runs as `postbuild`, after `astro build` produces dist/.
import { readdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { createHash } from "node:crypto";

const here = dirname(fileURLToPath(import.meta.url));
const distDir = resolve(here, "../dist");

// Matches inline <script> blocks: must have content and NOT be external (src=)
// or JSON-LD data (type="application/ld+json"). Captures the inner script text.
const inlineScript = /<script(?![^>]*\bsrc=)(?![^>]*type="application\/ld\+json")([^>]*)>([\s\S]*?)<\/script>/gi;

function walkHtml(dir) {
  const out = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const full = join(dir, entry.name);
    if (entry.isDirectory()) {
      out.push(...walkHtml(full));
    } else if (entry.name.endsWith(".html")) {
      out.push(full);
    }
  }
  return out;
}

function sha256B64(text) {
  return "sha256-" + createHash("sha256").update(text, "utf8").digest("base64");
}

function buildCsp(hashes) {
  const scriptSrc = ["'self'", ...hashes].join(" ");
  return [
    "default-src 'self'",
    "base-uri 'self'",
    "object-src 'none'",
    "frame-ancestors 'none'",
    "form-action 'self'",
    "connect-src 'self'",
    "img-src 'self' data: https:",
    "font-src 'self' data:",
    "style-src 'self' 'unsafe-inline'",
    `script-src ${scriptSrc}`,
  ].join("; ");
}

let files = 0;
let totalHashes = 0;

for (const file of walkHtml(distDir)) {
  let html = readFileSync(file, "utf8");
  const hashes = new Set();
  let match;
  inlineScript.lastIndex = 0;
  while ((match = inlineScript.exec(html)) !== null) {
    const body = match[2];
    if (body.trim().length === 0) continue;
    hashes.add(sha256B64(body));
  }

  // Remove any CSP meta we injected on a previous run (idempotent).
  html = html.replace(
    /<meta http-equiv="Content-Security-Policy"[^>]*>\s*/gi,
    "",
  );

  const csp = buildCsp([...hashes].sort());
  const meta = `<meta http-equiv="Content-Security-Policy" content="${csp}">`;

  // Inject immediately after <head> so it is parsed before other resources.
  if (/<head[^>]*>/i.test(html)) {
    html = html.replace(/(<head[^>]*>)/i, `$1${meta}`);
  } else {
    html = meta + html;
  }

  writeFileSync(file, html, "utf8");
  files++;
  totalHashes += hashes.size;
}

console.log(`[fix-csp] patched ${files} HTML file(s), ${totalHashes} inline script hash(es) total`);
