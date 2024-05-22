# `gh-not` 🔕

> GitHub rule-based notifications management

# Install

- Download a binary from the [release page.](https://github.com/nobe4/gh-not/releases/latest)

- Install via `gh`:
    ```shell
    gh extension install nobe4/gh-not
    ```

    Is used with `gh not`, while the others `gh-not`. The documentation uses
    `gh-not` exclusively.

- Install via `go`:
    ```shell
    go install github.com/nobe4/gh-not/cmd/gh-not@latest
    ```

- Build from sources

    ```shell
    go build ./cmd/gh-not
    ```

# Run

```shell
gh-not --help
gh-not sync
gh-not list --filter '.author.login | contains("4")'
...
```

# Configure

> [!NOTE]
> The tool requires a config file to exist for it to run, that's unfortunate and
> will be fixed soon.

The configuration file contains the rules to apply to the notifications. Each
rule contains three fields:

- `name`: the display name

- `action`: the action to perform on the notification

    The current list of action is found in [`actors.go`](./internal/actors/actors.go).

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

    See more at [`config.sample.yaml`](./config.sample.yaml).


[^gojq]: Technically, [`gojq`](https://github.com/itchyny/gojq) is used.
