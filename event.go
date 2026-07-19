package kero

import "time"

// Event is a value emitted by the terminal or runtime.
type Event interface {
	event()
}

// KeyEvent represents a keyboard input event.
type KeyEvent struct {
	Key  Key
	Rune rune
	Mod  Mod
}

func (KeyEvent) event() {}

// ResizeEvent represents a terminal size change.
type ResizeEvent struct {
	Width  int
	Height int
}

func (ResizeEvent) event() {}

// TickEvent represents a timed redraw tick.
type TickEvent struct {
	Time time.Time
}

func (TickEvent) event() {}

// MouseEvent represents a terminal mouse event.
type MouseEvent struct {
	X      int
	Y      int
	Button MouseButton
	Action MouseAction
	Mod    Mod
}

func (MouseEvent) event() {}

// Key identifies non-text keyboard keys.
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

// Mod stores keyboard or mouse modifiers.
type Mod uint8

const (
	ModNone Mod = 0
	ModCtrl Mod = 1 << iota
	ModAlt
	ModShift
)

// MouseButton identifies a mouse button.
type MouseButton int

const (
	MouseNone MouseButton = iota
	MouseLeft
	MouseMiddle
	MouseRight
	MouseWheelUp
	MouseWheelDown
)

// MouseAction identifies a mouse action.
type MouseAction int

const (
	MousePress MouseAction = iota
	MouseRelease
	MouseDrag
)
