# AGENTS.md — go-atomic-write

## Project

Single-package Go library providing TOCTOU-safe file writes via xxhash64 fingerprint verification, cross-platform file locking (`flock`/`LockFileEx`), atomic rename, and fsync for crash durability.

- **Module:** `github.com/larsartmann/go-atomic-write`
- **Go version:** 1.27
- **Main branch:** `master` (configured in `git-town.toml`)

## Commands

| Command                         | Purpose                                                               |
| ------------------------------- | --------------------------------------------------------------------- |
| `go test ./...`                 | Run all tests                                                         |
| `go test -race ./...`           | Run all tests with the race detector                                  |
| `go test -bench=. -benchmem`    | Run benchmarks (in `hash_bench_test.go`)                              |
| `go vet ./...`                  | Static analysis                                                       |
| `go build ./...`                | Verify compilation                                                    |
| `golangci-lint run ./...`       | Run all configured linters — MUST exit with `0 issues` before merging |
| `golangci-lint run ./... --fix` | Auto-fix formatting/imports (gci, gofumpt, goimports, golines)        |

`golangci-lint` config is the **gating quality check** of this project. The configuration lives in `.golangci.yml` and enables ~100 linters at their strictest defaults (`depguard` and `noinlineerr` are disabled by the converged config).

No Makefile; all commands are plain `go` toolchain. CI (`.github/workflows/ci.yml`) gates merges with test + race, lint, cross-compile (windows/amd64, darwin/amd64+arm64), `govulncheck`, and a drift guard (`go mod tidy -diff`, version-trio check, `pnpm install --frozen-lockfile`).

## Structure

Flat single-package layout — all source in the repository root:

| File                  | Purpose                                                                                                                            |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| `atomicwrite.go`      | Public API + staging: `Fingerprint`, `Write`, `WriteVerified`, `WriteIfChanged`, `FingerprintFile`, `writeAndSync`, `randomSuffix` |
| `rename_unix.go`      | POSIX `atomicRename` (single `rename(2)` + directory `fsync`)                                                                      |
| `rename_windows.go`   | Windows `atomicRename` (retry on `ERROR_ACCESS_DENIED`/`ERROR_SHARING_VIOLATION`)                                                  |
| `atomicwrite_test.go` | Unit + concurrency + integrity tests                                                                                               |
| `hash_bench_test.go`  | xxhash64 vs SHA-256 benchmarks                                                                                                     |

## Website

Marketing website and documentation built with Astro + Starlight + Tailwind v4. Deployed to Firebase Hosting.

- **Live URL:** `https://atomicwrite.lars.software` (DNS pending) / `https://atomicwrite.web.app` (live now)
- **Alt domain:** `https://go-atomic-write.lars.software` (DNS pending)
- **Firebase project:** `lars-software`
- **Firebase hosting site:** `atomicwrite`
- **Accent color:** Emerald (`#10b981`) — distinct from gogenfilter's cyan

### Website Commands

| Command                                                                | Purpose                                                                           |
| ---------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| `cd website && pnpm run dev`                                           | Local dev server                                                                  |
| `cd website && pnpm run build`                                         | Production build to `dist/` (`prebuild` syncs changelog, `postbuild` injects CSP) |
| `cd website && pnpm run typecheck`                                     | TypeScript + Astro type checking                                                  |
| `cd website && pnpm run preview`                                       | Preview production build locally                                                  |
| `cd website && pnpm run sync:changelog`                                | Regenerate changelog page from root `CHANGELOG.md`                                |
| `cd website && pnpm run brand:generate`                                | Regenerate all brand artifacts from the single geometry (`scripts/generate-brand-assets.mjs`) |
| `cd website && pnpm run brand:check`                                   | Exit 1 if any generated brand artifact drifted (CI gates on this)                 |
| `cd website && pnpm run lighthouse`                                    | Run Lighthouse CI against `dist/` (requires chromium in Nix devShell)             |
| `cd website && firebase deploy --only hosting --project lars-software` | Deploy to Firebase                                                                |

