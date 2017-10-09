package view

type Button struct {
	Text string
	Run  func()
	Displayable
	Focusable
}
