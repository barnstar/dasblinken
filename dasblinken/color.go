package dasblinken

import (
	"math"
	"math/rand/v2"
)

type HSV struct {
	h float64
	s float64
	v float64
}

type RGB struct {
	R float64
	G float64
	B float64
}

type Color interface {
	RGB_Fade(dim float64) uint32
	RGB() uint32
}

func (c *RGB) RGB_Fade(lum float64) uint32 {
	lum = math.Max(0, math.Min(1, lum))
	return uint32((min(c.R, 1.0)*255.0)*lum)<<16 +
		uint32((min(c.G, 1.0)*255.0)*lum)<<8 +
		uint32((min(c.B, 1.0)*255.0)*lum)
}

func (c *RGB) Faded(lum float64) RGB {
	return RGB{c.R * lum, c.G * lum, c.B * lum}
}

func (c *RGB) RGB() uint32 {
	return uint32(min(c.R, 1.0)*255.0)<<16 +
		uint32(min(c.G, 1.0)*255.0)<<8 +
		uint32(min(c.B, 1.0)*255.0)
}

func (color *HSV) RGB() uint32 {
	return color.RGB_Fade(1.0)
}

func (color *HSV) RGB_Fade(lum float64) uint32 {
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
	col := RGB{Rnot + m, Gnot + m, Bnot + m}
	return col.RGB_Fade(lum)
}

func RainbowPalette(wheelPos float64) RGB {
	if wheelPos > 1.0 || wheelPos < 0 {
		return RGB{1, 1, 1}
	}

	out := RGB{}

	if wheelPos <= 0.33 {
		out.R = wheelPos * 3
		out.G = 1.0 - wheelPos*3
		out.B = 0
	} else if wheelPos <= 0.66 && wheelPos > 0.33 {
		wheelPos -= .33
		out.R = 1.0 - wheelPos*3
		out.G = 0
		out.B = wheelPos * 3
	} else if wheelPos > 0.66 {
		wheelPos -= .66
		out.R = 0
		out.G = wheelPos * 3
		out.B = 1.0 - wheelPos*3
	}
	return out
}

func GreenFire(temperature float64) RGB {
	var heatcolor RGB

	// now figure out which third of the spectrum we're in:
	if temperature > 0.66 {
		// we're in the hottest third
		heatcolor.G = 1.0 - (temperature-0.66)/0.33     // full red
		heatcolor.B = 0.7 - (temperature-0.66)/0.33*0.7 // full green
		heatcolor.R = (temperature - 0.66) / 0.33       // ramp up blue

	} else if temperature > 0.33 && temperature <= 0.66 {
		// we're in the middle third
		heatcolor.G = 1.0
		heatcolor.B = (temperature - 0.33) / 0.33 * 0.7 // ramp up green
		heatcolor.R = 0

	} else {
		// we're in the coolest third
		heatcolor.G = temperature / 0.33 // ramp up red
		heatcolor.R = 0                  // no green
		heatcolor.B = 0
	}
	return heatcolor
}

func HeatPalette(temperature float64) RGB {
	var heatcolor RGB

	// now figure out which third of the spectrum we're in:
	if temperature > 0.66 {
		// we're in the hottest third
		heatcolor.R = 1.0 - (temperature-0.66)/0.33     // full red
		heatcolor.G = 0.7 - (temperature-0.66)/0.33*0.7 // full green
		heatcolor.B = (temperature - 0.66) / 0.33       // ramp up blue

	} else if temperature > 0.33 && temperature <= 0.66 {
		// we're in the middle third
		heatcolor.R = 1.0
		heatcolor.G = (temperature - 0.33) / 0.33 * 0.7 // ramp up green
		heatcolor.B = 0

	} else {
		// we're in the coolest third
		heatcolor.R = temperature / 0.33 // ramp up red
		heatcolor.G = 0                  // no green
		heatcolor.B = 0                  // no blue
	}
	return heatcolor
}

func ColdPalette(temperature float64) RGB {
	var heatcolor RGB

	if temperature > 0.66 {
		heatcolor.R = (temperature - 0.66) / 0.33
		heatcolor.G = (temperature - 0.66) / 0.33
		heatcolor.B = (temperature - 0.66) / 0.33
	} else if temperature > 0.33 && temperature <= 0.66 {
		heatcolor.R = 1.0 - (temperature-0.33)/0.33
		heatcolor.G = 1.0 - (temperature-0.33)/0.33
		heatcolor.B = 1.0
	} else {
		heatcolor.R = 0
		heatcolor.G = 0
		heatcolor.B = temperature / 0.33
	}
	return heatcolor
}

func RandomColor(in float64) float64 {
	return rand.Float64()
}

func Rotate(in float64) float64 {
	next := in + 0.01
	if next >= 1.0 {
		next = 0.01
	}
	return next
}
