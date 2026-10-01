package kero

// Context gives an app controlled access to runtime state.
type Context struct {
	Width  int
	Height int

	done         bool
	terminal     Terminal
	requestFrame func(fps int)
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

// RequestFrame requests one-shot frame.
// For example, RequestFrame(30) means "send me one FrameEvent about 1/30 second from now".
func (c *Context) RequestFrame(fps int) {
	if fps <= 0 {
		return
	}
	c.requestFrame(fps)
}
