package kero

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"testing"
)

func TestCopyToClipboard(t *testing.T) {
	var buf bytes.Buffer
	term := newTerminal(os.Stdin, &buf, DefaultOptions())

	text := "Hello, Kero Clipboard!"
	err := term.CopyToClipboard(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedEncoded := base64.StdEncoding.EncodeToString([]byte(text))
	expectedOutput := fmt.Sprintf("\x1b]52;c;%s\x07", expectedEncoded)

	if got := buf.String(); got != expectedOutput {
		t.Fatalf("got output %q, want %q", got, expectedOutput)
	}
}

func TestContextCopyToClipboard(t *testing.T) {
	var buf bytes.Buffer
	term := newTerminal(os.Stdin, &buf, DefaultOptions())
	ctx := Context{terminal: term}

	text := "Context Copy Test"
	if err := ctx.CopyToClipboard(text); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedEncoded := base64.StdEncoding.EncodeToString([]byte(text))
	expectedOutput := fmt.Sprintf("\x1b]52;c;%s\x07", expectedEncoded)

	if got := buf.String(); got != expectedOutput {
		t.Fatalf("got output %q, want %q", got, expectedOutput)
	}
}

func TestContextCopyToClipboardNilTerminal(t *testing.T) {
	ctx := Context{}
	if err := ctx.CopyToClipboard("test"); err != nil {
		t.Fatalf("expected nil error when terminal is nil, got %v", err)
	}
}
