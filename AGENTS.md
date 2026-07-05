# AGENTS.md

## Project Overview

`fetch` is a Go command-line API testing tool. It reads YAML/TOML/JSON-style case files, builds HTTP requests, runs them through the `fetch run` Cobra command, and verifies responses with assertions for status, headers, cookies, body content, JSON, HTML, XML, YAML, TOML, redirects, and errors.

The main entrypoint is `main.go`, CLI commands live in `cmd/`, and most implementation code lives under `pkg/`:

- `pkg/config`: case, request, assertion, config-file parsing, and template rendering hooks.
- `pkg/client`: HTTP request execution helpers.
- `pkg/assert`: primitive assertion operators used by response assertions.
- `pkg/assertion`: response-specific assertion implementations.
- `pkg/tpl`, `pkg/util`, `pkg/log`, `pkg/version`: shared support packages.
- `examples/`: executable case files used as reference material.
- `docs/`: Jekyll documentation site.

## Setup Commands

- Use Go 1.26 or newer for normal development; `go.mod` declares `go 1.26` and `toolchain go1.26.4`.
- Install formatting tools:
  ```bash
  make init
  ```
- Download module dependencies without vendoring:
  ```bash
  go mod download
  ```
- Sync vendored dependencies for the Makefile test command:
  ```bash
  go mod vendor
  ```
- Full dependency maintenance, when intentionally changing dependencies:
  ```bash
  make dep
  ```

`vendor/` is ignored by git, but `make test` uses `-mod=vendor`, so a fresh checkout needs `go mod vendor` before the canonical test command will run.

## Development Workflow

- Build the CLI binary:
  ```bash
  make build
  ```
- Run the CLI directly during development:
  ```bash
  go run . --help
  go run . run examples/http_get.yaml
  ```
- Run the built binary against an example:
  ```bash
  ./fetch run examples/http_get.yaml
  ```
- Run with verbose output when inspecting request/response behavior:
  ```bash
  ./fetch run examples/http_get.yaml -v
  ```
- Run multiple cases in parallel:
  ```bash
  ./fetch run examples/http_*.yaml -p 5
  ```

The CLI is built with Cobra. Add or edit commands in `cmd/`, register them from `init()`, and keep command behavior close to the existing `run`, `version`, `bootstrap`, and `generate` files.

## Testing Instructions

- Run the project test target:
  ```bash
  make test
  ```
- If `make test` fails with an inconsistent vendoring error before tests run, refresh local vendoring first:
  ```bash
  go mod vendor
  make test
  ```
- Run a focused package test while developing:
  ```bash
  go test ./pkg/assert -run TestStartsWith
  ```
- Run all packages without the Makefile coverage flags:
  ```bash
  go test ./...
  ```

Most tests are standard Go `testing` table tests in `*_test.go` files. `pkg/util` also has a Ginkgo/Gomega suite. Add focused tests next to the package you change, and prefer behavior-oriented cases that exercise the public package function or CLI behavior being modified.

## Code Style

- Format Go code with the project formatter target:
  ```bash
  make fmt
  ```
- `make fmt` runs `golines` with a 120-column limit and then `gofumpt`.
- Keep package names short and aligned with the existing directory names.
- Keep new CLI flags and environment variables documented in `docs/Usage/command-line.md` when they affect users.
- Do not introduce a new dependency if the standard library or an existing dependency already covers the need.
- If dependencies change, update `go.mod` and `go.sum` with `go mod tidy`; only refresh `vendor/` locally unless the repository policy changes.

## Documentation

The documentation site is under `docs/` and uses Jekyll with `github-pages`, `just-the-docs`, and `docs/Gemfile.lock`.

- Install the Bundler version required by the lockfile if your Ruby environment does not already have it:
  ```bash
  gem install bundler:2.3.26
  ```
- Install docs dependencies:
  ```bash
  bundle install
  ```
- Serve docs locally from `docs/`:
  ```bash
  bundle exec jekyll serve
  ```

Update README examples, `examples/`, and `docs/Usage/` together when changing user-visible case syntax or CLI behavior.

## Build and Release

- Local binary build:
  ```bash
  make build
  ```
- Docker image build:
  ```bash
  docker build -t fetch .
  ```
- Release automation is tag-driven in `.github/workflows/release.yml`.
- GoReleaser Homebrew tap settings are in `.goreleaser.yml`.
- GitHub Pages docs deploy from `.github/workflows/jekyll-gh-pages.yml` on pushes to `master`.

Keep the Dockerfile builder image, release workflow Go version, README, docs, and `go.mod` aligned when changing the Go version.

## Pull Request Checklist

- Keep diffs narrow and directly tied to the requested change.
- Run `make fmt` for Go code changes.
- Run `make test` after syncing local `vendor/` with `go mod vendor`.
- For docs-only changes, verify the relevant Markdown and Jekyll files; a Go test run is not required unless code or examples used by tests changed.
- For CLI behavior changes, update the corresponding docs and add or update an example case when useful.

## Agent Notes

- Do not fabricate commands, release behavior, or CI coverage. This repo has release and docs workflows, but no ordinary push/PR Go test workflow in `.github/workflows/`.
- Read existing code before editing; package boundaries are simple and should stay that way.
- Prefer small, direct changes over new abstractions.
- Preserve existing examples unless the requested change specifically updates case-file behavior.
- Treat `bootstrap` and `generate` as placeholder Cobra commands unless the task is explicitly about them.
