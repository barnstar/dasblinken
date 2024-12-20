package dasblinken

import (
	"image"
	"time"
)

// Runs the given render function.  Takes exactly duration ms or longer
func RenderFrame(duration time.Duration, render func()) {
	start := time.Now().UnixNano()
	render()
	now := time.Now().UnixNano()
	elapsed := time.Duration(now - start)
	time.Sleep(duration - elapsed)
}

type SpriteData []RGB

// A linear sequence of colors with a position and velocity
type LinearSprite struct {
	X    float32
	V    float32
	Lum  float64
	Data SpriteData
}

func Clear(e Effect) {
	e.Engine().Wait()
	ch := int(e.Opts().Channel)

	for i := 0; i < len(e.Engine().Leds(0)); i++ {
		e.Engine().Leds(ch)[i] = uint32(0x000000)
	}
	e.Engine().Render()
}

// Overlays the sprite s onto the strip at it's internal location
// This will wrap around in either direction
func WrappingOverlay(e Effect, s *LinearSprite) {
	ch := int(e.Opts().Channel)
	for i := 0; i < len(s.Data); i++ {
		ox := int(s.X) + i
		if ox < 0 || ox >= len(e.Engine().Leds(ch)) {
			continue
		}
		e.Engine().Leds(ch)[ox] = s.Data[i].RGB_Fade(s.Lum)
	}
}

type LedMatrix struct {
	Leds   []RGB
	Lum    []float64
	Width  int
	Height int
}

func NewLedMatrix(width, height int) *LedMatrix {
	return &LedMatrix{
		Leds:   make([]RGB, width*height),
		Lum:    make([]float64, width*height),
		Width:  width,
		Height: height,
	}
}

// 1->8
// 16->9
// 17->24
// etc... They are zigzagging
func (m *LedMatrix) SetPixel(x, y int, color RGB, lum float64) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return
	}
	var i int
	if x%2 == 1 {
		i = (m.Height - y - 1) + x*m.Height
	} else {
		i = y + x*m.Height
	}
	m.Leds[i] = color
	m.Lum[i] = lum
}

func (m *LedMatrix) ApplyLuminosity() {
	for i := range m.Leds {
		m.Leds[i].R = m.Leds[i].R * m.Lum[i]
		m.Leds[i].G = m.Leds[i].G * m.Lum[i]
		m.Leds[i].B = m.Leds[i].B * m.Lum[i]
	}
}

func RenderBuffer(e Effect, buffer []RGB) {
	count := e.Opts().Len()
	ch := int(e.Opts().Channel)

	e.Engine().Wait()
	for j := 0; j < count && j < len(buffer); j++ {
		e.Engine().Leds(ch)[j] = buffer[j].RGB()
	}
	e.Engine().Render()
}

func ClearBuffer(buffer []RGB) {
	for i := 0; i < len(buffer); i++ {
		buffer[i] = RGB{0, 0, 0}
	}
}

//	image := image.NewRGBA(image.Rect(0, 0, 8, 32))

func RenderImage(i *image.Image, e *Effect) {

}
