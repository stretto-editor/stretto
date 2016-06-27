package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/stretto-editor/gocui"
)

// create the file in the directory of the
func createFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		var file *os.File
		file, _ = os.Create(filename)
		file.Close()
	}
}

var newFileHandler = func() func(g *gocui.Gui, v *gocui.View) error {
	i := 0
	return func(g *gocui.Gui, v *gocui.View) error {
		i++
		name := fmt.Sprintf("file%d", i)
		newView, _ := newFileView(g, name)
		g.SetViewOnTop(newView.Name())
		g.SetCurrentView(newView.Name())
		g.SetWorkingView(newView.Name())
		return nil
	}
}()

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

func saveMain(v *gocui.View, filename string) error {
	if filename == "" {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
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
			if _, err = f.Write(p[:n]); err != nil {
				return err
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

func openFileHandler(g *gocui.Gui, v *gocui.View) error {
	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		return nil, openAndDisplayFile(g, input)
	}
	interactive(g, "Open File")
	return nil
}

func openAndDisplayFile(g *gocui.Gui, filename string) error {
	v, _ := newFileView(g, filename)
	g.SetWorkingView(v.Name())
	if g.CurrentMode().Name() != cmdMode {
		g.SetCurrentView(v.Name())
	}
	// g.SetViewOnTop(v.Name())
	err := openFile(v, filename)
	if err == nil {
		v.Title = filename
		return nil
	}
	return fmt.Errorf("Could not open file : %s", filename)
}

func saveAsHandler(g *gocui.Gui, v *gocui.View) error {
	currentDemonInput = func(g *gocui.Gui, filename string) (demonInput, error) {
		return nil, saveAs(g, filename)
	}
	interactive(g, "Save as")
	return nil
}

func saveAs(g *gocui.Gui, filename string) error {
	//v, _ := g.View("main")
	v := g.Workingview()
	createFile(filename)
	return saveMain(v, filename)
}
