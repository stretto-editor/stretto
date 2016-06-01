package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/stretto-editor/gocui"
)

// create the file in the directory of the
func createFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		var file *os.File
		file, _ = os.Create(filename)
		file.Close()
		currentFile = filename
	}
}

func saveMain(v *gocui.View, filename string) error {
	if filename == "" {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	p := make([]byte, 5)
	v.Rewind()
	var size int64 = -1
	for {
		n, err := v.Read(p)
		size += int64(n)
		if n > 0 {
			if _, err = f.Write(p[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			f.Truncate(size * int64(binary.Size(p[0])))
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func openAndDisplayFile(g *gocui.Gui, filename string) error {
	v, _ := g.View("main")
	err := openFile(v, filename)
	if err == nil {
		currentFile = filename
		v.Title = filename
		return nil
	}
	return fmt.Errorf("Could not open file : %s", filename)
}

func saveAs(g *gocui.Gui, filename string) error {
	v, _ := g.View("main")
	createFile(filename)
	return saveMain(v, filename)
}
