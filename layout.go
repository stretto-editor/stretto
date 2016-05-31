package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stretto-editor/gocui"
)

var requiredViewsInfo map[string]*struct {
	x, y, w, h int
	t          string
	e          bool
	f          string
	hi         bool
	wr         bool
}

func initRequiredViewsInfo(maxX, maxY int) {

	infoHeight := 2

	requiredViewsInfo = map[string]*struct {
		x, y, w, h int
		t          string // Title
		e          bool   // Editable
		f          string // Footer
		hi         bool   // Hidden
		wr         bool   // Wrap
	}{
		"main": {t: "undefined",
			e: true},
		"cmdline": {t: "Commandline",
			e: true},
		"inputline": {t: "Inputline for interactive actions",
			e: true},
		"infoline": {e: true,
			f: "INFO"},
		"error": {t: "Error :",
			e:  true,
			hi: true,
			wr: true},
	}

	// default geometries
	m, _ := requiredViewsInfo["main"]
	m.w = maxX + 1
	m.h = maxY - 1 - infoHeight
	m.x = -1
	m.y = 0

	c, _ := requiredViewsInfo["cmdline"]
	c.w = 30
	c.h = 2
	c.x = (maxX - c.w) / 2
	c.y = maxY - c.h - 10

	inp, _ := requiredViewsInfo["inputline"]
	inp.w = maxX * 80 / 100
	inp.h = 2
	inp.x = (maxX - inp.w) / 2
	inp.y = maxY - inp.h - 5

	inf, _ := requiredViewsInfo["infoline"]
	inf.w = maxX - 1
	inf.h = infoHeight
	inf.x = (maxX - inf.w) / 2
	inf.y = maxY - inf.h - 1

	e, _ := requiredViewsInfo["error"]
	e.w = maxX - 1
	e.h = 3
	e.x = (maxX - e.w) / 2
	e.y = maxY - inf.h - e.h - 1
}

func defaultLayout(g *gocui.Gui) error {
	var v *gocui.View
	var err error

	maxX, maxY := g.Size()

	initRequiredViewsInfo(maxX, maxY)

	for vname, settings := range requiredViewsInfo {
		v, err = g.SetView(vname, settings.x, settings.y, settings.x+settings.w, settings.y+settings.h)
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = settings.e
		v.Title = settings.t
		v.Footer = settings.f
		v.Hidden = settings.hi
		v.Wrap = settings.wr
	}

	// check if there is a second argument
	if len(os.Args) >= 2 {
		v, _ := g.View("main")
		v.Title = os.Args[1]
		if err := openFile(v, os.Args[1]); err != nil {
			return err
		}
		currentFile = os.Args[1]
	}

	info, _ := g.View("infoline")
	fmt.Fprintf(info, "Currently in edit mode \n"+"Cursor Position : 0,0")

	// main on top
	g.SetViewOnTop("main")
	g.SetCurrentView("main")

	return nil
}

func displayErrorView(g *gocui.Gui) {
	v, _ := g.View("error")
	if v.Hidden == true {
		m, _ := requiredViewsInfo["main"]
		e, _ := requiredViewsInfo["error"]
		m.h -= e.h
	}
	v.Hidden = false
	g.SetViewOnTop("error")
}

func hideErrorView(g *gocui.Gui) {
	v, _ := g.View("error")
	if v.Hidden == false {
		m, _ := requiredViewsInfo["main"]
		e, _ := requiredViewsInfo["error"]
		m.h += e.h
	}
	v.Hidden = true
	g.SetViewOnTop("main")
}

func layout(g *gocui.Gui) error {

	for vname, settings := range requiredViewsInfo {
		if _, err := g.SetView(vname, settings.x, settings.y, settings.x+settings.w, settings.y+settings.h); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
	}
	return nil
}

func openFile(v *gocui.View, name string) error {

	// inexisting view
	if v == nil {
		return gocui.ErrUnknownView
	}

	// get content of file
	f, err := ioutil.ReadFile(name)
	// inexisting file
	if err != nil {
		return err
	}

	v.Title = name
	v.Clear()
	fmt.Fprintf(v, "%s", f)
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	return nil
}
