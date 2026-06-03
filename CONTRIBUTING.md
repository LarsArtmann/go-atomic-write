# Contributing

Thanks for your interest in contributing to `go-atomic-write`!

## How to Contribute

1. Fork the repository
2. Create a feature branch from `master`
3. Make your changes
4. Ensure tests and benchmarks pass (see below)
5. Submit a pull request

## Development Setup

No Makefile, no flake.nix — plain Go toolchain only.

```bash
go test ./...          # Run all tests
go test -race ./...    # Run with race detector
go vet ./...           # Static analysis
go build ./...         # Verify compilation
go test -bench=. -benchmem -benchtime=100000x   # Run benchmarks
```

### Testing Conventions

- Standard `testing` package — no testify or other test utilities
- Use `t.Parallel()` for isolated tests, `t.TempDir()` for filesystem isolation
- Use `t.Helper()` for test helpers like `tempFile(t, content)`
- Concurrent tests should exercise the actual `flock` contention path, not just goroutine races

## Reporting Issues

Please use [GitHub Issues](https://github.com/larsartmann/go-atomic-write/issues) to report bugs or request features.
