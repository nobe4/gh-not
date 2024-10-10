# `gh-not` 🔕

> GitHub rule-based notifications management

[![Go Reference](https://pkg.go.dev/badge/github.com/nobe4/gh-not.svg)](https://pkg.go.dev/github.com/nobe4/gh-not)
[![CI](https://github.com/nobe4/gh-not/actions/workflows/ci.yml/badge.svg)](https://github.com/nobe4/gh-not/actions/workflows/ci.yml)
[![CodeQL](https://github.com/nobe4/gh-not/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/nobe4/gh-not/actions/workflows/github-code-scanning/codeql)

> [!IMPORTANT]
> Under heavy development, expect nothing.
> PRs/issues welcome!

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

# Configure

The configuration file contains the rules to apply to the notifications. Each
rule contains three fields:

- `name`: the display name

- `action`: the action to perform on the notification

    The current list of action is found in [`actions.go`](./internal/actions/actions.go).

- `filters`: a list of [`jq` filters](https://jqlang.github.io/jq/manual/#basic-filters)[^gojq]
    to filter notifications with.

    Each filter is inserted into the following patter: `.[] | select(%s)`.


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

# Automatic fetching

To automatically fetch new notifications and apply the rules, it is recommended
to setup an automated process to run `gh-not sync` regularly.

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
