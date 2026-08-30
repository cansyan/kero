package kero

// App is the user-provided terminal program.
type App interface {
	// Init runs once after the terminal is ready.
	Init(ctx *Context) error
	// Update receives input, resize events, ticks, and internal messages.
	Update(ctx *Context, ev Event) error
	// Draw renders the current application state onto the provided Frame.
	// It should avoid mutating application state.
	Draw(ctx *Context, f *Frame)
}
