# nothelp

> Helper tool for managing notes.

## Commands

- init
  - create config with davy defaults
- start
  - check for existance of today's note
  - if doesn't exist then:
    - make a copy of the template
    - update metadata, date/time etc
    - save file with date in the configured folder
  - open nvim on the line of the beginning of the start checklist (configurable?)
- stop
- today/inbox

## Requirements

- CLI first.
- Use simple markdown.
- Configuration in XDG_CONFIG_HOME.
  - directory to store notes.
  - template to use.
- Create a daily note based on a template.
- support start, stop, today and inbox workflows (davy knows).
- easy to add short snippets from CLI to daily note.
  - extension: should also support piping
- check for existance of daily note before doing anything.
- allow for 
- add tracking of habits?
- use vim as the editor of choice (makes it easier to go to specific lines)
- nvim plugin
  - command to open the note.
  - telescope searching.
