package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/stretto-editor/gocui"
)

// CmdHandler is the handler used for the action of a command
type CmdHandler func(g *gocui.Gui, cmd []string) error

// AutocompleteHandler is the handler called when autocompleting command arguments
type AutocompleteHandler func(prefix string, posArg int) string

// Command is the struct describing a command
type Command struct {
	name         string
	action       CmdHandler
	minArg       int
	maxArg       int
	errMin       error
	autocomplete AutocompleteHandler
}

var commands map[string]*Command

func validateCmd(g *gocui.Gui, v *gocui.View) error {
	var err error
	if v.Name() != "cmdline" {
		panic("Cmdline is not the current view")
	}
	cmdBuff := v.Buffer()
	if cmdBuff == "" {
		return nil
	}
	cmdBuff = cmdBuff[:len(cmdBuff)-1]
	cmd := strings.Fields(cmdBuff)
	if cmdCour := commands[cmd[0]]; cmdCour != nil {
		nbArgs := len(cmd) - 1
		if cmdCour.minArg > nbArgs {
			err = cmdCour.errMin
		} else if cmdCour.maxArg < nbArgs {
			err = ErrUnexpectedArgument
		} else {
			err = cmdCour.action(g, cmd)
		}
	} else {
		err = fmt.Errorf("unknown command : \"%s\"", cmd[0])
	}
	clearView(v)
	if err == gocui.ErrQuit {
		return err
	}
	if err != nil {
		displayError(g, err)
	}
	return nil
}

//AutocompleteCmd autocomplete the current command input by completing the command itself or the argument
func AutocompleteCmd(g *gocui.Gui, v *gocui.View) error {
	cmdBuff := v.Buffer()
	if cmdBuff == "" {
		return nil
	}
	ox, _ := v.Origin()
	cx, _ := v.Cursor()
	// the prefix ends at the cursor position
	cmdBuff = cmdBuff[:ox+cx]
	// if the prefix (i.e. the word behind the cursor) is not a space or nothing
	if ox+cx == 0 || cmdBuff[ox+cx-1] == ' ' {
		return nil
	}
	cmd := strings.Fields(cmdBuff)
	posInWord := 0
	posWord := 0
	i := 0
	/* the position of the cursor in the current word is needed to know exactly
	* the prefix. The for loop determines the position of the word and
	* the position of the cursor in the word
	 */
	for {
		//if the current position is a space we go to the next position
		if cmdBuff[posWord] == ' ' {
			posWord++
			continue
		}
		// here posWord is the position of the first letter of a word
		// if the cursor is between the beginning and the end of the current word
		// we found the word in which the prefix is, we can leave
		if posWord+len(cmd[i])+1 >= ox+cx {
			posInWord = ox + cx - posWord
			break
		}
		//else we moove to the end of the word
		posWord += len(cmd[i])
		i++
	}
	//if the prefix is contained in the first word, this word is a command
	if i == 0 {
		if cmdName := GetAutocompleteCmd(cmd[0][:posInWord], i); cmdName != "" {
			writeAutocomplete(v, cmd[0], cmdName)
		}
		return nil
	}
	command := commands[cmd[0]]
	if command == nil || command.maxArg < i || command.autocomplete == nil {
		return nil
	}
	if argumentName := command.autocomplete(cmd[i][:posInWord], i); argumentName != "" {
		writeAutocomplete(v, cmd[i][:posInWord], argumentName)
	}
	return nil
}

func writeAutocomplete(v *gocui.View, prefix, word string) {
	for i, c := range word {
		if i >= len(prefix) {
			v.EditWrite(c)
		}
	}
}

// GetAutocompleteCmd returns the command beginning by the prefix in argument
func GetAutocompleteCmd(prefix string, posArg int) string {
	count := 0
	output := ""
	for cmd := range commands {
		if strings.HasPrefix(cmd, prefix) && cmd == commands[cmd].name {
			count++
			if count > 1 {
				output = intersectionString(output, cmd)
			} else {
				output = cmd
			}
		}
	}
	return output
}

// GetAutocompleteFile returns the file beginning by the prefix in argument
func GetAutocompleteFile(prefix string, posArg int) string {
	currentdir := "."
	if index := strings.LastIndex(prefix, "/"); index != -1 {
		currentdir = prefix[:index+1]
	}
	files, err := ioutil.ReadDir(currentdir)
	if err != nil {
		return ""
	}
	output := ""
	count := 0
	isDir := true
	if currentdir == "." {
		currentdir = ""
	}
	for _, file := range files {
		if strings.HasPrefix(currentdir+file.Name(), prefix) {
			count++
			if count > 1 {
				output = intersectionString(output, currentdir+file.Name())
				isDir = false
			} else {
				output = currentdir + file.Name()
			}
		}
	}
	if isDir {
		f, err := os.Open(output)
		if err != nil {
			return output
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil || !fi.Mode().IsDir() {
			return output
		}
		output += "/"
	}
	return output
}

// GetAutocompleteBoolean returns the boolean beginning by the prefix in argument
func GetAutocompleteBoolean(prefix string, posArg int) string {
	if strings.HasPrefix("true", prefix) {
		return "true"
	}
	if strings.HasPrefix("false", prefix) {
		return "false"
	}
	return ""
}

func intersectionString(s1, s2 string) string {
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}
	output := ""
	for i := 0; i < len(s1) && s1[i] == s2[i]; i++ {
		output += string(s1[i])
	}
	return output
}
