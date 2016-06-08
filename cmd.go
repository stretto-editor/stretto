package main

import (
	"errors"
	"strconv"

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
	// ErrMissingLine raised when the number of the line to go to is missing
	ErrMissingLine = errors.New("the number of the line to go to is missing")
	// ErrGoToInWrapMode raised when the user try to use goto when wrap mode is active
	ErrGoToInWrapMode = errors.New("goto not available when wrap is active")
	// ErrNumberExpected raised when a number is expected in argument and other type was found
	ErrNumberExpected = errors.New("illegal parameter, number expected")
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
	commands["sc"] = &Command{"sc", saveAndClose, 0, 1, nil, GetAutocompleteFile}
	commands["replaceall"] = &Command{"replaceall", replaceAllCmd, 2, 2, ErrMissingPattern, nil}
	commands["repall"] = commands["replaceall"]
	commands["goto"] = &Command{"goto", goToCmd, 1, 2, ErrMissingLine, nil}
}

func quitCmd(g *gocui.Gui, cmd []string) error {
	return gocui.ErrQuit
}

func closeCmd(g *gocui.Gui, cmd []string) error {
	vMain := g.Workingview()
	closeView(g, vMain)
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
	replaceAll(g, cmd[1], cmd[2])
	return nil
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
	openAndDisplayFile(g, cmd[1])
	return nil
}

func saveAsCmd(g *gocui.Gui, cmd []string) error {
	saveAs(g, cmd[1])
	return nil
}

func setWrapCmd(g *gocui.Gui, cmd []string) error {
	vMain := g.Workingview()
	if cmd[1] == "true" {
		vMain.Wrap = true
	} else if cmd[1] == "false" {
		vMain.Wrap = false
	}
	return nil
}

func goToCmd(g *gocui.Gui, cmd []string) error {
	vMain := g.Workingview()
	if vMain.Wrap {
		return ErrGoToInWrapMode
	}
	var x, y int
	var err error
	if y, err = strconv.Atoi(cmd[1]); err != nil {
		return ErrNumberExpected
	}
	if len(cmd) > 2 {
		if x, err = strconv.Atoi(cmd[2]); err != nil {
			return ErrNumberExpected
		}
	}
	vMain.SetOrigin(0, 0)
	vMain.SetCursor(0, 0)
	_, cy := vMain.Cursor()
	cyPred := -1
	_, oy := vMain.Origin()
	oyPred := -1
	for cy+oy != cyPred+oyPred && oy+cy != y {
		_, cyPred = vMain.Cursor()
		_, oyPred = vMain.Origin()
		moveDown(g, vMain)
		_, cy = vMain.Cursor()
		_, oy = vMain.Origin()
	}
	vMain.MoveCursor(x, 0, false)
	vMain.Actions.Cut()
	switchModeHandlerFactory(editMode)(g, vMain)
	return nil
}
