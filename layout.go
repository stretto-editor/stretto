package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stretto-editor/gocui"
)

func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView("main", -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		// v.Wrap = true

		// check if there is a second argument
		if len(os.Args) >= 2 {
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
	var xcmd, ycmd int = (maxX - wcmd) / 2, maxY - hcmd - 5
	if v, err := g.SetView("cmdline", xcmd, ycmd, xcmd+wcmd, ycmd+hcmd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		fmt.Fprint(v, "cmdline")
		g.SetViewOnTop("main")
	}

	winput, hinput := maxX*80/100, 2
	var xinput, yinput int = (maxX - winput) / 2, maxY/2 - hinput/2
	if v, err := g.SetView("inputline", xinput, yinput, xinput+winput, yinput+hinput); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		fmt.Fprint(v, "input")
		g.SetViewOnTop("main")
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

	fmt.Fprintf(v, "%s", f)
	return nil
}
