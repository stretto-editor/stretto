package main

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretto-editor/gocui"
)

func initGui() *gocui.Gui {
	g := gocui.NewGui()
	g.Init()
	initModes(g)
	layout(g) // ! instead of g.SetLayout(layout)
	// since we do not enter gui's mainloop in any test
	initKeybindings(g)
	return g
}

func TestInitMode(t *testing.T) {

	g := gocui.NewGui()
	g.Init()
	layout(g)
	defer g.Close()

	expectedMode := []string{
		"file",
		"edit",
		"cmd",
	}

	initModes(g)

	for _, m := range expectedMode {
		_, e := g.Mode(m)
		if assert.Nil(t, e) {
			// should assert presence of the
			// functions associated with
			// entering and leaving modes
		}
	}
}

func TestDoSwitchMode2(t *testing.T) {
	var e error

	os.Args = []string{"main"}

	g := initGui()
	defer g.Close()

	// default behaviour
	e = doSwitchMode(g, "cmd")
	assert.NoError(t, e)
	m, _ := g.Mode("cmd")
	assert.Equal(t, g.CurrentMode(), m, "current mode should be cmd")

	// the rest : specific behaviours
	v, _ := g.View("cmdline")
	assert.Equal(t, g.CurrentView(), v, "current view should be cmdline")

	e = doSwitchMode(g, "file")
	assert.NoError(t, e)
	v, _ = g.View("main")
	assert.Equal(t, g.CurrentView(), v, "current view should be main")
	if assert.NotNil(t, g.CurrentView()) {
		assert.Equal(t, g.CurrentView().Editable, false, "current view should not be editable")
	}

	e = doSwitchMode(g, "edit")
	assert.NoError(t, e)
	v, _ = g.View("main")
	assert.Equal(t, g.CurrentView(), v, "current view should be main")
	if assert.NotNil(t, g.CurrentView()) {
		assert.Equal(t, g.CurrentView().Editable, true, "current view should be editable")
	}

	// unknown mode
	e = doSwitchMode(g, "notaknowmode")
	if assert.Error(t, e, "an error was expected") {
		assert.Equal(t, e, gocui.ErrUnknowMode)
	}
}

func TestDoSetTopView(t *testing.T) {
	var e error

	g := initGui()
	defer g.Close()

	e = doSetTopView(g, "cmdline")
	v, _ := g.View("cmdline")
	assert.NoError(t, e)
	assert.Equal(t, g.CurrentView(), v)

	e = doSetTopView(g, "notaknownview")
	if assert.Error(t, e, "an error was expected") {
		assert.Equal(t, e, gocui.ErrUnknownView)
	}
}

func TestValidateInput(t *testing.T) {
	var e error
	var v *gocui.View

	g := initGui()
	defer g.Close()

	// default behaviour :
	// 1 w/o a next demon
	v, _ = g.View("inputline")
	emptyDemon := func(g *gocui.Gui, input string) (demonInput, error) {
		return nil, nil
	}
	currentDemonInput = emptyDemon
	g.SetCurrentView("inputline")
	e = validateInput(g, v)
	assert.NoError(t, e)
	v, _ = g.View("inputline")
	assert.NotEqual(t, g.CurrentView(), v, "current view should not be the inputline")

	// 2 with a next demon
	// thats returning a error
	emptyDemon2 := func(g *gocui.Gui, input string) (demonInput, error) {
		return emptyDemon, errors.New("this is an error")
	}
	currentDemonInput = emptyDemon2
	g.SetCurrentView("inputline")
	e = validateInput(g, v)
	assert.NoError(t, e)
	v, _ = g.View("inputline")
	assert.Equal(t, g.CurrentView(), v, "current view should be the inputline")
	e = validateInput(g, v)
	assert.NoError(t, e, "this error should be handled in validateInput")
	assert.NotEqual(t, g.CurrentView(), v, "current view should not be the inputline")

	// 3 only gocui.ErrQuit should not be handled
	escapeDemon := func(g *gocui.Gui, input string) (demonInput, error) {
		return nil, gocui.ErrQuit
	}
	currentDemonInput = escapeDemon
	g.SetCurrentView("inputline")
	e = validateInput(g, v)
	if assert.Error(t, e, "an error was expected") {
		assert.Equal(t, e, gocui.ErrQuit)
	}

	// unauthorized calls :
	// 1 not from the inputline
	v, _ = g.View("main")
	assert.Panics(t, func() { validateInput(g, v) }, "Inputline is not the current view")

	// 2 no function to use the input
	v, _ = g.View("inputline")
	currentDemonInput = nil
	assert.Panics(t, func() { validateInput(g, v) }, "No Current Demon Input Available")
}

