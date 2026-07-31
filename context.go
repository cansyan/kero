package kero

// Context gives an app controlled access to runtime state.
type Context struct {
	Width  int
	Height int

	LastKey KeyEvent

	done bool
}

// Quit asks the program loop to stop after the current event is handled.
func (c *Context) Quit() {
	c.done = true
}

// Size returns the current terminal size.
func (c *Context) Size() Size {
	return Size{Width: c.Width, Height: c.Height}
}
