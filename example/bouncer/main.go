package main

import (
	"fmt"

	"github.com/cansyan/kero"
)

type Bouncer struct {
	x  int
	y  int
	dx int
	dy int
}

func (b *Bouncer) Init(ctx *kero.Context) error {
	b.x = ctx.Width / 2
	b.y = ctx.Height / 2
	b.dx = 1
	b.dy = 1
	return nil
}

func (b *Bouncer) Update(ctx *kero.Context, ev kero.Event) error {
	switch e := ev.(type) {
	case kero.KeyEvent:
		switch e.Key {
		case kero.KeyEsc:
			ctx.Quit()
		case kero.KeyRune:
			if e.Rune == 'q' {
				ctx.Quit()
			}
		}
	case kero.TickEvent:
		b.step(ctx)
	}

	return nil
}

func (b *Bouncer) View(ctx *kero.Context, f *kero.Frame) {
	title := kero.NewStyle().Foreground(kero.ColorCyan).Bold()
	ball := kero.NewStyle().Foreground(kero.ColorYellow).Bold()
	normal := kero.NewStyle()

	f.Write(2, 1, "kero bouncer", title)
	f.Write(2, 2, fmt.Sprintf("position: %d,%d", b.x, b.y), normal)
	f.Write(2, 3, "press q or esc to quit", normal)
	f.Set(b.x, b.y, 'o', ball)
}

func (b *Bouncer) step(ctx *kero.Context) {
	minX, minY := 0, 5
	maxX, maxY := ctx.Width-1, ctx.Height-1

	if maxX < minX || maxY < minY {
		return
	}

	b.x += b.dx
	b.y += b.dy

	if b.x <= minX {
		b.x = minX
		b.dx = 1
	} else if b.x >= maxX {
		b.x = maxX
		b.dx = -1
	}

	if b.y <= minY {
		b.y = minY
		b.dy = 1
	} else if b.y >= maxY {
		b.y = maxY
		b.dy = -1
	}
}

func main() {
	app := &Bouncer{}
	p := kero.New(app, kero.WithAltScreen(true), kero.WithFPS(20))

	if err := p.Run(); err != nil {
		panic(err)
	}
}
