package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/stretto-editor/gocui"
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

	initModes(g)

	g.SetLayout(layout)
	defaultLayout(g)
	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}
	g.Cursor = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		g.Close()
		log.Fatalln(err)
	}
}

func usage() {
	wiki := "Commands.md"
	fmt.Printf("Usage : \n\t stretto [file1]\n\n\n")
	flag.PrintDefaults()
	if f, err := ioutil.ReadFile(wiki); err != nil {
		fmt.Printf("\n Cannot load the documentation. Looking for %s\n", wiki)
		os.Exit(0)
	} else {
		fmt.Printf("%s", f)
	}
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
