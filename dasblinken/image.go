package dasblinken

import "image"

func (m *LedMatrix) renderImage(i *image.RGBA) {
	cr := float64(0xffff)
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			r, g, b, _ := i.At(x, y).RGBA()
			outC := RGB{float64(r) / cr, float64(g) / cr, float64(b) / cr}
			m.SetPixel(x, y, outC, 1.0)
		}
	}
}
