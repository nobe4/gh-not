# Changelog

## Next

## [v0.7.0](https://github.com/nobe4/gh-not/releases/tag/v0.7.0) - 2026-05-13

* doc: how to run gh-not as a systemd service by @Tethik in https://github.com/nobe4/gh-not/pull/305
* chore(lint): fix and update config by @nobe4 in https://github.com/nobe4/gh-not/pull/312
* Cache notification enrichment state by @parkerbxyz in https://github.com/nobe4/gh-not/pull/311
* Make notification enrichment atomic on failure by @parkerbxyz in https://github.com/nobe4/gh-not/pull/310
* Parallelize notification enrichment by @parkerbxyz in https://github.com/nobe4/gh-not/pull/309


## [v0.6.10](https://github.com/nobe4/gh-not/releases/tag/v0.6.10) - 2026-01-15

* Add merged by to enrichments by @Tethik in https://github.com/nobe4/gh-not/pull/304


## [v0.6.9](https://github.com/nobe4/gh-not/releases/tag/v0.6.9) - 2026-01-12

* auto(ln): update links by @github-actions[bot] in https://github.com/nobe4/gh-not/pull/300
* fix: parse $EDITOR arguments by @parkerbxyz in https://github.com/nobe4/gh-not/pull/303


## [v0.6.8](https://github.com/nobe4/gh-not/releases/tag/v0.6.8) - 2025-12-23

* feat: make config init safer by not overwriting an existing config file by @Tethik in https://github.com/nobe4/gh-not/pull/295
* feat: add devenv config by @nobe4 in https://github.com/nobe4/gh-not/pull/299
* feat: stricter validation of rules  by @Tethik in https://github.com/nobe4/gh-not/pull/297


## [v0.6.7](https://github.com/nobe4/gh-not/releases/tag/v0.6.7) - 2025-11-06

* auto(ln): update links by @github-actions[bot] in https://github.com/nobe4/gh-not/pull/279
* auto(ln): update links by @github-actions[bot] in https://github.com/nobe4/gh-not/pull/280
* chore(lint): format for the new golangci-lint version by @nobe4 in https://github.com/nobe4/gh-not/pull/283
* Update ln.yaml by @nobe4 in https://github.com/nobe4/gh-not/pull/287
* fix(lint): format by @nobe4 in https://github.com/nobe4/gh-not/pull/291
* fix: action-ln version by @nobe4 in https://github.com/nobe4/gh-not/pull/292
* fix: strip potential newlines from subject's title by @nobe4 in https://github.com/nobe4/gh-not/pull/290


## [v0.6.6](https://github.com/nobe4/gh-not/releases/tag/v0.6.6) - 2025-05-23

* Rename dependabot.yml to dependabot.yaml by @nobe4 in https://github.com/nobe4/gh-not/pull/267
* Update .ln-config.yaml by @nobe4 in https://github.com/nobe4/gh-not/pull/268
* auto(ln): update links by @github-actions in https://github.com/nobe4/gh-not/pull/270
* auto(ln): update links by @github-actions in https://github.com/nobe4/gh-not/pull/271
* refactor: improve logging for REPL view by @tebriel in https://github.com/nobe4/gh-not/pull/273
* chore(view.log_file): use XDG_STATE_DIR as default repl log location by @tebriel in https://github.com/nobe4/gh-not/pull/276
* fix(repl): always leave enough space for the status line by @nobe4 in https://github.com/nobe4/gh-not/pull/274


## [v0.6.5](https://github.com/nobe4/gh-not/releases/tag/v0.6.5) - 2025-05-15

* fix(lint): migrate and fix golangci-lint by @nobe4 in https://github.com/nobe4/gh-not/pull/261
* fix(ci): add permissions by @nobe4 in https://github.com/nobe4/gh-not/pull/263
* feat(ln): bootstrap ln by @nobe4 in https://github.com/nobe4/gh-not/pull/264
* auto(ln): update links by @github-actions in https://github.com/nobe4/gh-not/pull/266
* feat(json): allow to select for a single field from the REPL by @nobe4 in https://github.com/nobe4/gh-not/pull/262
* fix(open): add separating space by @nobe4 in https://github.com/nobe4/gh-not/pull/265


## [v0.6.4](https://github.com/nobe4/gh-not/releases/tag/v0.6.4) - 2025-03-06

* fix(repl): discard all output for `open` command by @nobe4 in https://github.com/nobe4/gh-not/pull/252


