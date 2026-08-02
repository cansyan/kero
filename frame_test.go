package kero

import "testing"

func TestFrameWriteClips(t *testing.T) {
	f := NewFrame(4, 2)
	f.Write(-2, 0, "abcd", NewStyle())

	if got := f.Cell(0, 0).Ch; got != 'c' {
		t.Fatalf("cell(0,0) = %q, want c", got)
	}
	if got := f.Cell(1, 0).Ch; got != 'd' {
		t.Fatalf("cell(1,0) = %q, want d", got)
	}
}

func TestFrameBox(t *testing.T) {
	f := NewFrame(4, 3)
	f.Box(Rect{X: 0, Y: 0, W: 4, H: 3}, NewStyle())

	checks := map[Point]rune{
		{X: 0, Y: 0}: '+',
		{X: 3, Y: 0}: '+',
		{X: 1, Y: 0}: '-',
		{X: 0, Y: 1}: '|',
	}
	for p, want := range checks {
		if got := f.Cell(p.X, p.Y).Ch; got != want {
			t.Fatalf("cell(%d,%d) = %q, want %q", p.X, p.Y, got, want)
		}
	}
}
