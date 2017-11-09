package view

import tb "github.com/nsf/termbox-go"

// Button is a UI component that has a text label,
// and performs an action when hit (using enter)
type Button struct {
	Text string
	Run  func()
	Displayable
}

func (b Button) HandleKey(e KeyEvent) {
	if e.event.Key == tb.KeyEnter || e.event.Key == tb.KeySpace {
		b.Run()
		e.consumed = true
	}
}

func NewButton(text string) Button {
	ret := Button{text, func() {}, defaultDisplayable()}
	ret.SetPalette(alternatePalette)
	ret.SetSize(Size{5, 1})
	return ret
}

func NewButtonWithAction(text string, action func()) Button {
	ret := NewButton(text)
	ret.Run = action
	return ret
}

// ButtonBar is a UI component that groups buttons
// together.
type ButtonBar struct {
	buttons []*Button
	Displayable
}

func (bar ButtonBar) Buttons() []*Button {
	return bar.buttons
}

func (bar ButtonBar) AddButton(b *Button) {
	bar.buttons = append(bar.buttons, b)
}

func (bar ButtonBar) RemoveButton(i int) {
	bar.buttons = append(bar.buttons[:i], bar.buttons[i+1:]...)
}
