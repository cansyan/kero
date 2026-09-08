package kero

// Context gives an app controlled access to runtime state.
type Context struct {
	Width  int
	Height int

	done     bool
	terminal Terminal
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
