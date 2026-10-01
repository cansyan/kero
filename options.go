package kero

// Options configures a Program.
type Options struct {
	AltScreen      bool
	Mouse          bool
	Kitty          bool // kitty keyboard protocol
	BracketedPaste bool
}

// Option changes Program options.
type Option func(*Options)

// DefaultOptions returns the default Program options.
func DefaultOptions() Options {
	return Options{
		AltScreen:      true,
		Mouse:          false,
		BracketedPaste: true,
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

// WithKitty enables or disables kitty keyboard protocol.
func WithKitty(v bool) Option {
	return func(o *Options) {
		o.Kitty = v
	}
}

// WithBracketedPaste enables or disables bracketed paste mode (\x1b[?2004h).
func WithBracketedPaste(v bool) Option {
	return func(o *Options) {
		o.BracketedPaste = v
	}
}
