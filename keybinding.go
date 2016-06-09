package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/stretto-editor/gocui"
)

// demonInput defines the prototype for functions that
// should be called later in validateInput
// A demonInput returns the next demonInput to be called,
// or nil if there is none
type demonInput func(g *gocui.Gui, input string) (demonInput, error)

// current is the next demonInput to be called
// In an handler (that use the inputline), it should be set.
// It is set back to nil once all the call have been made
// (in validateInput handler)
var currentDemonInput demonInput

// ErrViewCreated is use for interactive actions whose create a view
var ErrViewCreated = errors.New("A view was created")

func initKeybindings(g *gocui.Gui) error {

	var keyBindings = []struct {
		m string
		v string
		k interface{}
		h gocui.KeybindingHandler
	}{

		// ---------------------- COMMON COMMANDS ------------------------- //

		{m: fileMode, v: "main", k: gocui.KeyCtrlT, h: switchModeHandlerFactory(cmdMode)},
		{m: fileMode, v: "main", k: gocui.KeyF2, h: switchModeHandlerFactory(editMode)},
		{m: fileMode, v: "main", k: gocui.KeyCtrlQ, h: quitHandler},

		{m: editMode, v: "main", k: gocui.KeyCtrlT, h: switchModeHandlerFactory(cmdMode)},
		{m: editMode, v: "main", k: gocui.KeyF2, h: switchModeHandlerFactory(fileMode)},
		{m: editMode, v: "main", k: gocui.KeyCtrlQ, h: quitHandler},

		{m: cmdMode, v: "cmdline", k: gocui.KeyCtrlT, h: switchModeHandlerFactory(editMode)},

		// ---------------------- MAIN SECTION ---------------------------- //

		// ---------------------- NAVIGATION ------------------------------ //

		{m: fileMode, v: "main", k: gocui.KeyArrowLeft, h: moveLeft},
		{m: fileMode, v: "main", k: gocui.KeyArrowRight, h: moveRight},
		{m: fileMode, v: "main", k: gocui.KeyArrowUp, h: moveUp},
		{m: fileMode, v: "main", k: gocui.KeyArrowDown, h: moveDown},
		{m: fileMode, v: "main", k: gocui.KeyHome, h: cursorHome},
		{m: fileMode, v: "main", k: gocui.KeyEnd, h: cursorEnd},
		{m: fileMode, v: "main", k: gocui.KeyPgup, h: goPgUp},
		{m: fileMode, v: "main", k: gocui.KeyPgdn, h: goPgDown},

		{m: editMode, v: "main", k: gocui.KeyArrowLeft, h: moveLeft},
		{m: editMode, v: "main", k: gocui.KeyArrowRight, h: moveRight},
		{m: editMode, v: "main", k: gocui.KeyArrowUp, h: moveUp},
		{m: editMode, v: "main", k: gocui.KeyArrowDown, h: moveDown},
		{m: editMode, v: "main", k: gocui.KeyHome, h: cursorHome},
		{m: editMode, v: "main", k: gocui.KeyEnd, h: cursorEnd},
		{m: editMode, v: "main", k: gocui.KeyPgup, h: goPgUp},
		{m: editMode, v: "main", k: gocui.KeyPgdn, h: goPgDown},

		{m: editMode, v: "main", k: gocui.KeyF7, h: switchBufferForward},
		{m: fileMode, v: "main", k: gocui.KeyF7, h: switchBufferForward},
		{m: editMode, v: "main", k: gocui.KeyF8, h: switchBufferBackward},
		{m: fileMode, v: "main", k: gocui.KeyF8, h: switchBufferBackward},

		// ---------------------- USEFUL --- ------------------------------ //

		{m: fileMode, v: "main", k: 'o', h: openFileHandler},
		{m: fileMode, v: "main", k: 'w', h: closeFileHandler},
		{m: fileMode, v: "main", k: 's', h: saveHandler},
		{m: fileMode, v: "main", k: 'u', h: saveAsHandler},
		{m: fileMode, v: "main", k: 'f', h: searchHandler},
		{m: fileMode, v: "main", k: 'd', h: dirInfoHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlD, h: dirInfoHandler},

		{m: editMode, v: "main", k: gocui.KeyCtrlN, h: historicHandler},
		{m: editMode, v: "", k: gocui.KeyCtrlZ, h: undoHandler},
		{m: editMode, v: "", k: gocui.KeyCtrlY, h: redoHandler},

		{m: editMode, v: "main", k: gocui.KeyCtrlO, h: openFileHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlW, h: closeFileHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlS, h: saveHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlU, h: saveAsHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlF, h: searchHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlP, h: searchAndReplaceHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlC, h: copyHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlV, h: pasteHandler},
		{m: editMode, v: "main", k: gocui.KeyEnter, h: breaklineHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlJ, h: permutLinesUpHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlK, h: permutLinesDownHandler},

		{m: fileMode, v: "main", k: gocui.KeyF3, h: docHandler},
		{m: editMode, v: "main", k: gocui.KeyF3, h: docHandler},

		// ---------------------- INFO SECTION ---------------------------- //

		// ---------------------- NAVIGATION ------------------------------ //

		{m: fileMode, v: "tmp", k: gocui.KeyArrowUp, h: scrollUp},
		{m: fileMode, v: "tmp", k: gocui.KeyArrowDown, h: scrollDown},
		{m: fileMode, v: "tmp", k: gocui.KeyPgup, h: goPgUp},
		{m: fileMode, v: "tmp", k: gocui.KeyPgdn, h: goPgDown},
		{m: fileMode, v: "tmp", k: gocui.KeyEsc, h: quitTmpView},

		{m: editMode, v: "tmp", k: gocui.KeyArrowUp, h: scrollUp},
		{m: editMode, v: "tmp", k: gocui.KeyArrowDown, h: scrollDown},
		{m: editMode, v: "tmp", k: gocui.KeyPgup, h: goPgUp},
		{m: editMode, v: "tmp", k: gocui.KeyPgdn, h: goPgDown},
		{m: editMode, v: "tmp", k: gocui.KeyEsc, h: quitTmpView},

		// ---------------------- INPUT SECTION --------------------------- //

		// ---------------------- NAVIGATION ------------------------------ //

		{m: fileMode, v: "inputline", k: gocui.KeyHome, h: cursorHome},
		{m: fileMode, v: "inputline", k: gocui.KeyEnd, h: cursorEnd},
		{m: fileMode, v: "inputline", k: gocui.KeyArrowLeft, h: moveLeft},
		{m: fileMode, v: "inputline", k: gocui.KeyArrowRight, h: moveRight},

		{m: editMode, v: "inputline", k: gocui.KeyHome, h: cursorHome},
		{m: editMode, v: "inputline", k: gocui.KeyEnd, h: cursorEnd},
		{m: editMode, v: "inputline", k: gocui.KeyArrowLeft, h: moveLeft},
		{m: editMode, v: "inputline", k: gocui.KeyArrowRight, h: moveRight},

		// ---------------------- USEFUL --- ------------------------------ //

		{m: fileMode, v: "inputline", k: gocui.KeyEnter, h: validateInput},
		{m: fileMode, v: "inputline", k: gocui.KeyEsc, h: escapeInputHandler},

		{m: editMode, v: "inputline", k: gocui.KeyEnter, h: validateInput},
		{m: editMode, v: "inputline", k: gocui.KeyEsc, h: escapeInputHandler},

		{m: fileMode, v: "main", k: gocui.KeyEsc, h: escapeMainHandler},
		{m: editMode, v: "main", k: gocui.KeyEsc, h: escapeMainHandler},

		// ---------------------- CMD SECTION ---------------------------- //

		// ---------------------- NAVIGATION ------------------------------ //

		{m: cmdMode, v: "cmdline", k: gocui.KeyHome, h: cursorHome},
		{m: cmdMode, v: "cmdline", k: gocui.KeyEnd, h: cursorEnd},
		{m: cmdMode, v: "cmdline", k: gocui.KeyArrowLeft, h: moveLeft},
		{m: cmdMode, v: "cmdline", k: gocui.KeyArrowRight, h: moveRight},

		// CMDLINE
		{m: cmdMode, v: "cmdline", k: gocui.KeyEnter, h: validateCmd},
		{m: cmdMode, v: "cmdline", k: gocui.KeyTab, h: AutocompleteCmd},
	}

	for _, kb := range keyBindings {
		if err := g.SetKeybinding(kb.m, kb.v, kb.k, gocui.ModNone, kb.h); err != nil {
			return err
		}
	}
	return nil
}

