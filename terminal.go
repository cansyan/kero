package kero

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// Terminal controls the user's terminal.
type Terminal interface {
	Enter() error
	Leave() error
	ReadEvent() (Event, error)
	Size() (Size, error)
}

type ansiTerminal struct {
	in   *os.File
	out  io.Writer
	opts Options

	reader  *bufio.Reader
	state   string
	resizes chan os.Signal
	events  chan eventResult
	done    chan struct{}
}

func newTerminal(in *os.File, out io.Writer, opts Options) *ansiTerminal {
	return &ansiTerminal{
		in:      in,
		out:     out,
		opts:    opts,
		reader:  bufio.NewReader(in),
		resizes: make(chan os.Signal, 1),
		events:  make(chan eventResult, 16),
		done:    make(chan struct{}),
	}
}

func (t *ansiTerminal) Enter() error {
	if isTerminal(t.in) {
		state, err := t.sttyOutput("-g")
		if err != nil {
			return err
		}
		t.state = strings.TrimSpace(string(state))
		if err := t.sttyRun("raw", "-echo"); err != nil {
			return err
		}
	}

	if t.opts.AltScreen {
		if _, err := io.WriteString(t.out, "\x1b[?1049h"); err != nil {
			return err
		}
	}
	if t.opts.Mouse {
		if _, err := io.WriteString(t.out, "\x1b[?1000h\x1b[?1006h"); err != nil {
			return err
		}
	}
	_, err := io.WriteString(t.out, "\x1b[?25l\x1b[H\x1b[2J")
	if err != nil {
		return err
	}

	// Enable Kitty keyboard protocol (CSI > 1 u) to receive modified keys like Shift+Enter
	if t.opts.Kitty {
		_, err := io.WriteString(t.out, "\x1b[>1u")
		if err != nil {
			return err
		}
	}

	signal.Notify(t.resizes, syscall.SIGWINCH)

	// Start a dedicated goroutine that continuously reads runes from stdin
	// and publishes parsed Events to inputEvents or inputErr. This allows
	// ReadEvent to select on resizes and input concurrently without leaking
	// blocked readers.
	go func() {
		for {
			select {
			case <-t.done:
				return
			default:
			}
			r, _, err := t.reader.ReadRune()
			if err != nil {
				res := eventResult{ev: nil, err: err}
				select {
				case t.events <- res:
				case <-t.done:
				}
				return
			}
			ev, perr := t.parseRune(r)
			res := eventResult{ev: ev, err: perr}
			select {
			case t.events <- res:
			case <-t.done:
				return
			}
		}
	}()

	return nil
}

func (t *ansiTerminal) Leave() error {
	// Stop reader goroutine first.
	select {
	case <-t.done:
		// already closed
	default:
		close(t.done)
	}
	signal.Stop(t.resizes)

	var firstErr error
	if t.opts.Mouse {
		_, firstErr = io.WriteString(t.out, "\x1b[?1006l\x1b[?1000l")
	}
	if _, err := io.WriteString(t.out, "\x1b[?25h\x1b[0m"); firstErr == nil && err != nil {
		firstErr = err
	}
	// Restore keyboard mode (CSI < u) before leaving, only if it was enabled
	if t.opts.Kitty {
		_, err := io.WriteString(t.out, "\x1b[<u")
		if firstErr == nil && err != nil {
			firstErr = err
		}
	}

	if t.opts.AltScreen {
		if _, err := io.WriteString(t.out, "\x1b[?1049l"); firstErr == nil && err != nil {
			firstErr = err
		}
	}
	if t.state != "" {
		if err := t.sttyRun(t.state); firstErr == nil && err != nil {
			firstErr = err
		}
	}
	return firstErr
}

func (t *ansiTerminal) ReadEvent() (Event, error) {
	select {
	case <-t.resizes:
		size, err := t.Size()
		if err != nil {
			return nil, err
		}
		return ResizeEvent{Width: size.Width, Height: size.Height}, nil
	case res := <-t.events:
		return res.ev, res.err
	}
}

func (t *ansiTerminal) Size() (Size, error) {
	out, err := t.sttyOutput("size")
	if err != nil {
		return Size{Width: 80, Height: 24}, nil
	}
	fields := strings.Fields(string(out))
	if len(fields) != 2 {
		return Size{Width: 80, Height: 24}, nil
	}
	rows, err := strconv.Atoi(fields[0])
	if err != nil {
		return Size{Width: 80, Height: 24}, nil
	}
	cols, err := strconv.Atoi(fields[1])
	if err != nil {
		return Size{Width: 80, Height: 24}, nil
	}
	return Size{Width: cols, Height: rows}, nil
}

func (t *ansiTerminal) sttyOutput(args ...string) ([]byte, error) {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = t.in
	return cmd.Output()
}

func (t *ansiTerminal) sttyRun(args ...string) error {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = t.in
	return cmd.Run()
}

