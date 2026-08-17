package kero

import (
	"fmt"
	"strings"
	"time"
)

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

// String returns a human-readable representation of the key event,
// including modifiers. Examples: "ctrl+a", "shift+enter", "esc".
func (k KeyEvent) String() string {
	var mods []string
	if k.Mod&ModCtrl != 0 {
		mods = append(mods, "ctrl")
	}
	if k.Mod&ModAlt != 0 {
		mods = append(mods, "alt")
	}
	if k.Mod&ModMeta != 0 {
		mods = append(mods, "cmd")
	}
	if k.Mod&ModShift != 0 {
		mods = append(mods, "shift")
	}

	var keyName string
	switch k.Key {
	case KeyRune:
		// printable ASCII-ish
		if k.Rune >= 0x20 && k.Rune != 0x7f {
			r := k.Rune
			keyName = fmt.Sprintf("%c", r)
		} else {
			keyName = fmt.Sprintf("U+%04X", k.Rune)
		}
	case KeyEsc:
		keyName = "esc"
	case KeyEnter:
		keyName = "enter"
	case KeyBackspace:
		keyName = "backspace"
	case KeyTab:
		keyName = "tab"
	case KeyUp:
		keyName = "up"
	case KeyDown:
		keyName = "down"
	case KeyLeft:
		keyName = "left"
	case KeyRight:
		keyName = "right"
	case KeyHome:
		keyName = "home"
	case KeyEnd:
		keyName = "end"
	case KeyPgUp:
		keyName = "pgup"
	case KeyPgDown:
		keyName = "pgdown"
	case KeyDelete:
		keyName = "delete"
	default:
		keyName = "unknown"
	}

	if len(mods) == 0 {
		return keyName
	}
	return strings.Join(mods, "+") + "+" + keyName
}

// ResizeEvent represents a terminal size change.
type ResizeEvent struct {
	Width  int
	Height int
}

func (ResizeEvent) event() {}

// TickEvent represents a timed redraw tick.
// It only presents when program option FPS is set
type TickEvent struct {
	Time time.Time
}

func (TickEvent) event() {}

// MouseEvent represents a terminal mouse event.
type MouseEvent struct {
	X      int // starts from 0
	Y      int // starts from 0
	Button MouseButton
	Action MouseAction
	Mod    Mod
}

func (MouseEvent) event() {}

// String returns a human-readable representation of the mouse event.
func (m MouseEvent) String() string {
	var mods []string
	if m.Mod&ModCtrl != 0 {
		mods = append(mods, "ctrl")
	}
	if m.Mod&ModAlt != 0 {
		mods = append(mods, "alt")
	}
	if m.Mod&ModShift != 0 {
		mods = append(mods, "shift")
	}

	var btnName string
	switch m.Button {
	case MouseLeft:
		btnName = "left"
	case MouseMiddle:
		btnName = "middle"
	case MouseRight:
		btnName = "right"
	case MouseWheelUp:
		btnName = "wheel_up"
	case MouseWheelDown:
		btnName = "wheel_down"
	default:
		btnName = "none"
	}

	var actName string
	switch m.Action {
	case MousePress:
		actName = "press"
	case MouseRelease:
		actName = "release"
	case MouseDrag:
		actName = "drag"
	}

	desc := fmt.Sprintf("mouse_%s_%s(%d,%d)", btnName, actName, m.X, m.Y)
	if len(mods) == 0 {
		return desc
	}
	return strings.Join(mods, "+") + "+" + desc
}

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
	ModMeta // With Kitty keyboard protocol enabled, Kitty terminal reports Command key on macOS
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

type PasteStartEvent struct{}

func (PasteStartEvent) event() {}

type PasteEndEvent struct{}

func (PasteEndEvent) event() {}
