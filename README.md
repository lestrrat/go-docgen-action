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
| `version` | `latest` | docgen version to install: a tag, branch, commit, or `latest`. A tag or commit is cached between runs; a branch or `latest` is rebuilt every run. |
| `go-version` | | Go version to install. Empty reads it from `module-dir/go.mod`, so docgen is built with a toolchain at least as new as the module needs. |
| `module-dir` | `.` | Module root; must contain `go.mod`. |
| `out` | `docs` | Output directory. A relative path is joined with `module-dir`. |
| `check` | `false` | When `true`, write nothing and fail if committed docs differ. |
| `tags` | | Comma-separated build tags. |
| `goos` | | GOOS for the build context. Empty uses the docgen default (`linux`). |
| `goarch` | | GOARCH for the build context. Empty uses the docgen default (`amd64`). |
| `source-ref` | | Enable GitHub source links at this ref. |
| `source-base-url` | | Enable source links under this base URL. |
| `commit` | `false` | When `true` (and `check` is not `true`), commit changed docs under `out` and push them. See [Committing the docs](#committing-the-docs). |
| `commit-message` | `update generated docs` | Commit message. |
| `commit-user-name` | `github-actions[bot]` | Git author and committer name. |
| `commit-user-email` | `41898282+github-actions[bot]@users.noreply.github.com` | Git author and committer email. |

## Outputs

| Output | Description |
|--------|-------------|
| `changed` | `true` if generating changed files under the output directory, `false` if not. Empty in check mode, or when `module-dir` is not in a git repository. |
| `cache-hit` | `true` if the docgen binary came from the cache, `false` if it was built. Empty when `version` is not a tag or commit. |
| `pushed` | `true` if the action pushed a docs commit, `false` if not. Empty unless `commit` is `true`. |

## Committing the docs

Set `commit: "true"` to have the action commit and push changed files under `out`. The job
needs `contents: write`.

```yaml
on:
  push:
    branches: [main]
  pull_request:

permissions:
  contents: write

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5
        with:
          # On pull requests, check out the head branch, not the default merge commit.
          ref: ${{ github.head_ref || github.ref }}
      - uses: lestrrat/go-docgen-action@main
        with:
          commit: "true"
```

Where the commit goes:

- On a push, the action pushes to the current branch.
- On a pull request, it pushes to the pull request's head branch.

The action warns and continues, leaving the docs stale and `pushed` set to `false`, when:

- The pull request comes from a fork.
- A pull request job checked out the merge commit instead of the head branch.
- The ref is a tag.
- The push fails (read-only token, protected branch, or the branch moved).

Pushes made with the default `GITHUB_TOKEN` do not start other workflows. To run required checks
on the docs commit, give `actions/checkout` a personal access token or app token through its
`token` input.

## Caching

When `version` is a tag (such as `v0.1.0`) or a commit hash, the action caches the built
`docgen` binary with `actions/cache`. The cache key includes the runner OS and architecture,
the installed Go version, and `version`. Pin `version` to get the cache.
