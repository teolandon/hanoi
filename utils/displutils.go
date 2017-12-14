package utils

import (
	"github.com/teolandon/hanoi/areas"
	"github.com/teolandon/hanoi/utils/log"
	"github.com/teolandon/hanoi/view/displayable"
)

func LayoutedArea(size areas.Size, d displayable.Displayable) areas.Area {
	var width, height, x, y int
	layout := d.Layout()
	contentSize := d.Size()
	log.Log("ContentSize for", d, ":", contentSize)
	switch layout {
	case displayable.FitToParent:
		log.Log("Displayable", d, "is Fit To parent")
		width = size.Width()
		height = size.Height()
		x = 0
		y = 0
	case displayable.Centered:
		log.Log("Displayable", d, "is Centered")
		width = contentSize.Width()
		height = contentSize.Height()
		x = (size.Width() - width) / 2
		y = (size.Height() - height) / 2
	default:
		panic("nooo")
	}
	ret := areas.NewFromSize(x, y, areas.NewSize(width, height))
	log.Log("Layout of size", size, "is", layout, "and results in", ret)
	return ret
}
