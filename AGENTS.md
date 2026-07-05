# onepagepass — Agent Onboarding Guide

`onepagepass` is a small Go CLI that repacks a PDF scan of a document into a
"two-pages-per-A4" layout. It rasterizes the source PDF with MuPDF, merges the
top halves of consecutive pages, and writes a new A4 PDF.

## Layout

| Path                      | Purpose                                                            |
| ------------------------- | ----------------------------------------------------------------- |
| `cmd/main.go`             | CLI entry point (`urfave/cli/v3`), wires the `convert` command.    |
| `pkg/converter/`          | Conversion orchestration + page-merge logic (`image/draw`).        |
| `pkg/pdf/`                | PDF read (`go-fitz`/MuPDF) and write (`gopdf`) wrapper behind the `PdfFile` interface. |
| `pkg/converter/mocks/`    | Generated `PdfFile` mock (mockery) used by the converter tests.     |

## CLI Command Reference

All automation is managed via [`go-task`](https://taskfile.dev) (`Taskfile.yml`):

```bash
task setup          # download Go module dependencies
task lint           # run golangci-lint (v2 config)
task test           # run unit tests
task test:coverage  # tests with race detector + coverage profile
task format         # gofumpt/gci/golines + golangci-lint --fix
task generate       # regenerate mocks (go generate ./...)
task build          # compile bin/onepagepass with version ldflags
task clean          # remove build artifacts and test cache
```

## Task Rules

1. **Idiomatic Go**: standard library first, canonical error handling
   (`fmt.Errorf("...: %w", err)`), small interfaces. Follow Effective Go.
2. **Third-party dependencies require justification**: prefer the standard
   library. The direct dependency set is intentionally minimal — `go-fitz`
   (MuPDF rasterizer, no pure-Go equivalent), `gopdf` (PDF writer),
   `urfave/cli/v3` (CLI), and `testify` (tests). Do not add new dependencies
   without a clear reason the stdlib cannot cover.
3. **Testing with testify suites**: tests use `github.com/stretchr/testify`
   organised as suites. One top-level `func Test<Package>(t *testing.T)` entry
   point per package that only calls `suite.Run(...)`; all assertions live in
   suite methods (`s.Error`, `s.NoError`, `s.ErrorIs`, ...). The `PdfFile`
   collaborator is mocked (`pkg/converter/mocks`).
4. **Lint before finishing**: `task lint` must pass. It is the same gate CI
   enforces (`.github/workflows/ci.yml`).
5. **Commit messages**: no AI/agent attribution trailers of any kind.

## Notes

- Go version: **1.26** (see `go.mod`).
- `go-fitz` loads MuPDF at runtime via `purego` — building and testing need no
  CGO toolchain; only actual conversion at runtime touches the native library.
- Coverage is reported to Coveralls from CI.
