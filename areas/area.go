package areas

type Size struct {
	width  int
	height int
}

func NewSize(w, h int) Size {
	return Size{w, h}
}

func (s Size) Width() int {
	return s.width
}

func (s Size) Height() int {
	return s.height
}

type Area struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

func (a Area) Width() int {
	return a.x2 - a.x1
}

func (a Area) Height() int {
	return a.y2 - a.y1
}

func (a Area) X1() int {
	return a.x1
}

func (a Area) X2() int {
	return a.x2
}

func (a Area) Y1() int {
	return a.y1
}

func (a Area) Y2() int {
	return a.y2
}

func New(x1, x2, y1, y2 int) Area {
	return Area{x1, x2, y1, y2}
}

func NewFromSize(x, y int, size Size) Area {
	return Area{x, x + size.Width(), y, y + size.Height()}
}

func (a Area) SubArea(sub Area) Area {
	// TODO: Include more checks
	return Area{a.x1 + sub.x1, a.x1 + sub.x2, a.y1 + sub.y1, a.y1 + a.y2}
}

func (a Area) Padded(p Padding) Area {
	// TODO: Include more checks
	return Area{a.x1 + p.Left, a.x2 - p.Right, a.y1 + p.Up, a.y2 - p.Down}
}

type Padding struct {
	Up    int
	Down  int
	Left  int
	Right int
}
