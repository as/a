package cursor

import (
	"fmt"
	"image"
	"os/exec"
)

func widSelf() string {
	out, _ := exec.Command("xdotool", "getactivewindow").Output()
	if len(out) == 0{
		return ""
	}
	if out[len(out)-1] == '\n'{
		out=out[:len(out)-1]
	}
	return string(out)
}

func MoveRelativeTo(p image.Point) {
	id := widSelf()
	p = p.Add(image.Pt(6,6))
	exec.Command("xdotool", "mousemove", "-w", id, fmt.Sprint(p.X), fmt.Sprint(p.Y)).Run()
}

func MoveTo(p image.Point) bool {
	return moveTo(p)
}

func moveTo(p image.Point) bool {
	return false
}
