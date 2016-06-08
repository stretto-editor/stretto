package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stretto-editor/gocui"
)

// updateView is a fonction type used for updating view
type updateView func(int, int)

type viewInfo struct {
	x, y, w, h int
	c          string     // father container
	t          string     // Title
	e          bool       // Editable
	f          string     // Footer
	hi         bool       // Hidden
	wr         bool       // Wrap
	up         updateView // update
}

var requiredViewsInfo map[string]*viewInfo

var infoHeight = 2

func initRequiredViewsInfo(g *gocui.Gui) {
	updateCmdlineGeom := func(maxX, maxY int) {
		c, _ := requiredViewsInfo["cmdline"]
		c.w = 30
		c.h = 2
		c.x = (maxX - c.w) / 2
		c.y = maxY - c.h - 10
	}
	updateInputlineGeom := func(maxX, maxY int) {
		inp, _ := requiredViewsInfo["inputline"]
		inp.w = maxX * 80 / 100
		inp.h = 2
		inp.x = (maxX - inp.w) / 2
		inp.y = maxY - inp.h - 5
	}
	updateInfolineGeom := func(maxX, maxY int) {
		inf, _ := requiredViewsInfo["infoline"]
		inf.w = maxX - 1
		inf.h = infoHeight
		inf.x = (maxX - inf.w) / 2
		inf.y = maxY - inf.h - 1
	}
	updateErrorViewGeom := func(maxX, maxY int) {
		e, _ := requiredViewsInfo["error"]
		e.w = maxX - 1
		e.h = 3
		e.x = (maxX - e.w) / 2
		e.y = maxY - infoHeight - e.h - 1
	}
	updateHistoricView := func(maxX, maxY int) {
		h, _ := requiredViewsInfo["historic"]
		h.w = 20
		h.h = maxY - 15
		h.x = maxX - h.w - 5
		h.y = 5
	}

	requiredViewsInfo = map[string]*viewInfo{
		"cmdline": {
			t:  "Commandline",
			c:  "editable",
			e:  true,
			up: updateCmdlineGeom,
		},
		"inputline": {
			t:  "Inputline for interactive actions",
			c:  "editable",
			e:  true,
			hi: true,
			up: updateInputlineGeom,
		},
		"infoline": {
			e:  true,
			f:  "INFO",
			up: updateInfolineGeom,
		},
		"error": {
			t:  "Error :",
			e:  true,
			hi: true,
			wr: true,
			up: updateErrorViewGeom,
		},
		"historic": {
			t:  "History :",
			e:  true,
			hi: true,
			up: updateHistoricView,
		},
	}

	updateGeometry(g.Size())
}

func initTreeView(g *gocui.Gui) {
	g.SetViewNode("editable", "", 0, 0, 10, 10)
	g.SetViewNode("main", "editable", 0, 0, 10, 10)
}

func updateGeometry(maxX, maxY int) {
	for _, vi := range requiredViewsInfo {
		vi.up(maxX, maxY)
	}
}

func updateAllLayout(g *gocui.Gui) {
	updateGeometry(g.Size())

	if v, _ := g.View("error"); !v.Hidden {
		m, _ := requiredViewsInfo[g.Workingview().Name()]
		i, _ := requiredViewsInfo["inputline"]
		e, _ := requiredViewsInfo["error"]
		m.h -= e.h
		i.y -= e.h
		g.SetViewOnTop("error")
	}
	if v, _ := g.View("inputline"); !v.Hidden {
		g.SetViewOnTop("inputline")
	}
	if v, _ := g.View("historic"); !v.Hidden {
		g.SetViewOnTop("historic")
	}
}

func initView(g *gocui.Gui, vname string) (*gocui.View, error) {
	settings := requiredViewsInfo[vname]
	v, err := g.SetView(vname, settings.c, settings.x, settings.y, settings.x+settings.w, settings.y+settings.h)
	if err != gocui.ErrUnknownView && err != nil {
		return nil, err
	}
	v.Editable = settings.e
	v.Title = settings.t
	v.Footer = settings.f
	v.Hidden = settings.hi
	v.Wrap = settings.wr
	return v, nil
}

func defaultLayout(g *gocui.Gui) error {
	initTreeView(g)
	initRequiredViewsInfo(g)

	for vname := range requiredViewsInfo {
		initView(g, vname)
	}

	// check if there is a second argument
	if len(os.Args) >= 2 {
		openAndDisplayFile(g, os.Args[1])
	} else {
		newFileView(g, "file")
	}

	info, _ := g.View("infoline")
	maxX, _ := info.Size()
	mode := fmt.Sprintf("edit mode")
	pos := fmt.Sprintf("0:0")
	fmt.Fprintf(info, "%s", mode)
	fmt.Fprintf(info, "%[2]*.[2]*[1]s", pos, maxX-len(mode))

	// main on top
	c, _ := g.ViewNode("main")
	v := c.LastView()
	if v != nil {
		g.SetViewOnTop(v.Name())
		g.SetCurrentView(v.Name())
		g.SetWorkingView(v.Name())
	}
	return nil
}

func displayHistoric(g *gocui.Gui) {
	v, _ := g.View("historic")
	v.Hidden = false
	g.SetViewOnTop("historic")
}

func hideHistoric(g *gocui.Gui) {
	v, _ := g.View("historic")
	v.Hidden = true
}

func displayInputLine(g *gocui.Gui) {
	v, _ := g.View("inputline")
	v.Hidden = false
	g.SetViewOnTop("inputline")
}

func hideInputLine(g *gocui.Gui) {
	v, _ := g.View("inputline")
	v.Hidden = true
}

func displayErrorView(g *gocui.Gui) {
	v, _ := g.View("error")
	v.Hidden = false
	g.SetViewOnTop("error")
}

func hideErrorView(g *gocui.Gui) {
	v, _ := g.View("error")
	v.Hidden = true
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

func newFileView(g *gocui.Gui, filename string) (*gocui.View, error) {
	v, err := g.SetView(filename, "main", 0, 0, 100, 300)
	updateFileGeom := func(maxX, maxY int) {
		f, _ := requiredViewsInfo[filename]
		f.w = maxX + 1
		f.h = maxY - 1 - infoHeight
		f.x = -1
		f.y = 0
	}
	requiredViewsInfo[filename] = &viewInfo{
		t:  filename,
		c:  "main",
		e:  true,
		up: updateFileGeom,
	}
	updateFileGeom(g.Size())
	initView(g, filename)
	return v, err
}

func removeFileView(viewName string) {
	delete(requiredViewsInfo, viewName)
}

func createView(g *gocui.Gui, name string) (*gocui.View, error) {
	v, err := g.SetView(name, "", 0, 0, 100, 300)
	updateGeom := func(maxX, maxY int) {
		f, _ := requiredViewsInfo[name]
		f.w = 30
		f.h = maxY * 60 / 100
		f.x = 5
		f.y = 3
	}
	requiredViewsInfo[name] = &viewInfo{
		t:  name,
		c:  "",
		e:  false,
		wr: true,
		up: updateGeom,
	}
	initView(g, name)
	return v, err
}
