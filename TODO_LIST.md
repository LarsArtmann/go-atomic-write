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

## High Impact

| Task                                                                                          | Status    | Impact | Effort | Evidence                                                                                |
| --------------------------------------------------------------------------------------------- | --------- | ------ | ------ | --------------------------------------------------------------------------------------- |
| Remove dead website code: `comparisons` array, `ComparisonItem` interface, `fade-in-up` anim | 🔴 `TODO` | High   | 20min  | `website/src/data/sections.ts:34`, `website/src/data/types.ts:27`, `website/src/styles/global.css:39,95` — none referenced by components |
| Fix pulse-dot ignoring `prefers-reduced-motion`                                               | 🔴 `TODO` | High   | 10min  | `website/src/styles/global.css:123` — block targets only `[data-animate]`, not `.animate-pulse-dot` (accessibility: vestibular disorders) |
| Create CI/CD pipelines (Go CI + website deploy)                                               | 🔴 `TODO` | High   | 3h     | No `.github/workflows/` directory exists (verified)                                     |
| Re-add Content Security Policy to website build                                               | 🔴 `TODO` | High   | 2h     | `website/astro.config.mjs` has no CSP (verified); needs `scripts/fix-csp.mjs` patcher; build script update |
| Fix the recurring `404 was not found` build warning                                           | 🔴 `TODO` | High   | 30min  | No `website/src/pages/404.astro` exists; warning emitted every build                     |
| Keep AGENTS.md Go version in sync with `go.mod`                                               | 🔴 `TODO` | High   | 5min   | `go.mod` is `1.26.5`; `AGENTS.md` says `1.26.4` (drift recurs — was 1.26.3→1.26.4 in `9808ab1`) |

## Medium Impact

| Task                                                       | Status    | Impact | Effort | Evidence                                                                              |
| ---------------------------------------------------------- | --------- | ------ | ------ | ------------------------------------------------------------------------------------- |
| Add OG / social-share image generation                     | 🔴 `TODO` | Med    | 2h     | No `astro-og-canvas` dep, no `og:image`, no `website/src/pages/og/` endpoint          |
| Add `favicon.ico` for legacy browser support               | 🔴 `TODO` | Med    | 15min  | Only `website/public/favicon.svg` exists (verified)                                    |
| Add `aria-label`s to comparison matrix cells + WCAG audit  | 🔴 `TODO` | Med    | 1h     | `website/src/components/ComparisonSection.astro` — SVG check/X icons lack text alternatives; `text-muted` (`#78716c` on `#0a0908`) borderline AA |
| Run Lighthouse audit + add `lighthouserc.json` budgets     | 🔴 `TODO` | Med    | 1h     | No `website/lighthouserc.json`; no prior performance/accessibility scores             |
| Add `website/.editorconfig` for editor consistency         | 🔴 `TODO` | Med    | 5min   | No `website/.editorconfig` (verified); root `.editorconfig` exists                    |

## Low Impact

| Task                                                       | Status    | Impact | Effort | Evidence                                                                              |
| ---------------------------------------------------------- | --------- | ------ | ------ | ------------------------------------------------------------------------------------- |
| Verify sitemap + add `robots.txt` sitemap reference        | 🔴 `TODO` | Low    | 15min  | `website/public/robots.txt` exists; verify all 11 pages in sitemap                    |
| Add `firebase.json` long-term cache headers (CSS/JS)       | 🔴 `TODO` | Low    | 20min  | `website/firebase.json` has security headers but no 1-year immutable cache for assets |

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