## [v0.6.3](https://github.com/nobe4/gh-not/releases/tag/v0.6.3) - 2025-02-26

* fix(config): error when tilde is used in paths by @nobe4 in https://github.com/nobe4/gh-not/pull/247
* feat(config): enable $ENVIRONMENT-based path by @nobe4 in https://github.com/nobe4/gh-not/pull/250


## [v0.6.2](https://github.com/nobe4/gh-not/releases/tag/v0.6.2) - 2025-02-14

* feat(version): allow custom build info by @nobe4 in https://github.com/nobe4/gh-not/pull/245


## [v0.6.1](https://github.com/nobe4/gh-not/releases/tag/v0.6.1) - 2025-02-14

* fix: parallel test by @nobe4 in https://github.com/nobe4/gh-not/pull/239
* fix: move back to permission numbers by @nobe4 in https://github.com/nobe4/gh-not/pull/240
* feat: enable predeclared and nlreturn by @nobe4 in https://github.com/nobe4/gh-not/pull/241
* feat(version): use go 1.24.0's new version info by @nobe4 in https://github.com/nobe4/gh-not/pull/244


## [v0.6.0](https://github.com/nobe4/gh-not/releases/tag/v0.6.0) - 2025-01-08

* feat: enable all golangci-lint linters by @nobe4 in https://github.com/nobe4/gh-not/pull/229
* fix: enable bodyclose linter by @nobe4 in https://github.com/nobe4/gh-not/pull/230
* fix: enable errorlint linter by @nobe4 in https://github.com/nobe4/gh-not/pull/231
* fix: enable recvcheck linter by @nobe4 in https://github.com/nobe4/gh-not/pull/232
* fix: enable various linters by @nobe4 in https://github.com/nobe4/gh-not/pull/233
* fix: enable more linters by @nobe4 in https://github.com/nobe4/gh-not/pull/234
* feat(gh): allow multiple hosts by @nobe4 in https://github.com/nobe4/gh-not/pull/236
* feat: enable a couple more linters by @nobe4 in https://github.com/nobe4/gh-not/pull/237
* chore: enable yet more linters by @nobe4 in https://github.com/nobe4/gh-not/pull/238


## [v0.5.8](https://github.com/nobe4/gh-not/releases/tag/v0.5.8) - 2024-12-07

* typos suggestion by @ccoVeille in https://github.com/nobe4/gh-not/pull/222
* fix(config): init bootstrap the config folder by @nobe4 in https://github.com/nobe4/gh-not/pull/225
* ci: enable golangci-lint by @ccoVeille in https://github.com/nobe4/gh-not/pull/226
* fix(ci): make code compliant with golangci-lint by @nobe4 in https://github.com/nobe4/gh-not/pull/227


## [v0.5.7](https://github.com/nobe4/gh-not/releases/tag/v0.5.7) - 2024-11-12

* Add Assignees, Reviewers and ReviewersTeams by @monrad in https://github.com/nobe4/gh-not/pull/214


## [v0.5.6](https://github.com/nobe4/gh-not/releases/tag/v0.5.6) - 2024-11-04

* docs: add demo thanks to vhs by @nobe4 in https://github.com/nobe4/gh-not/pull/207
* docs(README): simplify the installation methods by @nobe4 in https://github.com/nobe4/gh-not/pull/208
* feat(gh): add error information in the output by @nobe4 in https://github.com/nobe4/gh-not/pull/213


## [v0.5.5](https://github.com/nobe4/gh-not/releases/tag/v0.5.5) - 2024-10-10

* fix(manager): relax enrichment failure by @nobe4 in https://github.com/nobe4/gh-not/pull/205


## [v0.5.4](https://github.com/nobe4/gh-not/releases/tag/v0.5.4) - 2024-10-08

* feat(tag): add tag action by @nobe4 in https://github.com/nobe4/gh-not/pull/202
* feat(root): display tags and filter by tag by @nobe4 in https://github.com/nobe4/gh-not/pull/203
* fix(gh): accept url.Error to retry by @nobe4 in https://github.com/nobe4/gh-not/pull/204


## [v0.5.3](https://github.com/nobe4/gh-not/releases/tag/v0.5.3) - 2024-10-03

* feat(json): implement json action for the REPL by @nobe4 in https://github.com/nobe4/gh-not/pull/193
* feat(type): add discussion type for pretty rendering by @nobe4 in https://github.com/nobe4/gh-not/pull/194
* fix(open): don't open if there's no HtmlUrl by @nobe4 in https://github.com/nobe4/gh-not/pull/195


