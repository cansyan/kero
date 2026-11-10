package kero

import "github.com/mattn/go-runewidth"

// Continuation is a sentinel rune used in wide-character continuation cells.
// Screen.Flush skips emitting output for cells containing Continuation.
const Continuation rune = -1

// Cell is one terminal cell.
type Cell struct {
	Ch    rune
	Style Style
}

// Frame is a drawable grid of terminal cells.
type Frame struct {
	width  int
	height int
	cells  []Cell
}

func NewFrame(w, h int) Frame {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	f := Frame{width: w, height: h, cells: make([]Cell, w*h)}
	f.Clear()
	return f
}

// Size returns the frame size.
func (f *Frame) Size() Size {
	return Size{Width: f.width, Height: f.height}
}

// Clear fills the frame with spaces using the default style.
func (f *Frame) Clear() {
	blank := Cell{Ch: ' ', Style: NewStyle()}
	for i := range f.cells {
		f.cells[i] = blank
	}
}

// Set writes one cell. Writes outside the frame are ignored.
// If ch is wide (e.g. CJK, width 2), the following cell(s) on the same row
// are marked as Continuation cells with the same style.
func (f *Frame) Set(x, y int, ch rune, s Style) {
	if x < 0 || y < 0 || x >= f.width || y >= f.height {
		return
	}
	if ch == 0 || ch == '\n' || ch == '\r' || ch == '\t' {
		ch = ' '
	}

	w := runewidth.RuneWidth(ch)
	if w <= 0 {
		w = 1
	}

	// If overwriting a continuation cell, clear the preceding wide character.
	if f.cells[y*f.width+x].Ch == Continuation {
		for prevX := x - 1; prevX >= 0; prevX-- {
			idx := y*f.width + prevX
			if f.cells[idx].Ch != Continuation {
				f.cells[idx] = Cell{Ch: ' ', Style: f.cells[idx].Style}
				break
			}
			f.cells[idx] = Cell{Ch: ' ', Style: f.cells[idx].Style}
		}
	}

	// If overwriting a lead wide character, clear its old continuation cells.
	oldRune := f.cells[y*f.width+x].Ch
	if oldRune != Continuation && oldRune != ' ' {
		oldW := runewidth.RuneWidth(oldRune)
		for i := 1; i < oldW && x+i < f.width; i++ {
			idx := y*f.width + (x + i)
			if f.cells[idx].Ch == Continuation {
				f.cells[idx] = Cell{Ch: ' ', Style: f.cells[idx].Style}
			}
		}
	}

	// If a wide character doesn't fit on the row, replace with a space.
	if ch != Continuation && x+w > f.width {
		ch = ' '
		w = 1
	}

	// If placing a wide character, clear any characters covered by its continuation cells.
	if ch != Continuation && w > 1 {
		for i := 1; i < w && x+i < f.width; i++ {
			idx := y*f.width + (x + i)
			cellRune := f.cells[idx].Ch
			if cellRune != Continuation {
				cellW := runewidth.RuneWidth(cellRune)
				for j := 1; j < cellW && x+i+j < f.width; j++ {
					cIdx := y*f.width + (x + i + j)
					if f.cells[cIdx].Ch == Continuation {
						f.cells[cIdx] = Cell{Ch: ' ', Style: f.cells[cIdx].Style}
					}
				}
			}
		}
	}

	f.cells[y*f.width+x] = Cell{Ch: ch, Style: s}

	if ch != Continuation {
		for i := 1; i < w; i++ {
			if x+i < f.width {
				f.cells[y*f.width+(x+i)] = Cell{Ch: Continuation, Style: s}
			}
		}
	}
}

// Write writes text starting at x/y. Text is clipped to the frame.
// Wide characters advance the column position by their display width.
func (f *Frame) Write(x, y int, text string, s Style) {
	if y < 0 || y >= f.height {
		return
	}
	for _, ch := range text {
		if x >= f.width {
			return
		}
		w := runewidth.RuneWidth(ch)
		if w <= 0 {
			w = 1
		}
		if x >= 0 {
			f.Set(x, y, ch, s)
		}
		x += w
	}
}

// Fill fills r with ch and s. Drawing is clipped to the frame.
func (f *Frame) Fill(r Rect, ch rune, s Style) {
	if ch == 0 {
		ch = ' '
	}
	for y := r.Y; y < r.Bottom(); y++ {
		for x := r.X; x < r.Right(); x++ {
			f.Set(x, y, ch, s)
		}
	}
}

// Box draws a single-line box around r. Drawing is clipped to the frame.
func (f *Frame) Box(r Rect, s Style) {
	if r.W <= 0 || r.H <= 0 {
		return
	}
	if r.W == 1 && r.H == 1 {
		f.Set(r.X, r.Y, '+', s)
		return
	}
	if r.H == 1 {
		for x := r.X; x < r.Right(); x++ {
			f.Set(x, r.Y, '-', s)
		}
		return
	}
	if r.W == 1 {
		for y := r.Y; y < r.Bottom(); y++ {
			f.Set(r.X, y, '|', s)
		}
		return
	}

	left := r.X
	right := r.Right() - 1
	top := r.Y
	bottom := r.Bottom() - 1

	f.Set(left, top, '+', s)
	f.Set(right, top, '+', s)
	f.Set(left, bottom, '+', s)
	f.Set(right, bottom, '+', s)

	for x := left + 1; x < right; x++ {
		f.Set(x, top, '-', s)
		f.Set(x, bottom, '-', s)
	}
	for y := top + 1; y < bottom; y++ {
		f.Set(left, y, '|', s)
		f.Set(right, y, '|', s)
	}
}

func (f *Frame) Cell(x, y int) Cell {
	if x < 0 || y < 0 || x >= f.width || y >= f.height {
		return Cell{Ch: ' ', Style: NewStyle()}
	}
	return f.cells[y*f.width+x]
}
