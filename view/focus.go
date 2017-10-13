package view

type KeyHandler interface {
	HandleKey(event KeyEvent)
}

type OkCancelDialog struct {
	container TitledContainer
	b         Button
}

func RequestFocus(f Displayable) {
	setFocused(f)
}

func executeEvent(event KeyEvent, f Displayable) {
	if event.consumed || f == nil {
		return
	}

	handler, ok := f.(KeyHandler)
	if ok {
		handler.HandleKey(event)
	}

	executeEvent(event, f.Parent())
}
