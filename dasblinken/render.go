package dasblinken

import (
	"time"
)

// Runs the given render function.  Takes exactly duration ms or longer
func doFrame(duration time.Duration, render func()) {
	start := time.Now().UnixNano()
	render()
	now := time.Now().UnixNano()
	elapsed := time.Duration(now - start)
	time.Sleep(duration - elapsed)
}

type spriteData []rgb

// A linear sequence of colours with a position and velocity
type sprite struct {
	x   float32
	v   float32
	lum float64
	s   spriteData
}

func clear(e Effect) {
	for i := 0; i < len(e.engine().Leds(0)); i++ {
		e.engine().Leds(e.Opts().Channel)[i] = uint32(0x000000)
	}
}

// Overlays the sprite s onto the strip at it's internal location
// This will wrap around in either direction
func wrappingOverlay(e Effect, s *sprite) {
	for i := 0; i < len(s.s); i++ {
		ox := int(s.x) + i
		if ox < 0 || ox >= len(e.engine().Leds(e.Opts().Channel)) {
			continue
		}
		e.engine().Leds(e.Opts().Channel)[ox] = s.s[i].toHex(s.lum)
	}
}

type LedMatrix struct {
	leds []uint32
	width int
	height int
}

//1->8
//16->9
//17->24
// etc... They are zigzagging
func (m *LedMatrix)setPixel(x, y int, color Color, lum float64) {	
	var i int
	if x%2 == 1 {
		i = (m.height - y) + x*m.height
	} else {
		i = y + x*m.height
	}
	m.leds[i] = color.toHex(lum)
}

