package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"testing"

	"github.com/jroimartin/gocui"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	g := initGui()
	defer g.Close()
	usr, _ := user.Current()
	configPath := usr.HomeDir + "/.stretto.json"
	copyPath := usr.HomeDir + "/r9w92W2Cn7MTtAhuCP5si2LxwP9UrmC6Y99.json"
	err := os.Rename(configPath, copyPath)
	f, _ := os.Create(configPath)
	w := bufio.NewWriter(f)
	s := "{\n" +
		"\"wrap\" : true,\n" +
		"\"cursor\" : true,\n" +
		"\"guibgcolor\" : \"black\",\n" +
		"\"guifgcolor\" : \"white\",\n" +
		"\"viewbgcolor\" : \"red\",\n" +
		"\"viewfgcolor\" : \"yellow\",\n" +
		"\"selbgcolor\" : \"blue\",\n" +
		"\"selfgcolor\" : \"green\",\n" +
		"\"highlight\" : true\n" +
		"}"
	fmt.Fprintf(w, s)
	w.Flush()
	initConfig(g)
	assert.Equal(t, gocui.Attribute(gocui.ColorBlack), gocui.Attribute(g.BgColor))
	assert.Equal(t, gocui.Attribute(gocui.ColorWhite), gocui.Attribute(g.FgColor))
	os.Remove(configPath)
	if err == nil {
		os.Rename(copyPath, configPath)
	}
}

func TestInitConfig2(t *testing.T) {
	g := initGui()
	defer g.Close()
	usr, _ := user.Current()
	configPath := usr.HomeDir + "/.stretto.json"
	copyPath := usr.HomeDir + "/r9w92W2Cn7MTtAhuCP5si2LxwP9UrmC6Y99.json"
	err := os.Rename(configPath, copyPath)
	f, _ := os.Create(configPath)
	w := bufio.NewWriter(f)
	s := "{\n" +
		"\"wrap\" : true,\n" +
		"\"cursor\" : true,\n" +
		"\"guibgcolor\" : \"magenta\",\n" +
		"\"guifgcolor\" : \"cyan\",\n" +
		"\"viewbgcolor\" : \"red\",\n" +
		"\"viewfgcolor\" : \"yellow\",\n" +
		"\"selbgcolor\" : \"blue\",\n" +
		"\"selfgcolor\" : \"green\",\n" +
		"\"highlight\" : true\n" +
		"}"
	fmt.Fprintf(w, s)
	w.Flush()
	initConfig(g)
	assert.Equal(t, gocui.Attribute(gocui.ColorMagenta), gocui.Attribute(g.BgColor))
	assert.Equal(t, gocui.Attribute(gocui.ColorCyan), gocui.Attribute(g.FgColor))
	os.Remove(configPath)
	if err == nil {
		os.Rename(copyPath, configPath)
	}
}

func TestNoConfig(t *testing.T) {
	g := initGui()
	defer g.Close()
	usr, _ := user.Current()
	configPath := usr.HomeDir + "/.stretto.json"
	copyPath := usr.HomeDir + "/r9w92W2Cn7MTtAhuCP5si2LxwP9UrmC6Y99.json"
	err := os.Rename(configPath, copyPath)
	initConfig(g)
	assert.Equal(t, gocui.Attribute(gocui.ColorBlack), gocui.Attribute(g.BgColor))
	assert.Equal(t, gocui.Attribute(gocui.ColorWhite), gocui.Attribute(g.FgColor))
	if err == nil {
		os.Rename(copyPath, configPath)
	}
}
