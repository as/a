package reg

import (
	"image"
	"sync/atomic"

	"github.com/as/frame"
)

var dreg = newReg()

func PutImage(v image.Image, path string) image.Image       { Put(v, path); return v }
func PutTheme(v frame.Color, path string) frame.Color       { Put(v, path); return v }
func PutPalette(v frame.Palette, path string) frame.Palette { Put(v, path); return v }
func GetImage(path string) image.Image                      { return Get(path).(image.Image) }
func GetTheme(path string) frame.Color                      { return Get(path).(frame.Color) }
func GetPalette(path string) frame.Palette                  { return Get(path).(frame.Palette) }

func Put(img interface{}, path string) interface{} {
	dreg.Register(img, path)
	return img
}
func Get(path string) interface{} {
	return dreg.Lookup(path)
}

type reg struct {
	atomic.Value
}

func newReg() *reg {
	a := &reg{}
	a.Store(make(map[string]interface{}))
	return a
}

func (r *reg) Register(img interface{}, path string) interface{} {
	m0 := r.Load().(map[string]interface{})
	m1 := make(map[string]interface{}, len(m0))
	for k, v := range m0 {
		m1[k] = v
	}
	m1[path] = img
	r.Store(m1)
	return img
}

func (r *reg) Lookup(path string) interface{} {
	return r.Load().(map[string]interface{})[path]
}
