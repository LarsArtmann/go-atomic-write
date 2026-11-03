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

> All 15 items from the 2026-09-17 harvest (report §f) were cleared on
> 2026-11-03 — CI build matrix, govulncheck, drift guard, GOTOOLCHAIN pin,
> version-trio gate, Dependabot npm root-ecosystem fix, npm vulnerability sweep
> (62 → 1 advisory), README benchmark refresh (ranges), release checklist,
> Node-24 action verification, Lighthouse against the deployed site, and the
> `go 1.27.1` investigation (fleet `go-version-auto-configure`, ADR-0001).
> See `CHANGELOG.md` `[Unreleased]`.

## High Impact

| Task                                                                                                                         | Status  | Impact | Effort | Evidence                                                   |
| ---------------------------------------------------------------------------------------------------------------------------- | ------- | ------ | ------ | ---------------------------------------------------------- |
| Cut release `v0.6.0` from `[Unreleased]` (temp-file leak fix, `go 1.27`, TS6 self-heal) — the fleet bump below depends on it | 🔴 TODO | High   | S      | `CHANGELOG.md` `[Unreleased]`; AGENTS.md Release Checklist |

## Medium Impact

| Task                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | Status  | Impact | Effort | Evidence                                                                                                         |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ | ------ | ---------------------------------------------------------------------------------------------------------------- |
| Fleet bump: 22 LarsArtmann repos import `go-atomic-write` (enumerated 2026-11-03 via `gh search code`, Go imports only: crm, BuildFlow, CV, .hermes, mr-sync, go-auto-upgrade, go-version-auto-configure, golangci-lint-auto-configure, samber-do-auditlog, go-workflow-auditlog, project-meta, library-policy, index, clean-wizard, ArtCode, securitymd, linter-autoconfigure-sdk, DiscordSync, dependabot-auto-configure, media-atlas, samber-linter, go-checker-helpers) — bump to the new release via the go-ecosystem-upgrade flow | 🔴 TODO | Medium | L      | `gh search code '"github.com/larsartmann/go-atomic-write"' --language go --owner LarsArtmann`                    |
| Find and configure the fleet npm sweep that force-bumps devDeps to latest majors — it re-bumped `typescript` to 7.x three times on 11-03 (breaks `astro check`) and regenerated `global.out.css` twice; the workspace-yaml override neutralizes the TS damage, but the sweep should learn the constraint or exclude this repo                                                                                                                                                                                                           | 🔴 TODO | Medium | M      | commits 3cf0176 + package.json mtimes 21:01/21:38/21:59 on 2026-11-03; `website/pnpm-workspace.yaml` `overrides` |

## Low Impact

| Task                                                                                                                                                                                                        | Status     | Impact | Effort | Evidence                                                                                    |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- | ------ | ------ | ------------------------------------------------------------------------------------------- |
| `minimatch@3.1.5` moderate advisory (glob@7 ← rimraf@3 ← chrome-launcher ← @lhci/cli): overriding to 11 breaks glob@7's CJS API at runtime — wait for rimraf/chrome-launcher upstream majors, then re-audit | 🔵 BLOCKED | Low    | S      | `website/pnpm-workspace.yaml` `overrides`; audit 2026-11-03 (1 advisory left)               |
| The deployed site serves the previous content until the next `Website` workflow deploy — verify `~23 GB/s` claims + changelog go live after push                                                            | 🔴 TODO    | Low    | S      | `https://atomicwrite.web.app` (fetched 2026-11-03); `pnpm run lighthouse:deployed` re-check |

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
