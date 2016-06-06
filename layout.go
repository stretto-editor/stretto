package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stretto-editor/gocui"
)

var requiredViewsInfo map[string]*struct {
	x, y, w, h int
	c          string // father container
	t          string // Title
	e          bool   // Editable
	f          string // Footer
	hi         bool   // Hidden
	wr         bool   // Wrap
}

func initRequiredViewsInfo(g *gocui.Gui) {

	requiredViewsInfo = map[string]*struct {
		x, y, w, h int
		c          string // father container
		t          string // Title
		e          bool   // Editable
		f          string // Footer
		hi         bool   // Hidden
		wr         bool   // Wrap
	}{
		"file": {t: "",
			c: "main",
			e: true},
		"cmdline": {t: "Commandline",
			c: "editable",
			e: true},
		"inputline": {t: "Inputline for interactive actions",
			c:  "editable",
			e:  true,
			hi: true},
		"infoline": {e: true,
			f: "INFO"},
		"error": {t: "Error :",
			e:  true,
			hi: true,
			wr: true},
	}

	setDefaultGeometry(g.Size())
}

func initTreeView(g *gocui.Gui) {
	g.SetViewNode("editable", "", 0, 0, 10, 10)
	g.SetViewNode("main", "editable", 0, 0, 10, 10)
}

func setDefaultGeometry(maxX, maxY int) {
	infoHeight := 2
	// default geometry
	f, _ := requiredViewsInfo["file"]
	f.w = maxX + 1
	f.h = maxY - 1 - infoHeight
	f.x = -1
	f.y = 0

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

func updateAllLayout(g *gocui.Gui) {
	setDefaultGeometry(g.Size())

	if v, _ := g.View("error"); !v.Hidden {
		m, _ := requiredViewsInfo["file"]
		i, _ := requiredViewsInfo["inputline"]
		e, _ := requiredViewsInfo["error"]
		m.h -= e.h
		i.y -= e.h
		g.SetViewOnTop("error")
	}
	if v, _ := g.View("inputline"); !v.Hidden {
		g.SetViewOnTop("inputline")
	}
}

func defaultLayout(g *gocui.Gui) error {
	var v *gocui.View
	var err error

	initTreeView(g)
	initRequiredViewsInfo(g)

	for vname, settings := range requiredViewsInfo {
		v, err = g.SetView(vname, settings.c, settings.x, settings.y, settings.x+settings.w, settings.y+settings.h)
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
		v, _ := g.View("file")
		// v := g.Workingview()
		v.Title = os.Args[1]
		if err := openFile(v, os.Args[1]); err != nil {
			return err
		}
		v.Title = os.Args[1]
	}

	info, _ := g.View("infoline")
	maxX, _ := info.Size()
	mode := fmt.Sprintf("edit mode")
	pos := fmt.Sprintf("0:0")
	fmt.Fprintf(info, "%s", mode)
	fmt.Fprintf(info, "%[2]*.[2]*[1]s", pos, maxX-len(mode))

	// main on top
	g.SetViewOnTop("file")
	g.SetCurrentView("file")
	g.SetWorkingView("file")

	return nil
}

func displayInputLine(g *gocui.Gui) {
	v, _ := g.View("inputline")
	v.Hidden = false
	g.SetViewOnTop("inputline")
}

func hideInputLine(g *gocui.Gui) {
	v, _ := g.View("inputline")
	v.Hidden = true
	// g.SetViewOnTop("main")
	g.SetViewOnTop(g.Workingview().Name())
}

func displayErrorView(g *gocui.Gui) {
	v, _ := g.View("error")
	v.Hidden = false
	g.SetViewOnTop("error")
}

func hideErrorView(g *gocui.Gui) {
	v, _ := g.View("error")
	v.Hidden = true
	// g.SetViewOnTop("main")
	g.SetViewOnTop(g.Workingview().Name())
}

func layout(g *gocui.Gui) error {
	updateAllLayout(g)

	for vname, settings := range requiredViewsInfo {
		if _, err := g.SetView(vname, settings.c, settings.x, settings.y, settings.x+settings.w, settings.y+settings.h); err != nil {
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

func displayDoc(g *gocui.Gui) {
	v, err := g.View("cmdinfo")
	if err != nil {
		v.Hidden = false
		g.SetViewOnTop("cmdinfo")
	}
}

func createDocView(g *gocui.Gui) (*gocui.View, error) {
	maxX, maxY := g.Size()
	wcmd, hcmd := maxX*70/100, maxY*70/100
	var xcmd, ycmd int = (maxX - wcmd) / 2, maxY/2 - hcmd/2
	var v *gocui.View
	var err error
	if v, err = g.SetView("cmdinfo", "", xcmd, ycmd, xcmd+wcmd, ycmd+hcmd); err != nil && err == gocui.ErrUnknownView {
		v.Editable = false
		v.Wrap = true
		v.Title = "Commands Summary"
		return v, nil
	}
	return nil, err

}
