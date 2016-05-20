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
		v.Wrap = true

		// check if there is a second argument
		if len(os.Args) >= 2 {
			openFile(v, os.Args[1])
			currentFile = os.Args[1]
		}
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}

	w_cmdl, h_cmdl := 30, 2
	var x_cmdl, y_cmdl int = (maxX - w_cmdl) / 2, maxY - h_cmdl - 5
	if v, err := g.SetView("cmdline", x_cmdl, y_cmdl, x_cmdl+w_cmdl, y_cmdl+h_cmdl); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		fmt.Fprint(v, "cmdline")
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
		panic(err)
	}

	fmt.Fprintf(v, "%s", f)
	return nil
}
