<!--
Copyright 2026 Darko Demić.
Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.
-->

# CI/CD for the fork

- **Status:** Accepted 2026-10-07 with Darko's answers in §7. CI, the SBOM with Grype, Dependabot and CODEOWNERS are built and green (§3–§6); the D-Bus integration test (§8) is next.
- **Date:** 2026-10-07
- **Owner:** Darko
- **Related:** `0001-module-path.md`; issue #5 (FreeBSD build); gpwebcam's `.github/workflows/ci.yml`, whose layout (actions pinned to commits, `contents: read`, oldest and newest Go) this plan follows.

## 1. Where we are

Checked 2026-10-07 on `main` at c773042, which is release `v1.13.0`.

| Area | State |
|---|---|
| CI | None. The fork has no `.github` directory, and fyne-io/systray has no workflows either. |
| Ruleset "Main protection" on the default branch, `main` (renamed from `master` on 2026-10-07) | Pull request with one approval, code owner review, squash only, linear history and signed commits. "Required status checks" is on, but its list is empty. Darko bypasses the ruleset. There is no `CODEOWNERS` file, so the code owner rule has nobody to ask. |
| Formatting | `gofmt -l` flags three files inherited from upstream: `example/icon/iconunix.go`, `example/icon/iconwin.go` and `systray_test.go`. |
| `go vet` | Passes for linux and freebsd. Fails for windows: `systray_windows_test.go:47` still calls `addOrUpdateMenuItem` with its old four arguments; the function takes seven. Broken upstream as well. |
| Tests | Seven tests. Linux runs five, `-race` passes, and they cover 34.9% of statements. The two Windows tests do not compile. Nobody has built or tested macOS (cgo, needs a Mac). |
| Uncovered on Linux | Everything that needs a session bus: `registerSystray`, `nativeStart`, the dbusmenu methods (`GetLayout`, `GetProperty`, `Event`, `AboutToShow`), `SetIcon`, `SetIconName`, `SetIconThemePath`, `hideMenuItem` and `showMenuItem`. The radio items have one test, which checks only the `toggle-type` and `toggle-state` they produce. |
| Cross-builds | linux, freebsd, openbsd, netbsd and windows build without cgo. darwin needs cgo and the macOS SDK. |
| Vulnerabilities | `govulncheck ./...`: "No vulnerabilities found" (godbus v5.2.2, x/sys v0.27.0). |
| Dependency watch | No `dependabot.yml`. Dependabot alerts and Dependabot security updates are off, and so is the dependency graph (its SBOM export API returns 404). Secret scanning and push protection are on. |
| Releases | `v1.13.0` was made by hand: a GPG-signed annotated tag and a GitHub release with written notes. |

## 2. What is missing

1. **Checks on pull requests.** Nothing runs on a PR or on `main`, which is why the ruleset's list of required checks is empty.
2. **Platform coverage.** Only Linux is ever built by a person. The Windows tests do not compile, macOS has never been built, and the BSD build broke unnoticed until issue #5.
3. **Dependency and vulnerability watch.** No version updates, no alerts and no vulnerability scan. A new advisory for godbus or x/sys would go unnoticed.
4. **Tests of the D-Bus side.** The code that registers the StatusNotifierItem and serves the menu, which is what gpwebcam depends on, runs in no test.

## 3. Pull request checks

`.github/workflows/ci.yml` runs on every pull request, on every push to `main`, every Monday at 05:17 UTC on `main` and on demand. Like gpwebcam, it pins actions to commits and runs with `permissions: contents: read`. Go comes from `actions/setup-go`: 1.20, the `go` directive in `go.mod` and the oldest supported Go, and `stable`. A library has no release toolchain to pin, so mise is not needed here.

| Job | Runner | Steps |
|---|---|---|
| Lint | ubuntu-24.04 | `gofmt -l`, then `go vet ./...` for linux, windows and freebsd (`GOOS=...`). |
| Test (Go 1.20), Test (Go stable) | ubuntu-24.04 | `go test -race` with coverage; the total goes into the job summary. |
| Cross-build | ubuntu-24.04 | `go build ./...` for freebsd, openbsd, netbsd and windows. |
| Windows | windows-2025 | `go vet ./...` and `go test ./...`. |
| macOS | macos-15 | `go vet`, `go build` and `go test` with cgo. |
| govulncheck | ubuntu-24.04 | `govulncheck ./...` for linux, windows and freebsd. |
| SBOM | ubuntu-24.04 | The CycloneDX SBOM from §5, uploaded as a workflow artifact and scanned with Grype. |

All of them are required checks (Darko, §7). Their names go into the ruleset once the workflow has run green on `main`.

The same PR fixes what the checks would fail on: it runs `gofmt` on the three upstream files (`systray_test.go` had CRLF line endings) and updates `systray_windows_test.go` to the current signatures of `addOrUpdateMenuItem`, `addSeparatorMenuItem` and `hideMenuItem`. While at it, it adds a checked radio item to that test.

The first run on the runners (2026-10-07) found two more problems, both in upstream tests that had never run on those platforms:

- **Windows:** `TestMenuItem_Remove`, the one test for all platforms, ran first, called `quit()` and returned without waiting for the tray to go away. The window class "SystrayClass" was still registered when `TestBaseWindowsTray` registered its own, which failed with "Class already exists". After that, `TestWindowsRun` hung until the 5-minute timeout. The test now waits until `Run` has returned.
- **macOS:** `TestMenuItem_Remove` calls `Run` on a goroutine, but the Cocoa event loop only works on the main thread, so the tray never became ready and the test hung. It is skipped on macOS with that reason. The macOS job therefore checks vet and the cgo build; it runs no test until there is one that can run off the main thread.

