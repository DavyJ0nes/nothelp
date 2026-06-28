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
