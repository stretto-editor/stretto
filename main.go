package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if e := createArgFileIfNotExists(); e != nil {
		log.Panicln(e)
	}

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}
	g.Cursor = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func usage() {
	fmt.Printf("stretto [file1]\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func createArgFileIfNotExists() (err error) {
	argsWithoutProg := os.Args[1:]
	var filename string
	for _, s := range argsWithoutProg {
		if !strings.HasPrefix(s, "-") {
			filename = s
			break
		}
	}
	if filename != "" {
		if _, err = os.Stat(filename); os.IsNotExist(err) {
			var file *os.File
			file, err = os.Create(filename)
			file.Close()
		}
	}
	return
}
