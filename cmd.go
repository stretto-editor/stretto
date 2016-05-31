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
		return quit(g, v)
	case "qs", "sq":
		return saveAndQuit(g, cmd)
	case "c!":
		vMain, _ := g.View("main")
		closeView(vMain)
	case "sc":
		saveAndClose(g, cmd)
	case "o", "open":
		openCmd(g, cmd)
	case "saveas", "sa":
		saveAsCmd(g, cmd)
	case "replaceall", "repall":
		replaceAll(g, cmd)
	case "setwrap":
		setWrapCmd(g, cmd)
	//TODO: go to the line specified
	default:
		displayError(g, ErrUnknownCommand)
	}
	clearView(v)
	return nil
}

func saveAndQuit(g *gocui.Gui, cmd []string) error {
	if currentFile == "" && len(cmd) == 1 {
		displayError(g, ErrMissingFilename)
		return nil
	}
	vMain, _ := g.View("main")
	filename := currentFile
	if filename == "" {
		filename = cmd[1]
	}
	createFile(filename)
	if err := saveMain(vMain, filename); err != nil {
		return err
	}
	return quit(g, vMain)
}

func replaceAll(g *gocui.Gui, cmd []string) {
	if len(cmd) == 2 {
		vMain, _ := g.View("main")
		for found, x, y := searchForward(vMain, cmd[1], 0, 0); found; found, x, y = searchForward(vMain, cmd[1], x, y) {
			replaceAt(vMain, x, y, cmd[1], cmd[2])
		}
	} else if len(cmd) == 1 {
		displayError(g, ErrMissingPattern)
	} else {
		displayError(g, ErrUnexpectedArgument)
	}
}

func saveAndClose(g *gocui.Gui, cmd []string) {
	if currentFile != "" || len(cmd) > 1 {
		vMain, _ := g.View("main")
		if currentFile == "" {
			createFile(cmd[1])
			saveMain(vMain, cmd[1])
		} else {
			saveMain(vMain, currentFile)
		}
		closeView(vMain)
	} else {
		displayError(g, ErrMissingFilename)
	}
}

func openCmd(g *gocui.Gui, cmd []string) {
	if len(cmd) == 2 {
		openAndDisplayFile(g, cmd[1])
	} else if len(cmd) == 1 {
		displayError(g, ErrMissingFilename)
	} else {
		displayError(g, ErrUnexpectedArgument)
	}
}

func saveAsCmd(g *gocui.Gui, cmd []string) {
	if len(cmd) == 2 {
		saveAs(g, cmd[1])
	} else if len(cmd) == 1 {
		displayError(g, ErrMissingFilename)
	} else {
		displayError(g, ErrUnexpectedArgument)
	}
}

func setWrapCmd(g *gocui.Gui, cmd []string) {
	if len(cmd) == 2 {
		vMain, _ := g.View("main")
		if cmd[1] == "true" {
			vMain.Wrap = true
		} else if cmd[1] == "false" {
			vMain.Wrap = false
		}
	} else if len(cmd) == 1 {
		displayError(g, ErrWrapArgument)
	} else {
		displayError(g, ErrUnexpectedArgument)
	}
}
