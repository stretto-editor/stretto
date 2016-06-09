package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAutocompleteCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "o")
	AutocompleteCmd(g, v)
	assert.Equal(t, "open\n", v.Buffer(), "\"o\" should be completed by \"open\"")
}

func TestAutocompleteCmdEmpty(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "")
	AutocompleteCmd(g, v)
	assert.Equal(t, "", v.Buffer(), "An empty buffer should not be completed")
}

func TestAutocompleteNothing(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "open ")
	AutocompleteCmd(g, v)
	assert.Equal(t, "open \n", v.Buffer(), "if there is nothing behind the cursor, tab does nothing")
}

func TestAutocompleteFile(t *testing.T) {
	g := initGui()
	defer g.Close()
	filename := "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"
	os.Create(filename)
	v, _ := g.View("cmdline")
	writeInView(v, "open r9w92W2Cn7MT")
	AutocompleteCmd(g, v)
	assert.Equal(t, "open r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt\n", v.Buffer(),
		"\"r9w92W2Cn7MT\" should be completed  by \"r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt\" "+
			"(to work, you need to don't have another file or"+
			" directory beginning with \"r9w92W2Cn7MT\")")
	os.Remove(filename)
}

func TestAutocompleteDir(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "open ./r9w92W2Cn7MT")
	os.Mkdir("r9w92W2Cn7MTtAhuCP5si2LxwP9UrmC6Y99", 0777)
	AutocompleteCmd(g, v)
	assert.Equal(t, "open ./r9w92W2Cn7MTtAhuCP5si2LxwP9UrmC6Y99/\n", v.Buffer(),
		"\"r9w92W2Cn7MT\" should be completed  by \"r9w92W2Cn7MTtAhuCP5si2LxwP9UrmC6Y99/\"")
	os.Remove("r9w92W2Cn7MTtAhuCP5si2LxwP9UrmC6Y99")
}

func TestAutocompleteTooMuchArg(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "open Commands t")
	AutocompleteCmd(g, v)
	assert.Equal(t, "open Commands t\n", v.Buffer(), "open doesn't expect a third argument, nothing to complete")
}

func TestAutocompleteNoAutocompletefunction(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "quit t")
	AutocompleteCmd(g, v)
	assert.Equal(t, "quit t\n", v.Buffer(), "quit doesn't expect any argument"+
		" so it has no autocomplete function and any arg should not be completed")
}

func TestAutocompleteArgOfUnknownCommand(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "command t")
	AutocompleteCmd(g, v)
	assert.Equal(t, "command t\n", v.Buffer(), "command is not a valid command"+
		" so it has no autocomplete function and any arg should not be completed")
}

func TestAutocompleteIntersectionFile(t *testing.T) {
	g := initGui()
	defer g.Close()
	filename := "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"
	filename2 := "r9w92W2Cn7MTtAhuCP5si2LxwP9UrmC6Y99"
	os.Create(filename)
	os.Create(filename2)
	v, _ := g.View("cmdline")
	writeInView(v, "open r9w92W2Cn7MT")
	AutocompleteCmd(g, v)
	assert.Equal(t, "open r9w92W2Cn7MTtAhuCP5si2L\n", v.Buffer(),
		"\"r9w92W2Cn7MT\" should be completed  by \"r9w92W2Cn7MTtAhuCP5si2L\" "+
			"(to work, you need to don't have another file or"+
			" directory beginning with \"r9w92W2Cn7MT\")")
	os.Remove(filename)
	os.Remove(filename2)
}

func TestAutocompleteIntersectionCommand(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "s")
	AutocompleteCmd(g, v)
	assert.Equal(t, "s\n", v.Buffer(), "should not be completed because \"s\" is the intersection of \"setwrap\" and \"saveas\"")
}

func TestAutocompleteBool(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "setwrap t")
	AutocompleteCmd(g, v)
	assert.Equal(t, "setwrap true\n", v.Buffer(), "setwrap expect a boolean so t should be completed by true")
	v.Clear()
	v.SetCursor(0, 0)
	writeInView(v, "setwrap f")
	AutocompleteCmd(g, v)
	assert.Equal(t, "setwrap false\n", v.Buffer(), "setwrap expect a boolean so f should be completed by false")

	v.Clear()
	v.SetCursor(0, 0)
	// weird argument
	writeInView(v, "setwrap a")
	AutocompleteCmd(g, v)
	assert.Equal(t, "setwrap a\n", v.Buffer(), "setwrap expect a boolean and a is not a prefix of true or false")
}
