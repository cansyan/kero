package kero

// Size stores a width and height in terminal cells.
type Size struct {
	Width  int
	Height int
}

// Point stores an x/y coordinate in terminal cells.
type Point struct {
	X int
	Y int
}

// Rect stores a rectangle in terminal cells.
type Rect struct {
	X int
	Y int
	W int
	H int
}

// Right returns the first x coordinate after r.
func (r Rect) Right() int {
	return r.X + r.W
}

// Bottom returns the first y coordinate after r.
func (r Rect) Bottom() int {
	return r.Y + r.H
}

// Inset returns r with every edge moved inward by n cells.
func (r Rect) Inset(n int) Rect {
	w := r.W - n*2
	h := r.H - n*2
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	return Rect{X: r.X + n, Y: r.Y + n, W: w, H: h}
}

// Contains reports whether p is inside r.
func (r Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.Right() && p.Y >= r.Y && p.Y < r.Bottom()
}

// SplitVertical splits r into left and right rectangles.
func SplitVertical(r Rect, leftWidth int) (left Rect, right Rect) {
	if leftWidth < 0 {
		leftWidth = 0
	}
	if leftWidth > r.W {
		leftWidth = r.W
	}
	left = Rect{X: r.X, Y: r.Y, W: leftWidth, H: r.H}
	right = Rect{X: r.X + leftWidth, Y: r.Y, W: r.W - leftWidth, H: r.H}
	return left, right
}

// SplitHorizontal splits r into top and bottom rectangles.
func SplitHorizontal(r Rect, topHeight int) (top Rect, bottom Rect) {
	if topHeight < 0 {
		topHeight = 0
	}
	if topHeight > r.H {
		topHeight = r.H
	}
	top = Rect{X: r.X, Y: r.Y, W: r.W, H: topHeight}
	bottom = Rect{X: r.X, Y: r.Y + topHeight, W: r.W, H: r.H - topHeight}
	return top, bottom
}

// Center returns a w by h rectangle centered inside r.
func Center(r Rect, w, h int) Rect {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	if w > r.W {
		w = r.W
	}
	if h > r.H {
		h = r.H
	}
	return Rect{
		X: r.X + (r.W-w)/2,
		Y: r.Y + (r.H-h)/2,
		W: w,
		H: h,
	}
}
