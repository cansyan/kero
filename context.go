package kero

// Context gives an app controlled access to runtime state.
type Context struct {
	Width  int
	Height int

	done         bool
	terminal     Terminal
	requestFrame func(fps int)
	post         func(func(*Context))
}

// Quit asks the program loop to stop after the current event is handled.
func (c *Context) Quit() {
	c.done = true
}

// Size returns the current terminal size.
func (c *Context) Size() Size {
	return Size{Width: c.Width, Height: c.Height}
}

// CopyToClipboard copies text to the system clipboard via the terminal.
func (c *Context) CopyToClipboard(text string) error {
	if c.terminal == nil {
		return nil
	}
	return c.terminal.CopyToClipboard(text)
}

// RequestFrame schedules a FrameEvent for a future animation frame.
//
// It must be called from the program's event loop, such as from Init,
// Update, or a callback passed to Post. To request a frame from
// another goroutine, call Post and request the frame from its callback.
//
// The application should call RequestFrame again from the FrameEvent
// handler if more frames are needed.
//
// For example, RequestFrame(30) means "send me one FrameEvent about 1/30 second from now".
func (c *Context) RequestFrame(fps int) {
	if fps <= 0 {
		return
	}
	c.requestFrame(fps)
}

// Post schedules fn to run on the program's event loop.
//
// Post may be called safely from any goroutine. The function runs later
// in the same context as Update and Draw, so it may safely modify
// application state.
func (c *Context) Post(fn func(*Context)) {
	if fn == nil {
		return
	}

	c.post(fn)
}
