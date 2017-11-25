package view

import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/view/colors"

// Button is a UI component that has a text label,
// and performs an action when hit (using enter)
type Button struct {
	Text string
	Run  func()
	displayable
}

func (b Button) String() string {
	return "Button with text \"" + b.Text + "\""
}

func (b Button) Draw() {
	// STUB
}

func (b Button) HandleKey(e KeyEvent) {
	if e.event.Key == tb.KeyEnter || e.event.Key == tb.KeySpace {
		b.Run()
		e.consumed = true
	}
}

func (b Button) Size() areas.Size {
	orig := b.displayable.Size()
	return areas.NewSize(orig.Width(), 1)
}

func NewButton(text string) Button {
	ret := Button{text, func() {}, defaultDisplayable()}
	ret.SetPalette(colors.AlternatePalette)
	ret.SetSize(areas.NewSize(5, 1))
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
