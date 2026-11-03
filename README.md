# go-docgen-action

A GitHub Action that generates Markdown API documentation for a Go module with
[docgen](https://github.com/lestrrat-go/docgen). Use it for modules that pkg.go.dev
cannot host, such as ones under a non-OSS license. The docs are plain Markdown
files that you commit to the module's own repository.

The action installs Go and `docgen`, then runs `docgen`. It never commits or
pushes. Your workflow decides what to do with the generated files.

## Generate docs

```yaml
- uses: actions/checkout@v5
- uses: lestrrat/go-docgen-action@main
  id: docgen
- run: git status --short docs
```

## Fail when committed docs are stale

```yaml
- uses: actions/checkout@v5
- uses: lestrrat/go-docgen-action@main
  with:
    check: "true"
```

The step prints the paths of files that differ and fails.

## Inputs

| Input | Default | Description |
|-------|---------|-------------|
| `version` | `latest` | docgen version to install: a tag, branch, commit, or `latest`. |
| `go-version` | `stable` | Go version to install. docgen needs a recent toolchain. |
| `module-dir` | `.` | Module root; must contain `go.mod`. |
| `out` | `docs` | Output directory. A relative path is joined with `module-dir`. |
| `check` | `false` | When `true`, write nothing and fail if committed docs differ. |
| `tags` | | Comma-separated build tags. |
| `goos` | | GOOS for the build context. Empty uses the docgen default (`linux`). |
| `goarch` | | GOARCH for the build context. Empty uses the docgen default (`amd64`). |
| `source-ref` | | Enable GitHub source links at this ref. |
| `source-base-url` | | Enable source links under this base URL. |

## Outputs

| Output | Description |
|--------|-------------|
| `changed` | `true` if generating changed files under the output directory, `false` if not. Empty in check mode, or when `module-dir` is not in a git repository. |
