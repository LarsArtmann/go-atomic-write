import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const websiteRoot = join(dirname(fileURLToPath(import.meta.url)), "..");
const repoRoot = join(websiteRoot, "..");

const failures = [];
const check = (condition, message) => {
	if (!condition) failures.push(message);
};

const decodeEntities = (value) =>
	value
		.replace(/&quot;/g, '"')
		.replace(/&#39;/g, "'")
		.replace(/&lt;/g, "<")
		.replace(/&gt;/g, ">")
		.replace(/&amp;/g, "&");

const stripTags = (value) => value.replace(/<[^>]+>/g, "");

const html = readFileSync(join(websiteRoot, "dist", "index.html"), "utf8");
const goMod = readFileSync(join(repoRoot, "go.mod"), "utf8");

const modulePath = goMod.match(/^module\s+(\S+)/m)?.[1];
check(Boolean(modulePath), "could not read the module path from go.mod");

const heroBlock = html.match(/<pre[^>]*><code[^>]*>([\s\S]*?)<\/code><\/pre>/);
check(Boolean(heroBlock), "hero code block not found in dist/index.html");

const visibleCode = heroBlock ? decodeEntities(stripTags(heroBlock[1])) : "";

const copyButton = html.match(/code-preview[\s\S]*?<button[^>]*data-copy[^>]*data-code="([^"]*)"/);
check(Boolean(copyButton), "hero copy button with data-code not found in dist/index.html");

const copiedCode = copyButton ? decodeEntities(copyButton[1]) : "";

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

const copyButtons = html.match(/<button[^>]*data-copy[^>]*>/g) ?? [];
check(copyButtons.length >= 2, "expected at least two copy buttons on the landing page");

for (const tag of copyButtons) {
	const code = tag.match(/data-code="([^"]*)"/)?.[1] ?? "";
	check(decodeEntities(code).trim().length > 0, `copy button has an empty data-code: ${tag}`);
}

if (failures.length > 0) {
	console.error("verify-build failed:");
	for (const failure of failures) console.error(`  - ${failure}`);
	process.exit(1);
}

console.log(
	`verify-build: hero code parity OK, module path OK (${modulePath}), ${copyButtons.length} copy buttons checked`,
);