## [v0.5.2](https://github.com/nobe4/gh-not/releases/tag/v0.5.2) - 2024-10-02

* feat(actions): allow passing string parameters to actions by @nobe4 in https://github.com/nobe4/gh-not/pull/188
* feat(config): enable passing arguments to the rule by @nobe4 in https://github.com/nobe4/gh-not/pull/189
* feat(assign): adds a new assign action by @nobe4 in https://github.com/nobe4/gh-not/pull/191


## [v0.5.1](https://github.com/nobe4/gh-not/releases/tag/v0.5.1) - 2024-09-29

* fix(cache): catch nil wrap expiration date by @nobe4 in https://github.com/nobe4/gh-not/pull/180
* fix(cache): initialize the wrapper explicitly by @nobe4 in https://github.com/nobe4/gh-not/pull/181
* fix(repl): items store their index in the list by @nobe4 in https://github.com/nobe4/gh-not/pull/182
* feat(repl): hide hidden notifications by @nobe4 in https://github.com/nobe4/gh-not/pull/183
* feat(command): use suggestion's text for command by @nobe4 in https://github.com/nobe4/gh-not/pull/184
* docs(README): update basic commands format by @nobe4 in https://github.com/nobe4/gh-not/pull/185


## [v0.5.0](https://github.com/nobe4/gh-not/releases/tag/v0.5.0) - 2024-09-28

### Breaking change, kinda 🐉

The format for the cache data changed slightly.
It now stores the `refreshed_at` time in the JSON object rather than using Filesystem data.
The next time you run `gh-not sync`, it will automatically update the format and you shouldn't see any issue.
Until then, the cache will be considered 54 years old.

### What's Changed

* feat: add a cache wrapper by @nobe4 in https://github.com/nobe4/gh-not/pull/177
* feat(cache): store RefreshedAt information in the Wrapper by @nobe4 in https://github.com/nobe4/gh-not/pull/178
* feat(cache): move the expiration logic out by @nobe4 in https://github.com/nobe4/gh-not/pull/179

## [v0.4.10](https://github.com/nobe4/gh-not/releases/tag/v0.4.10) - 2024-09-19

* feat(gh): add `per_page` parameter for API calls by @nobe4 in https://github.com/nobe4/gh-not/pull/172
* feat(sync): add the time of last sync by @nobe4 in https://github.com/nobe4/gh-not/pull/173


## [v0.4.9](https://github.com/nobe4/gh-not/releases/tag/v0.4.9) - 2024-09-14

* test(integration): add error handling for responses by @nobe4 in https://github.com/nobe4/gh-not/pull/154
* test: fix CI tests by @nobe4 in https://github.com/nobe4/gh-not/pull/167
* feat(notification): add LastCommentor field by @nobe4 in https://github.com/nobe4/gh-not/pull/168


## [v0.4.8](https://github.com/nobe4/gh-not/releases/tag/v0.4.8) - 2024-09-13

* fix(view): don't render an empty notifications list by @nobe4 in https://github.com/nobe4/gh-not/pull/165
* feat(open): add output for open action by @nobe4 in https://github.com/nobe4/gh-not/pull/166


## [v0.4.7](https://github.com/nobe4/gh-not/releases/tag/v0.4.7) - 2024-09-12

* feat(repl): add notification stats in the statusline by @nobe4 in https://github.com/nobe4/gh-not/pull/162
* fix(repl): Make `None` actually work by @nobe4 in https://github.com/nobe4/gh-not/pull/163
* fix(repl): select all only works on filtered items by @nobe4 in https://github.com/nobe4/gh-not/pull/164


## [v0.4.6](https://github.com/nobe4/gh-not/releases/tag/v0.4.6) - 2024-09-12

* refactor(api): drop `Do` method in the `Requestor` interface by @nobe4 in https://github.com/nobe4/gh-not/pull/147
* test(integration): First pass at creating integration tests by @nobe4 in https://github.com/nobe4/gh-not/pull/146
* test(integration): add header handling in integration tests by @nobe4 in https://github.com/nobe4/gh-not/pull/150
* feat(cache): get the refreshed date by @nobe4 in https://github.com/nobe4/gh-not/pull/161


## [v0.4.5](https://github.com/nobe4/gh-not/releases/tag/v0.4.5) - 2024-08-11

