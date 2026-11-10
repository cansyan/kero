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

func TestFrameWideCharacters(t *testing.T) {
	f := NewFrame(6, 2)
	f.Write(0, 0, "世界hi", NewStyle())

	if got := f.Cell(0, 0).Ch; got != '世' {
		t.Fatalf("cell(0,0) = %q, want '世'", got)
	}
	if got := f.Cell(1, 0).Ch; got != Continuation {
		t.Fatalf("cell(1,0) = %v, want Continuation", got)
	}
	if got := f.Cell(2, 0).Ch; got != '界' {
		t.Fatalf("cell(2,0) = %q, want '界'", got)
	}
	if got := f.Cell(3, 0).Ch; got != Continuation {
		t.Fatalf("cell(3,0) = %v, want Continuation", got)
	}
	if got := f.Cell(4, 0).Ch; got != 'h' {
		t.Fatalf("cell(4,0) = %q, want 'h'", got)
	}
	if got := f.Cell(5, 0).Ch; got != 'i' {
		t.Fatalf("cell(5,0) = %q, want 'i'", got)
	}
}

func TestFrameOverwriteContinuation(t *testing.T) {
	f := NewFrame(4, 1)
	f.Write(0, 0, "世", NewStyle())
	if got := f.Cell(0, 0).Ch; got != '世' {
		t.Fatalf("cell(0,0) = %q, want '世'", got)
	}
	if got := f.Cell(1, 0).Ch; got != Continuation {
		t.Fatalf("cell(1,0) = %v, want Continuation", got)
	}

	// Overwrite continuation cell at x=1
	f.Set(1, 0, 'a', NewStyle())
	if got := f.Cell(0, 0).Ch; got != ' ' {
		t.Fatalf("cell(0,0) after overwriting continuation = %q, want ' '", got)
	}
	if got := f.Cell(1, 0).Ch; got != 'a' {
		t.Fatalf("cell(1,0) after overwriting continuation = %q, want 'a'", got)
	}
}

func TestFrameOverwriteLeadWideCharacter(t *testing.T) {
	f := NewFrame(4, 1)
	f.Write(0, 0, "世", NewStyle())

	// Overwrite lead cell at x=0 with narrow character
	f.Set(0, 0, 'x', NewStyle())
	if got := f.Cell(0, 0).Ch; got != 'x' {
		t.Fatalf("cell(0,0) = %q, want 'x'", got)
	}
	if got := f.Cell(1, 0).Ch; got != ' ' {
		t.Fatalf("cell(1,0) = %q, want ' '", got)
	}
}

func TestFrameWideCharacterBoundaryClip(t *testing.T) {
	f := NewFrame(3, 1)
	f.Write(2, 0, "世", NewStyle())

	if got := f.Cell(2, 0).Ch; got != ' ' {
		t.Fatalf("cell(2,0) = %q, want ' ' when wide char does not fit", got)
	}
}

func TestFrameWideCharacterOverwritesWideContinuation(t *testing.T) {
	f := NewFrame(4, 1)
	f.Write(1, 0, "世", NewStyle()) // cells: [0: ' ', 1: '世', 2: Cont, 3: ' ']

	// Place wide char at 0, which takes cells 0 and 1 (overwriting '世' at 1, which has continuation at 2)
	f.Write(0, 0, "界", NewStyle())
	if got := f.Cell(0, 0).Ch; got != '界' {
		t.Fatalf("cell(0,0) = %q, want '界'", got)
	}
	if got := f.Cell(1, 0).Ch; got != Continuation {
		t.Fatalf("cell(1,0) = %v, want Continuation", got)
	}
	if got := f.Cell(2, 0).Ch; got != ' ' {
		t.Fatalf("cell(2,0) = %q, want ' '", got)
	}
}
