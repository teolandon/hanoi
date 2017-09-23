package menu

import "fmt"
import tb "github.com/nsf/termbox-go"
import "strconv"
import "unicode/utf8"
import "time"

var (
	initialized      bool                 = false
	itemRanges       map[string]itemRange = make(map[string]itemRange)
	mainMenu         Menu
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

func initMenu() {
	menuFlag = true

loop:
	for menuFlag {
		printMenu()
		switch ev := tb.PollEvent(); ev.Type {
		case tb.EventKey:
			if ev.Key == tb.KeyEsc || ev.Ch == 'q' {
				break loop
			} else if ev.Key == tb.KeyArrowLeft {
				navigateMenuLeft()
			} else if ev.Key == tb.KeyArrowRight {
				navigateMenuRight()
			} else if ev.Key == tb.KeyEnter {
				item := mainMenu[selectedMenuItem]
				switch item.Type {
				case FuncMenuItem:
					item.Function()()
				case IntMenuItem:
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
}

func editIntMenuItem() {
	if selectedMenuItem < 0 {
		return
	}
	item := mainMenu[selectedMenuItem]

	temp := selectedMenuItem
	deselectCurrentMenuItem()
	defer func() {
		selectMenuItem(temp)
	}()

	iR, ok := itemRanges[item.Name]
	if !ok {
		return
	}

	valueLength := intCharSize(item.Value())
	valuePos := iR.end - valueLength
	arrowPos := valuePos + (valueLength)/2

	printHelper("▲", arrowPos, menuY-1, selectedColor, selectedColor)
	printHelper(strconv.Itoa(item.Value()), valuePos, menuY,
		tb.ColorDefault|tb.AttrUnderline, tb.ColorDefault|tb.AttrUnderline)
	printHelper("▼", arrowPos, menuY+1, selectedColor, selectedColor)
	tb.Sync()

	defer func() {
		printHelper("▲", iR.end-1, menuY-1, tb.ColorDefault, tb.ColorDefault)
		printHelper(strconv.Itoa(item.Value()), iR.end-1, menuY,
			tb.ColorDefault, tb.ColorDefault)
		printHelper("▼", iR.end-1, menuY+1, tb.ColorDefault, tb.ColorDefault)
		tb.Sync()
	}()

intEditLoop:
	for {
		switch ev := tb.PollEvent(); ev.Type {
		case tb.EventKey:
			if ev.Key == tb.KeyEsc || ev.Key == tb.KeyEnter || ev.Ch == 'q' {
				break intEditLoop
			} else if ev.Key == tb.KeyArrowUp {
				incIntMenuItem(temp)
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

func incIntMenuItem(index int) {
	item := mainMenu[index]
	oldCharSize := intCharSize(item.Value())
	newVal := item.Value() + 1
	item.SetValue(newVal)
	newCharSize := intCharSize(newVal)
	if oldCharSize != newCharSize {
		populateRangesFrom(index)
		printItems(index, mainMenu.Size())
	} else {
		printMenuItem(item)
	}
}

func decIntMenuItem(index int) {
	item := mainMenu[index]
	oldCharSize := intCharSize(item.Value())
	newVal := item.Value() - 1
	item.SetValue(newVal)
	newCharSize := intCharSize(newVal)
	if oldCharSize != newCharSize {
		populateRangesFrom(index)
		printItems(index, mainMenu.Size())
	} else {
		printMenuItem(item)
	}
}

func printMenu() {
	populateRanges(0, mainMenu.Size())
	printItems(0, mainMenu.Size())
	selectMenuItem(selectedMenuItem)
}

func populateRangesFrom(start int) {
	populateRanges(start, mainMenu.Size())
}

func populateRanges(start, end int) {
	x, y := tb.Size()
	menuSize := menuCharSize(mainMenu)

	if menuSize > x || y < 3 {
		panic("terminal too small")
	}

	menuY = (y / 2) + 1

	offset := 0
	if start > 0 {
		offset += itemRanges[mainMenu[start-1].Name].end + 3
		fmt.Println("Increased offset by", offset)
	} else {
		offset += (x - menuSize) / 2
	}
	for _, item := range mainMenu[start:end] {
		nameSize := utf8.RuneCountInString(item.Name)

		tb.Sync()
		switch item.Type {
		case FuncMenuItem:
			itemRanges[item.Name] = itemRange{offset, offset + nameSize}
		case IntMenuItem:
			currItemSize := nameSize + intCharSize(item.Value()) + 1
			itemRanges[item.Name] = itemRange{offset, offset + currItemSize}
		}

		offset += itemRanges[item.Name].size() + 3
	}
}

func printItems(start, end int) {
	x1 := itemRanges[mainMenu[start].Name].start
	x2 := itemRanges[mainMenu[end-1].Name].end
	eraseRange(x1, x2, menuY)
	eraseRange(x1, x2, menuY-1)
	eraseRange(x1, x2, menuY+1)
	for _, item := range mainMenu {
		printMenuItem(item)
		printStr("|", itemRanges[item.Name].end+1, menuY)
	}
	final, _ := tb.Size()
	eraseRange(itemRanges[mainMenu[mainMenu.Size()-1].Name].end, final, menuY)
	tb.Sync()
}

func eraseRange(start, end, y int) {
	for ; start < end; start++ {
		eraseHelper(start, y)
	}
	tb.Sync()
}

// Does NOT Sync terminal
func eraseHelper(x, y int) {
	tb.SetCell(x, y, ' ', tb.ColorDefault, tb.ColorDefault)
}

func deselectCurrentMenuItem() {
	if selectedMenuItem == -1 {
		return
	}

	item := mainMenu[selectedMenuItem]

	unHighlightRange(itemRanges[item.Name])
	selectedMenuItem = -1
}

func printMenuItem(item *MenuItem) {
	printMenuItemHelper(item, tb.ColorDefault, tb.ColorDefault)
	tb.Sync()
}

func selectMenuItem(index int) {
	deselectCurrentMenuItem()
	item := mainMenu[index]

	highlightItem(item)
	selectedMenuItem = index
}

func printMenuItemHelper(item *MenuItem, fg, bg tb.Attribute) {
	name := item.Name
	nameSize := utf8.RuneCountInString(name)

	x := itemRanges[item.Name].start

	switch item.Type {
	case FuncMenuItem:
		printHelper(name, x, menuY, fg, bg)
	case IntMenuItem:
		valueStr := strconv.Itoa(item.Value())
		printHelper(name+" "+valueStr, x, menuY, fg, bg)
		valueLength := utf8.RuneCountInString(valueStr)
		arrowPos := x + nameSize + 1 + (valueLength)/2
		printHelper("▲", arrowPos, menuY-1, tb.ColorDefault, tb.ColorDefault)
		printHelper("▼", arrowPos, menuY+1, tb.ColorDefault, tb.ColorDefault)
	}
}

func highlightMenuRange(start, end int) {
	for _, item := range mainMenu[start:end] {
		highlightItem(item)
	}
}

func highlightItem(item *MenuItem) {
	iR, ok := itemRanges[item.Name]
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

// Does NOT Sync Terminal
func printHelper(s string, x, y int, fg, bg tb.Attribute) {
	for i, r := range s {
		tb.SetCell(x+i, y, r, fg, bg)
	}
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

func Init(menu *Menu) {
	if !initialized && menu.Size() > 0 {
		initialized = true
		mainMenu = *menu
		selectedMenuItem = 0
		tb.Flush()
		initMenu()
	}
}

func intCharSize(i int) int {
	return utf8.RuneCountInString(strconv.Itoa(i))
}

func menuCharSize(m Menu) (ret int) {
	for _, item := range m {
		name := item.Name
		nameSize := utf8.RuneCountInString(name)
		switch item.Type {
		case FuncMenuItem:
			ret += nameSize + 3
		case IntMenuItem:
			valSize := intCharSize(item.Value())
			ret += nameSize + valSize + 4
		}
	}
	ret -= 1
	return
}

func exitMenu() {
	menuFlag = false
}

func DefaultMenu() Menu {
	run := NewFuncMenuItem("Run", func() {})
	size := NewIntMenuItem("Size", 3)
	exit := NewFuncMenuItem("Exit", exitMenu)
	return []*MenuItem{&run, &size, &exit}
}
