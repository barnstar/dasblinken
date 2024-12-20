package dasblinken

import (
	"time"
)

type renderFunc func()

// Runs the given render function.  Takes exactly duration ms or longer
func doFrame(duration time.Duration, f renderFunc) {
	start := time.Now().UnixNano()
	f()
	now := time.Now().UnixNano()
	elapsed := time.Duration(now - start)
	time.Sleep(duration - elapsed)
}

type spriteData []rgb

// A linear sequence of colours with a position and velocity
type sprite struct {
	x float32
	v float32
	s spriteData
}

func clear(e Effect) {
	for i := 0; i < len(e.engine().Leds(0)); i++ {
		e.engine().Leds(0)[i] = uint32(0x000000)
	}
}

// Overlays the sprite s onto the strip at it's internal location
// This will wrap around in either direction
func wrappingOverlay(e Effect, s *sprite) {
	for i := 0; i < len(s.s); i++ {
		ox := int(s.x) + i
		if ox < 0 || ox >= len(e.engine().Leds(0)) {
			continue
		}
		e.engine().Leds(0)[ox] = s.s[i].toHex(1.0)
	}
}
