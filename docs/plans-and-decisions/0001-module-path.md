<!--
Copyright 2026 Darko Demić.
Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.
-->

# 0001 — Go module path: github.com/darkodemic/systray

- **Status:** Accepted 2026-10-07. The module is `github.com/darkodemic/systray` instead of `fyne.io/systray`.
- **Date:** 2026-10-07
- **Supersedes:** — / **Superseded by:** —
- **Owner:** Darko
- **Related:** `NOTICE` (provenance and copyright of the fork), `AGENTS.md` (branches for fyne-io PRs start from upstream `master`)

## Context

- The fork inherited `module fyne.io/systray` from fyne-io/systray. The `fyne.io` domain belongs to the Fyne project, so under that name the fork looks like their package.
- With that path, `go get github.com/darkodemic/systray` does not work: Go rejects the module because it declares itself as `fyne.io/systray`. The fork can therefore only be used through a `replace` directive. gpwebcam uses it exactly that way: `replace fyne.io/systray => github.com/darkodemic/systray v1.12.3-0.20261006205618-9c45f672f861` (checked 2026-10-07).
- The README described integration with the Fyne toolkit and `fyne package`, which is not relevant to this fork.

## Decision

1. `go.mod` declares `module github.com/darkodemic/systray`, and all internal imports, the example and the README use that path.
2. The Fyne-specific parts of the README are removed: the section on Fyne apps, the `fyne package` hints and the links to developer.fyne.io and pkg.go.dev/fyne.io.
3. The provenance stays on record. The README (intro and Credits), `NOTICE` and `AGENTS.md` still mention fyne-io/systray, because that is attribution, not branding.

## Consequences

Positive:

- The fork is used directly, `go get github.com/darkodemic/systray`, without a `replace` directive.
- pkg.go.dev shows the fork's documentation under its own name.

Negative:

- This is a breaking change. Anyone who imports `fyne.io/systray` with a `replace` to the fork has to change the imports to the new path and remove the `replace` when moving to a version after this change. An older version a project is already pinned to (gpwebcam at 9c45f67) keeps working.
- Pulling changes from upstream causes conflicts in the import lines of five files (`systray_unix.go`, `systray_menu_unix.go`, `systray_notifier_unix.go`, `systray_unix_test.go`, `example/main.go`) whenever upstream changes those lines.

Risks:

- A PR to fyne-io made from a branch based on the fork's `master` would carry the new path. The rule in `AGENTS.md` that such branches start from upstream `master` already prevents this.

## Alternatives considered

- **Keep `fyne.io/systray`.** Fewer conflicts with upstream, but the fork would stay usable only through a `replace` directive and would carry someone else's domain.
- **Vanity path (for example `darkodemic.com/systray`).** Independent of GitHub, but it requires hosting `go-import` meta tags, and there is no need for that now.

## Out of scope

- Moving gpwebcam to the new path.
- The `Makefile` target `tag-changelog`, which generates the changelog from getlantern/systray. Handled separately in issue #6, which removes the `Makefile` and the stale `CHANGELOG.md`.
