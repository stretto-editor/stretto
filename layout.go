package main

import (
	"fmt"

	"io/ioutil"
	"os"

	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView("main", -1, -1, maxX, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		// check if there is a second argument
		if len(os.Args) >= 2 {
			openFile(v, os.Args[1])
		}
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}

	if _, err := g.SetView("cmdline", -1, maxY-5, maxX, maxY); err != nil &&
		err != gocui.ErrUnknownView {
		return err
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
