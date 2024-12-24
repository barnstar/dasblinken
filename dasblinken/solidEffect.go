package dasblinken

import (
	"math/rand/v2"
)

type SolidEffect struct {
	opts SolidEffectOpts
	ws   wsEngine
	EffectControl

	colour     float64
	destColour float64
	colourStep float64
}

type SolidEffectOpts struct {
	base EffectsOpts `json:"-"`

	rotationFrames float64
	palletteFunc   func(float64) rgb `json:"-"`
	rotationFunc   func(float64) float64
}

func NewSolidEffect(opts SolidEffectOpts) *SolidEffect {
	effect := SolidEffect{}
	effect.opts = opts
	return &effect
}

func (e *SolidEffect) engine() wsEngine {
	return e.ws
}

func (e *SolidEffect) Start() error {
	return e.startEffect(e.opts.base, e)
}

func (e *SolidEffect) Stop() {
	clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.stopEffect(e.engine())
}

func (e *SolidEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *SolidEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *SolidEffect) animate(buffer []rgb) {
	e.colour = e.colour + e.colourStep
	if (e.colourStep > 0 && e.colour >= e.destColour) ||
		(e.colourStep < 0 && e.colour <= e.destColour) {
		e.pickNewDestColour()
	}

	for i := 0; i < e.opts.base.LedCount; i++ {
		buffer[i] = e.opts.palletteFunc(e.colour)
	}
}

func (e *SolidEffect) pickNewDestColour() {
	e.destColour = e.opts.rotationFunc(e.colour)
	e.colourStep = (e.destColour - e.colour) / e.opts.rotationFrames
}

func (e *SolidEffect) run(engine wsEngine) {
	e.ws = engine

	e.colour = rand.Float64()
	buffer := make([]rgb, e.opts.base.LedCount)
	for i := 0; i < e.opts.base.LedCount; i++ {
		buffer[i] = e.opts.palletteFunc(e.colour)
	}
	renderBuffer(e, buffer)
	e.pickNewDestColour()

	for e.running.Load() == true {
		doFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer)
			renderBuffer(e, buffer)
		})
	}
	e.done <- true
}
