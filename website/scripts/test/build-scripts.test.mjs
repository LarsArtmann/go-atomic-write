// Regression tests for the two load-bearing postbuild scripts:
//   - fix-csp.mjs  (idempotency + the 2026-11-04 quoted-hash incident class)
//   - sync-changelog.mjs (<-escaping so MDX cannot break the build)
// Run: pnpm run test:scripts  (node --test)
import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { mkdtempSync, readFileSync, rmSync, writeFileSync, mkdirSync } from "node:fs";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import test from "node:test";

const scriptsDir = resolve(import.meta.dirname, "..");

function runScript(script, env) {
  return spawnSync(process.execPath, [join(scriptsDir, script)], {
    env: { ...process.env, ...env },
    encoding: "utf8",
  });
}

function makeDist() {
  const dist = mkdtempSync(join(tmpdir(), "csp-test-"));
  writeFileSync(
    join(dist, "index.html"),
    `<!doctype html><html><head><title>t</title></head><body>
<script>console.log("inline one")</script>
<script src="/js/external.js"></script>
<script type="application/ld+json">{"@type":"Thing"}</script>
<script>console.log("inline two")</script>
</body></html>`,
    "utf8",
  );
  return dist;
}

test("fix-csp hashes only inline scripts, quotes every hash source, allows buttondown", () => {
  const dist = makeDist();
  try {
    const run = runScript("fix-csp.mjs", { CSP_DIST_DIR: dist });
    assert.equal(run.status, 0, run.stderr);

    const html = readFileSync(join(dist, "index.html"), "utf8");
    const metas = html.match(/<meta http-equiv="Content-Security-Policy"[^>]*>/giu) ?? [];
    assert.equal(metas.length, 1, "exactly one injected CSP meta");

    const scriptSrc = metas[0].match(/script-src ([^;"]+)/iu)?.[1] ?? "";
    const hashes = scriptSrc.split(" ").filter((source) => source.includes("sha256-"));
    assert.equal(hashes.length, 2, "both inline scripts hashed, external + ld+json skipped");
    for (const hash of hashes) {
      assert.match(hash, /^'sha256-[A-Za-z0-9+/=]{43}='$/u, `hash must be quoted: ${hash}`);
    }
    assert.equal((scriptSrc.match(/(^|\s)sha256-/gu) ?? []).length, 0, "no bare hash sources");
    assert.match(metas[0], /form-action 'self' https:\/\/buttondown\.email/u);
  } finally {
    rmSync(dist, { recursive: true, force: true });
  }
});

test("fix-csp is idempotent: a second run is a byte-for-byte no-op with one meta", () => {
  const dist = makeDist();
  try {
    assert.equal(runScript("fix-csp.mjs", { CSP_DIST_DIR: dist }).status, 0);
    const afterFirst = readFileSync(join(dist, "index.html"), "utf8");
    assert.equal(runScript("fix-csp.mjs", { CSP_DIST_DIR: dist }).status, 0);
    const afterSecond = readFileSync(join(dist, "index.html"), "utf8");
    assert.equal(afterSecond, afterFirst, "second run must not change the file");
    const metaCount = (afterSecond.match(/Content-Security-Policy/giu) ?? []).length;
    assert.equal(metaCount, 1, "exactly one CSP meta after double patching");
  } finally {
    rmSync(dist, { recursive: true, force: true });
  }
});

test("sync-changelog escapes < for MDX, strips the H1, and is idempotent", () => {
  const dir = mkdtempSync(join(tmpdir(), "changelog-test-"));
  const source = join(dir, "CHANGELOG.md");
  const target = join(dir, "changelog.mdx");
  try {
    writeFileSync(
      source,
      "# Changelog\n\n## [Unreleased]\n\n- generic `a < b` comparison and `<div>` markup survive\n",
      "utf8",
    );
    assert.equal(runScript("sync-changelog.mjs", { CHANGELOG_SOURCE: source, CHANGELOG_TARGET: target }).status, 0);
    const first = readFileSync(target, "utf8");

    assert.equal(first.match(/^#\s+Changelog/mu), null, "leading H1 stripped");
    assert.match(first, /&lt;div&gt;/u, "raw <div> escaped");
    assert.match(first, /`a &lt; b`/u, "inline-code < escaped");
    assert.equal((first.match(/</gu) ?? []).length, (first.match(/&lt;/gu) ?? []).length, "no unescaped < outside char refs");

    assert.equal(runScript("sync-changelog.mjs", { CHANGELOG_SOURCE: source, CHANGELOG_TARGET: target }).status, 0);
    assert.equal(readFileSync(target, "utf8"), first, "second run is byte-identical");
  } finally {
    rmSync(dir, { recursive: true, force: true });
  }
});

test("sync-changelog keeps a page without a leading H1 intact", () => {
  const dir = mkdtempSync(join(tmpdir(), "changelog-test-"));
  const source = join(dir, "CHANGELOG.md");
  const target = join(dir, "changelog.mdx");
  try {
    writeFileSync(source, "## [0.1.0] - 2026-06-03\n\n- entry\n", "utf8");
    assert.equal(runScript("sync-changelog.mjs", { CHANGELOG_SOURCE: source, CHANGELOG_TARGET: target }).status, 0);
    const out = readFileSync(target, "utf8");
    assert.match(out, /## \[0\.1\.0\]/u);
    mkdirSync(dir, { recursive: true });
  } finally {
    rmSync(dir, { recursive: true, force: true });
  }
});