func historicHandler(g *gocui.Gui, v *gocui.View) error {
	if v, _ := g.View("historic"); v.Hidden {
		displayHistoric(g)
	} else {
		hideHistoric(g)
	}
	return nil
}

func switchBufferForward(g *gocui.Gui, v *gocui.View) error {
	c, _ := g.ViewNode("main")
	newView := c.RoundRobinForward()
	if newView != nil {
		g.SetCurrentView(newView.Name())
		g.SetWorkingView(newView.Name())
	}
	return nil
}

func switchBufferBackward(g *gocui.Gui, v *gocui.View) error {
	c, _ := g.ViewNode("main")
	newView := c.RoundRobinBackward()
	if newView != nil {
		g.SetCurrentView(newView.Name())
		g.SetWorkingView(newView.Name())
	}
	return nil
}

func undoHandler(g *gocui.Gui, v *gocui.View) error {
	v.Actions.Undo()
	g.UpdateHistoric()
	return nil
}

func redoHandler(g *gocui.Gui, v *gocui.View) error {
	v.Actions.Redo()
	g.UpdateHistoric()
	return nil
}

/*
//TODO : deleted this (debug only)
func testCmdWrite(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	s := 'T'
	c := gocui.NewWriteCmd(v, cx, cy, s)
	v.Actions.Exec(c)
	g.UpdateHistoric()
	return nil
}
*/

