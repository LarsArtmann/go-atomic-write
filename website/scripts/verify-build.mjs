import { createHash } from "node:crypto";
import { join } from "node:path";
import { readFileSync } from "node:fs";

const websiteRoot = join(import.meta.dirname, "..");
const repoRoot = join(websiteRoot, "..");

const failures = [];
const check = (condition, message) => {
	if (!condition) {
		failures.push(message);
	}
};

const decodeEntities = (value) =>
	value
		.replaceAll("&quot;", '"')
		.replaceAll("&#39;", "'")
		.replaceAll("&lt;", "<")
		.replaceAll("&gt;", ">")
		.replaceAll("&amp;", "&");

const stripTags = (value) => value.replaceAll(/<[^>]+>/gu, "");

const html = readFileSync(join(websiteRoot, "dist", "index.html"), "utf8");
const goMod = readFileSync(join(repoRoot, "go.mod"), "utf8");

const modulePath = goMod.match(/^module\s+(?<modulePath>\S+)/mu)?.groups?.modulePath;
check(Boolean(modulePath), "could not read the module path from go.mod");

const heroBlock = html.match(/<pre[^>]*><code[^>]*>(?<codeBlock>[\s\S]*?)<\/code><\/pre>/u);
check(Boolean(heroBlock), "hero code block not found in dist/index.html");

const visibleCode = heroBlock ? decodeEntities(stripTags(heroBlock.groups.codeBlock)) : "";

const copyButton = html.match(
	/code-preview[\s\S]*?<button[^>]*data-copy[^>]*data-code="(?<codeValue>[^"]*)"/u,
);
check(Boolean(copyButton), "hero copy button with data-code not found in dist/index.html");

const copiedCode = copyButton ? decodeEntities(copyButton.groups.codeValue) : "";

check(
	visibleCode === copiedCode,
	"hero code shown on the page does not match the code the copy button copies",
);

if (modulePath) {
	check(
		html.includes(`go get ${modulePath}@latest`),
		`install command is missing "go get ${modulePath}@latest"`,
	);
	check(
		visibleCode.includes(`"${modulePath}"`),
		`hero import path does not match the go.mod module path (${modulePath})`,
	);
}

const copyButtons = html.match(/<button[^>]*data-copy[^>]*>/gu) ?? [];
check(
	copyButtons.length === 2,
	`expected exactly two copy buttons on the landing page, found ${copyButtons.length}`,
);

for (const tag of copyButtons) {
	const code = tag.match(/data-code="(?<codeValue>[^"]*)"/u)?.groups?.codeValue ?? "";
	check(decodeEntities(code).trim().length > 0, `copy button has an empty data-code: ${tag}`);
}

const metaDescription = html.match(/<meta name="description" content="(?<desc>[^"]*)"/u)?.groups
	?.desc;
check(
	Boolean(metaDescription && metaDescription.length >= 50),
	"meta description missing or suspiciously short",
);
check(
	!metaDescription?.toLowerCase().includes("toctou"),
	"meta description still uses the mechanism register (TOCTOU)",
);

const faqCount = (html.match(/<details[^>]*class="[^"]*group[^"]*"/gu) ?? []).length;
check(faqCount === 6, `expected 6 FAQ items, found ${faqCount}`);

const staleApi = ["atomicwrite.Write(path, newData, fp)", "WriteFunc(path, fn, Fingerprint{})"];
for (const pattern of staleApi) {
	check(!html.includes(pattern), `dist teaches the removed pre-0.4.0 API: ${pattern}`);
}

// Every inline <script> (Astro-inlined modules) must have its quoted hash in
// the meta CSP — the 2026-11-04 incident: bare hash tokens were invalid, so
// every inline script was silently blocked.
const csp = html.match(/<meta http-equiv="Content-Security-Policy" content="(?<csp>[^"]*)"/u)
	?.groups?.csp;
check(Boolean(csp), "CSP meta not found in dist/index.html");
const quotedHashes = new Set(csp?.match(/'sha256-[^']+'/gu));
const inlineScripts = html.matchAll(
	/<script(?![^>]*\bsrc=)(?![^>]*type="application\/ld\+json")(?<attrs>[^>]*)>(?<body>[\s\S]*?)<\/script>/gu,
);
for (const m of inlineScripts) {
	if (m.groups.body.trim().length === 0) {
		continue;
	}
	const hash = `'sha256-${createHash("sha256").update(m.groups.body, "utf8").digest("base64")}'`;
	check(
		quotedHashes.has(hash),
		`inline script (attrs: ${m.groups.attrs.slice(0, 40)}) is not allow-listed by a quoted hash in the CSP`,
	);
}

if (failures.length > 0) {
	console.error("verify-build failed:");
	for (const failure of failures) {
		console.error(`  - ${failure}`);
	}
	process.exit(1);
}

console.log(
	`verify-build: hero code parity OK, module path OK (${modulePath}), ${copyButtons.length} copy buttons checked`,
);
