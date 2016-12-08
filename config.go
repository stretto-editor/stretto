package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"github.com/stretto-editor/gocui"
)

type config struct {
	Wrap        bool
	Cursor      bool
	Guibgcolor  string
	Guifgcolor  string
	Viewbgcolor string
	Viewfgcolor string
	Selbgcolor  string
	Selfgcolor  string
	Highlight   bool
}

var userconfig config

func initConfig(g *gocui.Gui) {
	usr, _ := user.Current()
	file, e := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".stretto.json"))

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		return
	}
	json.Unmarshal(file, &userconfig)

	c, _ := g.ViewNode("main")
	v := c.LastView()

	v.Wrap = userconfig.Wrap
	v.Highlight = userconfig.Highlight
	v.BgColor = setColor(userconfig.Viewbgcolor)
	v.FgColor = setColor(userconfig.Viewfgcolor)
	v.SelBgColor = setColor(userconfig.Selbgcolor)
	v.SelFgColor = setColor(userconfig.Selfgcolor)
	v.Highlight = userconfig.Highlight

	g.Cursor = userconfig.Cursor
	g.BgColor = setColor(userconfig.Guibgcolor)
	g.FgColor = setColor(userconfig.Guifgcolor)
}

func setColor(s string) gocui.Attribute {
	switch s {
	case "black":
		return gocui.ColorBlack
	case "red":
		return gocui.ColorRed
	case "green":
		return gocui.ColorGreen
	case "yellow":
		return gocui.ColorYellow
	case "blue":
		return gocui.ColorBlue
	case "magenta":
		return gocui.ColorMagenta
	case "cyan":
		return gocui.ColorCyan
	case "white":
		return gocui.ColorWhite
	default:
		return gocui.ColorDefault

	}
}
