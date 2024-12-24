package dasblinken

import (
	"image"
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
	e.engine().Wait()
	for i := 0; i < len(e.engine().Leds(0)); i++ {
		e.engine().Leds(e.Opts().Channel)[i] = uint32(0x000000)
	}
	e.engine().Render()
}

// Overlays the sprite s onto the strip at it's internal location
// This will wrap around in either direction
func wrappingOverlay(e Effect, s *sprite) {
	for i := 0; i < len(s.s); i++ {
		ox := int(s.x) + i
		if ox < 0 || ox >= len(e.engine().Leds(e.Opts().Channel)) {
			continue
		}
		e.engine().Leds(e.Opts().Channel)[ox] = s.s[i].RGB_Fade(s.lum)
	}
}

type LedMatrix struct {
	leds   []rgb
	lum    []float64
	width  int
	height int
}

// 1->8
// 16->9
// 17->24
// etc... They are zigzagging
func (m *LedMatrix) setPixel(x, y int, color rgb, lum float64) {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return
	}
	var i int
	if x%2 == 1 {
		i = (m.height - y - 1) + x*m.height
	} else {
		i = y + x*m.height
	}
	m.leds[i] = color
	m.lum[i] = lum
}

func (m *LedMatrix) applyLuminosity() {
	for i := range m.leds {
		m.leds[i].r = m.leds[i].r * m.lum[i]
		m.leds[i].g = m.leds[i].g * m.lum[i]
		m.leds[i].b = m.leds[i].b * m.lum[i]
	}
}

func renderBuffer(e Effect, buffer []rgb) {
	ledCount := e.Opts().LedCount

	e.engine().Wait()
	for j := 0; j < ledCount && j < len(buffer); j++ {
		e.engine().Leds(e.Opts().Channel)[j] = buffer[j].RGB()
	}
	e.engine().Render()
}

func clearBuffer(buffer []rgb) {
	for i := 0; i < len(buffer); i++ {
		buffer[i] = rgb{0, 0, 0}
	}
}

//	image := image.NewRGBA(image.Rect(0, 0, 8, 32))

func renderImage(i *image.Image, e *Effect) {

}
