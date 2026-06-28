# CLAUDE.md

Guidance for AI assistants working in this repo. Keep changes minimal, match the
existing patterns, and verify with `mise run ci` before yielding.

## What this is

`nothelp` — a personal Go CLI (built on [cobra](https://github.com/spf13/cobra))
for daily notes plus a few habit/training trackers. It is macOS-oriented: it
opens files in `nvim`, controls apps with `osascript`, and writes into the local
Obsidian iCloud vault. Single Go module.

## Layout

- `main.go` — entrypoint; runs `cmd.RootCmd().Execute()` and exits non-zero on
  error.
- `cmd/` — one file per command. Each exposes `XxxCmd() *cobra.Command` whose
  `RunE` is a thin wrapper delegating to a `xxxRun()` function. Shared helpers
  live in `cmd/cmd.go`. Register new commands in `cmd/root.go` (keep the list
  alphabetical).
- `internal/templates/` — embedded note templates (`*.md`) and their rendering.
- `internal/config/` — filesystem paths for the notes vault.

## Conventions

- Command shape: `XxxCmd()` builder + `xxxRun()` logic. Keep `RunE` a one-liner;
  put real work in `xxxRun()`. Multi-part commands (see `pressup`, `training`)
  add subcommands via `cmd.AddCommand(...)`.
- Errors: return them up the stack (plain or wrapped); `main` prints and exits.
  Do not use `log.Fatal` or `panic` in command logic.
- Dates: the format is `dateFormat = "2006-01-02"`. Use the helpers in
  `cmd/cmd.go` — `todayDate()`, `yesterdayDate()`, `thisWeek()` (ISO week,
  `YYYY-Www`) — rather than formatting inline.
- Note file flow: `openNoteFile` (daily) and `openWeeklyNoteFile` (weekly)
  resolve the path (create from template if missing, else reuse; daily also
  checks the archive), then open it in `nvim` at the line matching a heading.
  Reuse these helpers; don't reimplement the create/open dance.
- Notes are written `0o600` under the directories from `internal/config`
  (`notes/daily`, `notes/weekly`, `notes/daily/archive`). Tracker state
  (`pressup`, `training`) is JSON in the vault.

## Embedded templates (read before editing)

- Templates are **data, not docs**. They are embedded via
  `//go:embed daily_template.md weekly_template.md` in
  `internal/templates/templates.go`.
- They are rendered with `text/template` using a per-call `FuncMap` closure, so
  placeholders are bare function calls: the daily template uses `{{date}}` and
  the weekly template uses `{{week}}` (NOT `{{.Date}}`/`{{.Week}}`). If you add
  a placeholder, register its function in the matching `Parse*` function.
- `.rumdl.toml` excludes `internal/templates/`. Never let the markdown formatter
  reflow these files — it changes the rendered note output. If you rename a
  template, update both the `//go:embed` directive and the `ParseFS` call.

## Tooling (mise)

- Everything runs through [mise](https://mise.jdx.dev). Run `mise install` once,
  then `mise run <task>`. See `mise.toml` for the full task list.
- **Always run `mise run ci` before finishing.** It is the exact CI suite:
  `mise:check`, `fmt:check`, `fmt:yaml:check`, `lint:scripts`, `mod:check`,
  `lint:go`, `build`, `test`. Fix failures at the source — never suppress.
- Linting is golangci-lint v2 with default linters (errcheck included). Handle
  or explicitly discard (`_ =`) errors; don't add `//nolint` without a reason.
- Useful local-only helpers: `fmt:md`, `fmt:yaml`, `fmt:go`, `lint:go:fix`,
  `mod:tidy`.

## Tests

- Run with `mise run test` (`go test -race ./...`). Keep tests beside the code
  (`*_test.go`).
- `internal/config` uses `gotest.tools/v3/assert`; other packages use plain
  `testing`. Match the package you're in.
- Use `t.TempDir()` for filesystem tests. Never touch the real vault and never
  spawn `nvim` from a test.

## Commits & releases

- [Conventional Commits](https://www.conventionalcommits.org) are required
  (`convco`; PR CI runs `lint:commits`). On push to `main`, `mise run release`
  (`.mise/scripts/release.sh`) computes the next version from commit history
  (`feat` → minor, `fix`/`perf` → patch, breaking → minor while 0.x), tags it,
  and publishes binaries + notes with GoReleaser (`.goreleaser.yml`). So the
  commit type drives the release — preview with `mise run release:dry`. Override
  the version with `RELEASE_VERSION=X.Y.Z` (env var, or the workflow_dispatch
  `version` input) to skip the auto bump.
- Commit signing (GPG) is enabled — prefer to let the human commit. If you must
  commit to validate, use a throwaway commit and `git reset --soft` afterward.

## Dependencies

- mise tools, Go modules, and GitHub Actions are updated by Renovate (one
  grouped PR, first Tuesday of each month). Don't hand-edit `mise.lock`;
  regenerate it with `mise run mise:lock`.
