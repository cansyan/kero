package kero

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

// Color identifies a basic terminal color.
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

// Attr stores terminal text attributes.
type Attr uint16

const (
	AttrNone Attr = 0
	AttrBold Attr = 1 << iota
	AttrUnderline
	AttrReverse
	AttrDim
)
