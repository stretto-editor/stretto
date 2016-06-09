package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCursorHomeEnd(t *testing.T) {

	g := initGui()
	defer g.Close()

	err := openAndDisplayFile(g, "LICENSE")
	assert.Nil(t, err, "No error should be found")

	wView := g.Workingview()
	xPos, yPos := wView.Cursor()

	assert.Equal(t, xPos, 0, "Cursor should be at first column")
	assert.Equal(t, yPos, 0, "Cursor should be at first line")

	err = cursorEnd(g, wView)
	assert.Nil(t, err, "cursorEnd shouldn't return any error")

	xPos, _ = wView.Cursor()
	line, _ := wView.Line(yPos)

	assert.Equal(t, xPos, len(line), "Cursor at wrong position")
	assert.Equal(t, yPos, 0, "Cursor at wrong position")

	err = cursorHome(g, wView)
	assert.Nil(t, err, "cursorEnd shouldn't return any error")

	xPos, _ = wView.Cursor()
	assert.Equal(t, xPos, 0, "Cursor at wrong position")
	assert.Equal(t, yPos, 0, "Cursor at wrong position")

}

func TestPageUpDown(t *testing.T) {

	g := initGui()
	defer g.Close()

	err := openAndDisplayFile(g, "LICENSE")
	assert.Nil(t, err, "No error should be found")

	wView := g.Workingview()
	_, ySize := wView.Size()
	xPos, yPos := wView.Cursor()

	err = goPgUp(g, wView)
	assert.Nil(t, err, "No error should happen")

	assert.Equal(t, xPos, 0, "Cursor at wrong position")
	assert.Equal(t, yPos, 0, "Cursor at wrong position")

	err = goPgDown(g, wView)
	assert.Nil(t, err, "No error should happen")

	xPos, yPos = wView.Origin()
	xT, yT := wView.Cursor()

	assert.Equal(t, xPos+xT, 0, "Cursor at wrong position")
	assert.Equal(t, yPos+yT, ySize, "Cursor at wrong position")
}

func TestMoveLeftRight(t *testing.T) {

	g := initGui()
	defer g.Close()

	err := openAndDisplayFile(g, "LICENSE")
	assert.Nil(t, err, "No error should be found")

	wView := g.Workingview()
	xPos, yPos := wView.Cursor()

	assert.Equal(t, xPos, 0, "Cursor at wrong position")
	assert.Equal(t, yPos, 0, "Cursor at wrong position")

	moveRight(g, wView)
	xPos, yPos = wView.Cursor()
	assert.Equal(t, xPos, 1, "Cursor at wrong position")
	assert.Equal(t, yPos, 0, "Cursor at wrong position")

	//wView.MoveCursor(-1, 0, false)
	moveLeft(g, wView)
	xPos, yPos = wView.Cursor()
	assert.Equal(t, xPos, 0, "Cursor at wrong position")
	assert.Equal(t, yPos, 0, "Cursor at wrong position")

	wView.MoveCursor(4, 0, false)
	xPos, yPos = wView.Cursor()
	assert.Equal(t, xPos, 4, "Cursor at wrong position")
	assert.Equal(t, yPos, 0, "Cursor at wrong position")

	wView.MoveCursor(0, 4, false)
	//xPos, yPos = wView.Cursor()
	xPos, _ = wView.Cursor()
	assert.Equal(t, xPos, 4, "Cursor at wrong position")
	// viewlines == nil
	//assert.Equal(t, yPos, 4, "Cursor at wrong position")

	wView.MoveCursor(0, -2, false)
	//xPos, yPos = wView.Cursor()
	xPos, _ = wView.Cursor()
	assert.Equal(t, xPos, 4, "Cursor at wrong position")
	//assert.Equal(t, yPos, 2, "Cursor at wrong position")

}
