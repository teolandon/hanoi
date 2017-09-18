package view

import "fmt"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/menu"
import "strconv"
import "unicode/utf8"
import "time"

var (
	initialized      bool                         = false
	itemRanges       map[*menu.MenuItem]itemRange = make(map[*menu.MenuItem]itemRange)
	menuY            int                          = -1
	selectedMenuItem *menu.MenuItem               = nil
)

var (
	selected tb.Attribute = tb.ColorDefault | tb.AttrReverse
)

type itemRange struct {
	start int
	end   int
}

func (iR itemRange) size() int {
	return iR.end - iR.start + 1
}

func (iR itemRange) String() string {
	return "{" + strconv.Itoa(iR.start) + ", " + strconv.Itoa(iR.end) + "}"
}

type choice uint16

const (
	exit choice = iota
	cont
)

func initMenu() choice {
	m := defaultMenu()
	printMenu(m)
	selectMenuItem(m.Items[0])

loop:
	for {
		switch ev := tb.PollEvent(); ev.Type {
		case tb.EventKey:
			if ev.Key == tb.KeyEsc || ev.Ch == 'q' {
				break loop
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	return exit
}

func printMenu(m menu.Menu) {
	x, y := tb.Size()
	menuSize := menuSize(m)

	if menuSize > x || y < 3 {
		panic("terminal too small")
	}

	menuStart := (x - menuSize) / 2
	menuY = (y / 2) + 1

	offset := 0
	for _, item := range m.Items {
		currX := menuStart + offset
		iR := printMenuItem(item, currX, menuY)
		itemRanges[&item] = iR
		print("|", iR.end+1, menuY)
		offset += itemRanges[&item].size() + 2
		fmt.Println("Range of menuitem", &item, ":", itemRanges[&item])
		fmt.Println("Offset of menuitem:", offset)
	}
	erase(menuStart+offset-2, menuY)
}

func erase(x, y int) {
	tb.SetCell(x, y, ' ', tb.ColorDefault, tb.ColorDefault)
	tb.Sync()
}

func deselectMenuItem(item menu.MenuItem) {
	iR, ok := itemRanges[&item]
	if !ok {
		return
	}

	printMenuItem(item, iR.start, menuY)
}

func printMenuItem(item menu.MenuItem, x, y int) itemRange {
	return printMenuItemHelper(item, x, y, tb.ColorDefault, tb.ColorDefault)
}

func selectMenuItem(item menu.MenuItem) {
	if selectedMenuItem != nil {
		deselectMenuItem(*selectedMenuItem)
	}

	iR, ok := itemRanges[&item]
	if !ok {
		return
	}

	printMenuItemHelper(item, iR.start, menuY, tb.ColorDefault, tb.ColorDefault|tb.AttrReverse)
	selectedMenuItem = &item
}

func printMenuItemHelper(item menu.MenuItem, x, y int, fg, bg tb.Attribute) itemRange {
	name := item.Name()
	nameSize := utf8.RuneCountInString(name)
	switch v := item.(type) {
	case menu.FuncMenuItem:
		printHelper(name, x, y, fg, bg)
		return itemRange{x, x + nameSize}
	case menu.IntMenuItem:
		valueStr := strconv.Itoa(v.Value())
		printHelper(name+" "+valueStr, x, y, fg, bg)
		currItemSize := nameSize + utf8.RuneCountInString(valueStr) + 1
		printHelper("▲", x+currItemSize-1, y-1, fg, bg)
		printHelper("▼", x+currItemSize-1, y+1, fg, bg)
		return itemRange{x, x + currItemSize}
	default:
		return itemRange{-1, -1}
	}
}

func print(s string, x, y int) {
	printHelper(s, x, y, tb.ColorDefault, tb.ColorDefault)
}

func printHelper(s string, x, y int, fg, bg tb.Attribute) {
	for i, r := range s {
		tb.SetCell(x+i, y, r, fg, bg)
	}
	tb.Sync()
}

func Init() {
	if !initialized {
		initialized = true
		tb.Flush()
		ch := initMenu()
		initChoiceLoop(ch)
	}
}

func initChoiceLoop(ch choice) {
	switch ch {
	case exit:
		return
	case cont:
		initHanoi()
	}
}

func initHanoi() {
	tb.Flush()
	// x, y := tb.Size()
}

func menuSize(menu menu.Menu) (ret int) {
	for _, item := range menu.Items {
		name := item.Name()
		ret += utf8.RuneCountInString(name)
		ret += 1
	}
	ret -= 1
	return
}

func defaultMenu() menu.Menu {
	run := menu.NewFuncMenuItem("Run", initHanoi)
	size := menu.NewIntMenuItem("Size", 3)
	return menu.Menu{[]menu.MenuItem{run, size}}
}
