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
	selectedColor tb.Attribute = tb.ColorDefault | tb.AttrReverse
)

type itemRange struct {
	start int
	end   int
}

func (iR itemRange) size() int {
	return iR.end - iR.start
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
				switch v := mainMenu[selectedMenuItem].(type) {
				case menu.FuncMenuItem:
					v.Function()()
				case menu.IntMenuItem:
					editIntMenuItem()
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

func editIntMenuItem() {
	if selectedMenuItem < 0 {
		return
	}
	curr := mainMenu[selectedMenuItem]
	item, ok := curr.(menu.IntMenuItem)
	if !ok {
		return
	}

	temp := selectedMenuItem
	defer func() {
		selectMenuItem(temp)
	}()
	deselectCurrentMenuItem()

	iR, ok := itemRanges[item]
	if !ok {
		return
	}

	printHelper("▲", iR.end-1, menuY-1, tb.ColorDefault|tb.AttrReverse, tb.ColorDefault|tb.AttrReverse)
	printHelper(strconv.Itoa(item.Value()), iR.end-1, menuY,
		tb.ColorDefault|tb.AttrReverse, tb.ColorDefault|tb.AttrReverse)
	printHelper("▼", iR.end-1, menuY+1, tb.ColorDefault|tb.AttrReverse, tb.ColorDefault|tb.AttrReverse)
	defer func() {
		printHelper("▲", iR.end-1, menuY-1, tb.ColorDefault, tb.ColorDefault)
		printHelper(strconv.Itoa(item.Value()), iR.end-1, menuY,
			tb.ColorDefault, tb.ColorDefault)
		printHelper("▼", iR.end-1, menuY+1, tb.ColorDefault, tb.ColorDefault)
	}()

intEditLoop:
	for {
		switch ev := tb.PollEvent(); ev.Type {
		case tb.EventKey:
			if ev.Key == tb.KeyEsc || ev.Key == tb.KeyEnter || ev.Ch == 'q' {
				break intEditLoop
			} else if ev.Key == tb.KeyArrowUp {
				incIntMenuItem(item)
			} else if ev.Key == tb.KeyArrowDown {
				decIntMenuItem(temp)
			} else {
				fmt.Println("Uncovered event:", ev)
			}

		default:
			fmt.Println("Uncovered event:", ev)
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func incIntMenuItem(item menu.MenuItem) {
	intItem, ok := item.(menu.IntMenuItem)
	if !ok {
		return
	}

	oldCharSize := intCharSize(intItem.Value())
	newVal := intItem.Value() + 1
	intItem.SetValue(intItem.Value() + 1)
	newCharSize := intCharSize(newVal)
	if oldCharSize != newCharSize {
		fmt.Println("Value size changed!")
	} else {
		fmt.Println("Reprinting only this place")
	}
}

func decIntMenuItem(index int) {

}

func printMenu() {
	printMenuRange(0, mainMenu.Size())
	printItems(0, mainMenu.Size())
}

func printMenuRange(start, end int) {
	populateRanges(start, end)
}

func populateRanges(start, end int) {
	x, y := tb.Size()
	menuSize := menuCharSize(mainMenu)

	if menuSize > x || y < 3 {
		panic("terminal too small")
	}

	menuStart := (x - menuSize) / 2
	menuY = (y / 2) + 1
	offset := 0
	for _, item := range mainMenu[start:end] {
		currX := menuStart + offset
		nameSize := utf8.RuneCountInString(item.Name())

		tb.Sync()
		switch v := item.(type) {
		case menu.FuncMenuItem:
			itemRanges[item] = itemRange{currX, currX + nameSize}
		case menu.IntMenuItem:
			currItemSize := nameSize + intCharSize(v.Value()) + 1
			itemRanges[item] = itemRange{currX, currX + currItemSize}
		}

		offset += itemRanges[item].size() + 3
	}
}

func printItems(start, end int) {
	for _, item := range mainMenu {
		printMenuItem(item)
		printStr("|", itemRanges[item].end+1, menuY)
	}
	erase(itemRanges[mainMenu[mainMenu.Size()-1]].end+1, menuY)
}

func erase(x, y int) {
	tb.SetCell(x, y, ' ', tb.ColorDefault, tb.ColorDefault)
	tb.Sync()
}

func deselectCurrentMenuItem() {
	if selectedMenuItem == -1 {
		return
	}

	item := mainMenu[selectedMenuItem]

	unHighlightRange(itemRanges[item])
	selectedMenuItem = -1
}

func printMenuItem(item menu.MenuItem) {
	printMenuItemHelper(item, tb.ColorDefault, tb.ColorDefault)
}

func selectMenuItem(index int) {
	deselectCurrentMenuItem()
	item := mainMenu[index]

	highlightItem(item)
	selectedMenuItem = index
}

func printMenuItemHelper(item menu.MenuItem, fg, bg tb.Attribute) {
	name := item.Name()
	nameSize := utf8.RuneCountInString(name)

	x := itemRanges[item].start

	switch v := item.(type) {
	case menu.FuncMenuItem:
		printHelper(name, x, menuY, fg, bg)
	case menu.IntMenuItem:
		valueStr := strconv.Itoa(v.Value())
		printHelper(name+" "+valueStr, x, menuY, fg, bg)
		currItemSize := nameSize + utf8.RuneCountInString(valueStr) + 1
		printHelper("▲", x+currItemSize-1, menuY-1, tb.ColorDefault, tb.ColorDefault)
		printHelper("▼", x+currItemSize-1, menuY+1, tb.ColorDefault, tb.ColorDefault)
	}
}

func highlightMenuRange(start, end int) {
	for _, item := range mainMenu[start:end] {
		highlightItem(item)
	}
}

func highlightItem(item menu.MenuItem) {
	iR, ok := itemRanges[item]
	if !ok {
		return
	}

	highlightRange(iR)
}

func highlightRange(iR itemRange) {
	for i := iR.start; i < iR.end; i++ {
		highlight(i, menuY)
	}
}

func unHighlightRange(iR itemRange) {
	for i := iR.start; i < iR.end; i++ {
		unHighlight(i, menuY)
	}
}

func highlight(x, y int) {
	highlightHelper(x, y, tb.ColorDefault|tb.AttrReverse)
}

func unHighlight(x, y int) {
	highlightHelper(x, y, tb.ColorDefault)
}

func highlightHelper(x, y int, color tb.Attribute) {
	termX, _ := tb.Size()
	cells := tb.CellBuffer()
	char := cells[x+y*termX].Ch
	tb.SetCell(x, y, char, color, color)
	tb.Sync()
}

func printStr(s string, x, y int) {
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

func intCharSize(i int) int {
	return utf8.RuneCountInString(strconv.Itoa(i))
}

func menuCharSize(m menu.Menu) (ret int) {
	for _, item := range m {
		name := item.Name()
		nameSize := utf8.RuneCountInString(name)
		switch v := item.(type) {
		case menu.FuncMenuItem:
			ret += nameSize + 3
		case menu.IntMenuItem:
			valSize := intCharSize(v.Value())
			ret += nameSize + valSize + 4
		}
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
