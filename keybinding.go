package main

import (
	"encoding/binary"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/stretto-editor/gocui"
)

var currentFile string

const fileMode = "file"
const editMode = "edit"
const cmdMode = "cmd"

type input struct {
	channel chan int
	content string
}

var in input

func initModes(g *gocui.Gui) {
	g.SetMode(cmdMode)
	g.SetMode(fileMode)
	g.SetMode(editMode)
}

func initKeybindings(g *gocui.Gui) error {

	var keyBindings = []struct {
		m string
		v string
		k interface{}
		h gocui.KeybindingHandler
	}{

		// NAVIGATION

		{m: fileMode, v: "", k: gocui.KeyHome, h: cursorHome},
		{m: fileMode, v: "", k: gocui.KeyEnd, h: cursorEnd},
		{m: fileMode, v: "", k: gocui.KeyPgup, h: goPgUp},
		{m: fileMode, v: "", k: gocui.KeyPgdn, h: goPgDown},

		{m: editMode, v: "", k: gocui.KeyHome, h: cursorHome},
		{m: editMode, v: "", k: gocui.KeyEnd, h: cursorEnd},
		{m: editMode, v: "", k: gocui.KeyPgup, h: goPgUp},
		{m: editMode, v: "", k: gocui.KeyPgdn, h: goPgDown},

		// USEFUL

		{m: fileMode, v: "", k: gocui.KeyCtrlQ, h: quitHandler},
		{m: fileMode, v: "main", k: gocui.KeyTab, h: switchModeTo(editMode)},
		{m: fileMode, v: "", k: gocui.KeyCtrlT, h: currTopViewHandler("cmdline")},
		{m: fileMode, v: "cmdline", k: gocui.KeyCtrlT, h: currTopViewHandler("main")},
		{m: fileMode, v: "main", k: 'o', h: openFileHandler},

		{m: editMode, v: "", k: gocui.KeyCtrlQ, h: quitHandler},
		{m: editMode, v: "main", k: gocui.KeyTab, h: switchModeTo(fileMode)},
		{m: editMode, v: "", k: gocui.KeyCtrlT, h: currTopViewHandler("cmdline")},
		{m: editMode, v: "cmdline", k: gocui.KeyCtrlT, h: currTopViewHandler("main")},
		{m: editMode, v: "main", k: gocui.KeyCtrlO, h: openFileHandler},

		// EDITION

		{m: fileMode, v: "main", k: 's', h: saveHandler},
		{m: fileMode, v: "main", k: 'u', h: saveAsHandler},
		{m: fileMode, v: "main", k: 'f', h: searchHandler},
		{m: fileMode, v: "main", k: 'p', h: searchAndReplaceHandler},
		{m: fileMode, v: "main", k: 'a', h: exampleInputFunc},

		{m: fileMode, v: "inputline", k: gocui.KeyEnter, h: validateInput},

		{m: editMode, v: "main", k: gocui.KeyCtrlS, h: saveHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlU, h: saveAsHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlF, h: searchHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlP, h: searchAndReplaceHandler},
		{m: editMode, v: "main", k: gocui.KeyCtrlA, h: exampleInputFunc},
		{m: editMode, v: "main", k: gocui.KeyCtrlC, h: copy},
		{m: editMode, v: "main", k: gocui.KeyCtrlV, h: paste},

		{m: editMode, v: "main", k: gocui.KeyEnter, h: breaklineHandler},
		{m: editMode, v: "inputline", k: gocui.KeyEnter, h: validateInput},
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

func searchHandler(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		v, _ := g.View("main")
		search(v, input)
		return nil, nil
	}

	g.SetCurrentView("inputline")
	g.SetViewOnTop("inputline")
	g.CurrentView().MoveCursor(0, 0, false)

	return nil
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

		g.SetCurrentView("inputline")
		g.SetViewOnTop("inputline")
		g.CurrentView().MoveCursor(0, 0, false)

		return nil
	}

	vmain, _ := g.View("main")
	if err := saveMain(vmain, currentFile); err != nil {
		return err
	}
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

	g.SetCurrentView("inputline")
	g.SetViewOnTop("inputline")
	g.CurrentView().MoveCursor(0, 0, false)

	return nil
}

func searchAndReplaceHandler(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		v, _ := g.View("main")

		if found := search(v, input); !found {
			return nil, nil
		}

		searched := input

		return func(g *gocui.Gui, input string) (demonInput, error) {
			v, _ := g.View("main")

			for i := 0; i < len(searched); i++ {
				v.EditDelete(false)
			}

			for _, c := range input {
				v.EditWrite(c)
			}
			return nil, nil
		}, nil

	}

	g.SetCurrentView("inputline")
	g.SetViewOnTop("inputline")
	g.CurrentView().MoveCursor(0, 0, false)

	return nil
}

