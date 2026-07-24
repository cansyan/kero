# Repository Guidelines

## Project Structure & Module Organization

`kero` is a small Go terminal UI library. Core implementation files live in the repository root as package `kero`, including `app.go`, `program.go`, `terminal.go`, `frame.go`, `screen.go`, `style.go`, `geometry.go`, and `textinput.go`. Tests sit beside the code as `*_test.go` files. Runnable examples belong under `example/<name>/main.go`; current examples include `example/counter` and `example/bouncer`. Design intent and architectural notes are in `DESIGN.md`.

## Build, Test, and Development Commands

- `go test ./...`: run all package and example tests.
- `go run ./example/counter`: launch the counter example.
- `go run ./example/bouncer`: launch the animation example.
- `go test ./... -run TestName`: run a focused test while iterating.

There is no separate build system; use standard Go tooling.

## Coding Style & Naming Conventions

Use idiomatic Go formatted with `gofmt`. Keep implementation in the root `kero` package unless there is clear pressure for a subpackage. Prefer plain structs, small interfaces, explicit control flow, and readable code over clever abstractions. Exported API names should be short and descriptive, such as `Frame`, `Rect`, `Program`, and `TextInput`. Examples should import the local module path with `import "kero"`.

## Testing Guidelines

Use Go’s standard `testing` package. Add focused tests for geometry math, frame clipping and drawing behavior, style serialization, input parsing, widgets, and program loop changes. Name tests by behavior, for example `TestRectInset` or `TestSplitVertical`. Prefer direct expected values and small fixtures.

## Commit & Pull Request Guidelines

Recent commits use short imperative messages, for example `Add simple animation example`. Follow that style: describe the change, not the process. Pull requests should include a concise summary, note any API or terminal behavior changes, list tests run, and include screenshots or terminal recordings when changing examples or rendering behavior.

## Security & Configuration Tips

Terminal cleanup is important. Changes in `program.go` or `terminal.go` must preserve restoration of raw mode, alternate screen, cursor visibility, and mouse state on errors. Avoid hidden global state and keep runtime options explicit in `Options`.
