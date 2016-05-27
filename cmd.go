package main

import (
	"strings"

	"github.com/stretto-editor/gocui"
)

func validateCmd(g *gocui.Gui, v *gocui.View) error {
	cmdBuff := v.Buffer()
	if cmdBuff == "" {
		return nil
	}
	cmdBuff = cmdBuff[:len(cmdBuff)-1]
	cmd := strings.Split(cmdBuff, " ")
	switch cmd[0] {
	case "quit", "q":
		return quit(g, v)
	case "qs", "sq":
		return saveAndQuit(g, cmd)
	case "sc":
		saveAndClose(g, cmd)
	case "o", "open":
		if len(cmd) > 1 {
			openAndDisplayFile(g, cmd[1])
		}
	case "saveas", "sa":
		if len(cmd) > 1 {
			saveAs(g, cmd[1])
		}
	case "replaceall", "repall":
		replaceAll(g, cmd)
	}
	v.Clear()
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	return nil
}

func saveAndQuit(g *gocui.Gui, cmd []string) error {
	if currentFile == "" && len(cmd) == 1 {
		return nil // print error command
	}
	vMain, _ := g.View("main")
	filename := currentFile
	if filename == "" {
		filename = cmd[1]
	}
	createFile(filename)
	saveMain(vMain, filename)
	return quit(g, vMain)
}

func replaceAll(g *gocui.Gui, cmd []string) {
	if len(cmd) > 2 {
		vMain, _ := g.View("main")
		for found, x, y := searchForward(vMain, cmd[1], 0, 0); found; found, x, y = searchForward(vMain, cmd[1], x, y) {
			replaceAt(vMain, x, y, cmd[1], cmd[2])
		}
	}
}

func saveAndClose(g *gocui.Gui, cmd []string) {
	if currentFile != "" || len(cmd) > 1 {
		vMain, _ := g.View("main")
		if currentFile == "" {
			createFile(cmd[1])
			saveMain(vMain, cmd[1])
		} else {
			saveMain(vMain, currentFile)
		}
		currentFile = ""
		vMain.Title = "undefined"
		vMain.Clear()
		vMain.SetCursor(0, 0)
	}
}
