# nothelp

A simple helper tool for managing daily notes.

It uses an embeded template that is the daily note along with my checklists that
help keep me organised.

## Commmands

- `start`: Kick off the day by opening the daily note to the start of the
  morning checklist.
- `stop`: End the day by running the evening checklist.
- `today`: Open the notes section for today's note.
- `yesterday`: Open the notes section for yesterday's note.
- `weekly`: Open (or create) this week's review note, named by ISO week (e.g.
  `2026-W27`).

## Local development

Tooling and tasks are managed with [mise](https://mise.jdx.dev): it pins the Go
toolchain and dev tools (see `mise.toml`) and runs every task CI runs.

```sh
brew install mise   # or see https://mise.jdx.dev/installing-mise.html
mise install        # install the pinned toolchain + tools
mise run ci         # run the full CI-parity suite locally
```

`mise run ci` is the strict suite CI enforces:

| Task | What it does |
| --- | --- |
| `mise:check` | Verify `mise.lock` is in sync |
| `fmt:check` | Markdown is formatted (rumdl) |
| `fmt:yaml:check` | YAML is formatted (yamlfmt) |
| `lint:scripts` | shellcheck on any shell scripts |
| `mod:check` | `go.mod` / `go.sum` are tidy |
| `lint:go` | golangci-lint |
| `build` | `go build ./...` |
| `test` | `go test -race ./...` |

Local-only helpers (not part of `ci`): `fmt:md` / `fmt:yaml` / `fmt:go` (format
in place), `lint:go:fix`, `mod:tidy`, `lint:commits` (convco), and `mise:lock`.
Run `mise tasks` to list them all.

Tool, Go-module, and GitHub-Actions updates are managed by Renovate, which opens
a single grouped PR on the first Tuesday of each month. `mise.lock` is committed
and regenerated on tool bumps via `mise lock`.

## Ideas

- easy to add short snippets from CLI to daily note.
- sync notes repo and push updates.
- add tracking of habits
- generate stats:
  - completed the checklist?
  - semantic analysis of what was learned.
- monthly / quarterly review notes following the same template approach.
- auto-link a weekly review to that week's daily notes.
- nvim plugin
  - command to open the note.
  - telescope searching.
