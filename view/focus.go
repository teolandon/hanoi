package view

import (
	"github.com/teolandon/hanoi/utils/log"
)

type KeyHandler interface {
	HandleKey(event KeyEvent)
}

type OkCancelDialog struct {
	container TitledContainer
	b         Button
}

func executeEvent(event KeyEvent, f Displayable) {
	log.Log("Executing handler on", f)
	if f != nil {
		log.Log("Its parent:", f.Parent())
	}
	log.Log()
	if event.consumed || f == nil {
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
