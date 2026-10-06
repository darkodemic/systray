<!--
Copyright 2026 Darko Demić.
Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.
-->

# Agent instructions

## Copyright header on changed files

This repository is a fork of [fyne-io/systray](https://github.com/fyne-io/systray),
which is itself a fork of [getlantern/systray](https://github.com/getlantern/systray).
It is licensed under the Apache License 2.0 (`LICENSE`), and the copyright of the
fork is recorded in `NOTICE`. Section 4(b) of the license requires every modified
file to carry a prominent notice that it was changed.

Whenever you change a file that came from upstream, make sure it starts with this
header, written in the file's comment syntax:

```go
// Modifications Copyright 2026 Darko Demić.
// Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.
```

Keep the word "Modifications": it is what states that the file was changed, so
dropping it would no longer meet section 4(b).

A file you create from scratch gets this header instead:

```go
// Copyright 2026 Darko Demić.
// Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.
```

Rules:

- Use the current year. If the file already has the header and the year is older,
  extend it to a range (`2026-2027`) instead of adding a second header.
- Comment syntax: `//` in `.go`, `.m`, `.h` and `.c` files, `#` in `Makefile` and
  shell scripts, `<!-- -->` in `.md` and `.xml` files.
- Go files, and any file with a `//go:build` line: the header goes on the very
  first lines, followed by a blank line, before the `//go:build` line and the
  `package` clause. The blank line keeps the header out of the package
  documentation.
- XML files: the header goes right after the `<?xml ... ?>` declaration, which
  must stay on the first line.
- These files never get a header: `LICENSE`, `NOTICE`, `CLAUDE.md` (it only
  imports this file), `go.mod`, `go.sum`, generated code under
  `internal/generated/` (regenerate it instead of editing it), and binary files
  such as images and icons.
- A branch for a pull request to fyne-io/systray starts from fyne-io's `master`
  and gets no headers, because upstream does not carry this notice.

## Language

Everything in this repository is written in English: code, comments, the README
and the plan and decision documents in `docs/plans-and-decisions/`, whatever
language the conversation is in.
