package main

import (
	"encoding/binary"
	"fmt"
	"github.com/stretto-editor/gocui"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var currentFile string

const fileMode = "file"
const editMode = "edit"
const cmdMode = "cmd"

func initModes(g *gocui.Gui) {
	openCmdMode := func(g *gocui.Gui) error {
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
		g.SetCurrentView("main")
		g.SetViewOnTop("main")
		return nil
	}
	openFileMode := func(g *gocui.Gui) error {
		var v *gocui.View
		var err error
		if v, err = g.View("main"); err != nil {
			return err
		}
		v.SetEditable(false)
		return nil
	}
	closeFileMode := func(g *gocui.Gui) error {
		var v *gocui.View
		var err error
		if v, err = g.View("main"); err != nil {
			return err
		}
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
}

func initKeybindings(g *gocui.Gui) error {

	var keyBindings = []struct {
		m string
		v string
		k interface{}
		h gocui.KeybindingHandler
	}{

		// ---------------------- COMMON COMMANDS ------------------------- //

		{m: fileMode, v: "", k: gocui.KeyCtrlT, h: switchModeHandlerFactory(cmdMode)},
		{m: fileMode, v: "", k: gocui.KeyTab, h: switchModeHandlerFactory(editMode)},
		{m: fileMode, v: "", k: gocui.KeyCtrlQ, h: quitHandler},

		{m: editMode, v: "", k: gocui.KeyCtrlT, h: switchModeHandlerFactory(cmdMode)},
		{m: editMode, v: "", k: gocui.KeyTab, h: switchModeHandlerFactory(fileMode)},
		{m: editMode, v: "", k: gocui.KeyCtrlQ, h: quitHandler},

		{m: cmdMode, v: "", k: gocui.KeyCtrlT, h: switchModeHandlerFactory(fileMode)},

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

		// ---------------------- USEFUL --- ------------------------------ //

		{m: fileMode, v: "main", k: 'o', h: openFileHandler},
		{m: fileMode, v: "main", k: 's', h: saveHandler},
		{m: fileMode, v: "main", k: 'u', h: saveAsHandler},
		{m: fileMode, v: "main", k: 'f', h: searchHandler},
		{m: fileMode, v: "main", k: 'b', h: commandInfoHandler},

		{m: editMode, v: "main", k: gocui.KeyCtrlO, h: openFileHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlS, h: saveHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlU, h: saveAsHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlF, h: searchHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlB, h: commandInfoHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlP, h: searchAndReplaceHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlC, h: copy},
		{m: editMode, v: "main", k: gocui.KeyCtrlV, h: paste},
		{m: editMode, v: "main", k: gocui.KeyEnter, h: breaklineHandler},

		// ---------------------- INFO SECTION ---------------------------- //

		// ---------------------- NAVIGATION ------------------------------ //

		{m: fileMode, v: "cmdinfo", k: gocui.KeyArrowUp, h: scrollUp},
		{m: fileMode, v: "cmdinfo", k: gocui.KeyArrowDown, h: scrollDown},
		{m: fileMode, v: "cmdinfo", k: gocui.KeyPgup, h: goPgUp},
		{m: fileMode, v: "cmdinfo", k: gocui.KeyPgdn, h: goPgDown},
		{m: fileMode, v: "cmdinfo", k: gocui.KeyEsc, h: quitInfo},

		{m: editMode, v: "cmdinfo", k: gocui.KeyArrowUp, h: scrollUp},
		{m: editMode, v: "cmdinfo", k: gocui.KeyArrowDown, h: scrollDown},
		{m: editMode, v: "cmdinfo", k: gocui.KeyPgup, h: goPgUp},
		{m: editMode, v: "cmdinfo", k: gocui.KeyPgdn, h: goPgDown},
		{m: editMode, v: "cmdinfo", k: gocui.KeyEsc, h: quitInfo},

		// ---------------------- INPUT SECTION --------------------------- //

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
	}

	for _, kb := range keyBindings {
		if err := g.SetKeybinding(kb.m, kb.v, kb.k, gocui.ModNone, kb.h); err != nil {
			return err
		}
	}

	g.SetCurrentMode(editMode)

	return nil
}

func breaklineHandler(g *gocui.Gui, v *gocui.View) error {
	v.EditNewLine()
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// demonInput defines the prototype for functions that
// should be called later in validateInput
// A demonInput returns the next demonInput to be called,
// or nil if there is noone
type demonInput func(g *gocui.Gui, input string) (demonInput, error)

// current is the next demonInput to be called
// In an handler (that use the inputline), it should be set.
// It is set back to nil once all the call have been made
// (in validateInput handler)
var currentDemonInput demonInput

func interactive(g *gocui.Gui, s string) {
	g.SetCurrentView("inputline")
	g.SetViewOnTop("inputline")
	g.CurrentView().Title = " " + s + " "
	g.CurrentView().MoveCursor(0, 0, false)
}

func saveHandler(g *gocui.Gui, v *gocui.View) error {

	if currentFile == "" {

		currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {

			createFile(input)

			v, _ := g.View("main")
			if err := saveMain(v, input); err != nil {
				return nil, err
			}

			return nil, nil
		}

		interactive(g, "Save")
		return nil
	}

	vmain, _ := g.View("main")
	if err := saveMain(vmain, currentFile); err != nil {
		return err
	}
	return nil
}

func commandInfoHandler(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	wcmd, hcmd := maxX*70/100, maxY*70/100
	var xcmd, ycmd int = (maxX - wcmd) / 2, maxY/2 - hcmd/2
	if v, err := g.SetView("cmdinfo", xcmd, ycmd, xcmd+wcmd, ycmd+hcmd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = false
		v.Wrap = true
		v.Title = "Commands Summary"
		openFile(v, "Commands.md")
		g.SetViewOnTop("cmdinfo")
		g.SetCurrentView("cmdinfo")
	}
	return nil
}

func quitInfo(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("main")
	g.SetViewOnTop("main")
	g.DeleteView("cmdinfo")
	return nil
}

// create the file in the directory of the
func createFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		var file *os.File
		file, _ = os.Create(filename)
		file.Close()
		currentFile = filename
	}
}

func quitHandler(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		if input != "n" {
			if currentFile == "" {
				interactive(g, "File name")
				return func(g *gocui.Gui, input string) (demonInput, error) {

					createFile(input)

					v, _ := g.View("main")
					if err := saveMain(v, input); err != nil {
						return nil, err
					}

					return nil, gocui.ErrQuit
				}, nil

			}

			v, _ := g.View("main")
			if err := saveMain(v, currentFile); err != nil {
				return nil, err
			}
		}
		return nil, gocui.ErrQuit
	}

	interactive(g, "Save Modifications (y/n)")
	return nil
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
		g.SetCurrentView("main")
		g.SetViewOnTop("main")
	}

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
	} else {
		fmt.Fprint(v, "ok")
	}
	displayErrorView(g)
}

