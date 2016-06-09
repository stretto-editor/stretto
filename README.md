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
wget .. && install.sh
sudo mv ~/stretto ~/Commands.md /usr/bin/
```

# Easy setup


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
tabulation and eventually and syntax highlighting.
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
