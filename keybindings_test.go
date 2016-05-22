package main

import (
	"fmt"
	"github.com/stretto-editor/gocui"
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
