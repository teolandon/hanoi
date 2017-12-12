package view

import (
	"github.com/teolandon/hanoi/utils/log"
	"github.com/teolandon/hanoi/view/displayable"
)

type KeyHandler interface {
	HandleKey(event KeyEvent)
}

func executeEvent(event KeyEvent, f displayable.Displayable) {
	log.Log("Executing handler on", f)
	if f != nil {
		log.Log("Its parent:", f.Parent())
	}
	log.Log()
	if event.consumed {
		return
	}

	_, ok := f.(mainContainer)
	if f == nil || ok { // If any displayable in the chain is parent-less, call event on main controller
		main.HandleKey(event)
		return
	}

	handler, ok := f.(KeyHandler)
	if ok {
		log.Log("Handling", handler, "with key", string(event.event.Ch))
		handler.HandleKey(event)
	} else {
		log.Log("Not a handler,", handler)
	}

	executeEvent(event, f.Parent())
}
