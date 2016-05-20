package main

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/jroimartin/gocui"
)

var currentFile string

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

	if err := g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, cursorHome); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, cursorEnd); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, goPgUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, goPgDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
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

func cursorHome(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, _ := v.Cursor()
		v.MoveCursor(-cx, 0, true)
	}
	return nil
}

func cursorEnd(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		l, _ := v.Line(cy)
		v.MoveCursor(len(l)-cx, 0, true)
	}
	return nil
}

func goPgUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, y := v.Size()
		yOffset := oy - y
		if yOffset < 0 {
			v.SetOrigin(ox, 0)
			if oy == 0 {
				_, cy := v.Cursor()
				v.MoveCursor(0, -cy, false)
			}
		} else {
			v.SetOrigin(ox, yOffset)
			v.MoveCursor(0, -yOffset, false)
		}
	}
	return nil
}

func goPgDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, y := v.Size()
		err := v.SetOrigin(ox, oy+y)
		if err != nil {
			return err
		}
		v.MoveCursor(0, 0, false)
	}
	return nil
}

func saveMain(g *gocui.Gui, v *gocui.View) error {
	if currentFile == "" {
		return nil
	}
	f, err := os.OpenFile(currentFile, os.O_WRONLY, 0666)
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
