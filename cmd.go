package main

import (
	"errors"
	"strings"

	"github.com/stretto-editor/gocui"
)

var (
	// ErrMissingFilename raised when you want to save a file but there is no filename specified
	ErrMissingFilename = errors.New("missing filename as argument")
	// ErrUnknownCommand The user has entered an unknown command
	ErrUnknownCommand = errors.New("unknown command")
	// ErrMissingPattern raised a word is missing for the search and replace
	ErrMissingPattern = errors.New("missing search or replace word")
	// ErrUnexpectedArgument argument found when it wasn't espected
	ErrUnexpectedArgument = errors.New("unexpected argument")
	// ErrWrapArgument raised when the true or false argument is missing in the wrap command
	ErrWrapArgument = errors.New("expected true or false argument")
)

func validateCmd(g *gocui.Gui, v *gocui.View) error {
	var err error
	if v.Name() != "cmdline" {
		panic("Cmdline is not the current view")
	}
	cmdBuff := v.Buffer()
	if cmdBuff == "" {
		return nil
	}
	cmdBuff = cmdBuff[:len(cmdBuff)-1]
	cmd := strings.Fields(cmdBuff)
	switch cmd[0] {
	case "quit", "q!":
		err = quit(g, v)
	case "qs", "sq":
		err = saveAndQuit(g, cmd)
	case "c!":
		// vMain, _ := g.View("main")
		vMain := g.Workingview()
		closeView(g, vMain)
	case "sc":
		err = saveAndClose(g, cmd)
	case "o", "open":
		err = openCmd(g, cmd)
	case "saveas", "sa":
		err = saveAsCmd(g, cmd)
	case "replaceall", "repall":
		err = replaceAllCmd(g, cmd)
	case "setwrap":
		err = setWrapCmd(g, cmd)
	//TODO: go to the line specified
	default:
		err = ErrUnknownCommand
	}
	clearView(v)
	if err == gocui.ErrQuit {
		return err
	}
	if err != nil {
		displayError(g, err)
	}
	return nil
}

func saveAndQuit(g *gocui.Gui, cmd []string) error {
	// vMain, _ := g.View("main")
	vMain := g.Workingview()
	if vMain.Title == "" && len(cmd) == 1 {
		return ErrMissingFilename
	}
	if vMain.Title == "" {
		vMain.Title = cmd[1]
	}
	createFile(vMain.Title)
	if err := saveMain(vMain, vMain.Title); err != nil {
		return err
	}
	return quit(g, vMain)
}

func replaceAllCmd(g *gocui.Gui, cmd []string) error {
	if len(cmd) == 3 {
		replaceAll(g, cmd[1], cmd[2])
		return nil
	}
	if len(cmd) == 1 {
		return ErrMissingPattern
	}
	return ErrUnexpectedArgument
}

func saveAndClose(g *gocui.Gui, cmd []string) error {
	// vMain, _ := g.View("main")
	vMain := g.Workingview()
	if vMain.Title != "" || len(cmd) > 1 {
		if vMain.Title == "" {
			createFile(cmd[1])
			vMain.Title = cmd[1]
		}
		saveMain(vMain, vMain.Title)
		closeView(g, vMain)
		return nil
	}
	return ErrMissingFilename
}

func openCmd(g *gocui.Gui, cmd []string) error {
	if len(cmd) == 2 {
		openAndDisplayFile(g, cmd[1])
		return nil
	}
	if len(cmd) == 1 {
		return ErrMissingFilename
	}
	return ErrUnexpectedArgument
}

func saveAsCmd(g *gocui.Gui, cmd []string) error {
	if len(cmd) == 2 {
		saveAs(g, cmd[1])
		return nil
	}
	if len(cmd) == 1 {
		return ErrMissingFilename
	}
	return ErrUnexpectedArgument
}

func setWrapCmd(g *gocui.Gui, cmd []string) error {
	if len(cmd) == 2 {
		//vMain, _ := g.View("main")
		vMain := g.Workingview()
		if cmd[1] == "true" {
			vMain.Wrap = true
		} else if cmd[1] == "false" {
			vMain.Wrap = false
		}
		return nil
	}
	if len(cmd) == 1 {
		return ErrWrapArgument
	}
	return ErrUnexpectedArgument
}