func breaklineHandler(g *gocui.Gui, v *gocui.View) error {
	v.EditNewLine()
	updateInfos(g)
	return nil
}

func interactive(g *gocui.Gui, s string) {
	g.SetCurrentView("inputline")
	displayInputLine(g)
	g.CurrentView().Title = " " + s + " "
	g.CurrentView().MoveCursor(0, 0, false)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func saveHandler(g *gocui.Gui, v *gocui.View) error {
	// vMain, _ := g.View("main")
	vMain := g.Workingview()
	if vMain.Title == "" {
		currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
			createFile(input)
			vMain.Title = input
			if err := saveMain(vMain, vMain.Title); err != nil {
				return nil, err
			}
			return nil, nil
		}

		interactive(g, "Save")
		return nil
	}

	if err := saveMain(vMain, vMain.Title); err != nil {
		return err
	}
	return nil
}

func docHandler(g *gocui.Gui, v *gocui.View) error {
	if v, err := newTmpView(g, "cmdinfo"); err != gocui.ErrUnknownView {
		displayError(g, err)
	} else {
		openFile(v, "Commands.md")
		g.SetViewOnTop(v.Name())
		g.SetCurrentView(v.Name())
	}
	return nil
}

func quitTmpView(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView(g.CurrentView().Name())
	removeInfoView(g.CurrentView().Name())
	g.SetCurrentView(g.Workingview().Name())
	g.SetViewOnTop(g.Workingview().Name())
	return nil
}

func dirInfoHandler(g *gocui.Gui, v *gocui.View) error {
	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		return nil, showDirectory(g, input)
	}
	interactive(g, "Directory Content")
	return nil
}

func showDirectory(g *gocui.Gui, directory string) error {
	if directory[len(directory)-1] != '/' {
		s := []string{}
		s = append(s, directory)
		s = append(s, "/")
		directory = strings.Join(s, "")
	}
	if v, err := newTmpView(g, "Directory Info"); err != gocui.ErrUnknownView {
		displayError(g, err)
	} else {
		files, err := ioutil.ReadDir(directory)
		if err != nil {
			return fmt.Errorf("%s is not a valid directory", directory)
		}
		displayDirectoryContent(v, files)
		g.SetViewOnTop(v.Name())
		g.SetCurrentView(v.Name())
		return ErrViewCreated
	}
	return nil
}

func displayDirectoryContent(v *gocui.View, files []os.FileInfo) {
	for _, file := range files {
		if file.IsDir() {
			fmt.Fprintln(v, " "+file.Name()+"/")
		} else if !file.IsDir() {
			fmt.Fprintln(v, " "+file.Name())
		}
	}
}

func quitHandler(g *gocui.Gui, v *gocui.View) error {
	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		if input != "n" {
			// vMain, _ := g.View("main")
			vMain := g.Workingview()
			if vMain.Title == "" {
				interactive(g, "File name")
				return func(g *gocui.Gui, input string) (demonInput, error) {

					createFile(input)
					vMain.Title = input
					if err := saveMain(vMain, vMain.Title); err != nil {
						return nil, err
					}

					return nil, gocui.ErrQuit
				}, nil

			}
			if err := saveMain(vMain, vMain.Title); err != nil {
				return nil, err
			}
		}
		return nil, gocui.ErrQuit
	}

	interactive(g, "Save Modifications (y/n)")
	return nil
}

func closeFileHandler(g *gocui.Gui, v *gocui.View) error {
	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		// vMain, _ := g.View("main")
		vMain := g.Workingview()
		if input != "n" {
			if vMain.Title == "" {
				interactive(g, "File name")
				return func(g *gocui.Gui, input string) (demonInput, error) {
					createFile(input)
					vMain.Title = input
					if err := saveMain(vMain, vMain.Title); err != nil {
						return nil, err
					}
					closeView(g, vMain)
					return nil, nil
				}, nil
			}
			if err := saveMain(vMain, vMain.Title); err != nil {
				return nil, err
			}
		}
		closeView(g, vMain)
		return nil, nil
	}

	interactive(g, "Save Modifications (y/n)")
	return nil
}

