package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stretto-editor/gocui"
)

func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView("main", 0, 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Title = "undefined"
		// v.Wrap = true

		// check if there is a second argument
		if len(os.Args) >= 2 {
			v.Title = os.Args[1]
			if err := openFile(v, os.Args[1]); err != nil {
				return err
			}
			currentFile = os.Args[1]
		}
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}

	wcmd, hcmd := 30, 2
	var xcmd, ycmd int = (maxX - wcmd) / 2, maxY - hcmd - 10
	if v, err := g.SetView("cmdline", xcmd, ycmd, xcmd+wcmd, ycmd+hcmd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "Commandline"
		g.SetViewOnTop("main")
	}

	winput, hinput := maxX*80/100, 2
	var xinput, yinput int = (maxX - winput) / 2, maxY/2 - hinput/2
	if v, err := g.SetView("inputline", xinput, yinput, xinput+winput, yinput+hinput); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "Inputline for interactive actions"
		g.SetViewOnTop("main")
	}

	winfo, hinfo := maxX-1, 4
	var xinfo, yinfo int = (maxX - winfo) / 2, maxY - hinfo - 1
	if v, err := g.SetView("infoline", xinfo, yinfo, xinfo+winfo, yinfo+hinfo); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Footer = "INFO"
		info, _ := g.View("infoline")
		fmt.Fprintf(info, "Currently in edit mode \n"+"Cursor Position : 0,0")
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

	v.Clear()
	fmt.Fprintf(v, "%s", f)
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	return nil
}
