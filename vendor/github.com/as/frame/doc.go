// Package frame provides plan9-like editable text images on a raster display. This implementation
// preserves NUL bytes, and uses a set of replacement characters for unrenderable text glyphs generated
// with a smaller sized font (hexadecimal or ascii representation).
//
// A frame's text is not addressable. Once the characters are written to the frame, there is no
// mechanism to retrieve their position from within the frame. Use a buffer to store text for reading
// and the range addresses of the frame to access bytes from that buffer.
//
// See github.com/as/ui/win for an example.
//
// A frame is created using the New function
//
//   img := image.NewRGBA(image.Rect(0,0,100,100))
//   fr := frame.New(img, img.Bounds(), frame.NewGoMono(), frame.Mono)
//
// A frame supports these common operations
//   Insert: Insert text
//   Delete: Delete text
//   IndexOf: Index for point
//   PointOf: Point for index
//   Select: Select range
//   Dot: Return selected range
//
// Insert and Delete
//
// Frames supports two operations for rendering text: Insert and Delete. Insert inserts text at the
// given index and moves existing characters after the index to the right. Delete deletes text in the
// given range (a range is a pair of indices) and moves existing character after the index to the
// left.
//
// The two operations are inverses of each other.
//
//   fr.Insert([]byte("hello world."), 0)
//   fr.Delete(0, 11)
//
// Insert and delete return the number of characters inserted or deleted.
//
// To delete the last insertion:
//   p0 := 0
//   n := fr.Insert([]byte("123"), p0)
//   fr.Delete(p0, p0+n)
//
// To execute a traditional "write" operation:
//   s := []byte("hello")
//   fr.Delete(0, int64(len(s)))
//   fr.Insert(s, 0)
//
// Projection
//
// Frames can translate between coordinates of the mouse and character offsets in the frame itself using
// IndexOf and PointOf.
//
//   p0  := fr.IndexOf(image.Pt(0, 0)) // Returns the index under the 2D point (0,0)
//   pt0 := fr.PointOf(5) // Returns the 2D point over the index
//
// Selection
//
// Frames support selecting ranges of text along with returning those selected ranges.
//
//   fr.Select(p0, p1)
//   fr.Dot()
//
// A more complicated facility exists for making a live selection. See example/basic for an example of
// how to use it.
//
//  fr.Sweep(...)
//
// Drawing
//
// No special operations are needed after a call to Insert, Delete, or Select. The frame's bitmap
// is updated. However, there are four functions that will redraw the frame on the bitmap if
// this is necessary.
//
// Recolor(pt image.Point, p0, p1 int64, cols Palette)
//   Recolor colors the range p0:p1 by redrawing the foreground, background, and font glyphs
//
// Redraw(pt image.Point, p0, p1 int64, issel bool)
//   Redraw redraws the characters between p0:p1. It accesses the cache of drawn glyph widths
//   to avoid remeasuring strings
//
// RedrawAt(pt image.Point, text, back image.Image)
//   RedrawAt refreshes the entire image to the right of the given pt. Everything below is redrawn.
//
// Refresh()
//   Refresh recomputes the state of the frame from scratch. This is an expensive operation compared
//   to redraw
//
// Display Sync
//
// After any operation that alters the frame, one can be sure that the changes can be written to
// the frame's bitmap. However, the same can not be said for the exp/shiny window. There currently
// exists an optimization (see github.com/as/drawcache) that caches rectangles that need to be
// redrawn to the screen. This is because shiny (or the native drivers for it) are too slow to
// refresh the entire window is that window's resolution is very high.
//
//
// This rendering pipeline is bottlenecked, so an optimization is located between the |*|
//
// insert | frame | shiny buffer |*| shiny window
//
//
package frame
