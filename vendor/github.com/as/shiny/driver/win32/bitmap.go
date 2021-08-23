package win32

const (
	BiRGB        = 0
	BiBitfields  = 3
	DibRGBColors = 0
	AcSrcOver    = 0x00
	AcSrcAlpha   = 0x01
	SrcCopy      = 0x00cc0020
)
const (
	GmCompatible = 1
	GmAdvanced   = 2
	MWTIdentity  = 1
)

type BitmapInfo struct {
	Header BitmapInfoV4
	Colors [1]RGBQuad
}

type BitmapInfoV4 struct {
	Size                    uint32
	Width                   int32
	Height                  int32
	Planes                  uint16
	BitCount                uint16
	Compression             uint32
	SizeImage               uint32
	XPelsPerMeter           int32
	YPelsPerMeter           int32
	ClrUsed                 uint32
	ClrImportant            uint32
	Red, Green, Blue, Alpha uint32
	Endpoints               [3]uint32
	Gamma                   [3]uint32
}

type BitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}
