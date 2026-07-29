package kero

// TextInput is a small editable single-line text widget.
type TextInput struct {
	Value  string
	Cursor int
}

// Update applies keyboard input to the text input.
func (t *TextInput) Update(ev Event) {
	e, ok := ev.(KeyEvent)
	if !ok {
		return
	}

	runes := []rune(t.Value)
	if t.Cursor < 0 {
		t.Cursor = 0
	}
	if t.Cursor > len(runes) {
		t.Cursor = len(runes)
	}

	switch e.Key {
	case KeyRune:
		if e.Mod&ModCtrl != 0 {
			return
		}
		runes = append(runes, 0)
		copy(runes[t.Cursor+1:], runes[t.Cursor:])
		runes[t.Cursor] = e.Rune
		t.Cursor++
	case KeyBackspace:
		if t.Cursor > 0 {
			runes = append(runes[:t.Cursor-1], runes[t.Cursor:]...)
			t.Cursor--
		}
	case KeyDelete:
		if t.Cursor < len(runes) {
			runes = append(runes[:t.Cursor], runes[t.Cursor+1:]...)
		}
	case KeyLeft:
		if t.Cursor > 0 {
			t.Cursor--
		}
	case KeyRight:
		if t.Cursor < len(runes) {
			t.Cursor++
		}
	case KeyHome:
		t.Cursor = 0
	case KeyEnd:
		t.Cursor = len(runes)
	}

	t.Value = string(runes)
}

// Draw renders the text input and its cursor.
func (t TextInput) Draw(f *Frame, r Rect, s Style) {
	cursorStyle := s.Reverse()
	runes := []rune(t.Value)

	if t.Cursor < 0 {
		t.Cursor = 0
	}
	if t.Cursor > len(runes) {
		t.Cursor = len(runes)
	}

	for i, ch := range runes {
		x := r.X + i
		y := r.Y

		if x >= r.Right() {
			break
		}

		if i == t.Cursor {
			f.Set(x, y, ch, cursorStyle)
			continue
		}

		f.Set(x, y, ch, s)
	}

	if t.Cursor == len(runes) {
		x := r.X + t.Cursor
		if x < r.Right() {
			f.Set(x, r.Y, ' ', cursorStyle)
		}
	}
}
