package kero

import "testing"

func TestRectInset(t *testing.T) {
	r := Rect{X: 2, Y: 1, W: 20, H: 5}
	got := r.Inset(1)
	want := Rect{X: 3, Y: 2, W: 18, H: 3}
	if got != want {
		t.Fatalf("Inset(1) = %#v, want %#v", got, want)
	}
}

func TestSplitVertical(t *testing.T) {
	root := Rect{X: 0, Y: 0, W: 80, H: 24}
	left, right := SplitVertical(root, 24)

	if left != (Rect{X: 0, Y: 0, W: 24, H: 24}) {
		t.Fatalf("left = %#v", left)
	}
	if right != (Rect{X: 24, Y: 0, W: 56, H: 24}) {
		t.Fatalf("right = %#v", right)
	}
}

func TestCenter(t *testing.T) {
	got := Center(Rect{X: 0, Y: 0, W: 80, H: 24}, 20, 6)
	want := Rect{X: 30, Y: 9, W: 20, H: 6}
	if got != want {
		t.Fatalf("Center = %#v, want %#v", got, want)
	}
}
