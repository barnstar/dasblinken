package dasblinken

import (
	"math"
	"math/rand/v2"
)

type StaticEffect struct {
	opts StaticEffectOpts
	ws   wsEngine
	EffectControl
}

type StaticEffectOpts struct {
	base EffectsOpts
}

func NewStaticEffect(opts StaticEffectOpts) *StaticEffect {
	effect := StaticEffect{}
	effect.opts = opts
	return &effect
}

func (e *StaticEffect) engine() wsEngine {
	return e.ws
}

func (e *StaticEffect) Start() error {
	return e.startEffect(e.opts.base, e)
}

func (e *StaticEffect) Stop() {
	clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.stopEffect(e.engine())
}

func (e *StaticEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *StaticEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *StaticEffect) animate(buffer []rgb) {
	for i := 0; i < e.opts.base.LedCount; i++ {
		if rand.Float64() < 0.6 {
			b := rand.Float64()
			b = math.Sqrt(b)
			buffer[i] = rgb{b, b, b}
		}
	}
}

func (e *StaticEffect) run(engine wsEngine) {
	e.ws = engine

	buffer := make([]rgb, e.opts.base.LedCount)
	for i := 0; i < e.opts.base.LedCount; i++ {
		b := rand.Float64()
		b = math.Sqrt(b)
		buffer[i] = rgb{b, b, b}
	}

	for e.running.Load() == true {
		doFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer)
			renderBuffer(e, buffer)
		})
	}
	e.done <- true
}
