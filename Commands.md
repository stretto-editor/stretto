# Keyboard Shortcuts and Commands Actions

**This topic is a draft and may contain wrong information. **

**Related version : 0.1.0 - 07/2016**

The stretto editor contains three modes :
 * File
 * Edition
 * Command

## Navigation

F2 : Switch between Edition and File modes.
Ctrl+T : Enter in or escape from Commandline.

Available in Edition and File modes:

Keypress              | Actions
----------------------| --------------------------------------
Up, Down, Left, Right | Move the cursor to the corresponding direction.
End, Home             | Move the cursor to the end or the beginning of the line.
PageUp, PageDown      | Move to the previous or the next PageUp.
Middle mouse move     | Scroll in the corresponding direction.
F7, F8                | Switch between opened files

## Useful

Edition   | File      | Actions
--------- | --------- | --------------------------------------
F3        | F3        | Open documentation
Ctrl+D    | D         | Display the content of a directory
Ctrl+O    | O         | Open a file
Ctrl+W    | W         | Close the current file
Ctrl+Q    | Ctrl+Q    | Quit

You can escape form interactive action at anytime with ESC.

## Edition

Edition   | File      | Actions
----------| --------- | --------------------------------------
Ctrl+S    | S         | Save
Ctrl+U    | U         | Save As
Ctrl+F    | F         | Search forward for next occurence
Ctrl+P    |           | Search and replace next occurence
Ctrl+C    |           | Copy (available on Linux with xclip installed)
Ctrl+V    |           | Paste (available on Linux with xclip installed)
Ctrl+Z    |           | Undo last action
Ctrl+Y    |           | Redo last undone action
Ctrl+N    |           | Display historic of the current view (Undo/Redo stack)
Ctrl+J    |           | Permute the current line with the previous one
Ctrl+K    |           | Permute the current line with the next one

## Command

Long       | Short      | Args               | Actions
-----------|------------|------------------- | ---------------
quit       | q!         |                    | Quit
           | qs, sq     | [filename]         | Save and Quit
close      | c!         |                    | Close file
           | sc         | [filename]         | Save and Close
open       | o          | filename           | Open file
saveas     | sa         | filename           | Save As
replaceall | repall     | findStr replaceStr | Replace all occurence
setwrap    |            | true|false         | Set/disable the wrap
goto       |            | [line [column]]    | Go to the specified location

There is an autocompletion on commands for long versions.
There is also an autocompletion on directories and files for action which
required a file or a directory.

Stay tuned for future releases ...
