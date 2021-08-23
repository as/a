## Path
Package path simulates a current working directory in a text editing environment. It's a way to manage multiple working directories in one process.


```
// Each call to Look is similar to a user clicking on a file to open it
x := path.New("/lib/ndb/local").Look("..")

// Display in text editor
fmt.Println(x.Name()) 

// Use in os.Open() 
fmt.Println(x.Abs())


```
