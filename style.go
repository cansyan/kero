package kero

import (
	"strconv"
	"strings"
)

// Style stores foreground, background, and text attributes.
type Style struct {
	Fg   Color
	Bg   Color
	Attr Attr
}

// NewStyle returns the default terminal style.
func NewStyle() Style {
	return Style{Fg: ColorDefault, Bg: ColorDefault, Attr: AttrNone}
}

// Foreground returns s with a foreground color.
func (s Style) Foreground(c Color) Style {
	s.Fg = c
	return s
}

// Background returns s with a background color.
func (s Style) Background(c Color) Style {
	s.Bg = c
	return s
}

// Bold returns s with bold text.
func (s Style) Bold() Style {
	s.Attr |= AttrBold
	return s
}

// Underline returns s with underlined text.
func (s Style) Underline() Style {
	s.Attr |= AttrUnderline
	return s
}

// Reverse returns s with foreground and background reversed by the terminal.
func (s Style) Reverse() Style {
	s.Attr |= AttrReverse
	return s
}

// Dim returns s with dim text.
func (s Style) Dim() Style {
	s.Attr |= AttrDim
	return s
}

// Italic returns s with italic text
func (s Style) Italic() Style {
	s.Attr |= AttrItalic
	return s
}

type ColorType uint8

const (
	ColorTypeDefault ColorType = iota
	ColorTypeANSI              // Standard 8/16 ANSI colors
	ColorType256               // 256-color palette (0-255)
	ColorTypeRGB               // 24-bit True Color
)

type Color struct {
	Type ColorType
	R    uint8 // Used for RGB (or ANSI index if preferred)
	G    uint8
	B    uint8
	Idx  uint8 // Used for 256-color index or basic ANSI index
}

// Basic ANSI Color Constants (0-7 offset logic for standard 30-37 / 40-47)
var (
	ColorDefault = Color{Type: ColorTypeDefault}
	ColorBlack   = Color{Type: ColorTypeANSI, Idx: 0}
	ColorRed     = Color{Type: ColorTypeANSI, Idx: 1}
	ColorGreen   = Color{Type: ColorTypeANSI, Idx: 2}
	ColorYellow  = Color{Type: ColorTypeANSI, Idx: 3}
	ColorBlue    = Color{Type: ColorTypeANSI, Idx: 4}
	ColorMagenta = Color{Type: ColorTypeANSI, Idx: 5}
	ColorCyan    = Color{Type: ColorTypeANSI, Idx: 6}
	ColorWhite   = Color{Type: ColorTypeANSI, Idx: 7}
)

func ColorANSI(idx uint8) Color {
	return Color{Type: ColorTypeANSI, Idx: idx}
}

func Color256(idx uint8) Color {
	return Color{Type: ColorType256, Idx: idx}
}

func ColorRGB(r, g, b uint8) Color {
	return Color{Type: ColorTypeRGB, R: r, G: g, B: b}
}

// ColorHex parses a hex color string (e.g., "#3E4451", "3E4451", "#3E4", "3E4")
// and returns an RGB Color. Returns ColorDefault if parsing fails.
func ColorHex(hex string) Color {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		hex = string([]byte{
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
		})
	}

	if len(hex) != 6 {
		return ColorDefault
	}

	val, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return ColorDefault
	}

	return ColorRGB(
		uint8((val>>16)&0xFF),
		uint8((val>>8)&0xFF),
		uint8(val&0xFF),
	)
}

// Attr stores terminal text attributes.
type Attr uint16

const (
	AttrNone Attr = 0
	AttrBold Attr = 1 << iota
	AttrUnderline
	AttrReverse
	AttrDim
	AttrItalic
)
