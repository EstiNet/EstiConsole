package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"fmt"
)

var (
	viewArr = []string{"v1", "v2", "v3"} //list of switchable views
	active  = 0
	cuiGUI  **gocui.Gui //static CUI object
)

/*
 * Start the CUI
 */

func attachCUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	cuiGUI = &g
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//CUI options
	g.Highlight = true
	g.SelFgColor = gocui.ColorWhite
	g.Mouse = true

	g.SetManagerFunc(layout)

	//set keybindings and mouse bindings

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error { //send command
		out, err := (**cuiGUI).View("v2")
		if err != nil {
			log.Fatal(err)
		}
		out.Rewind() //move buffer to beginning of line
		b := out.ViewBuffer()
		(**cuiGUI).Update(func (g *gocui.Gui) error { //clear the command view's text and move cursor to beginning async
			out, err := (**cuiGUI).View("v2")
			if err != nil {
				return err
			}
			out.Clear() //clear text
			out.SetCursor(0, 0) //move cursor
			return nil
		})

		SendCommand(string(b), procName) //send command over grpc to server
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error { //user click on view switch
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
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
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
	(**cuiGUI).Update(func(g *gocui.Gui) error {
		out, err := (**cuiGUI).View(view)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(out, str)
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
	if v, err := g.SetView("v1", 0, 0, maxX-21, maxY-4); err != nil {
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

	if v, err := g.SetView("v2", 0, maxY-3, maxX-21, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Command (press ENTER to send)"
		v.Wrap = true
		v.Autoscroll = true
		v.Editable = true
	}
	if v, err := g.SetView("v3", maxX-20, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v3 (editable)"
		v.Editable = true
		v.Wrap = true
	}
	return nil
}

/*
 * Quit CUI
 */

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
