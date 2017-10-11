package view

type KeyHandler interface {
	HandleKey(event KeyEvent)
}

type Focusable interface {
	Parent() Focusable
}

func RequestFocus(f Focusable) {
	setFocused(f)
}

func executeEvent(event KeyEvent, f Focusable) {
	if event.consumed || f == nil {
		return
	}

	handler, ok := f.(KeyHandler)
	if ok {
		handler.HandleKey(event)
	}

	executeEvent(event, f.Parent())
}
