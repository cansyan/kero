package kero

import "testing"

func TestTextInputUpdate(t *testing.T) {
	input := TextInput{}
	input.Update(KeyEvent{Key: KeyRune, Rune: 'a'})
	input.Update(KeyEvent{Key: KeyRune, Rune: 'b'})
	input.Update(KeyEvent{Key: KeyLeft})
	input.Update(KeyEvent{Key: KeyRune, Rune: 'X'})

	if input.Value != "aXb" {
		t.Fatalf("Value = %q, want aXb", input.Value)
	}
	if input.Cursor != 2 {
		t.Fatalf("Cursor = %d, want 2", input.Cursor)
	}
}

func TestTextInputDrawCursorAtEnd(t *testing.T) {
	input := TextInput{Value: "hi", Cursor: 2}
	f := newFrame(4, 1)

	input.Draw(&f, Rect{X: 0, Y: 0, W: 4, H: 1}, NewStyle())

	cell := f.cell(2, 0)
	if cell.Ch != ' ' {
		t.Fatalf("cursor cell = %q, want space", cell.Ch)
	}
	if cell.Style.Attr&AttrReverse == 0 {
		t.Fatalf("cursor cell is not reversed")
	}
}
