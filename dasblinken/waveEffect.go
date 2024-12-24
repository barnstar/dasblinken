package dasblinken

import (
	"fmt"
	"math"
)

type WaveEffect struct {
	opts WaveEffectOpts
	ws   wsEngine
	EffectControl

	offset float64
}

type WaveEffectOpts struct {
	base EffectsOpts

	palletteFunc func(float64) rgb `json:"-"`
}

func NewWaveEffect(opts WaveEffectOpts) *WaveEffect {
	effect := WaveEffect{}
	effect.opts = opts
	return &effect
}

func (e *WaveEffect) engine() wsEngine {
	return e.ws
}

func (e *WaveEffect) Start() error {
	return e.startEffect(e.opts.base, e)
}

func (e *WaveEffect) Stop() {
	clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.stopEffect(e.engine())
}

func (e *WaveEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *WaveEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *WaveEffect) animate(buffer *LedMatrix) {
	e.offset += 0.05
	clearBuffer(buffer.leds)
	for x := 0; x < e.opts.base.Width; x++ {
		o := (e.offset * 0.2) + float64(x)/float64(e.opts.base.Width)
		oi := int(o)
		c := o - float64(oi)
		y := ((math.Sin(5*float64(x)/float64(e.opts.base.Height)+e.offset+.3)*0.5 + 0.5) * float64(e.opts.base.Height))
		if y < 0 {
			y = 0
		}
		if y > float64(e.opts.base.Height-1) {
			y = float64(e.opts.base.Height - 1)
		}
		buffer.setPixel(x, int(math.Round(y-3)), e.opts.palletteFunc(c), 0.3)
		buffer.setPixel(x, int(math.Round(y-2)), e.opts.palletteFunc(c), 0.5)
		buffer.setPixel(x, int(math.Round(y-1)), e.opts.palletteFunc(c), 0.7)
		buffer.setPixel(x, int(math.Round(y)), e.opts.palletteFunc(c), 1.4)
		buffer.setPixel(x, int(math.Round(y+1)), e.opts.palletteFunc(c), 0.7)
		buffer.setPixel(x, int(math.Round(y-2)), e.opts.palletteFunc(c), 0.5)
		buffer.setPixel(x, int(math.Round(y-3)), e.opts.palletteFunc(c), 0.3)
	}
	buffer.applyLuminosity()
}

func (e *WaveEffect) run(engine wsEngine) {
	e.ws = engine

	ledCount := e.opts.base.LedCount
	fmt.Printf("New FireMatrixEffect with width %d, height %d\n", e.opts.base.Width, e.opts.base.Height)

	buffer := LedMatrix{make([]rgb, ledCount), make([]float64, ledCount), e.opts.base.Width, e.opts.base.Height}

	for e.running.Load() == true {
		doFrame(e.opts.base.FrameTime, func() {
			e.animate(&buffer)
			renderBuffer(e, buffer.leds)
		})
	}
	e.done <- true
}
