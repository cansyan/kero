package kero

import (
	"os"
	"time"
)

// Program owns the terminal lifecycle and event loop.
type Program struct {
	app  App
	opts Options

	terminal Terminal
	screen   *Screen
	ctx      Context
}

// New creates a Program.
func New(app App, opts ...Option) *Program {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	terminal := newTerminal(os.Stdin, os.Stdout, options)
	return &Program{
		app:      app,
		opts:     options,
		terminal: terminal,
		screen:   NewScreen(os.Stdout, 0, 0),
	}
}

// Run starts the terminal program and blocks until the app quits or an error occurs.
func (p *Program) Run() error {
	if err := p.terminal.Enter(); err != nil {
		return err
	}
	defer p.terminal.Leave()

	size, err := p.terminal.Size()
	if err != nil {
		return err
	}

	p.ctx.Width = size.Width
	p.ctx.Height = size.Height
	p.screen.Resize(size.Width, size.Height)

	if err := p.app.Init(&p.ctx); err != nil {
		return err
	}
	if err := p.render(); err != nil {
		return err
	}

	events := make(chan eventResult)
	go p.readEvents(events)

	var ticks <-chan time.Time
	if p.opts.FPS > 0 {
		ticker := time.NewTicker(time.Second / time.Duration(p.opts.FPS))
		defer ticker.Stop()
		ticks = ticker.C
	}

	for !p.ctx.done {
		var ev Event

		select {
		case result := <-events:
			if result.err != nil {
				return result.err
			}
			ev = result.ev
		case tm := <-ticks:
			ev = TickEvent{Time: tm}
		}

		if resize, ok := ev.(ResizeEvent); ok {
			p.ctx.Width = resize.Width
			p.ctx.Height = resize.Height
			p.screen.Resize(resize.Width, resize.Height)
		}

		if err := p.app.Update(&p.ctx, ev); err != nil {
			return err
		}

		if ke, ok := ev.(KeyEvent); ok {
			p.ctx.LastKeyEvent = ke
		}

		if err := p.render(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Program) render() error {
	f := p.screen.Frame()
	f.Clear()
	p.app.View(&p.ctx, f)
	return p.screen.Flush()
}

func (p *Program) readEvents(events chan<- eventResult) {
	for {
		ev, err := p.terminal.ReadEvent()
		events <- eventResult{ev: ev, err: err}
		if err != nil {
			return
		}
	}
}

type eventResult struct {
	ev  Event
	err error
}
