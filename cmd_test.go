package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretto-editor/gocui"
)

func TestValidateCmd(t *testing.T) {
	g := initGui()
	defer g.Close()

	// unauthorized calls : not from the cmdline
	v := g.Workingview()
	assert.Panics(t, func() { validateCmd(g, v) }, "Cmdline is not the current view")
}

func TestUnknownCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vError, _ := g.View("error")
	writeInView(v, "kl,sflk,f")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrUnknownCommand.Error(), "unknown command error expected")
}

func TestEmptyCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	validateCmd(g, v)
}

func TestOpenCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vError, _ := g.View("error")
	filename := "Commands.md"
	writeInView(v, "o "+filename)
	validateCmd(g, v)
	/*f, err := os.Open(filename)
	assert.Nil(t, err, err)
	content, _ := ioutil.ReadAll(f)*/
	content := getContentFile(filename)
	vMain := g.Workingview()
	assert.Equal(t, string(content)+"\n", vMain.Buffer(), "vMain should contains the content of "+filename)

	v.EditWrite('o')
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrMissingFilename.Error(), "missing argument error expected")

	writeInView(v, "o "+filename+" useless args")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrUnexpectedArgument.Error(), "unexpected argument error expected")
}

func TestCloseCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain := g.Workingview()
	writeInView(v, "c!")
	validateCmd(g, v)
	assert.Equal(t, "", vMain.Title, "Title of the main view should be empty")
	assert.Equal(t, "", vMain.Buffer(), "The buffer of the main view should be empty")

	//TODO : add an error when there is an unexpected argument
}

func TestReplaceAllCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vError, _ := g.View("error")
	// this test doesn't work because viewlines are empty so EditDelete doesn't work
	/*vMain, _ := g.View("main")
	text := " foo foo foo \n foo \n \n foo"
	expected := " bar bar bar \n bar \n \n bar\n"
	writeInView(vMain, text)
	//need to fill the viewlines
	writeInView(v, "repall foo bar")
	validateCmd(g, v)

	assert.Equal(t, expected, vMain.Buffer(), "all the words equal to the pattern should be replaced")
	*/
	//without arguments
	writeInView(v, "repall")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrMissingPattern.Error(), "missing pattern or replacement for search/replace")

	//with too many arguments
	writeInView(v, "repall 1 2 3")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrUnexpectedArgument.Error(), "unexpected third argument")
}

func TestSetWrapCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	// vMain, _ := g.View("main")
	vMain := g.Workingview()
	vError, _ := g.View("error")

	writeInView(v, "setwrap true")
	validateCmd(g, v)
	assert.Equal(t, true, vMain.Wrap, "wrap should be true")

	writeInView(v, "setwrap false")
	validateCmd(g, v)
	assert.Equal(t, false, vMain.Wrap, "wrap sould be false")

	writeInView(v, "setwrap ")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrWrapArgument.Error(), "missing argument error")

	writeInView(v, "setwrap useless args")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrUnexpectedArgument.Error(), "unexpected argument error")
}

func TestQuitCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	writeInView(v, "q!")
	err := validateCmd(g, v)
	assert.EqualError(t, err, gocui.ErrQuit.Error(), "Errquit should be returned from validatecmd when q! is executed")
}

func TestQuitAndSaveCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	// vMain, _ := g.View("main")
	vMain := g.Workingview()
	writeInView(v, "qs")
	vMain.Title = "6u8Y73wHm5QWmgRPcXk96y39cL.txt"

}

func TestSaveAsCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	// vMain, _ := g.View("main")
	vMain := g.Workingview()
	vError, _ := g.View("error")
	text := "This is a \n test on two lines"
	filename := "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"
	// First create the file
	writeInView(vMain, text)
	writeInView(v, "sa "+filename)
	validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")

	text2 := " another text"
	text += text2
	writeInView(vMain, text2)
	// Then write in the existing file
	writeInView(v, "sa "+filename)
	validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")

	os.Remove(filename)

	writeInView(v, "sa ")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrMissingFilename.Error(), "missing filename error")

	writeInView(v, "sa useless args")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrUnexpectedArgument.Error(), "missing argument error")
}

func TestSaveAndCloseCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	// vMain, _ := g.View("main")
	vMain := g.Workingview()
	vError, _ := g.View("error")
	vMain.Title = ""
	text := "This is a \n test on two lines"
	filename := "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"
	//save with no current file
	writeInView(vMain, text)
	writeInView(v, "sc "+filename)
	validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")
	_, err1 := g.View(filename)
	assert.Equal(t, err1, gocui.ErrUnknownView, "view not distroyed")

	//save with a current file
	vMain = g.Workingview()
	vMain.Title = filename
	text = "I'm trying to save \n and close an opened file"
	writeInView(vMain, text)
	writeInView(v, "sc")
	validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content")
	_, err2 := g.View(filename)
	assert.Equal(t, err2, gocui.ErrUnknownView, "view not distroyed")
	os.Remove(filename)

	//try to save without a current file name and without an argument
	vMain.Title = ""
	writeInView(v, "sc")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrMissingFilename.Error(), "missing filename error")
}

func TestSaveAndQuitCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain := g.Workingview()
	vError, _ := g.View("error")

	text := "This is a \n test on two lines"
	filename := "r9w92W2Cn7MTtAhuCP5si2LH356r8FrjV.txt"
	//save with no current file
	writeInView(vMain, text)
	writeInView(v, "sq "+filename)
	err := validateCmd(g, v)
	assert.Equal(t, text, getContentFile(filename), "the save file doesn't contain the right content ")
	assert.EqualError(t, err, gocui.ErrQuit.Error(), "Errquit should be returned from validatecmd when sq is executed")
	os.Remove(filename)

	//try to save without a current file name and without an argument
	vMain.Title = ""
	writeInView(v, "sq")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrMissingFilename.Error(), "missing filename error ")
}

func TestGoToCmd(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("cmdline")
	vMain := g.Workingview()
	vError, _ := g.View("error")
	filename := "Commands.md"
	//basic conditions
	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
	writeInView(vMain, getContentFile(filename))
	writeInView(v, "goto 10 10")
	validateCmd(g, v)
	//viewlines == nil --> moveDown does nothing
	cx, _ := vMain.Cursor()
	//cx, cy := vMain.Cursor()
	//assert.Equal(t, 10, cy, "the y coordonate of the cursor should be 10 "+vMain.Title)
	assert.Equal(t, 10, cx, "the x coordonate of the cursor should be 10")

	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
	//outside of the view, need to move down
	writeInView(v, "goto 40")
	validateCmd(g, v)
	cx, _ = v.Cursor()
	//cx, cy = v.Cursor()
	//viewlines == nil --> moveDown does nothing
	//assert.Equal(t, cy, 40, "the y coordonate of the cursor should be 10 "+vMain.Title)
	assert.Equal(t, cx, 0, "the x coordonate of the cursor should be 10")

	//invalide parameter
	writeInView(v, "goto a")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrNumberExpected.Error(), "number expected error should be raised")
	//invalide parameter
	writeInView(v, "goto 10 b")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrNumberExpected.Error(), "number expected error should be raised")

	vMain.Wrap = true
	writeInView(v, "goto 10 10")
	validateCmd(g, v)
	assert.Contains(t, vError.Buffer(), ErrGoToInWrapMode.Error(), "wrap is not allowed with goto")
}

func writeInView(v *gocui.View, s string) {
	for _, c := range s {
		v.EditWrite(c)
	}
}

func getContentFile(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	content, _ := ioutil.ReadAll(f)
	return string(content)
}
