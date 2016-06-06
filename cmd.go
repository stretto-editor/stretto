package main

import (
	"errors"

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

func initCommands() {
	commands = make(map[string]*Command)
	commands["quit"] = &Command{"quit", quitCmd, 0, 0, nil, nil}
	commands["q!"] = commands["quit"]
	commands["sq"] = &Command{"sq", saveAndQuit, 0, 1, ErrMissingFilename, GetAutocompleteFile}
	commands["qs"] = commands["sq"]
	commands["saveas"] = &Command{"saveas", saveAsCmd, 1, 1, ErrMissingFilename, GetAutocompleteFile}
	commands["sa"] = commands["saveas"]
	commands["setwrap"] = &Command{"setwrap", setWrapCmd, 1, 1, ErrWrapArgument, GetAutocompleteBoolean}
	commands["open"] = &Command{"open", openCmd, 1, 1, ErrMissingFilename, GetAutocompleteFile}
	commands["o"] = commands["open"]
	commands["close"] = &Command{"close", closeCmd, 0, 0, nil, nil}
	commands["c!"] = commands["close"]
	commands["replaceall"] = &Command{"replaceall", replaceAllCmd, 2, 2, ErrMissingPattern, nil}
	commands["repall"] = commands["replaceall"]
	//TODO: go to the line specified
}

func quitCmd(g *gocui.Gui, cmd []string) error {
	return gocui.ErrQuit
}

func closeCmd(g *gocui.Gui, cmd []string) error {
	vMain, _ := g.View("main")
	closeView(vMain)
	return nil
}

func saveAndQuit(g *gocui.Gui, cmd []string) error {
	vMain, _ := g.View("main")
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
	replaceAll(g, cmd[1], cmd[2])
	return nil
}

func saveAndClose(g *gocui.Gui, cmd []string) error {
	vMain, _ := g.View("main")
	if vMain.Title != "" || len(cmd) > 1 {
		vMain, _ := g.View("main")
		if vMain.Title == "" {
			createFile(cmd[1])
			vMain.Title = cmd[1]
		}
		saveMain(vMain, vMain.Title)
		closeView(vMain)
		return nil
	}
	return ErrMissingFilename
}

func openCmd(g *gocui.Gui, cmd []string) error {
	openAndDisplayFile(g, cmd[1])
	return nil
}

func saveAsCmd(g *gocui.Gui, cmd []string) error {
	saveAs(g, cmd[1])
	return nil
}

func setWrapCmd(g *gocui.Gui, cmd []string) error {
	vMain, _ := g.View("main")
	if cmd[1] == "true" {
		vMain.Wrap = true
	} else if cmd[1] == "false" {
		vMain.Wrap = false
	}
	return nil
}
