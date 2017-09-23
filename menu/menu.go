package menu

type Menu []*MenuItem

type MenuItemType uint8

const (
	FuncMenuItem MenuItemType = iota
	IntMenuItem
)

func (m Menu) Size() int {
	return len(m)
}

type MenuItem struct {
	Name     string
	Type     MenuItemType
	function func()
	value    int
}

func (m MenuItem) Function() func() {
	if m.Type != FuncMenuItem {
		panic("type error")
	}
	return m.function
}

func (m MenuItem) Value() int {
	if m.Type != IntMenuItem {
		panic("type error")
	}
	return m.value
}

func (m *MenuItem) SetValue(newVal int) {
	item := *m
	if m.Type != IntMenuItem {
		panic("type error")
	}
	item.value = newVal
	*m = item
}

func NewFuncMenuItem(name string, function func()) MenuItem {
	return MenuItem{name, FuncMenuItem, function, 0}
}

func NewIntMenuItem(name string, defaultVal int) MenuItem {
	return MenuItem{name, IntMenuItem, func() {}, defaultVal}
}
