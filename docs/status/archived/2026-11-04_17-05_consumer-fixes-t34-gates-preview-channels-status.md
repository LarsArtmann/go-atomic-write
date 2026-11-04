# Status: PR preview channels live, docs gates green, consumer fixes (2026-11-04 17:05)

> **Archived 2026-11-04 (cadence sweep).** Live remainder: none — T32–T34 closed, PR preview channels live (merged `011e0bc`), all four consumer fixes landed (samber#62 landed the same day; crm baseline moved to crm#7).

Session resumed from the 07:09 WAIT handoff with the mandate "NOW GET SHIT DONE! The WHOLE TODO LIST!". This report covers T32 completion (smoke suite + CI + PR previews), T33 (brand rasters), T34 (annotation sweep), the AGENTS.md memory pass, and 4 consumer fixes (3 landed, 1 mid-flight). Written **before** the v0.6.2 release + fleet bump, which are the next block.

## a) FULLY DONE

1. **T32 — Playwright smoke suite + CI job + PR preview channels (end-to-end verified)**. Config fixed (misindented `projects` line), `test:ui` script added. Local run needed two iterations: Playwright's downloaded chromium cannot run on NixOS (missing FHS libs) → `launchOptions.executablePath` from `CHROME_PATH` (nix chromium); the first attempt used a bare `use.executablePath`, which is **silently ignored** (cost one all-red local run). CI needed two iterations too: Playwright 1.56 **does not support ubuntu-26.04 runners** → smoke job now uses `browser-actions/setup-chrome` + `channel: "chrome"`; the preview job failed `Resource not accessible by integration` → added `checks: write`. **PR #22** (changelog batch as the live exercise) went fully green — build, smoke, Lighthouse budgets (the T31 landing-growth worry is resolved: budgets PASS), and the preview channel **actually deployed** (`https://atomicwrite--pr-22-pe8lu3zc.web.app`, fetched and serving the full landing page). Merged as `011e0bc`.
2. **Changelog batch** — `[Unreleased]` now documents the whole 11-03/11-04 website batch (T20–T32: gitleaks, regression tests, link checker, CopyButton, brand:png + tokens, when-not-to-use + FAQ, SEO hygiene, config-reloader use case, CONTRIBUTING/deploy docs, npm-sweep tracker, smoke suite + previews), correctly ordered sections, verified through a green build (MDX `<` escaping works).
3. **T33 — brand rasters + social preview + contact sheet + usage rules** (`a6fbb9e`). `brand:png` now also writes the five compositions as 2× PNGs (`assets/brand/png/`), the GitHub social preview (`social-preview.png`, 1280×640, og-image scaled to height 640 and padded on the card background — seamless), and a labeled `contact-sheet.png` (montage; labels need `-font DejaVu-Sans` explicitly and a single global `-fill` — per-cell fill recolors labels). Brand README gained the new Files rows, the **Usage rules** section (clear space = 25% mark height, min sizes 16/150/170 px, composition-by-slot, tokenized color only, variant-by-background), and the social-preview upload recipe (owner action, no API). `brand:check` still green (9 artifacts, 4 token mirrors).
4. **T34 — annotation sweep, battery ALL FIVE GATES GREEN** (`993a04f` + daemon states). docs-gates found 6 unharvested items + 8 annotate findings; all resolved: 8 verdict/path repairs (CONTRIBUTING spacing, `aria-label` plural, two `02-06` verdict markers, `0709` release.yml-route + FUNDING.yml won't-implement, `0918` build-gate done at `c5f36ff`/`939738b`/`bb3c47b`, `06-03` canonical NOT-DO); **43 items annotated in the 05:36 §f** + 4 in §b with real hashes (T20–T33 chain, live checks); `22-20` items 7/31/34, `00-14` row 12, `02-06` rows 38/39, `20-42` row 9 routed. TODO_LIST rebuilt (6 stale done rows deleted; routed rows added: keyboard tab-order, brand pixel-QA remainder, owner-actions bundle, deploy-secret WIF follow-through, Windows test run, visual-regression harness, seeded collision test, CV chronic red, lighthouse:deployed). ROADMAP pruned (5 shipped ideas removed). `.hermes` ghost reconciled in the plan doc (2 rows now point at the 06-15 fleet summary). Status index refreshed (+5 missing reports). Final `docs-gates run`: **freshness PASS, audit-all PASS, annotate PASS, routing PASS, cadence PASS**.
5. **Session verifications** — `go get @latest` resolves v0.6.1; `sum.golang.org` carries the v0.6.1 checksum (`h1:lL4arpBDRh70H1ktdS37bu4VswO5PekGgjfynwqcbrM=`); verify-build's FAQ assertion intact (landing 6, docs FAQ separate); post-release Dependabot = actions bumps only (PR #21 gitleaks-action, PR #7 svgo).
6. **AGENTS.md memory pass** (daemon-committed + pushed) — 8 new gotchas (Playwright env-split launch + the `use.executablePath` silent-ignore trap, preview `checks: write` + artifact-on-PR, docs-gates battery + annotate tooling + dry-run rule, `BUILDFLOW_NO_RESULT_CACHE=1`, chromium screenshot recipe, sitemap `lastmod` Date, brand:png font/label gotcha, npm-sweep guards), command table (+`test:ui`, `check:links`, `test:scripts`, docs-gates), structure table (+`smoke.spec.mjs`, `brand-tokens.json`, `fleet-bump.sh`, workflow descriptions).
7. **Consumer fixes landed (3 of 4)**:
   - **go-checker-helpers#3** (`d3f4dbe`, pushed) — reproduced, root cause: `lo.Filter` always allocates so no-match/empty-input returned non-nil; `FilterByRule` now returns nil for both per the tested contract. `go test ./...` + vet green.
   - **DiscordSync#18** (`487012f9`, pushed) — reproduced; `mime.TypeByExtension` defers to the host mime.types (bare CI images lack `.mov`); `mime.AddExtensionType(".mov", "video/quicktime")` pins the mapping in code. Package tests green, build green.
   - **project-meta#3** — **does not reproduce anymore** (build/vet/test green on current master, `go mod tidy` a no-op, clean worktree, go1.27.0). Closed as obsolete with the evidence in the closing comment.

## b) PARTIALLY DONE

1. **samber-do-auditlog#62 — undici override landed, gate still RED**. The `undici: ^8.12.2` override + lockfile regen landed (daemon race took the commit as `d337743`); the 11 undici error findings are gone from `website/`. But the pre-commit gate still fails: BuildFlow's pnpm-audit step audits **every npm tree in the repo**, and this repo has TWO lockfile-less npm trees (`website/video/`, `scripts/e2e-live-smoke/`) whose fresh-resolution audits produce high-severity findings the gate maps to errors. I generated their first lockfiles (daemon committed as `a662a26`/`923f844`); audits now read: website 5 (none error), **video 29 (10 high)**, **e2e-live-smoke 4 (4 high)**. The gate also ran against the root `npm init` stub — removed it AND its `npm directory: /` dependabot entry together (the exact fleet incident class from go-atomic-write 2026-09-22); that commit FAILED the gate twice before I understood the multi-tree scope, so it is **not yet pushed as a commit of mine** — the deletions are in the tree, the daemon has the state, and the gate-greening work (overrides/bumps for zod in website, the video/e2e trees) is the open remainder.
2. **05:36 report annotations** — §b2 (live-site verify post-deploy), §b7 (fleet retry), §f25 (consumer fixes: 2 of 4 closed), §f36 (BuildFlow cache lesson — the AGENTS.md half is done; the annotation hash waits for the final sweep), §f37/§f42/§f46 deliberately left open.
3. **v0.6.2 release** — everything batched in `[Unreleased]`, but T33's entry is NOT yet written, the version heading not moved, no tag. Release checklist step 1 (lint gate uncached) not yet re-run today.

## c) NOT STARTED

1. **v0.6.2 cut** (move `[Unreleased]` → `[0.6.2]`, CI green before tag, annotated tag, push, pkg.go.dev + proxy verification).
2. **Fleet bump v0.6.2** via `scripts/fleet-bump.sh` (17 repos + retry crm/BuildFlow/index/clean-wizard; SKIP_DIRTY guard).
3. **Final sweep** (full battery: `go test -race`, lint gate uncached, `pnpm run check`, docs-gates, live-site fetch, TODO_LIST/CHANGELOG final state).
4. **07:09 report annotation** (this session's own §f queue).
5. samber-do-auditlog pnpm-version question: the two new lockfiles were generated with **pnpm 12.9.0** (no corepack pin in that repo's trees) — CI compatibility unverified.

## d) TOTALLY FUCKED UP

1. **The undici commit rode a STALE BuildFlow result cache into the gate** — the pre-commit hook's pnpm-audit was a 0 ms cache hit reporting the OLD 11-error verdict, so my first commit attempt failed against a verdict that no longer matched reality. Re-ran with `BUILDFLOW_NO_RESULT_CACHE=1` (the lesson from the 05:36 report — which I had already recorded in AGENTS.md this session — and still walked into it because I assumed the hook environment was fresh). The daemon then won the commit race, so the fix landed under a `chore: auto-commit` message with none of the context.
2. **Two wasted CI round trips on T32 that planning would have prevented** — `use.executablePath` (silently ignored option → all-red local run; the correct home is `launchOptions`) and Playwright-on-ubuntu-26.04 (unsupported runner → red CI; `setup-chrome` + `channel: "chrome"` was already the repo's proven Lighthouse pattern I should have copied from the start). Same for the preview job's `checks: write` — the action's documented permission trio, missed once.
3. **samber-do-auditlog gate declared "fixed" in my head while red** — I treated the undici override as the fix (that is what the issue asked), but the gate's acceptance (green) depends on three npm trees I hadn't audited. The lesson: run the FULL gate, not the step that matches the issue text, before claiming an issue fixed. Two failed commit attempts burned before the multi-tree scope was clear.
4. **Daemon races hit 3× this session** (AGENTS.md, T34 intermediates, the undici fix) — each cost a verify-and-supplement cycle. My write-then-commit workflow loses to the daemon whenever a verification step sits between write and commit.

## e) WHAT WE SHOULD IMPROVE

1. **Read the tool's contract before its happy path** — `use.executablePath` vs `launchOptions.executablePath`, the Firebase action's documented permission trio, Playwright runner support matrix. All three were documented; I discovered them by red runs.
2. **Gate = whole gate.** Fixing the finding named in an issue ≠ gate green. Always run the full findings gate (uncached) before pushing any "fixes #" commit.
3. **Never trust the pre-commit cache in the hook itself** — `BUILDFLOW_NO_RESULT_CACHE=1` should be the default posture for hook-invoked commits, not a recovery step.
4. **Daemon coexistence**: commit immediately after write, or `git stash`-guard long verifications; every "verify then commit" gap is a daemon race window.
5. **Check the repo's pnpm/corepack pin before generating lockfiles** — `pnpm install` with whatever is on PATH (12.9.0 there) can produce a lockfileVersion the CI's pinned pnpm cannot read (frozen-lockfile).
6. **Scope-issue vs scope-gate**: when an issue's fix surfaces NEW gate failures, note them in the issue thread immediately instead of silently expanding the change.

## f) TOP 50 NEXT (impact-sorted; 1–12 are this session's direct continuations)

1. **samber-do-auditlog gate-greening**: zod override in `website/`, triage the 10 high in `website/video/` and 4 high in `scripts/e2e-live-smoke/` (overrides or bumps), then `BUILDFLOW_NO_RESULT_CACHE=1 buildflow -s pnpm-audit` until green.
2. **Verify the two new lockfiles against that repo's CI pnpm pin** (12.9.0-generated lockfileVersion vs CI's corepack pnpm@11.20.0; regenerate with 11.20.0 if frozen-lockfile breaks).
3. **Push the stub-removal + dependabot-entry commit** (currently only in the tree/daemon commits) once the gate is green, with the full context message.
4. **Comment + close the 3 fixed consumer issues** with fix hashes (go-checker-helpers#3 `d3f4dbe`, DiscordSync#18 `487012f9`, samber#62 after gate green), per the closing formula.
5. **T33 changelog entry** (2× lockup rasters, social preview + upload recipe, contact sheet, usage rules) into `[Unreleased]`.
6. **Cut v0.6.2**: lint gate uncached → CI green → move `[Unreleased]` → `[0.6.2]` (date 2026-11-04) → annotated tag → push → verify pkg.go.dev + proxy + `go get @v0.6.2`.
7. **Fleet bump v0.6.2** with `scripts/fleet-bump.sh` — 17 repos + retry crm/BuildFlow/index/clean-wizard; SKIP_DIRTY will protect against any in-flight storm state.
8. **Live-site verification post-deploy** (changelog page shows `[0.6.2]`, OG card live, CSP metas, newsletter CSP) + annotate 05:36 §b2.
9. **`pnpm run lighthouse:deployed`** against the live site after the deploy settles (TODO_LIST row; never observed).
10. **Annotate the 07:09 report** (this session's own queue: items 1–24) + close the 05:36 leftovers (§b7, §f25) with the release/bump hashes.
11. **Final sweep**: `BUILDFLOW_NO_RESULT_CACHE=1 buildflow -s golangci-lint` (0 issues), `go test -race ./...`, `pnpm run check`, `docs-gates run .` all green, TODO_LIST/CHANGELOG/AGENTS final state.
12. **gitleaks-action v3 decision** (dependabot PR #21): verify v3 semantics against the v2 SHA pin (env vars, allowlist behavior), then merge or close with reason.
13. Dependabot PR #7 (svgo 4.1.0 in /website, open since 10-04): review, merge or close.
14. Keyboard tab-order + theme first-paint Playwright asserts (TODO_LIST Medium row).
15. Brand pixel-QA: scripted glyph-ink assertion + footer-logo screenshot QA (TODO_LIST Medium row).
16. Real-Windows test leg (or one manual Windows run) for the fingerprint/temp-path suite.
17. Visual-regression screenshot harness for header/docs/OG surfaces.
18. `tempPathFor` deterministic (seeded) collision test.
19. File the CV Go-test chronic red after checking storm-session ownership.
20. Owner actions bundle: upload `assets/brand/social-preview.png` (Settings → Social preview), OG card in real previewers, Buttondown subscribe test, iOS apple-touch check, delete the 7 stale `firebase-adminsdk` keys.
21. Deploy-secret WIF migration (Option A execution per `docs/deploy-secret-posture.md`).
22. Check merged-master CI state after `a6fbb9e`/`993a04f` (T33 + T34 pushes were not individually observed; the drift guard + Website workflow verdicts pending).
23. Verify Lighthouse budgets once more on master post-T33 (raster additions should not affect budgets; confirm).
24. `brand:check` verdict on CI for the generator change (`a6fbb9e` touched the generator; drift guard runs it — observe once).
25. Re-check `git status` of all 4 consumer repos before the fleet bump (samber is dirty-prone while its gate work continues).
26. samber-do-auditlog: the dependabot-auto-configure info finding ("pnpm /website matches nothing") — re-verify after the new lockfiles exist; possibly needs pnpm entries for video/e2e dirs or it self-resolves.
27. samber `sharp` librsvg CVE (<0.35.5): check website's sharp version; bump or override if in range.
28. TODO_LIST re-harvest after the release + consumer fixes (rows for consumer fixes were never added — they'll be DONE before ever needing rows; verify).
29. CHANGELOG: add the `docs-gates` battery + annotation-sweep entry under a docs/infra heading for [0.6.2]? (Decide: repo-facing docs work may belong in the changelog's Changed (docs) section.)
30. Confirm `website/video` + `scripts/e2e-live-smoke` lockfiles don't break their non-pnpm workflows (CI jobs referencing those dirs).
31. Consider a `pnpm-lock.yaml presence check` for every `package.json` dir in BuildFlow (the lockfile-less-tree phantom-audit root cause) — file upstream in the BuildFlow repo (verify-before-filing first).
32. Review the session's daemon commits (d337743, a662a26, 923f844, b1c75b1, a906bd2) for content correctness — spot-verified so far, not exhaustively.
33. ROADMAP: `siteConfig` JSON projection idea — unchanged, confirm still raw (no action).
34. Archive sweep trigger: 11 unarchived reports (threshold 15) — healthy, revisit at the next docs-health pass.
35. ROADMAP prune follow-up: verify no TODO_LIST/ROADMAP split-brain after the release (features shipped in 0.6.2 must not linger as raw ideas).
36. `git town` state: the squash-merged PR #22 branch was deleted; confirm no lingering town-sync state (sync after the release).
37. `website/flake.lock` bump check (fleet nixpkgs drift; last bumped 11-01).
38. Cross-check `engines.node` vs the Node 24 CI pin after the video/e2e lockfiles (their package.json engines fields, if any).
39. Consider wiring `docs-gates run .` as a repo CI job (currently manual-only; the battery is fleet-standard).
40. Consider adding `BUILDFLOW_NO_RESULT_CACHE=1` to the pre-commit hook invocation itself (repo-side hooks config) so the cache trap cannot recur for any committer.
41. Confirm `docs/status/README.md` renders correctly (table row for the html LEAVE-ALONE intact post-edit).
42. Sample-check the contact sheet + social preview on a dark GitHub profile (social preview upload preview).
43. `minimatch@3.1.5` advisory (TODO_LIST BLOCKED row) — unchanged; revisit when rimraf/chrome-launcher ship majors.
44. Verify `go.sum`/vendor state of the 4 skipped fleet repos before the bump retry (crm/BuildFlow/index/clean-wizard).
45. After the fleet bump: spot-verify 3 random consumer repos build (the AFTER_RED guard covers failures, but a spot check is cheap).
46. pkg.go.dev freshness for v0.6.2 (proxy + pkg.go.dev both, not just one).
47. Retire `/tmp/chr-prof-*` profile dirs + the gaw-verify scratch module (throwaway artifacts).
48. `AGENTS.md`: record the samber root-stub + lockfile-less-tree incident as a gotcha (it is the same class the repo already documents, now with a second data point).
49. Cross-project lesson candidate for crush-config `references/lessons.md`: "audit-scope vs issue-scope" (the undici gate incident) — commit there, never in-session global writes.
50. Re-run `docs-gates run .` after the release lands (the 06:15/07:09/this report's annotations shift freshness again).

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. **samber-do-auditlog gate scope**: the pre-commit gate is red on 3 npm trees — website (zod high), website/video (29 vulns, 10 high after I generated its FIRST lockfile), scripts/e2e-live-smoke (4 high, same). Do you want the full gate-green sweep now (overrides/bumps across all three trees), or keep #62 scoped to undici (which IS fixed), file the tree-hygiene as new issues, and accept the red gate until then? Also: the two new lockfiles were generated with pnpm 12.9.0 — that repo's website CI pins 11.20.0; should I regenerate both with 11.20.0 before pushing, or is 12 acceptable?
2. **Dependabot PR #21** bumps gitleaks-action v2 (my pinned SHA) → v3.0.0. Verify-and-merge now, or hold v2 through the v0.6.2 cut and revisit? (T20 was "gitleaks job, v2 SHA" — merging a major bump right before a release mixes concerns, but the PR is green.)
3. **Fleet bump v0.6.2 for the 4 skipped repos** (crm, BuildFlow, index, clean-wizard): retry them all with SKIP_DIRTY as the guard — or is the storm session expected to absorb some of them again (if its lanes are still open, a bump commit could collide with in-flight work the SKIP_DIRTY check can't see because the tree is momentarily clean between their commits)?
