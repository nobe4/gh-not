# `gh-not` 🔕

> GitHub rule-based notifications management

[![Go Reference](https://pkg.go.dev/badge/github.com/nobe4/gh-not.svg)](https://pkg.go.dev/github.com/nobe4/gh-not)
[![CI](https://github.com/nobe4/gh-not/actions/workflows/ci.yml/badge.svg)](https://github.com/nobe4/gh-not/actions/workflows/ci.yml)
[![CodeQL](https://github.com/nobe4/gh-not/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/nobe4/gh-not/actions/workflows/github-code-scanning/codeql)

> [!IMPORTANT]
> The project is mostly "done" at this point. I won't be adding any new features.
>
> I will accept PRs/issues for bug/improvements.

![demo.gif](./docs/demo.gif)

# Install

- Install via `gh` (preferred method):
    ```shell
    gh extension install nobe4/gh-not
    ```

    It is then used with `gh not`.

- Download a binary from the [release page.](https://github.com/nobe4/gh-not/releases/latest)

- Build from sources
    ```shell
    go generate ./...
    go build ./cmd/gh-not
    ```

    See [`version.go`](./internal/version/version.go) for custom build info.

# Getting started

Run the following commands to get started and see the available commands and
options:

```shell
gh-not --help
gh-not config --init
gh-not sync
gh-not --filter '.author.login | contains("4")'
gh-not --repl
...
```

# How it works

`gh-not` fetches the notifications from GitHub and saves them in a local cache.

The synchronization between local and remote notifications is described in
[`sync.go`](./internal/notifications/sync.go).

The `sync` command applies the rules to the notifications and performs the
specified actions. It's recommended to run this regularly, see [this
section](#automatic-fetching).

The other commands are used to interact with the local cache. It uses `gh
api`-equivalent to modify the notifications on GitHub's side.

# Authentication

`gh-not` uses `gh`'s built-in authentication, meaning that if `gh` works,
`gh-not` should work.

If you want to use a specific PAT, you can do so with the environment variable
`GH_TOKEN`. The PAT requires the scopes: `notifications`, and `repo`.

E.g.:

```bash
# gh's authentication for github.com
gh-not ...

# Using a PAT for github.com.
GH_TOKEN=XXX gh-not ...
```

`gh-not` also respects `GH_HOST` and `GH_ENTERPRISE_TOKEN` if you need to use a
non-`github.com` host.

E.g.:

```bash
# gh's authentication for ghe.io
GH_HOST=ghe.io gh-not ...

# Using a PAT for ghe.io
GH_HOST=ghe.io GH_ENTERPRISE_TOKEN=XXX gh-not ...
```

See the [`gh` environment documentation](https://cli.github.com/manual/gh_help_environment).

> [!IMPORTANT]
> If you plan on using `gh-not` with more than one host, you might want to
> create a separate cache for it. See [cache](#cache).

# Configuration

## Cache

The cache is where the notifications are locally stored.

It contains 2 fields:

- `path`: the path to the JSON file.

- `TTLInHours`: how long before the cache needs to be refreshed.

If you use multiple hosts, you might want to have separate configurations and
caches to prevent overrides. Create one config file per host you want to use and
point the cache's path to a _different file_.

E.g.

- `config.github.yaml`
    ```yaml
    cache:
      path: cache.github.json
    ...
    ```

    Use it with `gh-not --config config.github.yaml`.

- `config.gheio.yaml `
    ```yaml
    cache:
      path: cache.gheio.json
    ...
    ```
    Use it with `gh-not --config config.gheio.yaml`.

## Rules

The configuration file contains the rules to apply to the notifications. Each
rule contains three fields and must contain an action and at least one filter.

- `name`: the display name

- `action`: the action to perform on the notification

    The current list of action is found in [`actions.go`](./internal/actions/actions.go).

- `filters`: a list of [`jq` filters](https://jqlang.github.io/jq/manual/#basic-filters)[^gojq]
    to filter notifications with.

    Each filter is inserted into the following pattern: `.[] | select(%s)`.


    Each filter in the list is run one after the other, making it similar to
    joining them with `and`.

    It means that if you want to specify conditions with `or`, you need to write
    them directly in the filter.

    E.g.
    ```yaml
    rules:
      - filters:
          - (.author.login == "dependabot[bot]") or (.author.login == "nobe4")
          - .repository.full_name == "nobe4/gh-not"
    ```

    Filters like:

    ```shell
    jq '.[] | select((.author.login == "dependabot[bot]") or (.author.login == "nobe4"))' | jq '.[] | select(.repository.full_name == "nobe4/gh-not")'
    ```

    See more at [`config.go`](./internal/config/config.go) and [`rule.go`](./internal/config/rule.go).

### Examples

```yaml
- name: mark closed dependabot PRs as done
  filters:
    - .author.login == "dependabot[bot]"
    - .subject.state == "closed"
  action: done
```

```yaml
- name: ignore a specific repo
  filters:
    - .repository.name == "greg-ci-tests"
  action: done
```

```yaml
- name: close read notifications
  filters:
    - .unread == false
  action: done
```

```yml
- name: mark notifications of PRs closed over a week ago as read
  filters:
    - .subject.state == "closed"
    - .updated_at | fromdate < now - 604800
  action: read
```

```yml
- name: close read notifications older than 2 weeks
  filters:
    - .unread == false    
    - .updated_at | fromdate < now - 1209600
  action: done
```

# Automatic fetching

To automatically fetch new notifications and apply the rules, it is recommended
to set up an automated process to run `gh-not sync` regularly.

E.g.

- `cron`

    ```shell
    0 * * * *  gh-not sync --config=/path/to/config.yaml --verbosity=4 >> /tmp/gh-not-sync.out 2>> /tmp/gh-not-sync.err
    ```

- [`launchd`](https://launchd.info/) (macOS)

    ```xml
    <?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
    <plist version="1.0">
      <dict>
        <key>EnvironmentVariables</key>
        <dict>
          <key>PATH</key>
          <string>/opt/homebrew/bin/:$PATH</string>
        </dict>

        <key>Label</key>
        <string>launched.gh-not-sync.hourly</string>

        <key>ProgramArguments</key>
        <array>
          <string>sh</string>
          <string>-c</string>
          <string>gh-not sync --config=/path/to/config.yaml --verbosity=4</string>
        </array>

        <key>StandardErrorPath</key>
        <string>/tmp/gh-not-sync.err</string>
        <key>StandardOutPath</key>
        <string>/tmp/gh-not-sync.out</string>

        <key>StartInterval</key>
        <integer>3600</integer>
      </dict>
    </plist>
    ```

    It is recommended to use https://launched.zerowidth.com/ for generating such Plist.

[^gojq]: Technically, [`gojq`](https://github.com/itchyny/gojq) is used.