func closeView(g *gocui.Gui, v *gocui.View) {
	//clearView(v)
	//v.Title = ""
	g.DeleteView(g.Workingview().Name())
	removeInfoView(g.Workingview().Name())
	c, _ := g.ViewNode("main")
	activeView := c.LastView()
	if activeView == nil {
		activeView, _ = newFileView(g, "file")
	}
	g.SetWorkingView(activeView.Name())
	if g.CurrentMode().Name() != cmdMode {
		g.SetCurrentView(activeView.Name())
		g.SetViewOnTop(activeView.Name())
	}
}

func clearView(v *gocui.View) {
	v.Clear()
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
}

func validateInput(g *gocui.Gui, v *gocui.View) error {

	if v.Name() != "inputline" {
		panic("Inputline is not the current view")
	}
	if currentDemonInput == nil {
		panic("No Current Demon Input Available")
	}

	input := v.Buffer()
	v.SetCursor(0, 0)
	v.Clear()

	if le := len(input); le < 1 {
		input = ""
	} else {
		input = input[:le-1]
	}
	var err error
	currentDemonInput, err = currentDemonInput(g, input)

	// if currentDemonInput is not nil,
	// the inputline is still open and
	// we are expecting some input
	// (see SearchAndReplace for instance)
	if currentDemonInput == nil {
		hideInputLine(g)
		updateInfos(g)
		if err == ErrViewCreated {
			return nil
		}
		g.SetCurrentView(g.Workingview().Name())
	}
	//g.SetViewOnTop(g.Workingview().Name())

	// ErrQuit should be the only error not handled
	if err != gocui.ErrQuit {
		displayError(g, err)
		return nil
	}

	return err
}

func displayError(g *gocui.Gui, e error) {
	v, _ := g.View("error")
	v.Clear()
	if e != nil {
		fmt.Fprint(v, e.Error())
		displayErrorView(g)
	} else {
		fmt.Fprint(v, "ok")
		hideErrorView(g)
	}
}

func escapeInputHandler(g *gocui.Gui, v *gocui.View) error {
	doEscapeInput(g, v)
	return nil
}

func doEscapeInput(g *gocui.Gui, v *gocui.View) {
	if v.Name() != "inputline" {
		panic("Inputline is not the current view")
	}
	if currentDemonInput == nil {
		panic("No Current Demon Input Available")
	}
	v.SetCursor(0, 0)
	v.Clear()
	currentDemonInput = nil
	g.SetCurrentView(g.Workingview().Name())
	hideInputLine(g)
	updateInfos(g)
}

func escapeMainHandler(g *gocui.Gui, v *gocui.View) error {
	hideErrorView(g)
	return nil
}

func switchModeHandlerFactory(modename string) gocui.KeybindingHandler {
	return func(g *gocui.Gui, v *gocui.View) error {
		doSwitchMode(g, modename)
		return nil
	}
}

func doSwitchMode(g *gocui.Gui, modename string) error {
	g.CurrentMode().CloseMode(g)
	if err := g.SetCurrentMode(modename); err != nil {
		return err
	}
	g.CurrentMode().OpenMode(g)
	return updateInfos(g)
}

func updateInfos(g *gocui.Gui) error {
	inMainView, x, y := cursorInfo(g)
	if inMainView {
		info, err := g.View("infoline")
		if err != nil {
			return err
		}
		info.Clear()
		maxX, _ := info.Size()
		mode := fmt.Sprintf("%s mode", g.CurrentMode().Name())
		pos := fmt.Sprintf("%d:%d", y, x)
		fmt.Fprintf(info, "%s", mode)
		fmt.Fprintf(info, "%[2]*.[2]*[1]s", pos, maxX-len(mode))
	}
	return nil
}

func cursorInfo(g *gocui.Gui) (bool, int, int) {
	v := g.CurrentView()
	if v.Name() == g.Workingview().Name() {
		x, y := v.Cursor()
		x1, y1 := v.Origin()
		return true, x + x1, y + y1
	}
	return false, 0, 0
}

func copyHandler(g *gocui.Gui, v *gocui.View) error {
	if err := copy(); err != nil {
		displayError(g, err)
	}
	return nil
}

func pasteHandler(g *gocui.Gui, v *gocui.View) error {
	if err := paste(v); err != nil {
		displayError(g, err)
	}
	updateInfos(g)
	return nil
}

func openFileHandler(g *gocui.Gui, v *gocui.View) error {
	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		return nil, openAndDisplayFile(g, input)
	}
	interactive(g, "Open File")
	return nil
}

func saveAsHandler(g *gocui.Gui, v *gocui.View) error {
	currentDemonInput = func(g *gocui.Gui, filename string) (demonInput, error) {
		return nil, saveAs(g, filename)
	}
	interactive(g, "Save as")
	return nil
}

func permutLinesUpHandler(g *gocui.Gui, v *gocui.View) error {
	v.EditPermutLines(true)
	return nil
}

func permutLinesDownHandler(g *gocui.Gui, v *gocui.View) error {
	v.EditPermutLines(false)
	return nil
}
