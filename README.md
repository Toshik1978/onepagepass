[![Build and Test](https://github.com/Toshik1978/onepagepass/actions/workflows/ci.yml/badge.svg)](https://github.com/Toshik1978/onepagepass/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/Toshik1978/onepagepass/badge.svg?branch=main)](https://coveralls.io/github/Toshik1978/onepagepass?branch=main)

# One-page passport

> CLI tool that repacks a PDF scan of a small document into a "two-pages-per-A4" layout that is handier to print.

If you have a passport (or any other document) that is smaller than A5, a plain scan wastes half the sheet
when printed. `onepagepass` rasterizes the source PDF and stacks the top half of each page on top of the next,
so two document pages share a single A4 sheet.

## How it works

1. The source PDF is rendered to images page by page via [MuPDF](https://mupdf.com/) (`go-fitz`).
2. Consecutive pages are merged — the top half of page _N+1_ is placed into the bottom half of page _N_.
3. The merged pages are written to a new A4 PDF via [`gopdf`](https://github.com/signintech/gopdf).

The result is saved next to the source as `<name>.converted.pdf`.

## Install

```bash
go install github.com/Toshik1978/onepagepass/cmd@latest
```

Or build from source (requires [Task](https://taskfile.dev)):

```bash
task build      # produces bin/onepagepass
```

## Usage

```bash
onepagepass convert --pdf passport.pdf
```

| Flag    | Required | Default | Description                         |
| ------- | -------- | ------- | ----------------------------------- |
| `--pdf` | yes      | —       | Path to the source PDF file.        |
| `--dpi` | no       | `300`   | Rasterization resolution in DPI.    |

Run `onepagepass --help` or `onepagepass convert --help` for the full reference.

## Development

All automation is driven by [`Taskfile.yml`](./Taskfile.yml):

```bash
task setup          # download Go module dependencies
task lint           # run golangci-lint
task test           # run unit tests
task test:coverage  # tests with race detector + coverage profile
task format         # gofumpt/gci/golines + golangci-lint --fix
task build          # compile bin/onepagepass
task clean          # remove build artifacts and test cache
```

See [`CLAUDE.md`](./CLAUDE.md) / [`AGENTS.md`](./AGENTS.md) for contributor and agent onboarding notes.

## License

See [`LICENSE`](./LICENSE).
