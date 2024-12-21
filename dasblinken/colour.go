package dasblinken

import (
	"math"
)

type hsv struct {
	h float64
	s float64
	v float64
}

type rgb struct {
	r float64
	g float64
	b float64
}

type Color interface {
	toHex(dim float64) uint32
}

func (c *rgb) toHex(lum float64) uint32 {
	lum = math.Max(0, math.Min(1, lum))
	return uint32((c.g*255.0)*lum)<<16 +
		uint32((c.r*255.0)*lum)<<8 +
		uint32((c.b*255.0)*lum)
}

func (c *hsv) toHex(dim float64) uint32 {
	c2 := c.toRGB()
	return c2.toHex(dim)
}

func (color *hsv) toRGB() rgb {
	if color.h < 0 || color.h >= 360 ||
		color.s < 0 || color.s > 1 ||
		color.v < 0 || color.v > 1 {
		return rgb{}
	}
	// When 0 ≤ h < 360, 0 ≤ s ≤ 1 and 0 ≤ v ≤ 1:
	C := color.v * color.s
	X := C * (1 - math.Abs(math.Mod(color.h/60, 2)-1))
	m := color.v - C
	var Rnot, Gnot, Bnot float64
	switch {
	case 0 <= color.h && color.h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= color.h && color.h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= color.h && color.h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= color.h && color.h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= color.h && color.h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= color.h && color.h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	return rgb{Rnot + m, Gnot + m, Bnot + m}
}
