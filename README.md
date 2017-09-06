# A
A is my text editor. 
- Written in Go (no dependencies)
- Closely resembles the Acme and Sam text editors.
- Optimized for editing huge binary files. 
  `The underlying frame implementation (see as/frame) does not eschew null bytes. UTF-8 is not implemented at this time.`
- Standard UNIX keyboard shortcuts
  `No crazy vi/emacs tricks`
- Mouse Chords
- Command Language
  `Edit command language is 80% implemented (slow)`


![paint](a.png)

This repository will change frequently, things will break unexpectedly. See issues.

# usage
a [file ...]

# hints
To reshape the windows and columns, click on the invisible 10x10px sizer that I haven't rendered yet with the middle mouse button. Hold the button down and move the window to the location. Release the button.

# edit
- 80% of the sam command language is implemented.

Edit ,x,the standard editor is any editor,x,any editor,c,ed,

# commands
- Currently only in CWD
- Put ```[go build]``` in the tag
- Double click inside ```[```
- Middle click to execute

# look
- Right click on a string
- If its a file, it will open it
- If win32, it will also move the mouse

# mouse
```
1 Select text/sweep
1-2 Snarf (cut)
1-3 Paste
2 Execute select
3 Look select
```

# keyboard
```
^U  Delete from cursor to start of line.
^W  Delete word before the cursor.
^H  Delete character before the cursor.
^A  Move cursor to start of the line.
^E  Move cursor to end of the line.
```

# purpose
- ACME SAC doesn't run on my computer
- The solution is to create a text editor from scratch then

# future
- Live multi-client editing
- Go specific ast/compiler stuff
- Better CMD exec
- File system interface to shiny events

# see also
History of good text editors

- `The Acme User Interface for Programmers` (Rob Pike)
- `A Tutorial for the Sam Command Language` (Rob Pike)
- Plan 9 
- Inferno
- Acme SAC

[![Go Report Card](https://goreportcard.com/badge/github.com/as/a)](https://goreportcard.com/badge/github.com/as/a)