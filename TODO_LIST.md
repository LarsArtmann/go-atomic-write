# TODO List

> Short-term, actionable, bounded work items, verified against the actual code.
> For long-term vision and unrefined ideas, use ROADMAP.md.
> Items are ranked by impact. Status is verified, not assumed.

## Status legend

| Status           | Meaning                                                     |
| ---------------- | ----------------------------------------------------------- |
| 🔴 `TODO`        | Not started. Needs doing.                                   |
| 🟡 `IN_PROGRESS` | Actively being worked on.                                   |
| 🔵 `BLOCKED`     | Cannot proceed, external dependency or decision needed.     |
| 🟢 `DONE`        | Completed. Remove from this list and log in `CHANGELOG.md`. |

> Rebuilt 2026-11-03 from a full docs-health sweep; re-harvested 2026-11-04
> after the release/quality execution session (v0.6.1 cut, 17-repo fleet bump,
> CSP-quoting unblock, contrast audit, Lighthouse 100s, gate hardening, Go
> edge tests — evidence in `docs/status/2026-11-04_05-36…` and `CHANGELOG.md`).
> The rows below are the surviving open work; Lane 4–5 batches execute from
> `docs/planning/2026-11-04_02-12_release-fleet-and-site-quality-pareto-plan.md`.

## Medium Impact

| Task                                                                                                                                                                                              | Status  | Impact | Effort | Evidence                                                                                         |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ | ------ | ------------------------------------------------------------------------------------------------ |
| Brand system batch (remainder): `brand:png` script (rasters in code, not README), SVG well-formedness parse inside `brand:check`, `brand-tokens.json` single palette source for generator + `global.css`, real-Space-Grotesk wordmark width measure, footer-logo screenshot QA, scripted pixel-QA (glyph-ink assertion) — the `aria-hidden` decision already landed (generator emits it) | 🔴 TODO | Medium | M      | `scripts/generate-brand-assets.mjs`; `docs/status/2026-11-03_23-29…` §2.3–2.5, §5.2–5.7          |
| Fleet npm sweep tracker: find/configure the sweep that force-bumps devDeps to latest majors — it re-bumped `typescript` to 7.x three times on 11-03 (breaks `astro check`) and regenerated `global.out.css` twice; the workspace-yaml override neutralizes the TS damage, but the sweep should learn the constraint or exclude this repo | 🔴 TODO | Medium | M      | commits 3cf0176 + package.json mtimes 21:01/21:38/21:59 on 2026-11-03; `website/pnpm-workspace.yaml` `overrides` |
| Interactive navigation verification: ClientRouter nav re-bind (`[data-copy]`, theme persistence, `astro:page-load` animations) + keyboard tab-order pass — verified by construction only so far; Playwright smoke suite closes it properly (plan T32) | 🔴 TODO | Medium | M      | `animations.js` `41f636e`; status report §b.1; plan T32                                          |

## Low Impact

| Task                                                                                                                                                        | Status     | Impact | Effort | Evidence                                                                      |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- | ------ | ------ | ----------------------------------------------------------------------------- |
| `minimatch@3.1.5` moderate advisory (glob@7 ← rimraf@3 ← chrome-launcher ← @lhci/cli): overriding to 11 breaks glob@7's CJS API at runtime — wait for rimraf/chrome-launcher upstream majors, then re-audit | 🔵 BLOCKED | Low    | S      | `website/pnpm-workspace.yaml` `overrides`; audit 2026-11-03 (1 advisory left) |
| OG card in real previewers (Slack/Twitter/Facebook — needs real clients, owner action) + post-deploy fetch assertion in the Website workflow (real-client subscribe test and iOS icon check are owner actions too) | 🔴 TODO    | Low    | S      | live OG verified at artifact level 2026-11-04 (28.3 KB, quoted-hash CSP live) |
| Docs/SEO hygiene batch: sitemap `lastmod`, `go get` install command in every docs entry point, `theme-color` light variant, `og:image:type`/`og:site_name`/`og:image:alt`+`twitter:image:alt`, `engines.node` field, `website/flake.lock` update, internal link checker for docs + README | 🔴 TODO    | Low    | M      | `astro.config.mjs` (no lastmod); grep 2026-11-03 (no og: extras, no link checker) |
| `CopyButton.astro` extraction — the ~300-char Tailwind class string is repeated per copyable snippet                                                         | 🔴 TODO    | Low    | S      | `HeroSection.astro`, docs pages repeat the button markup                      |
| `gitleaks` step in CI (secret scan; skipped in release reviews twice)                                                                                        | 🔴 TODO    | Low    | S      | no gitleaks anywhere in `.github/workflows/` (grep 2026-11-03)                |
| Deploy-secret posture: decide workload-identity federation vs SA-JSON rotation cadence; delete the 7 stale `firebase-adminsdk` keys (GCP console, owner action) | 🔴 TODO    | Low    | S      | `docs/status/archived/2026-09-17_20-57…` §a.11, §g.1; AGENTS.md SA-cap gotcha |
| CONTRIBUTING: release-cutting section for contributors + deploy/secret-minting recipe + "brand artifacts are generated — never hand-edit" note              | 🔴 TODO    | Low    | S      | `CONTRIBUTING.md` has no release/brand content (grep 2026-11-03)              |

---

<!-- Guidance for the builder filling this in:
  - Source of truth is the CODE. Verify each item before adding, many
    documented TODOs are already done.
  - One task per row. If it takes more than ~2 hours, split it into smaller
    tasks.
  - Cite evidence (file:line) so the next person can verify without re-deriving.
  - DONE items should be REMOVED, not kept. Use CHANGELOG.md for history.
  - If a task is vague ("improve X"), refine it into concrete steps or move
    it to ROADMAP.md.
  - Deduplicate by semantic intent, not by text match.
  - For 80/20 impact prioritization, use the pareto-planning skill AFTER
    building the list here.
-->
