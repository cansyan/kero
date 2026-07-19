package main

import (
	"fmt"

	"kero"
)

type Counter struct {
	n int
}

func (c *Counter) Init(ctx *kero.Context) error {
	return nil
}

func (c *Counter) Update(ctx *kero.Context, ev kero.Event) error {
	switch e := ev.(type) {
	case kero.KeyEvent:
		switch e.Key {
		case kero.KeyEsc:
			ctx.Quit()
		case kero.KeyRune:
			switch e.Rune {
			case 'q':
				ctx.Quit()
			case '+':
				c.n++
			case '-':
				c.n--
			}
		}
	}

	return nil
}

func (c *Counter) View(ctx *kero.Context, f *kero.Frame) {
	title := kero.NewStyle().Foreground(kero.ColorCyan).Bold()
	normal := kero.NewStyle()

	f.Clear()
	f.Box(kero.Rect{X: 2, Y: 1, W: ctx.Width - 4, H: 5}, normal)
	f.Write(4, 2, "kero counter", title)
	f.Write(4, 3, fmt.Sprintf("value: %d", c.n), normal)
	f.Write(4, 4, "press +, -, q, or esc", normal)
}

func main() {
	app := &Counter{}
	p := kero.New(app, kero.WithAltScreen(true))

	if err := p.Run(); err != nil {
		panic(err)
	}
}
