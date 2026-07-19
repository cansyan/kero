package kero

// App is the user-provided terminal program.
type App interface {
	Init(ctx *Context) error
	Update(ctx *Context, ev Event) error
	View(ctx *Context, f *Frame)
}
