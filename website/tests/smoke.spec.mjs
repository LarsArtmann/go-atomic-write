// Smoke tests for the built landing page. These are the gates that "verified
// by construction" could not cover: real rendered DOM, real copy payloads,
// real layout at 320 px, and the astro:page-load animation binding.
import { expect, test } from "@playwright/test";

const installCommand = "go get github.com/larsartmann/go-atomic-write@latest";

test("both copy buttons carry the exact payload they display", async ({ page }) => {
	await page.goto("/");
	const buttons = page.locator("[data-copy]");
	await expect(buttons).toHaveCount(2);

	const installCode = await buttons.nth(0).getAttribute("data-code");
	expect(installCode).toBe(installCommand);

	const heroCode = await buttons.nth(1).getAttribute("data-code");
	expect(heroCode).toContain("WriteVerified");
	expect(heroCode?.length ?? 0).toBeGreaterThan(80);
});

test("landing FAQ has exactly six questions", async ({ page }) => {
	await page.goto("/");
	await expect(page.locator("details")).toHaveCount(6);
});

test("no horizontal overflow at 320 px", async ({ page }) => {
	await page.setViewportSize({ width: 320, height: 800 });
	await page.goto("/");
	const overflow = await page.evaluate(
		() => document.documentElement.scrollWidth - document.documentElement.clientWidth,
	);
	expect(overflow).toBeLessThanOrEqual(0);
});

test("animations bind via astro:page-load (the navigation-rebind contract)", async ({ page }) => {
	await page.goto("/");
	await expect(page.locator("[data-animate-bound]").first()).toBeAttached();
	const boundCount = await page.locator("[data-animate-bound]").count();
	const total = await page.locator("[data-animate]").count();
	expect(boundCount).toBe(total);
});

test("keyboard tab order starts at the skip link, then header nav", async ({ page }) => {
	await page.goto("/");
	const seq = [];
	for (let i = 0; i < 4; i++) {
		await page.keyboard.press("Tab");
		seq.push(
			await page.evaluate(() => {
				const el = document.activeElement;
				if (!el) return "none";
				const label = el.getAttribute("aria-label") ?? el.textContent?.trim().slice(0, 24) ?? "";
				return `${el.tagName.toLowerCase()}:${label}`;
			}),
		);
	}
	expect(seq[0]).toContain("Skip to content");
	expect(seq[1]).toMatch(/^(a|button):/);
	expect(seq.some((s) => s.includes("Docs"))).toBe(true);
});

test("theme applies before first paint and the theme-color pair matches the tokens", async ({
	browser,
}) => {
	// Light stored: the .light class must be on <html> at domcontentloaded —
	// theme-init.js runs before paint, so no wrong-theme flash is possible.
	const lightCtx = await browser.newContext();
	await lightCtx.addInitScript(() => {
		localStorage.setItem("theme", "light");
	});
	const lightPage = await lightCtx.newPage();
	await lightPage.goto("/", { waitUntil: "domcontentloaded" });
	expect(await lightPage.evaluate(() => document.documentElement.classList.contains("light"))).toBe(
		true,
	);

	// Dark stored (system-pref default or explicit): no .light class.
	const darkCtx = await browser.newContext();
	await darkCtx.addInitScript(() => {
		localStorage.setItem("theme", "dark");
	});
	const darkPage = await darkCtx.newPage();
	await darkPage.goto("/", { waitUntil: "domcontentloaded" });
	expect(await darkPage.evaluate(() => !document.documentElement.classList.contains("light"))).toBe(
		true,
	);

	// The theme-color pair must mirror the token values exactly.
	const themeColors = await darkPage.evaluate(() =>
		[...document.querySelectorAll('meta[name="theme-color"]')].map((m) =>
			m.getAttribute("content"),
		),
	);
	expect(themeColors).toContain("#0a0f0d");
	expect(themeColors).toContain("#f4f3ee");

	await lightCtx.close();
	await darkCtx.close();
});
