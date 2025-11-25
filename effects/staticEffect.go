package effects

import (
	"math"
	"math/rand/v2"

	. "barnstar.com/dasblinken"
)

type StaticEffect struct {
	opts StaticEffectOpts
	ws   WSEngine
	EffectState
}

type StaticEffectOpts struct {
	base EffectsOpts
}

func NewStaticEffect(opts StaticEffectOpts) *StaticEffect {
	effect := StaticEffect{}
	effect.opts = opts
	return &effect
}

func (e *StaticEffect) Engine() WSEngine {
	return e.ws
}

func (e *StaticEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *StaticEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *StaticEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *StaticEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *StaticEffect) animate(buffer []RGB) {
	for i := 0; i < e.opts.base.Len(); i++ {
		if rand.Float64() < 0.6 {
			b := rand.Float64()
			b = math.Sqrt(b)
			buffer[i] = RGB{R: b, G: b, B: b}
		}
	}
}

func (e *StaticEffect) Run(engine WSEngine) {
	e.ws = engine

	buffer := make([]RGB, e.opts.base.Len())
	for i := 0; i < e.opts.base.Len(); i++ {
		b := rand.Float64()
		b = math.Sqrt(b)
		buffer[i] = RGB{R: b, G: b, B: b}
	}

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer)
			RenderBuffer(e, buffer)
		})
	}
	e.Done <- true
}
