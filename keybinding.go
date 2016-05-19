package main

import (
	"github.com/jroimartin/gocui"
)

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlT, gocui.ModNone, currTopViewHandler("cmdline")); err != nil {
		return err
	}

	if err := g.SetKeybinding("cmdline", gocui.KeyCtrlT, gocui.ModNone, currTopViewHandler("main")); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func currTopViewHandler(name string) gocui.KeybindingHandler {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := g.SetCurrentView(name); err != nil {
			return err
		}
		_, err := g.SetViewOnTop(name)
		return err
	}
}
