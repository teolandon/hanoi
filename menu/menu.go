package menu

type Menu []*MenuItem

func (m Menu) Size() int {
	return len(m)
}

type MenuItem interface {
	Name() string
}

type FuncMenuItem struct {
	name     string
	function func()
}

type IntMenuItem struct {
	name  string
	value int
}

func (m FuncMenuItem) Name() string {
	return m.name
}

func (m FuncMenuItem) Function() func() {
	return m.function
}

func (m IntMenuItem) Name() string {
	return m.name
}

func (m IntMenuItem) Value() int {
	return m.value
}

func (m IntMenuItem) SetValue(newVal int) {
	m.value = newVal
}

func NewFuncMenuItem(name string, function func()) MenuItem {
	return FuncMenuItem{name, function}
}

func NewIntMenuItem(name string, defaultVal int) MenuItem {
	return IntMenuItem{name, defaultVal}
}