* refactor: rename `actors` to `actions` and `actor` to `runner` by @nobe4 in https://github.com/nobe4/gh-not/pull/143
* feat(actions): automatically generate help for the actions by @nobe4 in https://github.com/nobe4/gh-not/pull/144


## [v0.4.4](https://github.com/nobe4/gh-not/releases/tag/v0.4.4) - 2024-08-11

* feat(hide): rework hidden state by @nobe4 in https://github.com/nobe4/gh-not/pull/140
* fix(notifications): use fallback string format by @nobe4 in https://github.com/nobe4/gh-not/pull/142


## [v0.4.3](https://github.com/nobe4/gh-not/releases/tag/v0.4.3) - 2024-08-10

* fix(config): rename command accept default key by @nobe4 in https://github.com/nobe4/gh-not/pull/139


## [v0.4.2](https://github.com/nobe4/gh-not/releases/tag/v0.4.2) - 2024-08-10

* docs: general improvements by @nobe4 in https://github.com/nobe4/gh-not/pull/135
* feat(sync): reset Meta.Done on newer remote notification by @nobe4 in https://github.com/nobe4/gh-not/pull/137


## [v0.4.1](https://github.com/nobe4/gh-not/releases/tag/v0.4.1) - 2024-08-09

* feat(script): add script to help tagging the next release by @nobe4 in https://github.com/nobe4/gh-not/pull/132
* refactor(colors): switch to lipgloss by @nobe4 in https://github.com/nobe4/gh-not/pull/133


## [v0.4.0](https://github.com/nobe4/gh-not/releases/tag/v0.4.0) - 2024-08-09

* feat(views): implement logging for bubbletea interface by @nobe4 in https://github.com/nobe4/gh-not/pull/123
* feat(notifications): render notifications ahead of time by @nobe4 in https://github.com/nobe4/gh-not/pull/124
* feat(repl): start reworking the repl to use the list bubble by @nobe4 in https://github.com/nobe4/gh-not/pull/125
* refactor(normal): move handlers into separate file by @nobe4 in https://github.com/nobe4/gh-not/pull/127
* refactor(normal): refactor more methods by @nobe4 in https://github.com/nobe4/gh-not/pull/128
* refactor(normal): use the list help view exclusively by @nobe4 in https://github.com/nobe4/gh-not/pull/129
* refactor(repl): remove old repl by @nobe4 in https://github.com/nobe4/gh-not/pull/130
* refactor(repl): rework mappings for consistency by @nobe4 in https://github.com/nobe4/gh-not/pull/131


## [v0.3.6](https://github.com/nobe4/gh-not/releases/tag/v0.3.6) - 2024-07-23

* fix(config): use the config path to edit by @nobe4 in https://github.com/nobe4/gh-not/pull/122


## [v0.3.5](https://github.com/nobe4/gh-not/releases/tag/v0.3.5) - 2024-07-22

* fix(tests): use fixed times to sort on by @nobe4 in https://github.com/nobe4/gh-not/pull/119
* feat(notifications): add Meta RemoteExists field by @nobe4 in https://github.com/nobe4/gh-not/pull/120
* refactor(rules): filter directly on notifications by @nobe4 in https://github.com/nobe4/gh-not/pull/121


## [v0.3.4](https://github.com/nobe4/gh-not/releases/tag/v0.3.4) - 2024-07-22

* fix: load go version from mod file by @nobe4 in https://github.com/nobe4/gh-not/pull/112
* feat(repl): highlight current line better by @nobe4 in https://github.com/nobe4/gh-not/pull/113
* feat(repl): highlight current line better by @nobe4 in https://github.com/nobe4/gh-not/pull/114
* feat(json): add JSON output by @nobe4 in https://github.com/nobe4/gh-not/pull/115
* feat(root): add --flag to allow listing all notifications by @nobe4 in https://github.com/nobe4/gh-not/pull/116
* feat(root): add --rule flag by @nobe4 in https://github.com/nobe4/gh-not/pull/117


## [v0.3.3](https://github.com/nobe4/gh-not/releases/tag/v0.3.3) - 2024-07-21

* docs: add development notice by @nobe4 in https://github.com/nobe4/gh-not/pull/109
* feat: add sync stats and refresh note by @nobe4 in https://github.com/nobe4/gh-not/pull/110
* feat: add last updated date by @nobe4 in https://github.com/nobe4/gh-not/pull/111


## [v0.3.2](https://github.com/nobe4/gh-not/releases/tag/v0.3.2) - 2024-07-20

