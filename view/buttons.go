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

func NewButton(text string, parent Displayable) Button {
	ret := Button{text, func() {}, defaultDisplayable()}
	return ret
}

func NewButtonWithAction(text string, action func()) Button {
	return *new(Button)
}
