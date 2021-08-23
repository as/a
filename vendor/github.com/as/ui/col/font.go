package col

import "github.com/as/font"

type facer interface {
	Face() font.Face
	SetFont(font.Face)
}
