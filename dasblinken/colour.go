package dasblinken

import (
	"math"
	"math/rand/v2"
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
	RGB_Fade(dim float64) uint32
	RGB() uint32
}

func (c *rgb) RGB_Fade(lum float64) uint32 {
	lum = math.Max(0, math.Min(1, lum))
	return uint32((min(c.r, 1.0)*255.0)*lum)<<16 +
		uint32((min(c.g, 1.0)*255.0)*lum)<<8 +
		uint32((min(c.b, 1.0)*255.0)*lum)
}

func (c *rgb) Faded(lum float64) rgb {
	return rgb{c.r * lum, c.g * lum, c.b * lum}
}

func (c *rgb) RGB() uint32 {
	return uint32(min(c.r, 1.0)*255.0)<<16 +
		uint32(min(c.g, 1.0)*255.0)<<8 +
		uint32(min(c.b, 1.0)*255.0)
}

func (color *hsv) RGB() uint32 {
	return color.RGB_Fade(1.0)
}

func (color *hsv) RGB_Fade(lum float64) uint32 {
	if color.h < 0 || color.h >= 1.0 ||
		color.s < 0 || color.s > 1 ||
		color.v < 0 || color.v > 1 {
		return 0
	}
	// When 0 ≤ h < 360, 0 ≤ s ≤ 1 and 0 ≤ v ≤ 1:
	C := color.v * color.s
	X := C * (1 - math.Abs(math.Mod(6*color.h, 2)-1))
	m := color.v - C
	var Rnot, Gnot, Bnot float64
	switch {
	case 0 <= color.h && color.h < 1/6:
		Rnot, Gnot, Bnot = C, X, 0
	case 1/6 <= color.h && color.h < 2/6:
		Rnot, Gnot, Bnot = X, C, 0
	case 2/6 <= color.h && color.h < 3/6:
		Rnot, Gnot, Bnot = 0, C, X
	case 3/6 <= color.h && color.h < 4/6:
		Rnot, Gnot, Bnot = 0, X, C
	case 4/6 <= color.h && color.h < 5/6:
		Rnot, Gnot, Bnot = X, 0, C
	case 5/6 <= color.h && color.h < 1:
		Rnot, Gnot, Bnot = C, 0, X
	}
	col := rgb{Rnot + m, Gnot + m, Bnot + m}
	return col.RGB_Fade(lum)
}

func rainbowPalette(wheelPos float64) rgb {
	if wheelPos > 1.0 || wheelPos < 0 {
		return rgb{1, 1, 1}
	}

	out := rgb{}

	if wheelPos <= 0.33 {
		out.r = wheelPos * 3
		out.g = 1.0 - wheelPos*3
		out.b = 0
	} else if wheelPos <= 0.66 && wheelPos > 0.33 {
		wheelPos -= .33
		out.r = 1.0 - wheelPos*3
		out.g = 0
		out.b = wheelPos * 3
	} else if wheelPos > 0.66 {
		wheelPos -= .66
		out.r = 0
		out.g = wheelPos * 3
		out.b = 1.0 - wheelPos*3
	}
	return out
}

func greenFire(temperature float64) rgb {
	var heatcolor rgb

	// now figure out which third of the spectrum we're in:
	if temperature > 0.66 {
		// we're in the hottest third
		heatcolor.g = 1.0 - (temperature-0.66)/0.33     // full red
		heatcolor.b = 0.7 - (temperature-0.66)/0.33*0.7 // full green
		heatcolor.r = (temperature - 0.66) / 0.33       // ramp up blue

	} else if temperature > 0.33 && temperature <= 0.66 {
		// we're in the middle third
		heatcolor.g = 1.0
		heatcolor.b = (temperature - 0.33) / 0.33 * 0.7 // ramp up green
		heatcolor.r = 0

	} else {
		// we're in the coolest third
		heatcolor.g = temperature / 0.33 // ramp up red
		heatcolor.r = 0                  // no green
		heatcolor.b = 0
	}
	return heatcolor
}

func heatPalette(temperature float64) rgb {
	var heatcolor rgb

	// now figure out which third of the spectrum we're in:
	if temperature > 0.66 {
		// we're in the hottest third
		heatcolor.r = 1.0 - (temperature-0.66)/0.33     // full red
		heatcolor.g = 0.7 - (temperature-0.66)/0.33*0.7 // full green
		heatcolor.b = (temperature - 0.66) / 0.33       // ramp up blue

	} else if temperature > 0.33 && temperature <= 0.66 {
		// we're in the middle third
		heatcolor.r = 1.0
		heatcolor.g = (temperature - 0.33) / 0.33 * 0.7 // ramp up green
		heatcolor.b = 0

	} else {
		// we're in the coolest third
		heatcolor.r = temperature / 0.33 // ramp up red
		heatcolor.g = 0                  // no green
		heatcolor.b = 0                  // no blue
	}
	return heatcolor
}

func coldPalette(temperature float64) rgb {
	var heatcolor rgb

	if temperature > 0.66 {
		heatcolor.r = (temperature - 0.66) / 0.33
		heatcolor.g = (temperature - 0.66) / 0.33
		heatcolor.b = (temperature - 0.66) / 0.33
	} else if temperature > 0.33 && temperature <= 0.66 {
		heatcolor.r = 1.0 - (temperature-0.33)/0.33
		heatcolor.g = 1.0 - (temperature-0.33)/0.33
		heatcolor.b = 1.0
	} else {
		heatcolor.r = 0
		heatcolor.g = 0
		heatcolor.b = temperature / 0.33
	}
	return heatcolor
}

func randomColor(in float64) float64 {
	return rand.Float64()
}

func rotate(in float64) float64 {
	next := in + 0.01
	if next >= 1.0 {
		next = 0.01
	}
	return next
}
