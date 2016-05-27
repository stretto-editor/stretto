package main

import (
	"strings"

	"errors"

	"github.com/stretto-editor/gocui"
)

var (
	// ErrPatternNotFound raised when the pattern is not foud
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
	v, _ := g.View("main")
	x, y := v.Cursor()

	if found, newx, newy := searchForward(v, input, x, y); found == true {
		moveTo(v, newx, newy)
		return nil
	}

	return ErrPatternNotFound

}

func replaceAt(v *gocui.View, x, y int, oldstring, newstring string) {
	moveTo(v, x, y)
	for i := 0; i < len(oldstring); i++ {
		v.EditDelete(false)
	}

	for _, c := range newstring {
		v.EditWrite(c)
	}
}

func searchAndReplaceHandler(g *gocui.Gui, v *gocui.View) error {

	currentDemonInput = func(g *gocui.Gui, input string) (demonInput, error) {
		v, _ := g.View("main")
		x, y := v.Cursor()
		var xnew, ynew int
		var found bool

		if found, xnew, ynew = searchForward(v, input, x, y); !found {
			return nil, nil
		}

		searched := input
		interactive(g, "Search and replace - Replace string")

		return func(g *gocui.Gui, input string) (demonInput, error) {
			v, _ := g.View("main")
			replaceAt(v, xnew, ynew, searched, input)
			return nil, nil
		}, nil

	}

	interactive(g, "Search and replace - Search string")
	return nil
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

// Moves the cursor relatively to the origin of the view
func moveTo(v *gocui.View, x int, y int) error {
	_, yOrigin := v.Origin()
	_, ySize := v.Size()

	if y <= ySize {

		v.SetCursor(x, y)
		return nil
	}

	// how many times we move from the size of the window
	var i int
	for i = 0; y > ySize; i++ {
		y -= ySize

	}
	v.SetOrigin(0, yOrigin+i*ySize)
	v.SetCursor(x, y)
	return nil
}
