package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView("main", -1, -1, maxX, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprint(v, "texte")

		v.Highlight = true
		v.Editable = true
		v.Wrap = true
	}

	if _, err := g.SetView("cmdline", -1, maxY-5, maxX, maxY); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}

	return nil
}
