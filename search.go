package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/stretto-editor/gocui"
)

var (
	// ErrPatternNotFound raised when the pattern is not found
	ErrPatternNotFound = errors.New("Unable to find")
)

func searchHandler(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		return nil, search(g, input)
	}

	interactive(g, "Search")
	return nil
}

func search(g *gocui.Gui, input string) error {
	v := g.Workingview()

	if found, x, y := v.SearchForward(input); found == true {
		v.AbsMoveCursor(x, y, false)
		return nil
	}

	return fmt.Errorf("Could not find pattern \"%s\" forward", input)
}

func searchAndReplaceHandler(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {

		v := g.Workingview()
		var x, y int
		var found bool

		if found, x, y = v.SearchForward(input); !found {
			return nil, fmt.Errorf("Could not find pattern \"%s\" forward", input)
		}

		searched := input
		interactive(g, "Search and replace - Replace string")

		return func(g *gocui.Gui, input string) (demonInput, error) {
			v := g.Workingview()
			replaceAt(v, x, y, searched, input)
			return nil, nil
		}, nil

	}

	interactive(g, "Search and replace - Search string")
	return nil
}

func replaceAll(g *gocui.Gui, pattern, replacement string) {
	v := g.Workingview()
	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)

	found, x, y := v.SearchForward(pattern)
	for found {
		replaceAt(v, x, y, pattern, replacement)
		found, x, y = v.SearchForward(pattern)
	}
}

func replaceAt(v *gocui.View, x, y int, oldstring, newstring string) {
	v.AbsMoveCursor(x, y, false)
	for i := 0; i < len(oldstring); i++ {
		v.EditDelete(false)
	}
	for _, c := range newstring {
		v.EditWrite(c)
	}
}

// func gives how to move from the current origin
func searchForward(v *gocui.View, pattern string, x int, y int) (bool, int, int) {

	if len(pattern) > 0 {

		var s string
		var err error
		var sameline = 1

		for i := 0; err == nil; i++ {
			s, err = v.Line(y + i)
			if err == nil {

				if x < len(s) {
					indice := strings.Index(s[x+sameline:], pattern)

					if indice >= 0 {
						if sameline == 0 {
							return true, indice + sameline, y + i
						}
						return true, indice + sameline + x, y + i
					}
				}
				x, sameline = 0, 0
			}
		}
	}
	return false, 0, 0
}