Node.js 24 required (use `nix shell nixpkgs#nodejs_24` if not in PATH).

### Website Structure

| Path                            | Purpose                                                                                                                            |
| ------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| `website/astro.config.mjs`      | Astro config: Starlight (incl. docs-header `logo`), fonts, sitemap                                                                 |
| `website/src/pages/index.astro` | Landing page                                                                                                                       |
| `website/src/components/`       | 16 Astro components (Hero, FeatureGrid, HowItWorks, Comparison, FAQ, etc.)                                                         |
| `website/src/data/`             | Typed content: config, features, sections, hero-code                                                                               |
| `website/src/assets/brand.svg`  | GENERATED bundled copy of the mark for the Starlight docs header logo                                                              |
| `website/src/content/docs/`     | 9 Starlight documentation pages                                                                                                    |
| `website/src/styles/`           | global.css (emerald theme) + starlight.css                                                                                         |
| `website/public/`               | favicon.svg, favicon.ico, og-image.svg (source) + og-image.png (1200×630), manifest, robots.txt, JS (theme, animations, copy-code) |
| `assets/brand/`                 | GENERATED standalone logo marks + lockups (`mark.svg`, horizontal/stacked lockups in light+dark) + brand README                    |
| `website/firebase.json`         | Hosting config with security headers                                                                                               |
| `website/.firebaserc`           | Firebase project + hosting target                                                                                                  |
| `website/lighthouserc.json`     | Lighthouse CI config: performance/accessibility/SEO budgets (desktop preset)                                                       |
| `website/.editorconfig`         | Standalone editor config for the website subtree (2-space indent, `root = true`)                                                   |
| `website/scripts/`              | Build-time tooling: `sync-changelog.mjs` (prebuild), `fix-csp.mjs` + `verify-build.mjs` (postbuild)                                |
| `scripts/generate-brand-assets.mjs` | GENERATED brand system: one geometry → favicon, Logo.astro, OG card, lockups; `--check` mode gates CI                          |
| `.github/workflows/`            | CI (`ci.yml`: test/race + lint + cross-compile + govulncheck + drift guard incl. brand:check) + website build/deploy (`website.yml`) |

### DNS

CNAME records for `atomicwrite` and `go-atomic-write` subdomains are defined in `/home/lars/projects/domains/lars.software.tf`. Both point to `atomicwrite.web.app`. Terraform apply requires a Namecheap API key that is not stored in this repo.

## Architecture & Data Flow

```
Caller reads file → computes Fingerprint → modifies data → calls Write()
  └─ Write() stages to unique .tmp → fsync .tmp → locks + verifies fingerprint → atomic rename → fsync dir
       └─ On mismatch: returns ErrConcurrentModification, cleans up .tmp
```

Key internal functions:

- `Write()` — entry point; generates unique temp path, stages + fsyncs data, branches on fingerprint
- `WriteIfChanged()` — idempotent write; skips if content matches disk, delegates to `Write` (first-write) or `WriteVerified` (content change)
- `writeAndSync()` — creates temp file, writes data, fsyncs, closes (with cleanup on error)
- `randomSuffix()` — generates random hex suffix for unique temp file names (crypto/rand)
- `commitWithVerification()` — acquires exclusive `flock`, re-reads target, verifies match, renames
- `atomicRename()` — single `os.Rename` + directory fsync (POSIX); retry loop (Windows)
- `syncDir()` — opens and fsyncs the target directory after rename (POSIX only)

## Dependencies

| Dependency          | Purpose                                                                        |
| ------------------- | ------------------------------------------------------------------------------ |
| `cespare/xxhash/v2` | Non-cryptographic hash for fingerprinting (chosen over SHA-256 for ~11× speed) |
| `gofrs/flock`       | Cross-platform file locking                                                    |

Both are intentional, minimal, and not candidates for replacement.

## Testing Conventions

- Standard `testing` package — no testify or other test utilities
- `t.Parallel()` on most tests
- `t.TempDir()` for filesystem isolation
- `tempFile(t, content)` helper creates a file in a temp dir
- Concurrency test (`TestConcurrentWriteRACE`) uses `sync.WaitGroup.Go` and `atomic.Int32` with divergent content per writer — verifies both `flock` contention and data integrity (no byte interleaving)

