# Draft: go-workflow-auditlog — lint debt unmasked by the golangci-lint v2.14.0 pin

**Repo:** LarsArtmann/go-workflow-auditlog
**Kind:** own-repo tracker issue, unsolicited (AI-run CI sweep)
**Status:** DRAFTED, not yet filed (awaiting review per banner)

---

```markdown
> [!NOTE]
> This filing was drafted by Crush:glm-5.3-flash via Crush from an AI-run investigation, not at my request. When the failure traces to an external report, that source is linked in the body.
>
> - [ ] MANUALLY REVIEWED by `@Lars Artmann` at `[<date-time>]`

## Problem

CI on `master` is red at `623fd38` (run 41926688745). The golangci-lint pin moved v2.13.2 → v2.14.0 (22:33 tonight) to fix the Go 1.27 export-data typecheck failure — the toolchain problem is fixed, and v2.14.0 now reports four real findings that v2.13.2 never surfaced because it died in typecheck first:
```

report_builder.go:201:3: missing cases in switch of type auditlog.EventType: auditlog.EventTypeSystem (exhaustive)
report_builder.go:233:3: missing cases in switch of type auditlog.EventType: auditlog.EventTypeSystem (exhaustive)
plugin.go:144:1: godoc should start with symbol name ("Record") (godoclint)
record_api_test.go:14:1: Function TestAuditorRecord_InjectsSyntheticEvent missing the call to method parallel (paralleltest)

```
## Why now

The typecheck failure masked everything since the repo's `go` directive crossed 1.27 (golangci-lint v2.13.2 is built with go1.26 and cannot read Go 1.27 stdlib export data — same incident as go-atomic-write's CI, fixed there by pinning v2.14.0). The pin upgrade is the right fix; these four findings are pre-existing debt it exposed, not regressions.

## Fix

The two `exhaustive` findings need a domain decision: `EventTypeSystem` exists in the enum but `report_builder.go:201` and `:233` don't handle it — decide whether system events get their own report line or an explicit ignore. The other two are mechanical (`plugin.go:144` doc comment should start with `Record`; `record_api_test.go:14` needs `t.Parallel()` or a `//nolint:paralleltest` with reason).

## Context

Found during the v0.6.2 consumer CI sweep (go-atomic-write fleet bump); vendorHash and lint-pin repairs for the same sweep are already on `master`. The warning `Found unknown linters in //nolint directives: erraudit` also appeared — the erraudit plugin isn't loaded in v2.14.0's config, so those `//nolint:erraudit` comments are currently no-ops.

💘 Generated with Crush
```
