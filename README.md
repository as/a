# A
A is my text editor. 

- Closely resembles the Acme and Sam text editors.
- Supports editing binary files. 
  `The underlying frame implementation (see as/frame) does not eschew null bytes. UTF-8 is not implemented at this time.`
- Standard UNIX keyboard shortcuts
  `No crazy vi/emacs tricks`
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
- The obvious solution was to create a text editor from scratch