* fix(caller): ensure the caller is appropriately defined by @nobe4 in https://github.com/nobe4/gh-not/pull/108


## [v0.3.1](https://github.com/nobe4/gh-not/releases/tag/v0.3.1) - 2024-07-20

* fix(gh): paginate the correct number of pages by @nobe4 in https://github.com/nobe4/gh-not/pull/104
* docs(README): remove old flag for forcing a refresh by @nobe4 in https://github.com/nobe4/gh-not/pull/105
* feat(sync): add --force flag to apply rules even on Done notifications by @nobe4 in https://github.com/nobe4/gh-not/pull/106
* feat(manager): force strategy by @nobe4 in https://github.com/nobe4/gh-not/pull/107


## [v0.3.0](https://github.com/nobe4/gh-not/releases/tag/v0.3.0) - 2024-07-20

### Breaking changes

See https://github.com/nobe4/gh-not/pull/102

### What's Changed

* fix(manager): rework refresh to run only when necessary by @nobe4 in https://github.com/nobe4/gh-not/pull/100
* refactor(cmd): reorganize the commands  by @nobe4 in https://github.com/nobe4/gh-not/pull/102


## [v0.2.7](https://github.com/nobe4/gh-not/releases/tag/v0.2.7) - 2024-07-18

* feat(build): simplify build/install process by pulling debug info by @nobe4 in https://github.com/nobe4/gh-not/pull/96
* fix: only enrich non-done notifications by @nobe4 in https://github.com/nobe4/gh-not/pull/99


## [v0.2.6](https://github.com/nobe4/gh-not/releases/tag/v0.2.6) - 2024-07-08

* test(config,jq): add more tests to make explicit the order of operations by @nobe4 in https://github.com/nobe4/gh-not/pull/94
* feat(config, jq): add rule testing during config parsing by @nobe4 in https://github.com/nobe4/gh-not/pull/93


## [v0.2.5](https://github.com/nobe4/gh-not/releases/tag/v0.2.5) - 2024-07-06

* docs(readme): add ci and docs badge by @nobe4 in https://github.com/nobe4/gh-not/pull/79
* feat(dependabot): create dependabot config by @nobe4 in https://github.com/nobe4/gh-not/pull/80
* feat(jq,rule): tests and advanced filters by @nobe4 in https://github.com/nobe4/gh-not/pull/87
* fix(cmd): run manager.save only on commands that load and change it by @nobe4 in https://github.com/nobe4/gh-not/pull/90


## [v0.2.4](https://github.com/nobe4/gh-not/releases/tag/v0.2.4) - 2024-06-23

* fix(config): write only config's `data` to disk by @nobe4 in https://github.com/nobe4/gh-not/pull/76
* refactor(commands): various small improvements for the commands organization by @nobe4 in https://github.com/nobe4/gh-not/pull/77
* fix(config): struct annotation correction by @nobe4 in https://github.com/nobe4/gh-not/pull/78


## [v0.2.3](https://github.com/nobe4/gh-not/releases/tag/v0.2.3) - 2024-06-23

* tests(gh): write some tests around HTTP handling by @nobe4 in https://github.com/nobe4/gh-not/pull/70
* feat(gh): protect against all decoding errors by @nobe4 in https://github.com/nobe4/gh-not/pull/71
* feat(view): fixed size view by @nobe4 in https://github.com/nobe4/gh-not/pull/72
* docs(config): add documentation for various types by @nobe4 in https://github.com/nobe4/gh-not/pull/73


## [v0.2.2](https://github.com/nobe4/gh-not/releases/tag/v0.2.2) - 2024-06-22

* fix(actors): add space after the actor's verb by @nobe4 in https://github.com/nobe4/gh-not/pull/68
* feat(config): use viper for configuration with default values by @nobe4 in https://github.com/nobe4/gh-not/pull/69


## [v0.2.1](https://github.com/nobe4/gh-not/releases/tag/v0.2.1) - 2024-06-18

* fix(views): add conditional preventing failure on missing config by @nobe4 in https://github.com/nobe4/gh-not/pull/60
* feat(gh): don't enrich Done notifications by @nobe4 in https://github.com/nobe4/gh-not/pull/65
* feat(manager): don't apply rules on Done notifications by @nobe4 in https://github.com/nobe4/gh-not/pull/66


## [v0.2.0](https://github.com/nobe4/gh-not/releases/tag/v0.2.0) - 2024-06-15

### Breaking Change

- `Meta.ToDelete` is now `Meta.Done`

### What's Changed

