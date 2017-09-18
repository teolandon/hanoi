package menu

type Menu []MenuItem

var funcMap map[FuncMenuItem]func() = make(map[FuncMenuItem]func())

func (m Menu) Size() int {
	return len(m)
}

type MenuItem interface {
	Name() string
}

type FuncMenuItem struct {
	name string
}

type IntMenuItem struct {
	name  string
	value int
}

func (m FuncMenuItem) Name() string {
	return m.name
}

func (m FuncMenuItem) Function() func() {
	ret, ok := funcMap[m]
	if !ok {
		return func() {}
	}
	return ret
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
	ret := FuncMenuItem{name}
	funcMap[ret] = function
	return FuncMenuItem{name}
}

func NewIntMenuItem(name string, defaultVal int) MenuItem {
	return IntMenuItem{name, defaultVal}
}
