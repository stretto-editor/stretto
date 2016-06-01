# Keyboard Shortcuts and Commands Actions

**This topic is a draft and may contain wrong information. **
**Related version : 0.0.1**

The stretto editor contains three modes :
 * File
 * Edition
 * Command

**The last one is not yet implemented.**

## Navigation

Available in Edition and File modes.

Keypress              | Actions
----------------------| ----------------------
Up, Down, Left, Right | Moves the cursor to the corresponding direction
End, Home             | Moves the cursor to the end or the beginning of the line
PageUp, PageDown      | Moves to the previous or the next PageUp

The scroll is also available.

## Useful

Edition     | File          | Actions
----------- | ------------- | -----------
Ctrl+O      | O             | Open a file
Ctrl+Q      | Ctrl+Q        | Quit
Tab         | Tab           | Switch between modes (no visual effects)
Ctrl+T      | Ctrl+T        | Open/close command line
Ctrl+B      | B             | Open commands information
Ctrl+W      | W             | Close the current file

Escpace commands view with escape button.

## Edition

Edition     | File          | Actions
----------- | ------------- | -----------
Ctrl+S      | S             | Save
Ctrl+U      | U             | Save As
Ctrl+F      | F             | Search forward for next occurence
Ctrl+P      |               | Search and replace next occurence
Ctrl+C      |               | Copy
Ctrl+V      |               | Paste

## Command

Command    | Shortcuts       | Actions
-----------|-----------------| -----------
quit       | q!              | Quit
           | qs, sq          | Save and Quit
           | c!              | Close file
           | sc              | Save and Close
           | open            | Open file
saveas     | sa              | Save As
replaceall | repall          | Replace all occurence
setwrap    |                 | Set the Wrapper


Save and Quit needs a file in argument if no file is currently open
Save and Close needs a file in argument if no file is currently open
Open needs a file in argument
Save as needs a file in argument
Replaceall needs two arguments : what you replace followed by the replacement
SetWrap needs "true" or "false" in argument

Stay tuned for future releases ...
