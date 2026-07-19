# kero Design Document

`kero` is a small terminal UI framework for Go.

The goal is not to hide how terminal UIs work. The goal is to give the user a
clear structure for writing them: receive events, update state, render cells,
repeat.

## Values

- Simple enough to read in one sitting.
- Explicit control flow.
- Small interfaces over large inheritance trees.
- Plain Go structs over clever abstractions.
- Good defaults, but few hidden behaviors.
- Easy to debug by printing state or inspecting frames.

## Non-Goals

- A full desktop-style widget system.
- CSS-like styling.
- A retained-mode scene graph.
- Complex async runtime machinery.
- Supporting every terminal capability in the first version.

## Mental Model

Every kero program is a loop:

1. Read input from the terminal.
2. Convert input into an `Event`.
3. Send the event to the application.
4. Let the application update its state.
5. Ask the application to draw itself into a frame.
6. Flush the frame to the terminal.

Applications own their state. kero owns terminal setup, event reading, frame
diffing, and drawing.

## Package Shape

Initial module:

```text
kero
```

Initial directory layout:

```text
.
├── go.mod
├── app.go
├── context.go
├── event.go
├── frame.go
├── geometry.go
├── program.go
├── screen.go
├── style.go
├── terminal.go
└── example
    └── counter
        └── main.go
```

Implementation code should live in the repository root as package `kero`.
Examples should live under `example/`.

Example programs should import the local module path:

```go
import "kero"
```

The first version should stay in one root package until clear pressure appears.
Future packages like `input`, `ansi`, `layout`, or `widgets` should wait until
the root package becomes meaningfully crowded.

## Core Types

### App

An `App` is the user-provided program.

```go
type App interface {
    Init(ctx *Context) error
    Update(ctx *Context, ev Event) error
    View(ctx *Context, f *Frame)
}
```

`Init` runs once after the terminal is ready.

`Update` receives input, resize events, ticks, and internal messages.

`View` draws the current state. It should avoid mutating application state.

### Context

`Context` gives the app controlled access to the runtime.

```go
type Context struct {
    Width  int
    Height int

    done bool
}

func (c *Context) Quit()
func (c *Context) Size() Size
```

Early versions should keep this tiny. Later, it may expose timers, commands,
clipboard access, logging, or focus state.

### Program

`Program` owns the terminal lifecycle and event loop.

```go
type Program struct {
    app App
    opts Options

    screen *Screen
    ctx    Context
}

func New(app App, opts ...Option) *Program
func (p *Program) Run() error
```

The program should:

- Enter raw terminal mode.
- Switch to the alternate screen when configured.
- Hide the cursor while running.
- Restore the terminal on exit, even after errors.
- Read events.
- Render frames.

Pseudo-code:

```go
func (p *Program) Run() error {
    if err := p.terminal.Enter(); err != nil {
        return err
    }
    defer p.terminal.Leave()

    size, err := p.terminal.Size()
    if err != nil {
        return err
    }

    p.ctx.Width = size.Width
    p.ctx.Height = size.Height
    p.screen.Resize(size.Width, size.Height)

    if err := p.app.Init(&p.ctx); err != nil {
        return err
    }

    p.render()

    for !p.ctx.done {
        ev, err := p.terminal.ReadEvent()
        if err != nil {
            return err
        }

        switch e := ev.(type) {
        case ResizeEvent:
            p.ctx.Width = e.Width
            p.ctx.Height = e.Height
            p.screen.Resize(e.Width, e.Height)
        }

        if err := p.app.Update(&p.ctx, ev); err != nil {
            return err
        }

        p.render()
    }

    return nil
}

func (p *Program) render() error {
    f := p.screen.Frame()
    f.Clear()

    p.app.View(&p.ctx, f)

    return p.screen.Flush()
}
```

### Options

Options should be boring and explicit.

```go
type Options struct {
    AltScreen bool
    Mouse     bool
    // FPS limits timed redraws, such as spinners, clocks, and animations.
    // Event-only apps can ignore it and render only after input or resize.
    FPS       int
}

type Option func(*Options)

func WithAltScreen(v bool) Option
func WithMouse(v bool) Option
func WithFPS(fps int) Option
```