func search(v *gocui.View, pattern string) bool {
	if len(pattern) > 0 {

		var s string
		var err error
		var sameline = 1

		x, y := v.Cursor()

		for i := 0; err == nil; i++ {
			s, err = v.Line(y + i)
			if err == nil {
				// size of line is long enough to move the cursor
				if x < len(s) {
					indice := strings.Index(s[x+sameline:], pattern)

					// existing element on this line
					if indice >= 0 {
						if sameline == 0 {
							x, y = v.Cursor()
							v.MoveCursor(indice+sameline-x, i, false)
						} else {
							v.MoveCursor(indice+sameline, i, false)
						}
						return true
					}
				}
				x = 0
				sameline = 0
			}
		}
	}
	return false
}

func exampleInputFunc(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		vmain, _ := g.View("main")
		for _, ch := range input {
			vmain.EditWrite(ch)
		}
		return nil, nil
	}

	g.SetCurrentView("inputline")
	g.SetViewOnTop("inputline")
	g.CurrentView().MoveCursor(0, 0, false)
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

	return err
}

func switchModeTo(name string) gocui.KeybindingHandler {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := g.SetCurrentMode(name); err != nil {
			return err
		}
		if v.Name() == "cmdline" {
			v.SetCursor(0, 0)
			v.Clear()
		}
		if name == "cmd" {
			g.SetCurrentView("cmdline")
		}

		v.Editable = true

		if name == "file" {
			v.Editable = false
		}
		return nil
	}
}

func currTopViewHandler(name string) gocui.KeybindingHandler {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := g.SetCurrentView(name); err != nil {
			return err
		}
		if _, err := g.SetViewOnTop(name); err != nil {
			return err
		}
		return nil
	}
}

func cursorHome(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		v.SetOrigin(0, oy)
		v.SetCursor(0, cy)
	}
	return nil
}

func cursorEnd(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		x, _ := v.Size()
		l, _ := v.Line(cy)
		if len(l) > x {
			v.SetOrigin(len(l)-x+1, oy)
			v.SetCursor(x-1, cy)
		} else {
			v.SetCursor(len(l), cy)
		}
	}
	return nil
}

func goPgUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, y := v.Size()
		yOffset := oy - y
		if yOffset < 0 {
			if err := v.SetOrigin(ox, 0); err != nil {
				return err
			}
			if oy == 0 {
				_, cy := v.Cursor()
				v.MoveCursor(0, -cy, false)
			} else {
				v.MoveCursor(0, 0, false)
			}
		} else {
			v.SetOrigin(ox, yOffset)
			v.MoveCursor(0, 0, false)
		}
	}
	return nil
}

func goPgDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, y := v.Size()
		_, cy := v.Cursor()
		if y >= v.BufferSize() {
			v.MoveCursor(0, y-cy, false)
		} else if oy >= v.BufferSize()-y {
			v.MoveCursor(0, y, false)
		} else if oy+2*y >= v.BufferSize() {
			v.SetOrigin(ox, v.BufferSize()-y+1)
			v.MoveCursor(0, 0, false)
		} else {
			v.SetOrigin(ox, oy+y)
			v.MoveCursor(0, 0, false)
		}
	}
	return nil
}

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

/*
func searchInteractive(g *gocui.Gui, v *gocui.View) (bool, error) {
	currTopViewHandler("inputline")(g, v)
	g.CurrentView().MoveCursor(0, 0, false)
	in.channel <- 1

	if len(in.content) > 0 {

		var s string
		var err error
		var sameline = 1

		x, y := v.Cursor()

		for i := 0; err == nil; i++ {
			s, err = v.Line(y + i)
			if err == nil {
				// size of line is long enough to move the cursor
				if x < len(s)-1 {
					indice := strings.Index(s[x+sameline:], in.content)

					// existing element on this line
					if indice >= 0 {
						if sameline == 0 {
							x, y = v.Cursor()
							v.MoveCursor(indice+sameline-x, i, false)
						} else {
							v.MoveCursor(indice+sameline, i, false)
						}
						return true, nil
					}
				}
				x = 0
				sameline = 0
			}
		}
	}
	return false, nil
}
*/

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
		openFile(v, input)
		currentFile = input
		return nil, nil
	}

	g.SetCurrentView("inputline")
	g.SetViewOnTop("inputline")
	g.CurrentView().MoveCursor(0, 0, false)

	return nil
}

/*

func replace(g *gocui.Gui, v *gocui.View) error {
	go replaceInteractive(g, v)
	return nil
}

func replaceInteractive(g *gocui.Gui, v *gocui.View) error {
	if found, err := searchInteractive(g, v); err != nil || found == false {
		if err != nil {
			return err
		}
		return nil
	}

	searched := in.content

	currTopViewHandler("inputline")(g, v)
	g.CurrentView().MoveCursor(0, 0, false)
	in.channel <- 1

	for i := 0; i < len(searched); i++ {
		v.EditDelete(false)
	}

	replaced := in.content

	for _, c := range replaced {
		v.EditWrite(c)
	}

	return nil
}
*/
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

	g.SetCurrentView("inputline")
	g.SetViewOnTop("inputline")
	g.CurrentView().MoveCursor(0, 0, false)

	return nil
}
