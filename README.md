# Stretto
**An editor right from the 22nd century**

See [GitHub Pages](http://stretto-editor.github.io/) to join the project and [Facebook](https://www.facebook.com/Stretto-841328495972117/?fref=ts) to follow our news.

## Presentation

Stretto is a text-based editor halfway beetween modal editors such as Emacs or Vim,
and user-friendly ones such as Atom or SublimeText. Indeed, Stretto has tree modes
but also use nice shortcuts like Ctrl+s to simply save your file.

The stretto editor contains three modes :
 * File
 * Edition
 * Command

The **Edition** mode is appropriate to perform modification on files. The **File**
mode is for file visualization. You cannot perform actions which alter the file.
The **Command** mode is made to do complex actions.

Please see the [user documentation](Commands.md) for usage.

## Features & Screeshots


# Installation

```
wget https://raw.githubusercontent.com/stretto-editor/stretto/master/install.sh
chmod +x install.sh
./install.sh
```

# Easy setup

After the installation you will find stretto executable and Commands.md file in
your home directory. If you want to move them to another directory such as
/usb/bin, keep these two files together to get documentation in the application.

You'll find a hidden configuration file in your home directory named
stretto.json after the installation which will allow you to configure some
features such as :

* View color
* Background color
* Visible cursor
* Highlighting of current line
* Activate wrap

Warning : you can change the parameters but do not delete one of the
attributes in the configuration file.


# Road-Map

Ideas for futures releases:
 * Syntax highlighting
 * Screen split for multi-files visualization
 * Autocompletion on words
 * Search with regex
 * Configurable indentation
 * Multi-views on the same file
 * Minimap of the current file
 * (Un)comment lines

# About
## Contribution

Stretto is developped in Golang. See https://golang.org/doc/install or for
more details.

The architecture model is MVC. Please respect this model as much as possible.

TODO :
 * Increase the gap between the buffer and the view in order to have a real
tabulation and eventually syntax highlighting.
 * For the View, split SetView and create to functions. One for creation and
initialization and an other one for updating fields.
 * Create logs and bug reports to easing the developpement
 * Take all actions into account in Undo/Redo
 * Use the highlighting for next occurence or matching piece
 * Migrate to a graphical interface
 * Use a cache in case of crash
 * Separate the navigation from the cursor moves

Do not hesitate to contribute and share about the code, the architecture or new features.

## Licence

Please see the [licence](LICENCE) for more information.
