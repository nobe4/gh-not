# Contributing

## Setup

Check `go.mod` for the required Go version.

## Scripts

- `script/build`
- `script/test`
- `script/lint`

## Commits

Using [conventional commits](https://www.conventionalcommits.org/) is preferred.

```
feat(repl): add search filter
fix(cache): handle nil expiration
docs(README): update install steps
```

## Pull requests

- Keep PRs as small as possible. One concern per PR. If a
  PR can be split, I will ask you to split it.
- Enable "Allow edits from maintainers" on your PR.
- Must use conventional commit format for the PR title. PRs are squash-merged
  using the title as the commit message.
  The maintainer will rename the PRs if needed.

## Breaking changes

Add a `!` after the type/scope in your commit message:

```
feat(config)!: rename cache path
```

## Changelog

When you add a change, put it in the `## Next` section at
the top of `CHANGELOG.md`. Use this format:

```md
## Next

* short description of the change by @you in https://github.com/nobe4/gh-not/pull/NNN
```

Don't move entries to a version heading. The maintainer
handles versioning at release time.
