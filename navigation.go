package main

import (
	"github.com/stretto-editor/gocui"
)

func cursorHome(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		v.SetOrigin(0, oy)
		v.SetCursor(0, cy)
	}
	updateInfos(g)
	g.CurrentView().Actions.Cut()
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
	updateInfos(g)
	g.CurrentView().Actions.Cut()
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
	updateInfos(g)
	g.CurrentView().Actions.Cut()
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
	updateInfos(g)
	g.CurrentView().Actions.Cut()
	return nil
}

func moveLeft(g *gocui.Gui, v *gocui.View) error {
	moveAndInfo(g, v, -1, 0, false)
	return nil
}

func moveRight(g *gocui.Gui, v *gocui.View) error {
	moveAndInfo(g, v, 1, 0, false)
	return nil
}

func moveUp(g *gocui.Gui, v *gocui.View) error {
	moveAndInfo(g, v, 0, -1, false)
	return nil
}

func moveDown(g *gocui.Gui, v *gocui.View) error {
	moveAndInfo(g, v, 0, 1, false)
	return nil
}

func scrollUp(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	if oy != 0 {
		v.SetOrigin(0, oy-1)
	}
	updateInfos(g)
	g.CurrentView().Actions.Cut()
	return nil
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	// allowed infinite scroll
	v.SetOrigin(0, oy+1)
	updateInfos(g)
	g.CurrentView().Actions.Cut()
	return nil
}

func moveAndInfo(g *gocui.Gui, v *gocui.View, x int, y int, b bool) {
	v.MoveCursor(x, y, b)
	updateInfos(g)
	g.CurrentView().Actions.Cut()
}
