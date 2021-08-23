# Etch
[![Go Report Card](https://goreportcard.com/badge/github.com/as/etch)](https://goreportcard.com/badge/github.com/as/etch)

Package etch provides a simple facility to write graphical regression tests.
The `Assert` function handles the common case. Give it the test variable, the
images you have and want, and it will fail your case if they differ.

# Synopsis

```
have = image.NewRGBA(r)
want = image.NewRGBA(r)

// fail the test if the images differ, and write the delta to a png
etch.Assert(t, have, want, "Test.png")
```

# Visualize

Optionally, provide a filename to store the graphical difference as an uncompressed PNG if the test fails.

![paint](img/delta.png)

The `Extra` data in the image (have but don't want) is represented in `Red`.
The `Missing` data (`want`, but dont `have`) is represented in `Blue`. 
These can be changed by modifying `Extra` and `Missing` package variables

# Example

I observed a bug in A where the text on the last line wasn't cleaned up unless that last line ended in a newline character.
This means if the frame displays `^a\nb\nc$` and `b\n` is deleted, the user would see `^a\nc\nc$`. Nasty.

We can programatically check for any defect as long as we know how to reproduce it. 

# Step 1: Find Reproduction and Expected Result
Find the reproduction. In this case I also found steps that generate the expected result. You can also use a cached expected result from a previously known good configuration.

- Insert the multi-line text containing no trailing newline (good: insert a trailing newline)

![paint](img/1.png)

- Select any line but the last 

![paint](img/2.png)

- Delete the selection

![paint](img/3.png)

Above you can see the result of the middle line's deletion for both sessions. The window that
did not have the trailing newline did not clean up the last line after copying it up toward the
top of the frame

# Step 2: Create Images

Create two images

```
	have = image.NewRGBA(r)
	want = image.NewRGBA(r)
```

Now for the test case specific stuff. Your steps will replace mine
below depending on what you're actually doing to the images. In my case
the frame draws on them directly, so we really don't care about
its inner workings too much, just that there's a bug and we're
going to test for its existence using these two images: `have`
and `want`.


```
	// Create two frames
	h = New(r, font.NewBasic(fsize), have, A)
	w = New(r, font.NewBasic(fsize), want, A)
	
	// Insert some text with and without trailing newlines
	w.Insert([]byte("1234\ncccc\ndddd\n"), 0)
	h.Insert([]byte("1234\ncccc\ndddd"), 0)
	
	// Delete the second line
	h.Delete(5, 10)
	w.Delete(5, 10)
```

By this point, `want` will be an image with the _defect-free_
state and `have` will be an image with the _defective_ state

```
	etch.Assert(t, have, want, "TestDeleteLastLineNoNL.png")
```

# Step 4: Go Test

We run `go test`

```
--- FAIL: TestDeleteLastLineNoNL (0.03s)
	etch.go:64: delta: TestDeleteLastLineNoNL.png
FAIL
exit status 1
FAIL	github.com/as/frame	0.125s
```

We can look at the image to see what went wrong: `TestDeleteLastLineNoNL.png`

![paint](img/delta.png)

Although it looks obvious, remember that this test would fail if any of the pixels differ. It's not easy to compare images visually, and you shouldn't avoid automating tests for it. Automating the tests helps prevent regressions from going undetected and speeds up the edit/compile/test cycle. 

# Step 5: Apply the Fix

```
	f.Draw(f.b, image.Rect(pt0.X, pt0.Y, pt0.X+(f.r.Max.X-pt1.X), q0), f.b, pt1, f.op)
	f.Draw(f.b, image.Rect(f.r.Min.X, q0, f.r.Max.X, q0+(q2-q1)), f.b, image.Pt(f.r.Min.X, q1), f.op)
	// f.Paint(image.Pt(pt2.X, pt2.Y-(pt1.Y-pt0.Y)), pt2, f.Color.Back)

```

The bug is the commented line above. Once the comment is removed, the test passes. Because `go test`
can be run automatically on file changes, this eliminates the manual step of checking the image. The
test passes once `have` and `want` are the same, and when they're not, just open the delta in an image
viewer to see what went wrong.




