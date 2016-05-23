package main

import (
	"fmt"
	"github.com/stretto-editor/gocui"
	"strings"
	"testing"
)

func TestQuit(t *testing.T) {
	if err := quit(&gocui.Gui{}, &gocui.View{}); err != gocui.ErrQuit {
		t.Error("quit should return ErrQuit")
	}
}

func TestCursor(t *testing.T) {
	g := gocui.NewGui()
	g.Init()
	defer g.Close()

	maxX, maxY := g.Size()

	v, _ := g.SetView("testA", 0, 0, maxX, maxY)
	fmt.Fprint(v, "foo")
	cursorEnd(g, v)
	if x, y := v.Cursor(); x != 3 || y != 0 {
		t.Errorf("Cursor is not at the end of the line. Current position : %d %d", x, y)
	}
	cursorHome(g, v)
	if x, y := v.Cursor(); x != 0 || y != 0 {
		t.Error("Cursor is not at the beginning of the line")
	}

}

func TestInitMode(t *testing.T) {

	var testModeNames = []string{
		"cmd",
		"file",
		"edit",
	}

	g := gocui.NewGui()
	g.Init()
	defer g.Close()

	initModes(g)

	for _, mn := range testModeNames {
		if _, err := g.Mode(mn); err != nil {
			t.Errorf("mode %s doesnt exist", mn)
		}
	}
}

func TestSwitchMode(t *testing.T) {

	g := gocui.NewGui()
	g.Init()
	defer g.Close()

	g.SetMode("testA")
	g.SetMode("testB")
	g.SetCurrentMode("testA")

	f := switchModeTo("testB")
	f(g, &gocui.View{})

	if m := g.CurrentMode(); m.Name() != "testB" {
		t.Error("Wrong current mode")
	}
}

func TestPageDownUp(t *testing.T) {
	g := gocui.NewGui()
	g.Init()
	defer g.Close()

	maxX, maxY := g.Size()
	v, _ := g.SetView("test", 0, 0, maxX, maxY)
	g.SetCurrentView("test")
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(v, "%d\n", i)
	}

	goPgDown(g, v)
	goPgDown(g, v)

	if fl := strings.Split(v.ViewBuffer(), "\n"); fl[0] != string(2*maxY-1) {
		t.Errorf("Found line no %s instead of line no %d", fl[0], 2*maxY-1)
	}

	goPgUp(g, v)
	if fl := strings.Split(v.ViewBuffer(), "\n"); fl[0] != string(maxY-1) {
		t.Errorf("Found line no %s instead of line no %d", fl[0], 2*maxY-1)
	}
}

func TestCurrTopViewHandler(t *testing.T) {

	g := gocui.NewGui()
	g.Init()
	defer g.Close()

	maxX, maxY := g.Size()
	vA, _ := g.SetView("testA", 0, 0, maxX, maxY)
	fmt.Fprintf(vA, "foo")
	vB, _ := g.SetView("testB", 0, 0, maxX, maxY)
	fmt.Fprintf(vB, "bar")
	g.SetCurrentView("testA")

	f := currTopViewHandler("testB")
	f(g, vA)
	if g.CurrentView() != vB {
		t.Error("testB was expected to be the current view")
	}
	if g.CurrentView().ViewBuffer() != vB.ViewBuffer() {
		t.Error("testB was expected to be on top")
	}
}

func TestInitKeybindings(t *testing.T) {
	g := gocui.NewGui()
	g.Init()
	defer g.Close()

	initModes(g)
	g.SetLayout(layout)

	if err := initKeybindings(g); err != nil {
		t.Error("bad keybindings initialization")
	}

}
