package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"fmt"
	"strings"
	"time"
)

var (
	viewArr         = []string{"v1", "v2", "v3", "v5", "modetoggle", "v6"} //list of switchable views
	active          = 0
	cuiGUI          *gocui.Gui //static CUI object
	curCommandIndex = -1
	prevCommands    []string
	lightMode       = false
	prevBottomLine  = 0
	scrollPos       = 0
)

/*
 * Start the CUI
 */

func attachCUI() {
	g, err := gocui.NewGui(gocui.Output256)
	cuiGUI = g
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	//CUI options
	g.Highlight = true
	if lightMode {
		g.SelBgColor = gocui.ColorWhite
		g.BgColor = gocui.ColorWhite
		g.SelFgColor = 210
		g.FgColor = 240
	} else {
		g.SelBgColor = gocui.ColorBlack
		g.BgColor = gocui.ColorBlack
		g.SelFgColor = 31
		g.FgColor = 249
	}
	g.Mouse = true

	g.SetManagerFunc(layout)

	//set keybindings and mouse bindings

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, enterClick); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, mouseClick); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("v2", gocui.KeyArrowUp, gocui.ModNone, prevCommand); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("v2", gocui.KeyArrowDown, gocui.ModNone, forwardCommand); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("v1", gocui.MouseWheelUp, gocui.ModNone, v1ScrollUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("v1", gocui.KeyArrowUp, gocui.ModNone, v1ScrollUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("v1", gocui.MouseWheelDown, gocui.ModNone, v1ScrollDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("v1", gocui.KeyArrowDown, gocui.ModNone, v1ScrollDown); err != nil {
		log.Panicln(err)
	}

	g.Update(func(g *gocui.Gui) error { //write slice to view once the main cui loop starts
		writeSliceToView(attachLog, "v1") //TODO only write screen height size
		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

/*
 * Command scroll up
 */
func prevCommand(gui *gocui.Gui, view *gocui.View) error {
	if curCommandIndex == -1 && len(prevCommands) > 0 {
		curCommandIndex = len(prevCommands) - 1
		view.Clear()
		fmt.Fprintln(view, prevCommands[curCommandIndex])
		view.SetCursor(len(prevCommands[curCommandIndex]), 0)
	} else if curCommandIndex != 0 && curCommandIndex != -1 {
		curCommandIndex--
		view.Clear()
		fmt.Fprintln(view, prevCommands[curCommandIndex])
		view.SetCursor(len(prevCommands[curCommandIndex]), 0)
	}
	return nil
}

/*
 * Command scroll down
 */
func forwardCommand(gui *gocui.Gui, view *gocui.View) error {
	if curCommandIndex == len(prevCommands)-1 && len(prevCommands) > 0 {
		curCommandIndex = -1
		clearCommandView()
	} else if curCommandIndex > -1 && curCommandIndex < len(prevCommands)-1 {
		curCommandIndex++

		view.Clear()
		fmt.Fprintln(view, prevCommands[curCommandIndex])
		view.SetCursor(len(prevCommands[curCommandIndex]), 0)
	}
	return nil
}

/*
 * v1 scroll up
 */
func v1ScrollUp(gui *gocui.Gui, view *gocui.View) error {
	if x, y := view.Origin(); y != 0 {
		_, sy := view.Size()
		if y+4-sy >= 0 && y+4-sy < len(attachLog) && attachLog[y+4-sy] == "" { // y+4 is the y size of v1 (height), catch error with y+4-sy >= 0
			ObtainLogAtIndex(procName, y+4-sy)
		}
		if view.Autoscroll {
			view.Autoscroll = false
			prevBottomLine = y
			scrollPos = prevBottomLine
		}
		scrollPos--
		view.SetOrigin(x, y-1)
	}
	return nil
}

/*
 * v1 scroll down
 */

func v1ScrollDown(gui *gocui.Gui, view *gocui.View) error {
	if x, y := view.Origin(); scrollPos != prevBottomLine {
		if scrollPos == prevBottomLine-1 {
			view.Autoscroll = true
		}
		scrollPos++
		view.SetOrigin(x, y+1)
	}
	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

/*
 * Write text to view box async (from different goroutines)
 */

func writeToView(str string, view string) {
	(*cuiGUI).Update(func(g *gocui.Gui) error {
		out, err := (*cuiGUI).View(view)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(out, str)
		return nil
	})
}
func writeSliceToView(slice []string, view string) {
	(*cuiGUI).Update(func(g *gocui.Gui) error {
		out, err := (*cuiGUI).View(view)
		if err != nil {
			log.Fatal(err)
		}
		for _, str := range slice {
			fmt.Fprintln(out, ""+str+"\u001b[0m")
		}
		return nil
	})
}

/*
 * Function to switch focus to next view
 */

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	//set cursor to appear on command view
	if nextIndex == 1 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

/*
 * Setup CUI layout + views
 */

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", 0, 0, maxX-21, maxY-5); err != nil { // y 4
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = procName + " Console"
		v.Wrap = true
		v.Autoscroll = true
		if _, err = setCurrentViewOnTop(g, "v1"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("v2", 0, maxY-4, maxX-21, maxY-1); err != nil { // y 3
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Command (press ENTER to send)"
		v.Wrap = true
		v.Autoscroll = true
		v.Editable = true
	}
	if v, err := g.SetView("v3", maxX-20, 0, maxX-1, maxY-12); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Info"
		v.Editable = true
		v.Autoscroll = true
		v.Wrap = true
	}
	if v, err := g.SetView("v4", maxX-20, maxY-11, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Buttons"
		v.Wrap = true
	}
	if v, err := g.SetView("v5", maxX-19, maxY-10, maxX-2, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, "      About")
	}
	if v, err := g.SetView("modetoggle", maxX-19, maxY-7, maxX-2, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, "  Color Toggle")
	}
	if v, err := g.SetView("v6", maxX-19, maxY-4, maxX-2, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, "      Exit")
	}
	return nil
}

/*
 * Mouse click event
 */

func mouseClick(gui *gocui.Gui, view *gocui.View) error { //user click on view switch
	var index int
	for i := 0; i < len(viewArr); i++ {
		if viewArr[i] == view.Name() {
			index = i
			break
		}
	}
	var maxi, cur int
	cur = active
	if cur < 0 {
		cur = len(viewArr) - 1
	}
	if index < cur {
		maxi = (len(viewArr) - cur) + index
	} else {
		maxi = index - cur
	}
	for i := 0; i < maxi; i++ { //cycle through views
		nextView(gui, view)
	}
	(*cuiGUI).Update(func(g *gocui.Gui) error { //move cursor to beginning async
		out, err := (*cuiGUI).View("v2")
		if err != nil {
			return err
		}
		out.SetCursor(0, 0) //move cursor
		return nil
	})
	return nil
}

/*
 * Enter button click event
 */
func enterClick(gui *gocui.Gui, view *gocui.View) error { //send command
	if view.Name() == "v6" {
		return quit(gui, view)
	} else if view.Name() == "v5" {
		aboutPopup(gui)
		return nil
	} else if view.Name() == "modetoggle" {
		toggleMode(gui)
		return nil
	}

	//otherwise, send command

	out, err := (*cuiGUI).View("v2")
	if err != nil {
		log.Fatal(err)
	}
	out.Rewind() //move buffer to beginning of line
	b := out.ViewBuffer()
	b = strings.Replace(b, "\n", "", -1)
	clearCommandView()

	prevCommands = append(prevCommands, b)
	curCommandIndex = -1
	SendCommand(b, procName) //send command over grpc to server
	return nil
}

func clearCommandView() {
	(*cuiGUI).Update(func(g *gocui.Gui) error { //clear the command view's text and move cursor to beginning async
		out, err := (*cuiGUI).View("v2")
		if err != nil {
			return err
		}
		out.Clear()         //clear text
		out.SetCursor(0, 0) //move cursor
		return nil
	})
}

func changeViewColour(view string, ansicolour string) {
	(*cuiGUI).Update(func(g *gocui.Gui) error {
		out, err := (*cuiGUI).View(view)
		if err != nil {
			return err
		}
		str := ansicolour + out.ViewBuffer() + "\u001b[0m"
		out.Clear() //clear text
		fmt.Fprintln(out, str)
		return nil
	})
}

/*
 * Toggles the colours of the cui
 */

func toggleMode(gui *gocui.Gui) {
	lightMode = !lightMode
	if lightMode {
		gui.SelBgColor = gocui.ColorWhite
		gui.BgColor = gocui.ColorWhite
		gui.SelFgColor = 210
		gui.FgColor = 240
	} else {
		gui.SelBgColor = gocui.ColorBlack
		gui.BgColor = gocui.ColorBlack
		gui.SelFgColor = 31
		gui.FgColor = 249
	}
	for _, v := range gui.Views() {
		v.BgColor, v.FgColor = gui.BgColor, gui.FgColor
		v.SelBgColor, v.SelFgColor = gui.SelBgColor, gui.SelFgColor
	}
}

func aboutPopup(gui *gocui.Gui) { //show popup for extra info
	maxX, maxY := gui.Size()
	if v, err := gui.SetView("popup", maxX/2-10, maxY/2, maxX/2+10, maxY/2+4); err != nil {
		fmt.Fprintln(v, "EstiCli "+version)
		fmt.Fprintln(v, "――――――――――――――――――――")
		fmt.Fprintln(v, "EspiDev approves")
		go func() { //delete popup after 1.5 s
			t, _ := time.ParseDuration("1500ms")
			time.Sleep(t)
			(*cuiGUI).Update(func(g *gocui.Gui) error {
				return gui.DeleteView("popup")
			})
		}()
	}
}

/*
 * Quit CUI
 */

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
