package main

import "github.com/stretto-editor/gocui"

const fileMode = "file"
const editMode = "edit"
const cmdMode = "cmd"

func initModes(g *gocui.Gui) {
	openCmdMode := func(g *gocui.Gui) error {
		g.SetWorkingView(g.CurrentView().Name())
		if err := g.SetCurrentView("cmdline"); err != nil {
			return err
		}
		g.SetViewOnTop("cmdline")
		v, _ := g.View("cmdline")
		v.Clear()
		v.SetOrigin(0, 0)
		v.SetCursor(0, 0)
		return nil
	}
	closeCmdMode := func(g *gocui.Gui) error {
		g.SetCurrentView(g.Workingview().Name())
		g.SetViewOnTop(g.Workingview().Name())
		return nil
	}
	openFileMode := func(g *gocui.Gui) error {
		v := g.Workingview()
		v.SetEditable(false)
		return nil
	}
	closeFileMode := func(g *gocui.Gui) error {
		v := g.Workingview()
		v.SetEditable(true)
		return nil
	}
	openEditMode := func(g *gocui.Gui) error {
		return nil
	}
	closeEditMode := func(g *gocui.Gui) error {
		return nil
	}
	g.AddMode(cmdMode, openCmdMode, closeCmdMode)
	g.AddMode(fileMode, openFileMode, closeFileMode)
	g.AddMode(editMode, openEditMode, closeEditMode)

	g.SetCurrentMode(editMode)
}
