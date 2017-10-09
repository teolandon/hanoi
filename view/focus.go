package view

type Focusable struct {
	parent      *Focusable
	handleEvent func(event keyEvent)
}

func (f Focusable) RequestFocus() {
	setFocused(&f)
}

func (f Focusable) AcceptInput() {
	go func() {
		for {
			select {
			case event := <-eventChannel:
				f.executeEvent(event)
			case <-stopChannel:
				return
			}
		}
	}()
}

func (f Focusable) executeEvent(event keyEvent) {
	if event.consumed {
		return
	}

	f.handleEvent(event)
	f.parent.handleEvent(event)
}
