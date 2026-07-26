# AGENTS.md — go-atomic-write

## Project

Single-package Go library providing TOCTOU-safe file writes via xxhash64 fingerprint verification, cross-platform file locking (`flock`/`LockFileEx`), atomic rename, and fsync for crash durability.

- **Module:** `github.com/larsartmann/go-atomic-write`
- **Go version:** 1.26.5
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

`golangci-lint` config is the **gating quality check** of this project. The configuration lives in `.golangci.yml` and enables ~100 linters at their strictest defaults. Adding a new third-party import requires updating the `depguard.rules.main.allow` list.

No Makefile, no CI config. All commands are plain `go` toolchain.

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

| Command                                                                | Purpose                          |
| ---------------------------------------------------------------------- | -------------------------------- |
| `cd website && npm run dev`                                            | Local dev server                 |
| `cd website && npm run build`                                          | Production build to `dist/` (`prebuild` syncs changelog, `postbuild` injects CSP) |
| `cd website && npm run typecheck`                                      | TypeScript + Astro type checking |
| `cd website && npm run preview`                                        | Preview production build locally |
| `cd website && npm run sync:changelog`                                 | Regenerate changelog page from root `CHANGELOG.md` |
| `cd website && firebase deploy --only hosting --project lars-software` | Deploy to Firebase               |

Node.js 24 required (use `nix shell nixpkgs#nodejs_24` if not in PATH).

### Website Structure

| Path                            | Purpose                                                               |
| ------------------------------- | --------------------------------------------------------------------- |
| `website/astro.config.mjs`      | Astro config: Starlight, fonts, sitemap                               |
| `website/src/pages/index.astro` | Landing page                                                          |
| `website/src/components/`       | 14 Astro components (Hero, FeatureGrid, HowItWorks, Comparison, etc.) |
| `website/src/data/`             | Typed content: config, features, sections, hero-code                  |
| `website/src/content/docs/`     | 9 Starlight documentation pages                                       |
| `website/src/styles/`           | global.css (emerald theme) + starlight.css                            |
| `website/public/`               | favicon, manifest, robots.txt, JS (theme, animations, copy-code)      |
| `website/firebase.json`         | Hosting config with security headers                                  |
| `website/.firebaserc`           | Firebase project + hosting target                                     |
| `website/scripts/`              | Build-time tooling: `sync-changelog.mjs` (prebuild), `fix-csp.mjs` (postbuild) |
| `.github/workflows/`            | CI (`ci.yml`: Go gate) + website build/deploy (`website.yml`)         |

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

## Gotchas

- **Keep the Go version above in sync with `go.mod`** — this drift has recurred three times (1.26.3 → 1.26.4 → 1.26.5); whenever `go.mod`'s `go` directive bumps, update the `**Go version:**` line in the Project section above
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
- The `//nolint:gosec` comments on `os.ReadFile`/`os.OpenFile`/`os.Open` calls are intentional — `path` is caller-controlled, not user input
- The `//nolint:gosec` comments on `os.ReadFile`/`os.WriteFile` calls in `atomicwrite_test.go` are intentional — all paths are `t.TempDir()`-rooted and never user input
- Lint is gating: a single `golangci-lint` issue blocks merging. Common pitfalls:
  - Single-letter variables (`fp`, `wg`, `h`, `mb`) trigger `varnamelen` — use full words
  - Inline `if err := …; err != nil` triggers `noinlineerr` — assign first, then check
  - `b.Helper()` is required on every benchmark helper that takes `*testing.B` (thelper)
  - 0o644 in tests triggers `gosec G306` — annotate with `//nolint:gosec` and rationale
  - `os.ReadFile(path)`/`os.OpenFile` with a variable path triggers `gosec G304` — annotate with `//nolint:gosec` and rationale
  - Blank line between error assignment and `if err != nil` triggers `wsl_v5` — put the blank line above the assignment instead
  - New third-party imports require updating `.golangci.yml` `depguard.rules.main.allow`
- **Website changelog page is GENERATED** — `website/src/content/docs/changelog.mdx` is produced from root `CHANGELOG.md` by `scripts/sync-changelog.mjs` (runs on `prebuild`/`predev`). Edit `CHANGELOG.md`, never the `.mdx`. If the build breaks with an MDX "unexpected character" error, a changelog entry likely contains a raw `<`; the sync script escapes `<` but check new entries.
- **Website CSP is hash-based and post-build** — `scripts/fix-csp.mjs` (`postbuild`) injects a per-file CSP `<meta>` from inline-script SHA-256 hashes. There is **no `'unsafe-inline'` for `script-src`**. If a new inline script breaks the site under CSP, either move it to an external `/js/*.js` file (loaded via `<script is:inline src=…>`) or confirm `fix-csp.mjs` hashed it.
- **Do NOT add `website/src/pages/404.astro`** — Starlight ships its own `404.html`. A custom one causes a route collision (warning today, hard error in future Astro). The `404 was not found` line during build is a benign Starlight route log, not a warning.
- The website's `npm run build` requires a clean `.astro` cache when content files are renamed/extension-changed: run `rm -rf .astro dist node_modules/.cache` if the content layer fails to resolve a renamed doc.
