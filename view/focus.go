package view

import (
	"fmt"
)

type KeyHandler interface {
	HandleKey(event KeyEvent)
}

type OkCancelDialog struct {
	container TitledContainer
	b         Button
}

func RequestFocus(f Displayable) {
	SetFocused(f)
}

func executeEvent(event KeyEvent, f Displayable) {
	if event.consumed || f == nil {
		return
	}

	handler, ok := f.(KeyHandler)
	if ok {
		fmt.Println("Handling", handler)
		handler.HandleKey(event)
	}
	fmt.Println("Executing handler on", f.Parent())

	executeEvent(event, f.Parent())
}
