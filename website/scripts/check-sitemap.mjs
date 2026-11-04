import { join } from "node:path";
import { readFileSync } from "node:fs";

// Sitemap sanity gate (postbuild): the sitemap must exist, cover the built
// pages, and carry real lastmod dates. Born from the 07:09 f.17 finding that
// nothing asserted the sitemap after the `lastmod` crash fix — a silently
// empty or lastmod-less sitemap would ship unnoticed.
const MIN_URLS = 10;

const distDir = process.env.DIST_DIR_OVERRIDE ?? join(import.meta.dirname, "..", "dist");
const failures = [];
const check = (condition, message) => {
	if (!condition) {
		failures.push(message);
	}
};

let sitemapXml = "";
try {
	const indexXml = readFileSync(join(distDir, "sitemap-index.xml"), "utf8");
	const loc = indexXml.match(/<loc>(?<loc>[^<]+)<\/loc>/u)?.groups?.loc;
	check(Boolean(loc), "sitemap-index.xml has no <loc> entry");
	if (loc) {
		sitemapXml = readFileSync(join(distDir, loc.split("/").pop()), "utf8");
	}
} catch (error) {
	failures.push(`could not read the sitemap from dist/: ${error.message}`);
}

if (sitemapXml) {
	const urlCount = sitemapXml.match(/<url>/gu)?.length ?? 0;
	check(
		urlCount >= MIN_URLS,
		`sitemap lists ${urlCount} urls (expected >= ${MIN_URLS}) - pages went missing from the build`,
	);
	check(
		Boolean(sitemapXml.match(/<lastmod>\d{4}-\d{2}-\d{2}/u)),
		"sitemap has no <lastmod> dates (the lastmod Date-object regression)",
	);
}

if (failures.length > 0) {
	console.error("check-sitemap failed:");
	for (const failure of failures) {
		console.error(`  - ${failure}`);
	}
	process.exit(1);
}

const listed = sitemapXml.match(/<url>/gu)?.length ?? 0;
console.log(`check-sitemap: OK (${listed} urls, lastmod present)`);