* refactor(log): simplify logging and notification debug rendering by @nobe4 in https://github.com/nobe4/gh-not/pull/56
* refactor(notifications): rename ToDelete to Done by @nobe4 in https://github.com/nobe4/gh-not/pull/57


## [v0.1.5](https://github.com/nobe4/gh-not/releases/tag/v0.1.5) - 2024-06-15

* feat(version): improve version display by @nobe4 in https://github.com/nobe4/gh-not/pull/52
* feat(ci): check for deadcode by @nobe4 in https://github.com/nobe4/gh-not/pull/53
* feat(notifications): implement a better `sync`'ing logic  by @nobe4 in https://github.com/nobe4/gh-not/pull/55


## [v0.1.4](https://github.com/nobe4/gh-not/releases/tag/v0.1.4) - 2024-06-12

Thanks @offbyone for reporting #42

* feat(repl): add paginator view by @nobe4 in https://github.com/nobe4/gh-not/pull/40
* Update issue templates by @nobe4 in https://github.com/nobe4/gh-not/pull/43
* refactor(keymap): simplify the keymap by using the name as help by @nobe4 in https://github.com/nobe4/gh-not/pull/45
* feat(ci): create ci workflow by @nobe4 in https://github.com/nobe4/gh-not/pull/46
* fix(ci): fix CI using GitHub's editor by @nobe4 in https://github.com/nobe4/gh-not/pull/47
* feat(ci): run on ubuntu only by @nobe4 in https://github.com/nobe4/gh-not/pull/49
* feat(config,cache): ensure parent folder exists by @nobe4 in https://github.com/nobe4/gh-not/pull/44


## [v0.1.3](https://github.com/nobe4/gh-not/releases/tag/v0.1.3) - 2024-06-12

* docs(README): add section about running sync automatically by @nobe4 in https://github.com/nobe4/gh-not/pull/30
* feat(config): add `--init` flag by @nobe4 in https://github.com/nobe4/gh-not/pull/32
* feat(notifications): add read status in output by @nobe4 in https://github.com/nobe4/gh-not/pull/33
* feat(gh,config): config for all notifications and retry/page count by @nobe4 in https://github.com/nobe4/gh-not/pull/34
* feat(config, repl): configurable keymap by @nobe4 in https://github.com/nobe4/gh-not/pull/37
* docs(config): simplify documentation and init flag by @nobe4 in https://github.com/nobe4/gh-not/pull/39


## [v0.1.2](https://github.com/nobe4/gh-not/releases/tag/v0.1.2) - 2024-06-11

* refactor: use io.Writer in actors by @nobe4 in https://github.com/nobe4/gh-not/pull/26
* feat(actor,notifications): open notification from the subject URL by @nobe4 in https://github.com/nobe4/gh-not/pull/29


## [v0.1.1](https://github.com/nobe4/gh-not/releases/tag/v0.1.1) - 2024-06-11

* docs: add a how it works 101 section in the README by @nobe4 in https://github.com/nobe4/gh-not/pull/18
* feat: support no-config run by @nobe4 in https://github.com/nobe4/gh-not/pull/22
* refactor(manager): rework how the refresh condition is handled by @nobe4 in https://github.com/nobe4/gh-not/pull/25
* feat(config): add edit flag by @nobe4 in https://github.com/nobe4/gh-not/pull/24

Thanks to @williammartin and @alondahari for feedback

## [v0.1.0](https://github.com/nobe4/gh-not/releases/tag/v0.1.0) - 2024-06-11

* refactor: use native bubbletea help and keybinding by @nobe4 in https://github.com/nobe4/gh-not/pull/10
* refactor: split views into their own packages by @nobe4 in https://github.com/nobe4/gh-not/pull/11
* feat: create an interface for the API Caller by @nobe4 in https://github.com/nobe4/gh-not/pull/12
* refactor(notifications): use pointers for Notifications by @nobe4 in https://github.com/nobe4/gh-not/pull/14
* refactor: rename DeleteNil and ToInterface by @nobe4 in https://github.com/nobe4/gh-not/pull/16
* feat: notifications manager by @nobe4 in https://github.com/nobe4/gh-not/pull/17


## [v0.0.4](https://github.com/nobe4/gh-not/releases/tag/v0.0.4) - 2024-05-27

* feat: add version in release by @nobe4 in https://github.com/nobe4/gh-not/pull/5
* feat: implement simple REPL by @nobe4 in https://github.com/nobe4/gh-not/pull/9
