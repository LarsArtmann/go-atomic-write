# Roadmap

> Long-term direction and raw ideas. Items here are NOT actionable tasks.
> When an idea is refined into bounded work, it moves to TODO_LIST.md.

## Themes

### 1. Hardening & cross-platform correctness

The core write path is solid on POSIX. The remaining durability and concurrency
questions are edge platforms and crash recovery.

Raw ideas:

- Real-Windows testing of `rename_windows.go` retry loop and `FlushFileBuffers`
  durability (currently build-tag-compiled, never run on Windows hardware)
- Stale `.tmp` file sweeper — crashed writes leave uniquely-named temp files;
  decide whether a discovery/cleanup helper belongs in the library or the caller
- Fuzz testing of the fingerprint + commit path to harden against pathological inputs
- NFS / network-filesystem durability caveats — document and, where possible, test

### 2. Documentation depth

The website ships 9 reference pages. The next layer is the "why" and "how-to"
content that turns a reference into a learning resource.

Raw ideas:

- "Migration from `os.WriteFile`" guide with before/after patterns
- "Understanding TOCTOU races" deep-dive essay
- "Why xxhash64?" design rationale page
- "Design decisions" as a dedicated docs page (currently a README section)

### 3. Website maturity

The landing page and docs are live and functional. Maturity means SEO,
performance budgets, social proof, and visual richness.

Raw ideas:

- Visual polish: alternating section backgrounds, full-width stat bar, mobile how-it-works flow connector, comparison-matrix hover-expand + clearer partial glyph, animated write-pipeline diagram, animated swap arrows on scroll
- Code syntax highlighting — deliberately declined twice (manual token list is the single source shared with the copy payload); revisit only if a third surface needs highlighting
- Surface `WriteIfChanged`/`WriteFunc`/`WriteWithPerm` on the landing page; lead the hero metrics row with user benefit instead of raw throughput
- Demo material: short problem/solution video or animated GIF; persona section ("who is this for?")
- Dependents page (GitHub code search for importers) once the library has users
- Social proof: testimonials, "used by" once adoption exists; GitHub Sponsors/FUNDING if ever requested
- Promotion push: the library had zero importers until the 2026-11-03 enumeration (22 internal repos) — public promotion (HN/Reddit showcase, profile README) still untouched
- Benchmark dashboard (à la GitHub Pages) for tracking xxhash64 vs SHA-256 over time; benchmark numbers inside the comparison matrix
- Search Console submission + custom-domain canonicalization follow-through
- RSS/Atom feed, print stylesheet for docs, accessibility statement page, prefetch strategy tuning
- `siteConfig` JSON projection so build scripts consume the module path without importing TypeScript
- DNS: `atomicwrite.lars.software` and `go-atomic-write.lars.software` custom
  domains (external dependency — Namecheap DNS via Terraform)

### 4. Toward v1.0.0 stability

The API surface (`Write`, `WriteVerified`, `WriteIfChanged`, `WriteFunc*`,
`Fingerprint`) has settled. The path to a v1.0.0 compatibility commitment.

Raw ideas:

- Lock the public API after real-world usage feedback
- API-surface compat test (golden signatures snapshot) so breaking changes are explicit, not accidental
- Document the v1.0.0 compatibility guarantee and what counts as a breaking change
- Release automation (`release.yml` tag-based GitHub releases) before cutting v1.0.0
- Decide the release convention: keep tagging 0.x as GitHub `--prerelease` (skill convention) or switch to stable releases at v1.0

## Non-goals

Things we are deliberately NOT pursuing and why:

- **General file management:** this library does one thing — atomic, race-free,
  crash-durable writes. Copy, move, transactional multi-file writes are out of scope.
- **Analytics/tracking:** no telemetry on the website. Intentional privacy stance
  (gogenfilter removed Plausible for the same reason).
- **Test framework dependencies:** testing stays on the standard `testing` package.
  No testify or assertion libraries — keeps the dependency surface at two.
- **Cryptographic hashing:** fingerprints detect changes, not attackers. xxhash64
  stays; SHA-256 is benchmark-only. Not a security boundary and never will be.

---

<!-- Guidance for the builder:
  - NO bounded actionable tasks here. If it has a clear scope and effort
    estimate, it belongs in TODO_LIST.md.
  - NO status indicators on individual items. This is vision, not inventory.
  - Ideas should be raw and unrefined by design.
  - Non-goals are as important as goals: they prevent scope creep.
  - Revisit quarterly to prune stale directions.
-->
