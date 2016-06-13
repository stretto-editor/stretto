# Stretto
**An editor right from the 22nd century**
![cover logo](/assets/cover_1000.png)

See [GitHub Pages](http://stretto-editor.github.io/) to join the project and [Facebook](https://www.facebook.com/Stretto-841328495972117/?fref=ts) to follow our news.

## Presentation

Stretto is a text-based editor halfway between modal editors such as Emacs or Vim,
and user-friendly ones such as Atom or SublimeText. Indeed, Stretto has three modes
but also use nice shortcuts like Ctrl+s to simply save your file.

The stretto editor contains three modes :
 * File
 * Edition
 * Command
 
**You can switch mode by pressing F2, and get some help about the shortcuts by pressing F3.**

The **Edition** mode is appropriate to perform modification on files. The **File**
mode is for file visualization. You cannot perform actions which alter the file.
The **Command** mode is made to do complex actions.

Please see the [user documentation](Commands.md) for usage.

## Some features & Screenshots

#### User Interface
![ihm](https://cloud.githubusercontent.com/assets/17803473/15940756/c37d8d48-2e7d-11e6-82fd-08162dabb0d1.png)
 * 1: Main area where files are displayed
 * 2: Line for interactive action in Edit and File mode
 * 3: Error section for feedbacks
 * 4: Info view

#### Do-Undo-Redo
![undo](https://cloud.githubusercontent.com/assets/17803473/15940784/daddb526-2e7d-11e6-9d44-b9e1e0fd7cc2.png)
A representation of actions done, so that you don't have to press `ctrl+z` like a savage beast everytime you want to cancel something.

#### Theme
![theme](https://cloud.githubusercontent.com/assets/17803473/15940781/da6f5018-2e7d-11e6-9ba5-83220996e43e.png)
A personalization thanks to stretto.json

There are many more features. Just be curious !


# Installation

```bash
wget https://raw.githubusercontent.com/stretto-editor/stretto/master/install.sh
chmod +x install.sh
./install.sh
```
After install, you can move the `stretto` binary to a more appropriate destination :
```bash
sudo mv stretto /usr/bin
```
# Easy setup

After the installation you will find stretto executable and **Commands.md** file in
your home directory. If you want to move them to another directory such as
/usb/bin, keep these two files together to get documentation in the application.

You'll find a hidden configuration file in your home directory named
**stretto.json** after the installation which will allow you to configure some
features such as :

* View color
* Background color
* Visible cursor
* Highlighting of current line
* Activate wrap

***Warning : you can change the parameters but do not delete one of the
attributes in the configuration file.***


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

Stuffs that should be done :
 * Strengthen the test base (we love tests at Stretto)
 * Get rid of Gocui dependency by developing a new GUI architecture
 * Increase the gap between the buffer and the view in order to have a real
tabulation and eventually syntax highlighting.
 * For the View, split SetView and create to functions. One for creation and
initialization and an other one for updating fields.
 * Create logs and bug reports to easing the development
 * Take all actions into account in Undo/Redo
 * Use the highlighting for next occurrence or matching piece
 * Migrate to a graphical interface
 * Use a cache in case of crash
 * Separate the navigation from the cursor moves

# About
## Contribution
**Please read before contributing**
### Code structure
Stretto is developed in Golang. See https://golang.org/doc/install or for
more details.

The architecture model is MVC. Please respect this model as much as possible.

Do not hesitate to contribute and share about the code, the architecture or new features. 
**Please contact @eric-burel for any further information : [eb@lebrun-burel.com](mailto:eb@lebrun-burel.com)**

### Advice

- Please merge your work on the `dev` branch for new features (except for bug fixes). `master` is only for new versions, patches and fixes.
- Please unit test your code, **especially for bug fixes**. It would be sad to lost your work because of an unnoticed regression !
- Please add advice on how to solve the problem when opening an issue if you can, or at least as much details as possible. That would be of great help.
- Please always use English

Since we are a very small team, most pull requests will be merged without much checking. So please take your own responsibility and test your code before sending PRs ;)


## Licence

Please see the [licence](LICENCE) for more information.
