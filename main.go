package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

func main() {
	flag.Usage = usage
	flag.Parse()

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
	fmt.Printf("stretto [file1]")
	flag.PrintDefaults()
	os.Exit(1)
}
