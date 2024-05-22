# `gh-not`

> GitHub rule-based notifications managements

# Installation

> TBD when a release is done

# Configuration

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

    ```shell
    jq '.[] | select((.author.login == "dependabot[bot]") or (.author.login == "nobe4"))' | jq '.[] | select(.repository.full_name == "nobe4/gh-not")'
    ```


[^gojq]: Technically, [`gojq`](https://github.com/itchyny/gojq) is used.
