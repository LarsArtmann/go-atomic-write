> [!NOTE]
> This filing was drafted by GLM via Crush from an AI-run investigation, not at my request.
>
> - [ ] MANUALLY REVIEWED by `@Lars Artmann` at `[<date-time>]`

**TL;DR:** `go build ./...` refuses to run at HEAD because go.mod drifted from the module graph, and once tidied, the generated templ files call two `templ` functions the resolved templ version does not define. Ask: regenerate the templ files against the resolved templ version and commit a tidy go.mod. Found during the go-atomic-write v0.6.2 fleet bump — its SKIP_INCONSISTENT guard refused to touch the repo, which is how the breakage surfaced.

## Problem

At HEAD (`13ae851`), a clean tree fails immediately:

```
$ go build ./...
go: updates to go.mod needed; to update it:
    go mod tidy
```

Running `go mod tidy` aligns `templ-components` v1.26.0 → v1.27.0 (the module graph already resolves 1.27.0; go.mod is the stale side). After that, compilation fails on the generated files:

```
# github.com/larsartmann/crm/internal/httpapi
internal/httpapi/companies_templ.go:185:53: undefined: templ.ResolveAttributeValue
internal/httpapi/companies_templ.go:203:53: undefined: templ.JoinURLErrs
internal/httpapi/companies_templ.go:277:54: undefined: templ.JoinURLErrs
internal/httpapi/company_detail_templ.go:146:53: undefined: templ.JoinURLErrs
```

So the generated `*_templ.go` files were produced by a templ runtime newer than the one go.mod resolves. Two layers of the same drift: stale go.mod, then generated-code-ahead-of-runtime.

## Impact

- Every CI run on master is red at the build step.
- The go-atomic-write v0.6.2 fleet bump skipped this repo (its baseline guard refuses to bump an inconsistent module), so `go-atomic-write` sits at v0.6.1 here until the baseline is fixed.
- Reproduction is trivial: fresh clone → `go build ./...`.

## Fix shape

1. Commit the tidy (or realign the graph — whichever direction is intended).
2. Run `templ generate` with the resolved templ version so the generated files match the runtime again.
3. Then the fleet bump for go-atomic-write v0.6.2 can land (`go get github.com/larsartmann/go-atomic-write@v0.6.2` verified working once the tree builds).

💘 Generated with Crush
