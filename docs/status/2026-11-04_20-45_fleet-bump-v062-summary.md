# Fleet Bump Summary — go-atomic-write v0.6.2

**Created:** 2026-11-04 20:45 CET
**Session:** v0.6.2 release lane, fleet leg (continuation of the 19:35 in-flight report)
**Method:** `scripts/fleet-bump.sh <repo> v0.6.2` over the 21-repo queue (17 v0.6.1 consumers + 4 retry repos), with the daemon-commit carrier pattern: the auto-commit daemon frequently lands the `go get` result before the runner's commit step, so several bumps below cite the daemon commit that carries them. The runner is restart-safe — the first pass died silently (interrupted tool call) and the rerun handled already-bumped repos as no-ops.

## Bumped (17 repos, all verified at v0.6.2 on their canonical remote)

| Repo                         | Commit      | Notes                                                                                                                                                   |
| ---------------------------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| go-checker-helpers           | `aa95247`   | first-pass bump; runner rerun found NO_CHANGES. (The 19:35 report misattributed this hash to mr-sync.)                                                   |
| go-version-auto-configure    | `a94e836`   |                                                                                                                                                          |
| golangci-lint-auto-configure | `6d51fa3`   | daemon-carried                                                                                                                                           |
| samber-do-auditlog           | `e662ffb`   | first pass; repo additionally got the gate-greening work (`a7ee760`)                                                                                     |
| go-workflow-auditlog         | `291f4b4`   | runner commit + push                                                                                                                                     |
| project-meta                 | `612b58e`   | daemon-carried in `linter/go.mod`; loop's own pass SKIP_INCONSISTENT'd (go.sum untidy at baseline) — tidied and pushed after (`2b2c336`), origin verified |
| CV                           | `38e24a141` | first pass, daemon-carried; 14-module repo (two ~60 s integration suites) — the "45-minute library-policy hang" in the 19:35 report was actually this grind. BASELINE_RED on career-pipeline/vision + root integration, same-failure-set pre/post → commit proceeded per guard policy |
| library-policy               | `72ed2d4`   | NOT hung — completed with NO_CHANGES (consumes via flake master-pin; go.mod already at v0.6.2)                                                           |
| securitymd                   | `69e7173`   |                                                                                                                                                          |
| linter-autoconfigure-sdk     | `b566268`   |                                                                                                                                                          |
| DiscordSync                  | `ea6c444c`  | BASELINE_RED same-failure-set (its QuickTime MIME fix `487012f9` landed in the same arc)                                                                 |
| dependabot-auto-configure    | `30a5701`   |                                                                                                                                                          |
| ArtCode                      | `4f34657`   |                                                                                                                                                          |
| samber-linter                | `5795e92`   |                                                                                                                                                          |
| clean-wizard                 | `b061a7c`   | first pass, daemon-carried; the rerun's verification pass was killed after ~1 h stuck in an uninterruptible btrfs read inside `go test` (origin already had v0.6.2 — nothing lost) |
| mr-sync                      | `1891fc1`   | MANUAL (see vendor bug below) — surgical commit of go.mod/go.sum only; concurrent dashboard session's working tree left untouched                        |
| media-atlas                  | `f7f4a8a`   | MANUAL — tracked vendor/ synced in the same commit; pushed to `github` (this repo's branch.*.remote is `github`; its `origin` is a lagging Forgejo mirror) |

## Skipped (3 + 1 absent, with reasons)

| Repo        | Reason                                                                                                                                                                            |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| crm         | Baseline-red at HEAD, pre-existing: root go.mod drifted from the module graph (templ-components v1.26.0 vs graph v1.27.0), and after tidy the generated templ files call `templ.ResolveAttributeValue`/`templ.JoinURLErrs`, undefined in the resolved templ runtime. Source-verified and filed as [crm#7](https://github.com/LarsArtmann/crm/issues/7); local experimental tidy+bump reverted. |
| BuildFlow   | SKIP_DIRTY — concurrent session (matching `.githooks/pre-commit` processes observed live)                                                                                          |
| index       | SKIP_DIRTY — docs-stats TSV rewritten periodically (docs-archive-stats timer); tree flaps dirty                                                                                    |
| BDSM-World  | Repo absent from `/home/lars/projects/` entirely — nothing to bump (was v0.5.1-consumer in the 06:15 v0.6.1 pass)                                                                  |

## Runner bug found + fixed this pass

`fleet-bump.sh` ran `[vendor sync]` AFTER the after-build/after-test phase. Any vendored consumer repo (mr-sync, media-atlas) failed its after-build with "inconsistent vendoring" (go.mod bumped, vendor/modules.txt not), tripping a **false AFTER_RED** → `git restore` reverted the bump. Fix: vendor sync moved to directly after `go mod tidy`, before the verification phases. Both victims re-bumped manually with `go mod vendor` in the loop.

## Cross-pass facts worth keeping

- **Daemon-commit carrier**: 5 of 17 bumps landed via daemon commits (`golangci-lint-auto-configure`, CV, project-meta, clean-wizard + the first-pass set). The runner's own `git commit` then finds a clean tree and exits NO_CHANGES without pushing — the daemon commits ride the next push. Verified per-repo that origin carries v0.6.2.
- **Dual-remote topology**: media-atlas pushes to `github` (canonical) with `origin` pointing at a Forgejo mirror that lags. Pin checks that read `origin/master` are valid for the other 16 repos only.
- **`go mod tidy` prunes non-importers**: the runner's per-module `go get` + `tidy` added go-atomic-write to 9 CV sub-modules and tidy immediately removed it (they don't import it) — visible in the logs as empty `[verify go.mod]` grep output. Not a defect.
- **pkg.go.dev**: v0.6.2 published and rendered (checked 20:0x; the "not in the latest version" banner is the 0.x-prerelease labeling quirk).

## Deferred verification

- Retry crm after its baseline is repaired (issue filed); retry BuildFlow/index in a clean window (next bump pass is fine — both are v0.6.1/v0.6.0 consumers with live concurrent writers).
- CV chronic-red filing still pending storm-session-ownership confirmation (TODO_LIST row).