Default options:

```go
Options{
    AltScreen: true,
    Mouse:     false,
    FPS:       30,
}
```

## Events

Events are values. They should be easy to switch on.

```go
type Event interface {
    event()
}
```

Core event types:

```go
type KeyEvent struct {
    Key  Key
    Rune rune
    Mod  Mod
}

type ResizeEvent struct {
    Width  int
    Height int
}

type TickEvent struct {
    Time time.Time
}

type MouseEvent struct {
    X      int
    Y      int
    Button MouseButton
    Action MouseAction
    Mod    Mod
}
```

Keys:

```go
type Key int

const (
    KeyUnknown Key = iota
    KeyRune
    KeyEsc
    KeyEnter
    KeyBackspace
    KeyTab
    KeyUp
    KeyDown
    KeyLeft
    KeyRight
    KeyHome
    KeyEnd
    KeyPgUp
    KeyPgDown
    KeyDelete
)
```

Modifiers:

```go
type Mod uint8

const (
    ModNone Mod = 0
    ModCtrl Mod = 1 << iota
    ModAlt
    ModShift
)
```

## Drawing

kero draws a grid of cells.

```go
type Cell struct {
    Ch    rune
    Style Style
}

type Frame struct {
    width  int
    height int
    cells  []Cell
}
```

Frame API:

```go
func (f *Frame) Size() Size
func (f *Frame) Clear()
func (f *Frame) Set(x, y int, ch rune, s Style)
func (f *Frame) Write(x, y int, text string, s Style)
func (f *Frame) Fill(r Rect, ch rune, s Style)
func (f *Frame) Box(r Rect, s Style)
```

Coordinates are zero-based:

```text
(0,0) is the top-left cell.
x increases to the right.
y increases downward.
```

Drawing outside the frame should be safely clipped.

## Geometry

```go
type Size struct {
    Width  int
    Height int
}

type Point struct {
    X int
    Y int
}

type Rect struct {
    X int
    Y int
    W int
    H int
}

func (r Rect) Right() int
func (r Rect) Bottom() int
func (r Rect) Inset(n int) Rect
func (r Rect) Contains(p Point) bool
```

## Style

Styles should be compact and composable.

```go
type Style struct {
    Fg   Color
    Bg   Color
    Attr Attr
}

func NewStyle() Style
func (s Style) Foreground(c Color) Style
func (s Style) Background(c Color) Style
func (s Style) Bold() Style
func (s Style) Underline() Style
```

Colors:

```go
type Color int

const (
    ColorDefault Color = iota
    ColorBlack
    ColorRed
    ColorGreen
    ColorYellow
    ColorBlue
    ColorMagenta
    ColorCyan
    ColorWhite
)
```

Attributes:

```go
type Attr uint16

const (
    AttrNone Attr = 0
    AttrBold Attr = 1 << iota
    AttrUnderline
    AttrReverse
    AttrDim
)
```

## Screen

`Screen` is the low-level terminal renderer.

```go
type Screen struct {
    out io.Writer

    width  int
    height int

    prev Frame
    next Frame
}

func NewScreen(out io.Writer, w, h int) *Screen
func (s *Screen) Resize(w, h int)
func (s *Screen) Frame() *Frame
func (s *Screen) Flush() error
```

The first renderer can redraw the full screen every frame. A later version can
diff `prev` and `next` to write only changed cells.

## Terminal

Terminal handling should be isolated behind a small interface.

```go
type Terminal interface {
    Enter() error
    Leave() error
    ReadEvent() (Event, error)
    Size() (Size, error)
}
```

The concrete implementation can use `golang.org/x/term` for raw mode and manual
ANSI parsing for input.

## Layout

The first version should not include a full layout engine.

Instead, provide small helpers:

```go
func SplitVertical(r Rect, leftWidth int) (left Rect, right Rect)
func SplitHorizontal(r Rect, topHeight int) (top Rect, bottom Rect)
func Center(r Rect, w, h int) Rect
```

This keeps layout visible in user code.

Example usage:

```go
func (a *App) View(ctx *kero.Context, f *kero.Frame) {
    normal := kero.NewStyle()
    title := kero.NewStyle().Bold()

    root := kero.Rect{
        X: 0,
        Y: 0,
        W: ctx.Width,
        H: ctx.Height,
    }

    sidebar, content := kero.SplitVertical(root, 24)
    header, body := kero.SplitHorizontal(content, 3)

    f.Box(sidebar, normal)
    f.Write(sidebar.X+2, sidebar.Y+1, "Menu", title)
    f.Write(sidebar.X+2, sidebar.Y+3, "Inbox", normal)
    f.Write(sidebar.X+2, sidebar.Y+4, "Archive", normal)

    f.Box(header, normal)
    f.Write(header.X+2, header.Y+1, "Messages", title)

    f.Box(body, normal)
    f.Write(body.X+2, body.Y+1, "Select a message from the sidebar.", normal)
}
```

For an 80x24 terminal, the rough rectangles would be:

```text
root    = {X:0,  Y:0, W:80, H:24}
sidebar = {X:0,  Y:0, W:24, H:24}
content = {X:24, Y:0, W:56, H:24}
header  = {X:24, Y:0, W:56, H:3}
body    = {X:24, Y:3, W:56, H:21}
```

## Widgets

Widgets are optional helpers, not the center of the framework.

A widget should be a plain struct with draw and update methods:

```go
type TextInput struct {
    Value  string
    Cursor int
}

func (t *TextInput) Update(ev Event)
func (t TextInput) Draw(f *Frame, r Rect, s Style)
```

Widgets should not need their own runtime or hidden global state.

The cursor can be drawn with reverse video. If the cursor is over a character,
that character is drawn with `AttrReverse`. If the cursor is at the end of the
input, the widget draws a reversed space.

## Example

```go
package main

import (
    "fmt"

    "kero"
)

type Counter struct {
    n int
}

func (c *Counter) Init(ctx *kero.Context) error {
    return nil
}

func (c *Counter) Update(ctx *kero.Context, ev kero.Event) error {
    switch e := ev.(type) {
    case kero.KeyEvent:
        switch e.Key {
        case kero.KeyEsc:
            ctx.Quit()
        case kero.KeyRune:
            switch e.Rune {
            case 'q':
                ctx.Quit()
            case '+':
                c.n++
            case '-':
                c.n--
            }
        }
    }

    return nil
}

func (c *Counter) View(ctx *kero.Context, f *kero.Frame) {
    title := kero.NewStyle().Foreground(kero.ColorCyan).Bold()
    normal := kero.NewStyle()

    f.Clear()
    f.Box(kero.Rect{X: 2, Y: 1, W: ctx.Width - 4, H: 5}, normal)
    f.Write(4, 2, "kero counter", title)
    f.Write(4, 3, fmt.Sprintf("value: %d", c.n), normal)
    f.Write(4, 4, "press +, -, q, or esc", normal)
}

func main() {
    app := &Counter{}
    p := kero.New(app, kero.WithAltScreen(true))

    if err := p.Run(); err != nil {
        panic(err)
    }
}
```

## Error Handling

`Run` should return errors.

Terminal cleanup should still happen if:

- `Init` fails.
- `Update` fails.
- Rendering fails.
- Event reading fails.
- The app panics.

Panic recovery can be optional, but terminal restoration must be reliable.

## First Implementation Milestone

The first useful version should include:

- Raw mode setup and restoration.
- Alternate screen support.
- Keyboard input for common keys.
- Resize events.
- A frame buffer.
- Full-screen rendering.
- Basic styles.
- The `App`, `Context`, `Program`, and `Frame` APIs.
- One example app.

Mouse support, widgets, frame diffing, and advanced color can wait.

## Design Questions

Open questions to answer during implementation:

- Should `View` return an error?
- Should `Update` be allowed to emit commands?
- Should ticks be opt-in per app?
- Should kero support Unicode grapheme clusters from day one?
- Should `example/` contain one folder per example, or flat example files?

## Bias For Version 0

When unsure, choose the boring option:

- Redraw more instead of diffing early.
- Use fewer interfaces.
- Add fewer widgets.
- Keep terminal-specific code isolated.
- Make examples better before making abstractions bigger.
