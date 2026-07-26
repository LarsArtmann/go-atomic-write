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

> The previously-High items (dead website code, pulse-dot reduced-motion, CSP,
> CI/CD pipelines, the 404 warning) were completed on 2026-07-26 — see
> `CHANGELOG.md` `[Unreleased]`. The recurring `404 was not found` build line is a
> **benign Starlight route-generation log** (Starlight ships its own `404.html`);
> adding a custom `404.astro` causes a route collision and is NOT the fix.

## High Impact

| Task                                                        | Status | Impact | Effort | Evidence |
| ----------------------------------------------------------- | ------ | ------ | ------ | -------- |
| _(none open — previously-High items cleared on 2026-07-26)_ |        |        |        |          |

## Medium Impact

| Task                                                      | Status    | Impact | Effort | Evidence                                                                                                                                         |
| --------------------------------------------------------- | --------- | ------ | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| Add OG / social-share image generation                    | 🔴 `TODO` | Med    | 2h     | No `astro-og-canvas` dep, no `og:image`, no `website/src/pages/og/` endpoint                                                                     |
| Add `favicon.ico` for legacy browser support              | 🔴 `TODO` | Med    | 15min  | Only `website/public/favicon.svg` exists (verified)                                                                                              |
| Add `aria-label`s to comparison matrix cells + WCAG audit | 🔴 `TODO` | Med    | 1h     | `website/src/components/ComparisonSection.astro` — SVG check/X icons lack text alternatives; `text-muted` (`#78716c` on `#0a0908`) borderline AA |
| Run Lighthouse audit + add `lighthouserc.json` budgets    | 🔴 `TODO` | Med    | 1h     | No `website/lighthouserc.json`; no prior performance/accessibility scores                                                                        |
| Add `website/.editorconfig` for editor consistency        | 🔴 `TODO` | Med    | 5min   | No `website/.editorconfig` (verified); root `.editorconfig` exists                                                                               |

## Low Impact

| Task                                                        | Status    | Impact | Effort | Evidence                                                                                       |
| ----------------------------------------------------------- | --------- | ------ | ------ | ---------------------------------------------------------------------------------------------- |
| Verify sitemap + add `robots.txt` sitemap reference         | 🔴 `TODO` | Low    | 15min  | `website/public/robots.txt` exists; verify all 11 pages in sitemap                             |
| `website/flake.nix`: pin `nodejs_24` (not `pkgs.nodejs`)    | 🔴 `TODO` | Low    | 5min   | `website/flake.nix:43-57` uses `pkgs.nodejs` (Node 22 on unstable); AGENTS.md requires Node 24 |
| `website/flake.nix`: add `meta.description` to `deploy` app | 🔴 `TODO` | Low    | 2min   | `nix flake check` warns: `app 'apps.x86_64-linux.deploy' lacks attribute 'meta.description'`   |

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
