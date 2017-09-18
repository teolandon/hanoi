package view

import "fmt"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/menu"
import "strconv"
import "unicode/utf8"
import "time"

var (
	initialized      bool                        = false
	itemRanges       map[menu.MenuItem]itemRange = make(map[menu.MenuItem]itemRange)
	mainMenu         menu.Menu
	menuFlag         bool = false
	menuY            int  = -1
	selectedMenuItem int  = -1
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
	mainMenu = defaultMenu()
	printMenu()
	selectMenuItem(0)
	menuFlag = true

loop:

	for menuFlag {
		switch ev := tb.PollEvent(); ev.Type {
		case tb.EventKey:
			if ev.Key == tb.KeyEsc || ev.Ch == 'q' {
				break loop
			} else if ev.Key == tb.KeyArrowLeft {
				navigateMenuLeft()
			} else if ev.Key == tb.KeyArrowRight {
				navigateMenuRight()
			} else if ev.Key == tb.KeyEnter {
				fmt.Println("lololo")
				switch v := mainMenu[selectedMenuItem].(type) {
				case menu.FuncMenuItem:
					v.Function()()
				case menu.IntMenuItem:
					editIntMenuItem(v)
				}
			} else {
				fmt.Println("Uncovered event:", ev)
			}

		default:
			fmt.Println("Uncovered event:", ev)
			time.Sleep(10 * time.Millisecond)
		}
	}

	return exit
}

func editIntMenuItem(item menu.IntMenuItem) {
	// STUB
}

func printMenu() {
	x, y := tb.Size()
	menuSize := menuSize(mainMenu)

	if menuSize > x || y < 3 {
		panic("terminal too small")
	}

	menuStart := (x - menuSize) / 2
	menuY = (y / 2) + 1

	offset := 0
	for _, item := range mainMenu {
		currX := menuStart + offset
		iR := printMenuItem(item, currX, menuY)
		itemRanges[item] = iR
		print("|", iR.end+2, menuY)
		fmt.Println("Size of", item.Name(), ":", itemRanges[item].size())
		offset += itemRanges[item].size() + 3
	}
	erase(menuStart+offset-2, menuY)
}

func erase(x, y int) {
	tb.SetCell(x, y, ' ', tb.ColorDefault, tb.ColorDefault)
	tb.Sync()
}

func deselectMenuItem(index int) {
	item := mainMenu[index]
	iR, ok := itemRanges[item]
	if !ok {
		return
	}

	printMenuItem(item, iR.start, menuY)
}

func printMenuItem(item menu.MenuItem, x, y int) itemRange {
	return printMenuItemHelper(item, x, y, tb.ColorDefault, tb.ColorDefault)
}

func selectMenuItem(index int) {
	item := mainMenu[index]
	if selectedMenuItem != -1 {
		deselectMenuItem(selectedMenuItem)
	}

	iR, ok := itemRanges[item]
	if !ok {
		fmt.Println("Exiting")
		return
	}

	fmt.Println("Selecting menuItem", item.Name(), "whose range is", itemRanges[item])
	printMenuItemHelper(item, iR.start, menuY, tb.ColorDefault|tb.AttrReverse, tb.ColorDefault|tb.AttrReverse)
	selectedMenuItem = index
}

func printMenuItemHelper(item menu.MenuItem, x, y int, fg, bg tb.Attribute) itemRange {
	name := item.Name()
	nameSize := utf8.RuneCountInString(name)
	switch v := item.(type) {
	case menu.FuncMenuItem:
		printHelper(name, x, y, fg, bg)
		return itemRange{x, x + nameSize - 1}
	case menu.IntMenuItem:
		valueStr := strconv.Itoa(v.Value())
		printHelper(name+" "+valueStr, x, y, fg, bg)
		currItemSize := nameSize + utf8.RuneCountInString(valueStr) + 1
		printHelper("▲", x+currItemSize-1, y-1, tb.ColorDefault, tb.ColorDefault)
		printHelper("▼", x+currItemSize-1, y+1, tb.ColorDefault, tb.ColorDefault)
		return itemRange{x, x + currItemSize - 1}
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

func navigateMenuLeft() {
	menuSize := mainMenu.Size()
	newIndex := (selectedMenuItem - 1) % menuSize
	if newIndex < 0 {
		newIndex += menuSize
	}
	selectMenuItem(newIndex)
}

func navigateMenuRight() {
	newIndex := (selectedMenuItem + 1) % mainMenu.Size()
	selectMenuItem(newIndex)
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
	menuFlag = false
	tb.Flush()
	// x, y := tb.Size()
}

func menuSize(menu menu.Menu) (ret int) {
	for _, item := range menu {
		name := item.Name()
		ret += utf8.RuneCountInString(name)
		ret += 1
	}
	ret -= 1
	return
}

func exitMenu() {
	menuFlag = false
}

func defaultMenu() menu.Menu {
	run := menu.NewFuncMenuItem("Run", initHanoi)
	size := menu.NewIntMenuItem("Size", 3)
	exit := menu.NewFuncMenuItem("Exit", exitMenu)
	return []menu.MenuItem{run, size, exit}
}
