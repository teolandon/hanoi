package areas

// Size describes a property of an object
// using two integers, the width and the height.
type Size struct {
	width  int
	height int
}

// NewSize returns a new Size struct with the
// given attributes.
func NewSize(w, h int) Size {
	return Size{w, h}
}

// Width returns the width of a Size s.
func (s Size) Width() int {
	return s.width
}

// Height returns the height of a Size s.
func (s Size) Height() int {
	return s.height
}

// Padded returns an area that will represent
// an area starting at p.Left and p.Up, and will
// have the original width and height minus the
// width and height of the padding.
func (s Size) Padded(p Padding) Area {
	return Area{p.Left, s.width - p.Right, p.Up, s.height - p.Down}
}

// Area describes the property of an object that occupies
// a cetrain rectangular area, from point (x1, y1) to point
// (x2-1, y2-1).
type Area struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

// Size returns that size of the area a.
func (a Area) Size() Size {
	return Size{a.Width(), a.Height()}
}

// Width returns the width of the area a.
func (a Area) Width() int {
	return a.x2 - a.x1
}

// Height returns the height of the area a.
func (a Area) Height() int {
	return a.y2 - a.y1
}

// X1 is the getter for the x-coordinate x1 of area a.
func (a Area) X1() int {
	return a.x1
}

// X2 is the getter for the x-coordinate x2 of area a.
func (a Area) X2() int {
	return a.x2
}

// Y2 is the getter for the y-coordinate y1 of area a.
func (a Area) Y1() int {
	return a.y1
}

// Y1 is the getter for the y-coordinate y2 of area a.
func (a Area) Y2() int {
	return a.y2
}

// New returns a new Area with the given coordinates.
func New(x1, x2, y1, y2 int) Area {
	return Area{x1, x2, y1, y2}
}

// NewFromSize returns a new Area starting at (x, y), and has Size size.
func NewFromSize(x, y int, size Size) Area {
	return Area{x, x + size.Width(), y, y + size.Height()}
}

// SubArea returns a new Area such that the Area sub is
// transposed into the Area a.
func (a Area) SubArea(sub Area) Area {
	// TODO: Include more checks
	return Area{a.x1 + sub.x1, a.x1 + sub.x2, a.y1 + sub.y1, a.y1 + sub.y2}
}

// Padding describes 4 paddings to be applied to each side
// of a Size or Area.
type Padding struct {
	Left  int
	Right int
	Up    int
	Down  int
}
