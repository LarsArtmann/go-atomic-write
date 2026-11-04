import { defineConfig, devices } from "@playwright/test";

// Smoke suite over the PRODUCTION build (`pnpm run build && pnpm run test:ui`).
// The webServer target is `astro preview` - it serves dist/, so a failing build
// cannot mask as a passing test.
export default defineConfig({
	testDir: "tests",
	timeout: 30_000,
	fullyParallel: true,
	reporter: [["list"]],
	use: {
		...devices["Desktop Chrome"],
		baseURL: "http://localhost:4321",
	},
	webServer: {
		command: "pnpm run preview",
		url: "http://localhost:4321",
		reuseExistingServer: true,
		timeout: 60_000,
	},
	projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
});
