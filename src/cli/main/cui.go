package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"fmt"
)

var cuiGUI **gocui.Gui

func attachCUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	cuiGUI = &g
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

var (
	viewArr = []string{"v1", "v2", "v3"}
	active  = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func writeToView(str string, view string) {
	out, err := (**cuiGUI).View(view)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(out, str)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

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
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
