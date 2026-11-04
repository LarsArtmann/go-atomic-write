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