## 4. Vulnerabilities and dependencies

- **govulncheck** is the gate for vulnerabilities in Go code. It reads the Go vulnerability database and reports a vulnerability only when the code reaches an affected function, so it is run per platform.
- **Grype** scans the SBOM (§5) and reports every advisory for a module version, whether or not this module calls the affected code. Every finding goes into the job summary; only high and critical ones fail the job. On 2026-10-07 it reports GO-2026-5024 (CVE-2026-39824, low): an integer overflow in `windows.NewNTUnicodeString` in `golang.org/x/sys` before v0.44.0. systray does not call that function, which is why govulncheck reports nothing.
- **`.github/dependabot.yml`** updates `gomod` and `github-actions` monthly, like gpwebcam (Darko, §7). The action updates also keep the commit pins current.
- **x/sys stays below v0.31.0.** v0.31.0 and later require Go 1.23, and v0.44.0, the fix for GO-2026-5024, requires Go 1.25 (proxy.golang.org, checked 2026-10-07). Moving there would raise the `go` directive of this module, and with it of gpwebcam, which builds with Go 1.22. `dependabot.yml` therefore ignores x/sys from v0.31.0 on. This gets lifted when the oldest supported Go moves; a high or critical x/sys advisory that the code reaches would force that decision earlier.
- **Repository settings**, which need no PR: the dependency graph, Dependabot alerts and Dependabot security updates are on. Security updates open a PR right away for a vulnerable dependency, outside the monthly schedule.

## 5. SBOM

Decided (Darko, §7): a full CycloneDX SBOM. CI makes one on every run and keeps it as a workflow artifact, and `.github/workflows/release.yml` attaches one to every published release as `systray-<tag>.cdx.json`.

Darko suggested Syft. Both tools were tried on `v1.13.0` on 2026-10-07:

| | Syft 1.54.1 (`syft scan dir:.`) | cyclonedx-gomod v1.12.0 (`mod -type library -licenses -assert-licenses`) |
|---|---|---|
| Main component | Type `file`, the module itself listed as a separate component with version `UNKNOWN` | Type `library`, version `v1.13.0` taken from the git tag, purl with the version |
| Licenses | Found only with remote or module cache license search turned on | Apache-2.0, BSD-2-Clause and BSD-3-Clause, from the module cache |
| Hashes | None | SHA-256 for each dependency, from `go.sum` |
| Dependency graph | One edge, from the main module | Full: systray → godbus → x/sys |

For a Go library, cyclonedx-gomod gives the more complete SBOM, so this repository uses it. Syft fits gpwebcam, which ships binaries and packages: GoReleaser's `sboms` section runs Syft for every archive and package, and Syft reads a Go binary's build information. Grype reads the SBOMs of both tools.

## 6. Release workflow

The tag stays manual, because it is signed with Darko's key, and so do the release notes. `.github/workflows/release.yml` runs when a release is published, or by hand with a tag for an older release:

1. It checks that the tag looks like `vX.Y.Z`.
2. It generates the SBOM at that tag and attaches it to the release.
3. It requests the version from proxy.golang.org, so pkg.go.dev lists it right away.

After the merge, a manual run with `v1.13.0` attaches the SBOM to the first release.

## 7. Decisions

Answered by Darko on 2026-10-07:

1. **Windows and macOS jobs** are required checks.
2. **SBOM:** full CycloneDX: a workflow artifact on every run and a release asset (§5).
3. **Linter:** `go vet` is enough for now.
4. **Dependabot:** monthly.
5. **CODEOWNERS:** `* @darkodemic`, so the ruleset's code owner rule applies.

## 8. Tests

- **D-Bus integration test.** In CI it runs under `dbus-run-session`, which starts a private session bus. It starts the tray with `RunWithExternalLoop`, then reads it back through godbus as a client, the way a panel does:
  - the StatusNotifierItem properties (title, tooltip, icon, icon name and theme path);
  - the dbusmenu layout: plain, checkbox and radio items, submenus, disabled and hidden items, separators and shortcuts;
  - a click sent with `Event`, which must arrive on the item's `ClickedCh`;
  - `LayoutUpdated` and `ItemsPropertiesUpdated` after the code changes an item.

  Without a session bus (`DBUS_SESSION_BUS_ADDRESS` unset), the test is skipped, so `go test ./...` keeps working everywhere.
- **One tray per process.** On Linux `quit` closes a package-level channel that is never recreated, so `Run` works once per test binary: `TestMenuItem_Remove` panics under `go test -count=2`, on `main` too (checked 2026-10-07). The integration test must take that into account, for example by running in its own test binary or by resetting that state in the test.

## 9. Order of work

1. One PR: the Go fixes from §3, `ci.yml`, `release.yml`, `dependabot.yml`, `CODEOWNERS`, this plan, and the one mention of the fork's `master` in `0001-module-path.md`, which is now `main`. The fixes and the checks travel together, so the PR shows green on its own and the merge order of separate PRs cannot leave `main` red.
2. Settings: the dependency graph, alerts and security updates right away; the required checks in the ruleset after the first green run on `main`; then a manual run of the release workflow for `v1.13.0`.
3. The D-Bus integration test (§8), as its own PR.

## Where we are and what is next

2026-10-07: decisions recorded in §7. The PR from §9.1 (#10) is green on all eight checks after the test fixes in §3; Grype reports only GO-2026-5024 (low). The dependency graph, Dependabot alerts and security updates are on. Next, after the merge: the required checks in the ruleset and the release workflow run for `v1.13.0` (§9.2), then the D-Bus integration test (issue #11).
