## What
Manipulate a mouse cursor

## How
Move to cursor to the point (10,10)

```
package main

import(
	"github.com/as/cursor"
	"image"
	"fmt"
)

func main(){
	if !cursor.MoveTo(image.Pt(10,10)){
		fmt.Println("failed to cursor")
	}
}
```

## Why
Acme implementation

## Todo
```
1. Plan9
2. Linux
3. Darwin
```
