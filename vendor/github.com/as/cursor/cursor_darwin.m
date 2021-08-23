
// +build darwin
// +build 386 amd64
// +build !ios

#include "_cgo_export.h"
#include <pthread.h>
#include <stdio.h>

#import <Cocoa/Cocoa.h>
#import <Foundation/Foundation.h>

NSPoint
scalepoint(NSPoint pt, CGFloat scale)
{
    pt.x *= scale;
    pt.y *= scale;
    return pt;
}

CGFloat
pixelScale()
{
	NSWindow *win = [[[NSApplication sharedApplication] windows] objectAtIndex:0];
	NSSize size, ptsize;

	// Compute factors
	ptsize = [win.contentView bounds].size;
	size = [win.contentView convertSizeToBacking: ptsize];
	return size.width/ptsize.width;	
}
CGFloat
pointScale()
{
	return 1.0f / pixelScale();
}

NSRect
winbounds()
{
	NSWindow *win = [[[NSApplication sharedApplication] windows] objectAtIndex:0];
	NSRect r = win.frame;
	//CGFloat s = pixelScale();
	//r.origin.x *= s;
	//r.origin.y *= s;
	//r.size.width *= s;
	//r.size.height *= s;
	return r;
}

NSPoint
goodY(NSPoint pt )
{
	pt.y = 1800 - pt.y;
	return pt;
}

NSPoint
setmouse(float x, float y)
{
	NSPoint q;
	q.x=x;
	q.y=y;
	CGWarpMouseCursorPosition(NSPointToCGPoint(q));
	CGAssociateMouseAndMouseCursorPosition(true);

    return q;
}

NSRect
bounds()
{
	NSRect r = [[[NSScreen screens] objectAtIndex:0] frame];
	CGFloat s = pixelScale();
	r.origin.x *= s;
	r.origin.y *= s;
	r.size.width *= s;
	r.size.height *= s;
	return r;
}

void
acmesetmouse(int x, int y)
{
	NSWindow *win = [[[NSApplication sharedApplication] windows] objectAtIndex:0];
	NSSize size, ptsize;
	NSPoint q, pt, mpos;
	pt.x = x;
	pt.y = y;

	// Compute factors
	ptsize = [win.contentView bounds].size;
	size = [win.contentView convertSizeToBacking: ptsize];
	CGFloat topixels = size.width/ptsize.width;
	CGFloat topoints = 1.0f / topixels;

	mpos = scalepoint(NSMakePoint(pt.x, pt.y), topoints);

	q = [win.contentView convertPoint:mpos toView:nil];
	q = [win convertRectToScreen:NSMakeRect(q.x, q.y, 0, 0)].origin;

	NSRect r = bounds();
	//q.y = r.size.height - q.y;

	setmouse(q.x, q.y);
}
