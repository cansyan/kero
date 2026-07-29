package kero

// Options configures a Program.
type Options struct {
	AltScreen bool
	Mouse     bool
	// FPS limits timed redraws, such as spinners, clocks, and animations.
	// Event-only apps can ignore it and render only after input or resize.
	FPS int
}

// Option changes Program options.
type Option func(*Options)

// DefaultOptions returns the default Program options.
func DefaultOptions() Options {
	return Options{
		AltScreen: true,
		Mouse:     false,
	}
}

// WithAltScreen enables or disables the alternate terminal screen.
func WithAltScreen(v bool) Option {
	return func(o *Options) {
		o.AltScreen = v
	}
}

// WithMouse enables or disables terminal mouse reporting.
func WithMouse(v bool) Option {
	return func(o *Options) {
		o.Mouse = v
	}
}

// WithFPS sets the maximum timed redraw rate.
func WithFPS(fps int) Option {
	return func(o *Options) {
		if fps > 0 {
			o.FPS = fps
		}
	}
}
