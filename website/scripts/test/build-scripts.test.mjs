// Regression tests for the two load-bearing postbuild scripts:
//   - fix-csp.mjs  (idempotency + the 2026-11-04 quoted-hash incident class)
//   - sync-changelog.mjs (<-escaping so MDX cannot break the build)
// Run: pnpm run test:scripts  (node --test)
import { join, resolve } from "node:path";
import { mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import test from "node:test";
import { tmpdir } from "node:os";

const scriptsDir = resolve(import.meta.dirname, "..");

function runScript(script, env) {
  return spawnSync(process.execPath, [join(scriptsDir, script)], {
    env: { ...process.env, ...env },
    encoding: "utf8",
  });
}

const FIXTURE_PAGE = [
  "<!doctype html><html><head><title>t</title></head><body>",
  '<script>console.log("inline one")</script>',
  '<script src="/js/external.js"></script>',
  '<script type="application/ld+json">{"@type":"Thing"}</script>',
  '<script>console.log("inline two")</script>',
  "</body></html>",
].join("\n");

function makeDist() {
  const dist = mkdtempSync(join(tmpdir(), "csp-test-"));
  writeFileSync(join(dist, "index.html"), FIXTURE_PAGE, "utf8");
  return dist;
}

function cspMeta(html) {
  const metas = html.match(/<meta http-equiv="Content-Security-Policy"[^>]*>/giu) ?? [];
  return metas;
}

function scriptSrcFrom(meta) {
  const match = meta.match(/(?<src>script-src [^;"]+)/iu);
  return match?.groups.src ?? "";
}

function withDist(work) {
  const dist = makeDist();
  try {
    return work(dist);
  } finally {
    rmSync(dist, { recursive: true, force: true });
  }
}

function withTempSource(content, work) {
  const dir = mkdtempSync(join(tmpdir(), "changelog-test-"));
  const source = join(dir, "CHANGELOG.md");
  const target = join(dir, "changelog.mdx");
  writeFileSync(source, content, "utf8");
  try {
    return work({ source, target });
  } finally {
    rmSync(dir, { recursive: true, force: true });
  }
}

test("fix-csp hashes only inline scripts, quotes every hash source, allows buttondown", () => {
  withDist((dist) => {
    const run = runScript("fix-csp.mjs", { CSP_DIST_DIR: dist });
    assert.equal(run.status, 0, run.stderr);

    const metas = cspMeta(readFileSync(join(dist, "index.html"), "utf8"));
    assert.equal(metas.length, 1, "exactly one injected CSP meta");

    const sources = scriptSrcFrom(metas[0]).split(" ").filter((source) => source.includes("sha256-"));
    assert.equal(sources.length, 2, "both inline scripts hashed, external + ld+json skipped");
    for (const hash of sources) {
      assert.match(hash, /^'sha256-[A-Za-z0-9+/=]{43}='$/u, `hash must be quoted: ${hash}`);
    }
    assert.equal((scriptSrcFrom(metas[0]).match(/(?<boundary>^|\s)sha256-/gu) ?? []).length, 0, "no bare hash sources");
    assert.match(metas[0], /form-action 'self' https:\/\/buttondown\.email/u);
  });
});

test("fix-csp is idempotent: a second run is a byte-for-byte no-op with one meta", () => {
  withDist((dist) => {
    assert.equal(runScript("fix-csp.mjs", { CSP_DIST_DIR: dist }).status, 0);
    const afterFirst = readFileSync(join(dist, "index.html"), "utf8");
    assert.equal(runScript("fix-csp.mjs", { CSP_DIST_DIR: dist }).status, 0);
    const afterSecond = readFileSync(join(dist, "index.html"), "utf8");
    assert.equal(afterSecond, afterFirst, "second run must not change the file");
    assert.equal(cspMeta(afterSecond).length, 1, "exactly one CSP meta after double patching");
  });
});

test("sync-changelog escapes < for MDX, strips the H1, and is idempotent", () => {
  withTempSource("# Changelog\n\n## [Unreleased]\n\n- generic `a < b` comparison and `<div>` markup survive\n", ({ source, target }) => {
    const sync = () => runScript("sync-changelog.mjs", { CHANGELOG_SOURCE: source, CHANGELOG_TARGET: target });
    assert.equal(sync().status, 0);
    const first = readFileSync(target, "utf8");

    assert.equal(first.match(/^#\s+Changelog/mu), null, "leading H1 stripped");
    assert.match(first, /`a &lt; b`/u, "inline-code < escaped");
    assert.match(first, /`&lt;div>` markup/u, "< escaped; > kept by design (harmless in MDX)");
    assert.equal((first.match(/(?<!&)<(?!div>)/gu) ?? []).length, 0, "no unescaped < outside char refs");

    assert.equal(sync().status, 0);
    assert.equal(readFileSync(target, "utf8"), first, "second run is byte-identical");
  });
});

test("sync-changelog keeps a page without a leading H1 intact", () => {
  withTempSource("## [0.1.0] - 2026-06-03\n\n- entry\n", ({ source, target }) => {
    assert.equal(runScript("sync-changelog.mjs", { CHANGELOG_SOURCE: source, CHANGELOG_TARGET: target }).status, 0);
    assert.match(readFileSync(target, "utf8"), /## \[0\.1\.0\]/u);
  });
});

function writeFixture(dist, indexHtml) {
  writeFileSync(join(dist, "index.html"), indexHtml, "utf8");
  mkdirSync(join(dist, "docs"), { recursive: true });
  writeFileSync(join(dist, "docs", "index.html"), '<section id="setup">s</section>', "utf8");
  mkdirSync(join(dist, "home"), { recursive: true });
  writeFileSync(join(dist, "home", "other.html"), '<a href="../docs/">up</a>', "utf8");
}

function expectCheckLinks(dist, status, pattern) {
  const run = runScript("check-links.mjs", { DIST_DIR: dist });
  assert.equal(run.status, status, run.stderr);
  if (pattern) {
    assert.match(run.stderr, pattern);
  }
}

test("check-links accepts a clean fixture and rejects broken links and fragments", () => {
  const dist = mkdtempSync(join(tmpdir(), "links-test-"));
  try {
    writeFixture(dist, '<span id="install"></span><a href="/docs/">docs</a><a href="/docs/#setup">setup</a><a href="#install">install</a>');
    expectCheckLinks(dist, 0);

    writeFixture(dist, '<a href="/missing-page/">broken</a>');
    expectCheckLinks(dist, 1, /missing-page/u);

    writeFixture(dist, '<a href="/docs/#no-such-id">broken anchor</a>');
    expectCheckLinks(dist, 1, /no-such-id/u);
  } finally {
    rmSync(dist, { recursive: true, force: true });
  }
});
