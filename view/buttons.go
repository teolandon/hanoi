package view

import tb "github.com/nsf/termbox-go"

type Button struct {
	Text string
	Run  func()
	Parent
	Displayable
}

type Parent struct {
	p Focusable
}

func (par Parent) Parent() Focusable {
	return par.p
}

func (b Button) HandleKey(e KeyEvent) {
	if e.event.Key == tb.KeyEnter || e.event.Key == tb.KeySpace {
		b.Run()
		e.consumed = true
	}
}

func NewButton(text string, parent Focusable) Button {
	ret := Button{text, func() {}, Parent{parent}, defaultDisplayable()}

	return ret
}

func NewButtonWithAction(text string, action func()) Button {
	return *new(Button)
}
