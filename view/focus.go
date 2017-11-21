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

func executeEvent(event KeyEvent, f Displayable) {
	fmt.Println("Executing handler on", f)
	if f != nil {
		fmt.Println("Its parent:", f.Parent())
	}
	fmt.Println()
	if event.consumed || f == nil {
		return
	}

	handler, ok := f.(KeyHandler)
	if ok {
		fmt.Println("Handling", handler)
		handler.HandleKey(event)
	} else {
		fmt.Println("Not a handler,", handler)
	}

	executeEvent(event, f.Parent())
}