func TestDoEscapeInput(t *testing.T) {
	var v *gocui.View

	g := initGui()
	defer g.Close()

	// default behaviour
	g.SetCurrentView("inputline")
	emptyDemon := func(g *gocui.Gui, input string) (demonInput, error) {
		return nil, nil
	}
	currentDemonInput = emptyDemon
	v, _ = g.View("inputline")
	doEscapeInput(g, v)
	v, _ = g.View("main")
	assert.Equal(t, g.CurrentView(), v, "current view should be the main view")
	assert.Nil(t, currentDemonInput, "there should not be any demon waiting")

	// unauthorized calls :
	// 1 not from the inputline
	v, _ = g.View("main")
	assert.Panics(t, func() { doEscapeInput(g, v) }, "Inputline is not the current view")

	// 2 no function to use the input
	v, _ = g.View("inputline")
	currentDemonInput = nil
	assert.Panics(t, func() { doEscapeInput(g, v) }, "No Current Demon Input Available")
}

func TestDoCursorInfo(t *testing.T) {
	g := initGui()
	defer g.Close()

	g.SetCurrentView("main")
	g.CurrentView().SetCursor(4, 5)
	inMainView, _, _ := cursorInfo(g)
	//rx, ry := g.CurrentView().Cursor()

	assert.Equal(t, inMainView, true, "Current view should be main")
	//assert.Equal(t, x, rx, "X value should be equal")
	//assert.Equal(t, y, ry, "Y value should be equal")

	g.SetCurrentView("inputline")
	inMainView, _, _ = cursorInfo(g)

	assert.Equal(t, inMainView, false, "Current view shouldn't be main")

}

func TestDoInfoView(t *testing.T) {
	var v2 *gocui.View

	g := initGui()
	defer g.Close()

	v := g.CurrentView()
	e := commandInfoHandler(g, v)
	assert.Nil(t, e)

	v = g.CurrentView()
	v2, e = g.View("cmdinfo")

	assert.Equal(t, v, v2, "The view should be command info")
	assert.Nil(t, e)

	e = quitInfo(g, v)
	assert.Nil(t, e)

	_, e = g.View("cmdinfo")
	assert.Equal(t, e, errors.New("unknown view"), "No view should be found")

}

func TestDoQuitHandler(t *testing.T) {
	g := initGui()
	defer g.Close()
	v, _ := g.View("inputline")
	currentFile = "unknownfile"

	// espace with an empty input
	e := quitHandler(g, v)
	assert.Nil(t, e, "No error should be found")

	e = validateInput(g, v)
	assert.Nil(t, e, "Input shoud be valid")

	currentFile = ""

	e = quitHandler(g, v)
	assert.Nil(t, e, "No error should be found")
	e = validateInput(g, v)
	assert.Nil(t, e, "Input shoud be valid")

	e = validateInput(g, v)
	assert.NotNil(t, e, "Input already left")
}

func TestDoCopy(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}
	g := initGui()
	defer g.Close()
	teststring := "testinput"
	v, e := g.View("main")

	c := exec.Command("xclip", "-i")
	c.Stdin = strings.NewReader(teststring)
	c.Start()

	e = copy(g, v)
	assert.Nil(t, e)

	out, _ := exec.Command("xclip", "-o", "-selection", "c").Output()
	assert.Equal(t, string(out), teststring, "Should be equal")

}

func TestDoPaste(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}

	g := initGui()
	v, _ := g.View("main")
	defer g.Close()
	teststring := "testinput"

	c := exec.Command("xclip", "-i", "-selection", "c")
	c.Stdin = strings.NewReader(teststring)

	paste(g, v)
	assert.Equal(t, teststring+"\n", v.Buffer(), "Content shoud be the same")
}

func TestDoSaveAsHandler(t *testing.T) {
	// possible errors of called functions already tested in test_cmd
	g := initGui()
	defer g.Close()

	v, e := g.View("inputline")
	assert.Nil(t, e, "There should be no error")
	e = saveAsHandler(g, v)
	assert.Nil(t, e, "There should be no error")
	v.EditWrite('a')

	e = validateInput(g, v)
	assert.Nil(t, e, "There should be no error")
}

func TestDoOpenHandler(t *testing.T) {
	// possible errors of called functions already tested in test_cmd
	g := initGui()
	defer g.Close()

	v, _ := g.View("inputline")
	openFileHandler(g, v)
	v.EditWrite('a')

	validateInput(g, v)

}

func TestDoSaveHandler(t *testing.T) {
	// possible errors of called functions already tested in test_cmd
	g := initGui()
	defer g.Close()

	currentFile = ""

	v, _ := g.View("inputline")
	saveHandler(g, v)
	v.EditWrite('c')

	currentFile = "c"
	v, _ = g.View("main")
	v.EditWrite('k')
	v, _ = g.View("inputline")
	saveHandler(g, v)

	validateInput(g, v)

}
