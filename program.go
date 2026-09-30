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

	frameTimer    *time.Timer
	frameC        <-chan time.Time
	frameRequests chan struct{}
}

// New creates a Program.
func New(app App, opts ...Option) *Program {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	terminal := newTerminal(os.Stdin, os.Stdout, options)
	return &Program{
		app:           app,
		opts:          options,
		terminal:      terminal,
		screen:        NewScreen(os.Stdout, 0, 0),
		frameRequests: make(chan struct{}, 1),
		ctx: Context{
			terminal: terminal,
		},
	}
}

// CopyToClipboard copies text to the system clipboard via the terminal.
func (p *Program) CopyToClipboard(text string) error {
	return p.terminal.CopyToClipboard(text)
}

// Run starts the terminal program and blocks until the app quits or an error occurs.
func (p *Program) Run() error {
	p.ctx.requestFrame = p.requestFrame

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

	var lastFrame time.Time

	for !p.ctx.done {
		var ev Event

		select {
		case result := <-events:
			if result.err != nil {
				return result.err
			}
			ev = result.ev

		case <-p.frameRequests:
			p.scheduleFrame()
			continue

		case tm := <-p.frameC:
			delta := time.Duration(0)
			if !lastFrame.IsZero() {
				delta = tm.Sub(lastFrame)
			}
			lastFrame = tm

			// The timer has fired, so another frame isn't
			// currently scheduled.
			p.frameC = nil

			ev = FrameEvent{
				Time:  tm,
				Delta: delta,
			}
		}

		if resize, ok := ev.(ResizeEvent); ok {
			p.ctx.Width = resize.Width
			p.ctx.Height = resize.Height
			p.screen.Resize(resize.Width, resize.Height)
		}

		if err := p.app.Update(&p.ctx, ev); err != nil {
			return err
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
	p.app.Draw(&p.ctx, f)
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

func (p *Program) requestFrame() {
	select {
	case p.frameRequests <- struct{}{}:
	default:
	}
}

func (p *Program) scheduleFrame() {
	fps := p.opts.FPS
	if fps <= 0 {
		return
	}

	interval := time.Second / time.Duration(fps)

	// A frame is already scheduled.
	// Don't create another one.
	if p.frameC != nil {
		return
	}

	if p.frameTimer == nil {
		p.frameTimer = time.NewTimer(interval)
	} else {
		p.frameTimer.Reset(interval)
	}

	p.frameC = p.frameTimer.C
}

type eventResult struct {
	ev  Event
	err error
}