func (t *ansiTerminal) parseRune(r rune) (Event, error) {
	switch r {
	case 0x09:
		return KeyEvent{Key: KeyTab}, nil
	case 0x0d, 0x0a:
		return KeyEvent{Key: KeyEnter}, nil
	case 0x7f, 0x08:
		return KeyEvent{Key: KeyBackspace}, nil
	case 0x1b:
		return t.parseEscape()
	case 0x01:
		return KeyEvent{Key: KeyRune, Rune: 'a', Mod: ModCtrl}, nil
	case 0x02:
		return KeyEvent{Key: KeyRune, Rune: 'b', Mod: ModCtrl}, nil
	case 0x03:
		return KeyEvent{Key: KeyRune, Rune: 'c', Mod: ModCtrl}, nil
	case 0x04:
		return KeyEvent{Key: KeyRune, Rune: 'd', Mod: ModCtrl}, nil
	case 0x05:
		return KeyEvent{Key: KeyRune, Rune: 'e', Mod: ModCtrl}, nil
	case 0x06:
		return KeyEvent{Key: KeyRune, Rune: 'f', Mod: ModCtrl}, nil
	case 0x0e:
		return KeyEvent{Key: KeyRune, Rune: 'n', Mod: ModCtrl}, nil
	case 0x10:
		return KeyEvent{Key: KeyRune, Rune: 'p', Mod: ModCtrl}, nil
	case 0x11:
		return KeyEvent{Key: KeyRune, Rune: 'q', Mod: ModCtrl}, nil
	case 0x13:
		return KeyEvent{Key: KeyRune, Rune: 's', Mod: ModCtrl}, nil
	// Ctrl+\ (0x1c) -> represent as Rune '\' with ModCtrl
	case 0x1c:
		return KeyEvent{Key: KeyRune, Rune: '\\', Mod: ModCtrl}, nil
	default:
		if r < 0x20 {
			return KeyEvent{Key: KeyUnknown, Rune: r}, nil
		}
		return KeyEvent{Key: KeyRune, Rune: r}, nil
	}
}

func (t *ansiTerminal) parseEscape() (Event, error) {
	if t.reader.Buffered() == 0 {
		return KeyEvent{Key: KeyEsc}, nil
	}

	next, _, err := t.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return KeyEvent{Key: KeyEsc}, nil
		}
		return nil, err
	}
	if next != '[' && next != 'O' {
		return KeyEvent{Key: KeyRune, Rune: next, Mod: ModAlt}, nil
	}

	var seq []rune
	for {
		r, _, err := t.reader.ReadRune()
		if err != nil {
			return nil, err
		}
		seq = append(seq, r)
		if (r >= '@' && r <= '~') || len(seq) > 32 {
			break
		}
	}

	return keyFromEscape(next, string(seq)), nil
}

func keyFromEscape(prefix rune, seq string) Event {
	// Handle SGR ("u") sequences from modern terminals, e.g. "13;2u" = Enter with Shift
	if before, ok := strings.CutSuffix(seq, "u"); ok {
		s := before
		parts := strings.Split(s, ";")
		if len(parts) >= 1 {
			code, err := strconv.Atoi(parts[0])
			if err == nil {
				var mod Mod = ModNone
				if len(parts) >= 2 {
					if m, err := strconv.Atoi(parts[1]); err == nil {
						switch m {
						case 2:
							mod |= ModShift
						case 3:
							mod |= ModAlt
						case 4:
							mod |= ModShift | ModAlt
						case 5:
							mod |= ModCtrl
						case 6:
							mod |= ModShift | ModCtrl
						case 7:
							mod |= ModAlt | ModCtrl
						case 8:
							mod |= ModShift | ModAlt | ModCtrl
						}
					}
				}
				// Map known control keys; otherwise treat as a rune with modifiers.
				switch code {
				case 13:
					return KeyEvent{Key: KeyEnter, Mod: mod}
				case 9:
					return KeyEvent{Key: KeyTab, Mod: mod}
				case 27:
					return KeyEvent{Key: KeyEsc, Mod: mod}
				default:
					return KeyEvent{Key: KeyRune, Rune: rune(code), Mod: mod}
				}
			}
		}
	}

	// Also handle CSI sequences that end with '~', e.g. "13;2~"
	if before, ok := strings.CutSuffix(seq, "~"); ok {
		s := before
		parts := strings.Split(s, ";")
		if len(parts) >= 1 {
			code, err := strconv.Atoi(parts[0])
			if err == nil {
				var mod Mod = ModNone
				if len(parts) >= 2 {
					if m, err := strconv.Atoi(parts[1]); err == nil {
						switch m {
						case 2:
							mod |= ModShift
						case 3:
							mod |= ModAlt
						case 4:
							mod |= ModShift | ModAlt
						case 5:
							mod |= ModCtrl
						case 6:
							mod |= ModShift | ModCtrl
						case 7:
							mod |= ModAlt | ModCtrl
						case 8:
							mod |= ModShift | ModAlt | ModCtrl
						}
					}
				}
				// Map known numeric CSI '~' codes to keys.
				switch code {
				case 13:
					return KeyEvent{Key: KeyEnter, Mod: mod}
				case 9:
					return KeyEvent{Key: KeyTab, Mod: mod}
				case 27:
					return KeyEvent{Key: KeyEsc, Mod: mod}
				case 3:
					return KeyEvent{Key: KeyDelete, Mod: mod}
				case 5:
					return KeyEvent{Key: KeyPgUp, Mod: mod}
				case 6:
					return KeyEvent{Key: KeyPgDown, Mod: mod}
				default:
					return KeyEvent{Key: KeyRune, Rune: rune(code), Mod: mod}
				}
			}
		}
	}

	if prefix == 'O' {
		switch seq {
		case "H":
			return KeyEvent{Key: KeyHome}
		case "F":
			return KeyEvent{Key: KeyEnd}
		}
	}

	switch seq {
	case "A":
		return KeyEvent{Key: KeyUp}
	case "B":
		return KeyEvent{Key: KeyDown}
	case "C":
		return KeyEvent{Key: KeyRight}
	case "D":
		return KeyEvent{Key: KeyLeft}
	case "H":
		return KeyEvent{Key: KeyHome}
	case "F":
		return KeyEvent{Key: KeyEnd}
	case "3~":
		return KeyEvent{Key: KeyDelete}
	case "5~":
		return KeyEvent{Key: KeyPgUp}
	case "6~":
		return KeyEvent{Key: KeyPgDown}
	}
	return KeyEvent{Key: KeyUnknown}
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
