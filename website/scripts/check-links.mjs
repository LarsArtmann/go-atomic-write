// Internal link checker: every href/src in the built site must resolve inside
// dist/, and every fragment must point at an existing id. Also checks the
// README's relative asset references. Run against a real build:
//   pnpm run build && pnpm run check:links
// Override the dist dir with DIST_DIR (used by the regression tests).
import { dirname, join, relative, resolve } from "node:path";
import { readFileSync, readdirSync, statSync } from "node:fs";

const repoRoot = resolve(import.meta.dirname, "../..");
const distDir = resolve(process.env.DIST_DIR ?? resolve(import.meta.dirname, "../dist"));

const failures = [];

function* walk(dir) {
	for (const entry of readdirSync(dir, { withFileTypes: true })) {
		const full = join(dir, entry.name);
		if (entry.isDirectory()) {
			yield* walk(full);
		} else {
			yield full;
		}
	}
}

function pageHasId(html, id) {
	return new RegExp(String.raw`id=["']${id.replaceAll(/["'\\]/gu, String.raw`\$&`)}["']`, "u").test(
		html,
	);
}

function resolveTarget(fromFile, href) {
	const [path] = href.split("#");
	const target = path.startsWith("/") ? join(distDir, path) : resolve(dirname(fromFile), path);
	if (path === "" || path.endsWith("/")) {
		return `${target}/index.html`;
	}
	return target;
}

for (const file of [...walk(distDir)].filter((f) => f.endsWith(".html"))) {
	const html = readFileSync(file, "utf8");
	const refs = [...html.matchAll(/(?:href|src)="(?<ref>[^"]+)"/gu)].map((m) => m.groups.ref);
	const unique = [...new Set(refs)];

	for (const ref of unique) {
		if (ref.startsWith("http://") || ref.startsWith("https://") || ref.startsWith("mailto:")) {
			continue;
		}

		if (ref.startsWith("#")) {
			const id = ref.slice(1);
			if (id && !pageHasId(html, id)) {
				failures.push(`${relative(distDir, file)}: fragment ${ref} has no matching id on the page`);
			}
			continue;
		}

		const target = resolveTarget(file, ref);
		let exists = false;
		try {
			exists = statSync(target).isFile();
		} catch {
			exists = false;
		}
		if (!exists) {
			failures.push(`${relative(distDir, file)}: ${ref} -> missing ${relative(distDir, target)}`);
			continue;
		}

		const [, hash] = ref.split("#");
		if (hash) {
			const targetHtml = readFileSync(target, "utf8");
			if (!pageHasId(targetHtml, hash)) {
				failures.push(`${relative(distDir, file)}: ${ref} -> target page has no id "${hash}"`);
			}
		}
	}
}

const readme = readFileSync(join(repoRoot, "README.md"), "utf8");
for (const ref of [...readme.matchAll(/\]\((?<ref>[^)#]+)(?:#[^)]*)?\)/gu)].map(
	(m) => m.groups.ref,
)) {
	if (ref.startsWith("http://") || ref.startsWith("https://") || ref.startsWith("mailto:")) {
		continue;
	}
	try {
		statSync(resolve(repoRoot, ref)).isFile();
	} catch {
		failures.push(`README.md: ${ref} -> missing file in repo`);
	}
}

if (failures.length > 0) {
	console.error(`check-links: ${failures.length} broken link(s):`);
	for (const failure of failures) {
		console.error(`  - ${failure}`);
	}
	process.exit(1);
}
console.log(`check-links: OK (all internal links and fragments resolve)`);
