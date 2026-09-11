package kero

import (
	"fmt"
	"io"
	"strings"
)

// Screen renders frames to a terminal output stream.
type Screen struct {
	out io.Writer

	width  int
	height int

	prev Frame
	next Frame
}

// NewScreen creates a screen with the given size.
func NewScreen(out io.Writer, w, h int) *Screen {
	return &Screen{
		out:    out,
		width:  w,
		height: h,
		prev:   NewFrame(w, h),
		next:   NewFrame(w, h),
	}
}

// Resize changes the screen size and clears stored frames.
func (s *Screen) Resize(w, h int) {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	s.width = w
	s.height = h
	s.prev = NewFrame(w, h)
	s.next = NewFrame(w, h)
}

// Frame returns the next drawable frame.
func (s *Screen) Frame() *Frame {
	return &s.next
}

// Flush writes the full frame to the output stream.
func (s *Screen) Flush() error {
	var b strings.Builder
	b.WriteString("\x1b[H")

	current := NewStyle()
	b.WriteString(styleANSI(current))

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			cell := s.next.Cell(x, y)
			if cell.Style != current {
				current = cell.Style
				b.WriteString(styleANSI(current))
			}
			ch := cell.Ch
			if ch == 0 {
				ch = ' '
			}
			b.WriteRune(ch)
		}
		if y != s.height-1 {
			b.WriteString("\r\n")
		}
	}

	b.WriteString(styleANSI(NewStyle()))
	_, err := io.WriteString(s.out, b.String())
	if err != nil {
		return err
	}

	s.prev, s.next = s.next, s.prev
	s.next.Clear()
	return nil
}

func styleANSI(s Style) string {
	codes := []string{"0"}

	if s.Attr&AttrBold != 0 {
		codes = append(codes, "1")
	}
	if s.Attr&AttrDim != 0 {
		codes = append(codes, "2")
	}
	if s.Attr&AttrItalic != 0 {
		codes = append(codes, "3")
	}
	if s.Attr&AttrUnderline != 0 {
		codes = append(codes, "4")
	}
	if s.Attr&AttrReverse != 0 {
		codes = append(codes, "7")
	}

	if code := colorANSICode(s.Fg, false); code != "" {
		codes = append(codes, code)
	}
	if code := colorANSICode(s.Bg, true); code != "" {
		codes = append(codes, code)
	}

	return "\x1b[" + strings.Join(codes, ";") + "m"
}

func colorANSICode(c Color, isBg bool) string {
	switch c.Type {
	case ColorTypeANSI:
		offset := 30
		if isBg {
			offset = 40
		}
		return fmt.Sprintf("%d", offset+int(c.Idx))
	case ColorType256:
		layer := "38"
		if isBg {
			layer = "48"
		}
		return fmt.Sprintf("%s;5;%d", layer, c.Idx)
	case ColorTypeRGB:
		layer := "38"
		if isBg {
			layer = "48"
		}
		return fmt.Sprintf("%s;2;%d;%d;%d", layer, c.R, c.G, c.B)
	default:
		return ""
	}
}
