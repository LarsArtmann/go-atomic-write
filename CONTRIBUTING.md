# Contributing

Thanks for your interest in contributing to `go-atomic-write`!

## How to Contribute

1. Fork the repository
2. Create a feature branch from `master`
3. Make your changes
4. Ensure tests, lint, and (if touching the website) the build pass (see below)
5. Submit a pull request

## Library Development Setup

The Go library needs only the plain Go toolchain — no Makefile required.

```bash
go test ./...                       # Run all tests
go test -race ./...                 # Run with race detector
go vet ./...                        # Static analysis
go build ./...                      # Verify compilation
golangci-lint run ./...             # Gating linter (~100 linters, strict) — MUST be 0 issues
go test -bench=. -benchmem          # Run benchmarks
```

> `golangci-lint` (config in `.golangci.yml`) is the **gating** quality check.
> A single issue blocks merging. New third-party imports must be added to
> `depguard.rules.main.allow` in `.golangci.yml`.

### Testing Conventions

- Standard `testing` package — no testify or other test utilities
- Use `t.Parallel()` for isolated tests, `t.TempDir()` for filesystem isolation
- Use `t.Helper()` for test helpers like `tempFile(t, content)`
- Concurrent tests should exercise the actual `flock` contention path, not just goroutine races

## Website Development Setup

The marketing/docs site lives in `website/` and is built with Astro + Starlight +
Tailwind v4, deployed to Firebase Hosting. Node.js 24 is required
(use `nix shell nixpkgs#nodejs_24` if it is not in your PATH).

```bash
cd website
pnpm install                       # First-time setup
pnpm run dev                       # Local dev server
pnpm run build                     # Production build (regenerates the changelog page first)
pnpm run typecheck                 # TypeScript + Astro type checking
pnpm run preview                   # Preview the production build locally
```

### Website notes

- **Changelog page is generated.** `website/src/content/docs/changelog.mdx` is produced
  from the repo-root `CHANGELOG.md` by `website/scripts/sync-changelog.mjs`, which runs
  automatically via the `prebuild`/`predev` hooks. Edit `CHANGELOG.md` — never the
  generated `.mdx`.
- **Deploy** with `firebase deploy --only hosting --project lars-software` (from `website/`).

## Reporting Issues

Please use [GitHub Issues](https://github.com/larsartmann/go-atomic-write/issues) to report bugs or request features.
