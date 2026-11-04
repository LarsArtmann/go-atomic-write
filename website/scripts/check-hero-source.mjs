import { join } from "node:path";
import { readFileSync } from "node:fs";

// Source-level parity: the token list in hero-code.ts is the single source of
// truth for both the rendered hero code and the copy payload. This check runs
// against the SOURCE (pre-build), so a hero drift fails before dist exists.
const source = readFileSync(join(import.meta.dirname, "../src/data/hero-code.ts"), "utf8");

const tokens = [
	...source.matchAll(/text: (?<t>"[^"]*"|`[^`]*`),\s*className: "(?<c>[^"]+)"/gu),
].map((m) => `${m.groups.t},${m.groups.c}`);
const failures = [];

if (tokens.length < 8) {
	failures.push(`expected a full hero snippet (>= 8 tokens), found ${tokens.length}`);
}
const dupes = tokens.filter((l, i, all) => all.indexOf(l) !== i);
if (dupes.length > 0) {
	failures.push(`duplicate tokens in hero-code.ts: ${dupes.join("; ")}`);
}
if (!source.includes("WriteVerified")) {
	failures.push("hero snippet no longer teaches WriteVerified");
}
if (!source.includes("FingerprintFromBytes")) {
	failures.push("hero snippet no longer teaches FingerprintFromBytes");
}
for (const banned of ["atomicwrite.Write(path, newData, fp)", "Fingerprint{}"]) {
	if (tokens.some((t) => t.includes(banned))) {
		failures.push(`hero snippet teaches the removed API: ${banned}`);
	}
}

// The rendered markup must be derivable: every className in the token list
// must appear in the source as a literal string (Tailwind scans this file).
const classNames = [...source.matchAll(/className: "(?<c>[^"]+)"/gu)].map((m) => m.groups.c);
for (const c of classNames) {
	for (const cls of c.split(/\s+/u)) {
		if (!source.includes(cls)) {
			failures.push(`className ${cls} is not a literal in hero-code.ts (Tailwind would drop it)`);
		}
	}
}

if (failures.length > 0) {
	console.error("check-hero-source failed:");
	for (const failure of failures) {
		console.error(`  - ${failure}`);
	}
	process.exit(1);
}

console.log(`check-hero-source: ${tokens.length} tokens OK, API current, classes literal`);