## Release Checklist

1. `golangci-lint run ./...` exits with `0 issues` — the gate; fix before anything else
2. CI green on `master` BEFORE tagging (test + race, lint, cross-compile, govulncheck, drift guard) — v0.5.1 shipped a Windows compile break because Linux-only CI passed while the tag was already cut
3. `go build ./...` and `go test -race ./...` locally
4. Move `CHANGELOG.md` `[Unreleased]` content into a new `[X.Y.Z]` section
5. If the toolchain changed, bump the version trio together (`go.mod`, `.golangci.yml` `run.go`, this file) — CI now fails on drift
6. Annotated tag `vX.Y.Z`, push the tag — never re-tag/delete tags (module proxy poisoning)
7. Verify `pkg.go.dev` shows the new version and `go get github.com/larsartmann/go-atomic-write@vX.Y.Z` resolves
8. Website: changelog page syncs on build; deploy rides the `Website` workflow or `firebase deploy --only hosting:atomicwrite`

## Gotchas

- **Website build (pnpm 11): build-script approvals live in `website/pnpm-workspace.yaml` under `allowBuilds:`** (`esbuild: true`) — pnpm v11 ignores `pnpm.*` in `package.json` and silently skips unapproved postinstall scripts, so `astro build` then fails on a missing esbuild binary. A placeholder value (e.g. `esbuild: set this to true or false`) silently disables the whole key (cmdguard incident, fixed 2026-09-19).
- **pnpm v11 reads `overrides` from `website/pnpm-workspace.yaml`, NOT `package.json`** — the package.json `overrides` block was silently ignored (even `pnpm install --force` no-op'd against it; proven 2026-11-03 while fixing pnpm-audit findings). CI is pnpm-only (corepack pnpm@11.20.0 in both workflows), so workspace.yaml is the single source of truth; the package.json block is deleted. Edit overrides there, then `pnpm install` in the same change (frozen-lockfile drift guard).
- **pnpm-audit baseline (2026-11-03): 1 moderate accepted** — `minimatch@3.1.5` (DoS via globparts amplification) rides `@lhci/cli@0.16.0` → `chrome-launcher@0.13.4` → `rimraf@3` → `glob@7`. Every fix path verified broken at source level: rimraf v6 removes the callback API chrome-launcher 0.13.4 uses (`this.rimraf(dir, cb)`), minimatch 11 breaks glob 7, and latest @lhci/cli still pins chrome-launcher ^0.13.4. The moderate is dev-only, non-attacker-reachable, and does not gate (`--fail-on=error`). All other audit findings were cleared via workspace.yaml overrides (`basic-ftp ^6.2.1`, `glob@>=13 <14 → ^14.0.2`, `minimatch@>=10 <11 → ^11.0.1`). Revisit when @lhci/cli ships chrome-launcher ^1.
- **BuildFlow scanning policy lives in `.buildflow.yml` + `.jscpd.json` (both repo-root, 2026-11-03)**: `exclude: website/dist` keeps gitignored build output out of file-scanning tools; `skip_steps: [type-check]` because BuildFlow's type-check is generic `tsc --noEmit`, which structurally cannot resolve Astro virtual modules (`astro:content`) — it reported 4 errors while `astro check` was green (0 errors/33 files); the honest TS gate is website CI's `pnpm run typecheck`. jscpd's dist exclusion must live in `.jscpd.json` as `**/dist/**` — the `website/dist/**` form does NOT match jscpd 4.x globs, jscpd ignores nested `website/.gitignore` rules, and the root `.gitignore`'s `/dist/` is root-anchored.
- **Keep the Go version trio in sync (`go.mod` / `.golangci.yml` `run.go` / this file)** — drift has recurred four times (1.26.3 → 1.26.4 → 1.26.5 → 1.27); the CI drift-guard job now fails on mismatch, so bump all three together
- `.tmp` files use unique names (`path + "." + randomHex + ".tmp"`) to prevent concurrent writers from corrupting a shared staging file
- `.tmp` files are created alongside the target file (same directory) — callers need write permissions on the directory, not just the file
- `.tmp` files are in `.gitignore`; `.bak` files are no longer created (the `.bak` pattern was removed — it broke POSIX atomicity)
- `fsync` is called on both the temp file (before rename) and the target directory (after rename, POSIX only) for crash durability
- On Windows, directory fsync is a no-op (no POSIX-equivalent directory sync concept); file data durability is handled by `FlushFileBuffers` via `file.Sync()`
- On Windows, `os.Rename` is retried up to 5 times with exponential backoff on `ERROR_ACCESS_DENIED`/`ERROR_SHARING_VIOLATION` (antivirus, open handles)
- `Fingerprint` is `[8]byte` (xxhash64), stored big-endian — do not compare with `string` or `[]byte` representations of hashes
- `FingerprintFile` returns zero-value (not an error) for nonexistent files — this is the "first write" sentinel
- File permissions are preserved from the existing file, defaulting to `0644` for new files
- `ErrConcurrentModification` is a sentinel `errors.New` value — always check with `errors.Is`, not string matching
- Only `//nolint:makezero` remains in `hash_bench_test.go` — the `//nolint:gosec` comments were removed in 0.5.2 when the linter config was realigned; do not re-add them reflexively
- **golangci-lint GitHub Action must be v7+** for golangci-lint v2 (action v6 refuses `version: v2.x`), and the lint version must be pinned explicitly (`version: v2.13.2`) — `version: latest` resolves to the ancient v1.x line built with go1.24
- Lint is gating: a single `golangci-lint` issue blocks merging. Common pitfalls:
  - Single-letter variables (`fp`, `wg`, `h`, `mb`) trigger `varnamelen` — use full words
  - Inline `if err := …; err != nil` triggers `noinlineerr` — assign first, then check
  - `b.Helper()` is required on every benchmark helper that takes `*testing.B` (thelper)
  - 0o644 in tests triggers `gosec G306` — annotate with `//nolint:gosec` and rationale
  - `os.ReadFile(path)`/`os.OpenFile` with a variable path triggers `gosec G304` — annotate with `//nolint:gosec` and rationale
  - Blank line between error assignment and `if err != nil` triggers `wsl_v5` — put the blank line above the assignment instead
  - `depguard` is disabled in the converged config — new third-party imports need no allow-list edit
  - `exhaustruct` → `exhaustruct_v5` migration is DONE (config enables `exhaustruct_v5`; do not re-add the old name)
- **Website changelog page is GENERATED** — `website/src/content/docs/changelog.mdx` is produced from root `CHANGELOG.md` by `scripts/sync-changelog.mjs` (runs on `prebuild`/`predev`). Edit `CHANGELOG.md`, never the `.mdx`. If the build breaks with an MDX "unexpected character" error, a changelog entry likely contains a raw `<`; the sync script escapes `<` but check new entries.
- **Website CSP is hash-based and post-build** — `scripts/fix-csp.mjs` (`postbuild`) injects a per-file CSP `<meta>` from inline-script SHA-256 hashes. There is **no `'unsafe-inline'` for `script-src`**. If a new inline script breaks the site under CSP, either move it to an external `/js/*.js` file (loaded via `<script is:inline src=…>`) or confirm `fix-csp.mjs` hashed it.
- **Do NOT add `website/src/pages/404.astro`** — Starlight ships its own `404.html`. A custom one causes a route collision (warning today, hard error in future Astro). The `404 was not found` line during build is a benign Starlight route log, not a warning.
- The website's `pnpm run build` requires a clean `.astro` cache when content files are renamed/extension-changed: run `rm -rf .astro dist node_modules/.cache` if the content layer fails to resolve a renamed doc.
- **OG image is static, not dynamic** — `public/og-image.svg` is GENERATED by the brand generator (its embedded logo slot can never drift from the mark); `public/og-image.png` is the generated 1200×630 output. Regenerate the PNG with `magick -background none public/og-image.svg -resize 1200x630! public/og-image.png` (ImageMagick 7). No `astro-og-canvas` dependency needed.
- **Lighthouse CI needs chromium** — `pnpm run lighthouse` requires Chrome/Chromium. The Nix devShell provides it and sets `CHROME_PATH`. Run via `nix develop` or `nix shell nixpkgs#nodejs_24 nixpkgs#chromium`.
- **Brand assets are GENERATED — never hand-edit them** (2026-11-03): `assets/brand/*.svg`, `website/src/components/Logo.astro`, `website/public/favicon.svg`, `website/src/assets/brand.svg`, and `website/public/og-image.svg` are all written by `scripts/generate-brand-assets.mjs` (identity: "the atomic swap" — emerald squircle, new bytes in →, old bytes out ←). Edit the generator, then `cd website && pnpm run brand:generate`. `brand:check` (exact-bytes compare) runs as the first CI drift-guard step and in the Website workflow, so hand-edits fail CI. `favicon.ico` and `og-image.png` are ImageMagick-rasterized and NOT byte-gated — regenerate them after changing the SVGs (commands in `assets/brand/README.md`).
- **Starlight's `logo.src` resolves through the Astro asset pipeline, not public URLs** — `logo: { src: "/favicon.svg" }` fails the build with `ImageNotFound`; point it at a bundled copy instead (`logo: { src: "./src/assets/brand.svg" }`, the generator writes that copy). Same pattern typespec-asyncapi uses.
- **`website/.editorconfig` uses `root = true`** — this intentionally decouples the website from the Go-rooted `.editorconfig` (which uses tabs). The website uses 2-space indentation for all files.
- **`astro check` requires TypeScript 6.x** — TS7's native compiler dropped the programmatic API that `astro check` needs. The fleet npm sweep re-bumps the manifest to `^7.0.2` (three times on 2026-11-03 alone); the `"typescript@>=7": ^6.0.3` override in `pnpm-workspace.yaml` self-heals any 7.x demand back to 6.0.3 (verified: manifest at `^7.0.2` resolves 6.0.3). Do NOT fight the manifest — leave it, the override wins. Website CI's `pnpm run typecheck` is the alarm if the guard is ever removed.
- **The auto-commit daemon commits manifest changes without lockfiles** — this caused four `pnpm install --frozen-lockfile` CI failures (09-09, 09-13, 09-15, 09-17) and another four-package drift on 11-03. The CI drift-guard job now fails fast on it. After any `package.json`/`go.mod` change, regenerate the lockfile (`pnpm install`, `go mod tidy`) in the same change or CI breaks.
- **No npm manifest at the repo root** — a stray `npm init -y` stub was auto-committed 2026-09-22 (commit 5a6302f) and the `npm directory: /` Dependabot entry then failed on every run (no lockfile, no deps). The stub is deleted and the Dependabot entry removed (11-03); it recurred the SAME day (a second stub committed alongside a lockfile-less `npm directory: /` entry, also removed 11-03) — if a JS tool ever needs a root-level package.json, add the lockfile AND a matching Dependabot entry together, never either alone.
- **CI pins `GOTOOLCHAIN: local`** — no rogue toolchain download can silently pick a different Go; raising `go.mod`'s directive relies on setup-go installing that version from `go-version-file`.
- **A fleet npm sweep force-bumps devDeps to latest majors on a schedule** — on 2026-11-03 within one session it bumped `typescript` ×3, regenerated `global.out.css` ×2, and dprint re-quoted `pnpm-workspace.yaml`, daemon-committing each pass. Guards in place: `typescript@>=7` override (self-heal), `src/styles/*.out.css` gitignore, CI drift-guard + Website typecheck. When unexplained artifacts appear after a sweep, add the missing gitignore/override guard rather than re-deleting by hand. The `basic-ftp: ^6.2.1` override is deliberate: get-uri declares ^5 but 6.2.5 was verified CJS-loadable (2026-11-03) — if Lighthouse browser download ever breaks, remove that override first.
- **Website deploy needs `FIREBASE_SERVICE_ACCOUNT_LARS_SOFTWARE`** — a repo secret holding a JSON key for `github-website-deploy@lars-software.iam.gserviceaccount.com` (the shared GitHub-deploy SA; templ-components uses the same SA pattern). Without it every CI deploy fails at the final step and the live site silently freezes (it was stale 2026-07-26 → 2026-09-17 for exactly this reason). Local deploys work without the secret: `nix shell nixpkgs#nodejs nixpkgs#firebase-tools -c firebase deploy --only hosting:atomicwrite --project lars-software` (firebase-tools is logged in as `lartyhd@gmail.com`). IAM for `lars-software` is only visible to `lartyhd@gmail.com`, not `lars@helpless.ai`.
- **The `firebase-adminsdk` SA has hit GCP's 10-key cap** — minting new keys for it fails with `FAILED_PRECONDITION`. Use `github-website-deploy@lars-software` for new CI keys instead.
- **The landing-page hero snippet has one source of truth** — `website/src/data/hero-code.ts` holds a token list; `heroCode` (clipboard payload) and `highlightedHeroCode()` (rendered markup) are both derived from it, so displayed and copied code cannot diverge (they did: the markup showed a 3-arg `Write` that does not compile). Edit the token list, never the markup in `HeroSection.astro`. Tailwind scans this `.ts` file for the syntax classes (`text-amber`, `text-code-comment`, ...), so keep those as literal strings.
- **Copy buttons are declarative and delegation-based** — `website/public/js/copy-code.js` installs one `document` click listener guarded by `document.documentElement.dataset.copyBound`, so it works for every `[data-copy]` element (including ones added later) and cannot double-bind across ClientRouter navigations. A new copyable snippet needs only `data-copy` plus `data-code={...}`. Success/failure is announced via the `#copy-status` live region in `LandingLayout.astro`, not by the button label alone.
- **The module path is a single constant** — `siteConfig.modulePath` in `website/src/data/config.ts` is the only place the module path is written; `github`, `installCommand`, and `pkgGoDev` derive from it, `hero-code.ts` imports it for the sample's import path, and `astro.config.mjs` imports `siteConfig.description` for the Starlight docs meta description. Do not rebuild these by string-replacing `siteConfig.github`.
- **`postbuild` fails the build on landing-page drift** — `website/scripts/verify-build.mjs` asserts that the hero code rendered in `dist/index.html` equals the payload its copy button carries, that the install command matches the `go.mod` module path, and that every `[data-copy]` button has a non-empty `data-code`. If a copy/label regex stops matching, the build fails loudly rather than shipping a stale gate.
- **`pnpm run validate`** runs `html-validate dist/index.html` (config in `website/.htmlvalidate.json`); the Website CI workflow runs it after `pnpm run build`.
- **The `i18n` Starlight collection is intentionally empty** — `website/src/content/i18n/en.json` contains `{}` and is registered in `src/content.config.ts` only to silence Starlight's `The collection "i18n" does not exist or is empty` build warning. Starlight falls back to its built-in UI strings, so there is nothing to translate. Deleting the file brings the warning back.
- **`website/src/styles/global.out.css` was a stale generated artifact** — removed in 0.5.x; it had no references. Do not recreate a `*.out.css` in `src/styles/`; Tailwind v4 runs through the Vite plugin, not a CLI `-o` output.
- **No linter sub-module (decided 2026-11-03)** — evaluated against `go-branded-id/linter` (brandid-lint) and rejected. The only candidate fleet-wide rule ("use this library instead of hand-rolled `os.CreateTemp`+rename") has no syntactic signal and an unstateable false-positive budget; API-misuse rules are already covered by the compiler or `errorlint`/`errcheck`. Consumer-side bans (e.g. CRM's `internal/` CreateTemp+rename ban, CRM TODO #32) belong in the consumer's golangci config or an arch-test in the consumer repo. brandid-lint earned its sub-module by porting a proven rule with a real corpus — nothing equivalent exists here. Revisit only if ≥2 consumers accumulate a real misuse corpus with a nameable two-signal detection.
