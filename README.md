# nothelp

A small personal CLI for managing daily notes — plus a few habit and training
trackers — from the terminal. Notes are generated from embedded templates and
opened in Neovim.

> macOS-oriented: `start` / `stop` drive desktop apps via `osascript`, notes are
> written under the local Obsidian iCloud vault, and files open in `nvim`.

## Commands

### Daily notes

- `start`: Quit distracting apps, launch work apps, and open today's note at the
  Start checklist.
- `stop`: Quit work apps and open today's note at the Shutdown checklist.
- `today` (alias `inbox`): Open today's note at the Log section.
- `notes`: Open today's note at the Log section (same target as `today`).
- `yesterday`: Open yesterday's note at the Log section.
- `weekly` (alias `week`): Open (or create) this week's review note, named by
  ISO week (e.g. `2026-W27`).
- `archive`: Move daily notes into the archive folder.

Any of these create the note from the embedded template if it doesn't exist yet.

### Trackers

- `pressup log <morning|evening> <rounds>` / `pressup list`: Log and review the
  daily 50-pressup challenge (state kept as JSON in the vault).
- `training list [week|today|tomorrow]` / `training log [day]`: View and tick
  off the 19-week Vätternrundan training plan.

### Review

- `review [-o <path>]`: Generate a 7-day review of recent daily notes using an
  LLM. Requires `OPENAI_API_KEY`; optionally set `NOTHELP_OPENAI_MODEL` (default
  `gpt-4o-mini`) and `NOTHELP_OPENAI_BASE_URL` (default the OpenAI API).

## How it works

- Paths come from `internal/config`: notes live under
  `~/Library/Mobile Documents/iCloud~md~obsidian/Documents/notes` in `daily/`,
  `weekly/`, and `daily/archive/`.
- Daily notes are named `YYYY-MM-DD.md`, weekly notes `YYYY-Www.md` (ISO week).
- Templates are embedded (`internal/templates/*.md`) and rendered with Go's
  `text/template`; the command then opens the file in `nvim` at the relevant
  heading.

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

Commits follow [Conventional Commits](https://www.conventionalcommits.org)
(validated by `convco`); releases are cut from commit history by
go-semantic-release on push to `main`.

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