func escapeInputHandler(g *gocui.Gui, v *gocui.View) error {
	doEscapeInput(g, v)
	return nil
}

func escapeMainHandler(g *gocui.Gui, v *gocui.View) error {
	hideErrorView(g)
	return nil
}

func setTopViewHandlerFactory(viewname string) gocui.KeybindingHandler {
	return func(g *gocui.Gui, v *gocui.View) error {
		doSetTopView(g, viewname)
		return nil
	}
}
func switchModeHandlerFactory(modename string) gocui.KeybindingHandler {
	return func(g *gocui.Gui, v *gocui.View) error {
		doSwitchMode(g, modename)
		return nil
	}
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
	g.SetCurrentView("main")
	g.SetViewOnTop("main")
}

func doSetTopView(g *gocui.Gui, viewname string) error {
	if err := g.SetCurrentView(viewname); err != nil {
		return err
	} else {
		g.SetViewOnTop(viewname)
	}
	return nil
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
		fmt.Fprintf(info, "Currently in "+g.CurrentMode().Name()+" mode \n"+"Cursor Position : "+y+":"+x)
	}
	return nil
}

func cursorInfo(g *gocui.Gui) (bool, string, string) {
	v := g.CurrentView()
	if v.Name() == "main" {
		x, y := v.Cursor()
		x1, y1 := v.Origin()
		xstring, ystring := strconv.Itoa(x+x1), strconv.Itoa(y+y1)
		return true, xstring, ystring
	}
	return false, "", ""
}

// func currTopViewHandler(name string) gocui.KeybindingHandler {
// 	return func(g *gocui.Gui, v *gocui.View) error {
// 		if err := g.SetCurrentView(name); err != nil {
// 			return err
// 		}
// 		if _, err := g.SetViewOnTop(name); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	g.CurrentMode().OpenMode(g)
// 	return nil
// }

func saveMain(v *gocui.View, filename string) error {
	if filename == "" {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		if strings.HasSuffix(err.Error(), "permission denied") {
			return nil
		}
		return err
	}
	defer f.Close()

	p := make([]byte, 5)
	v.Rewind()
	var size int64 = -1
	for {
		n, err := v.Read(p)
		size += int64(n)
		if n > 0 {
			if _, er := f.Write(p[:n]); err != nil {
				return er
			}
		}
		if err == io.EOF {
			f.Truncate(size * int64(binary.Size(p[0])))
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func copy(g *gocui.Gui, v *gocui.View) error {
	//http://stackoverflow.com/questions/10781516/how-to-pipe-several-commands-in-go
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		return nil
	}
	c1 := exec.Command("xclip", "-o")
	c2 := exec.Command("xclip", "-i", "-selection", "c")
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	if err := c1.Start(); err != nil {
		// print : error can't find xclip
		return nil
	}
	if err := c2.Start(); err != nil {
		return err
	}
	if err := c1.Wait(); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if err := c2.Wait(); err != nil {
		return err
	}
	return nil
}

func paste(g *gocui.Gui, v *gocui.View) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	out, err := exec.Command("xclip", "-o", "-selection", "c").Output()
	s := string(out)
	if err != nil {
		//print error : can't find xclip
		return nil
	}
	for _, r := range s {
		if rune(r) == '\n' {
			v.EditNewLine()
		} else {
			v.EditWrite(rune(r))
		}
	}
	return nil
}

func openFileHandler(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		v, _ := g.View("main")
		err := openFile(v, input)
		if err == nil {
			currentFile = input
			return nil, nil
		}
		return nil, fmt.Errorf("Could not open file : %s", input)
	}

	interactive(g, "Open File")
	return nil
}

func saveAsHandler(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, filename string) (demonInput, error) {
		v, _ := g.View("main")
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			var file *os.File
			file, _ = os.Create(filename)
			file.Close()
		}
		return nil, saveMain(v, filename)
	}

	interactive(g, "Save as")
	return nil
}
