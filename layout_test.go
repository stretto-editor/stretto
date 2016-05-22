package main

import (
	"github.com/stretto-editor/gocui"
	"testing"
)

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

/*
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
*/
