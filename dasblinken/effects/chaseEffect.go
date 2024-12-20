package effects

import (
	. "barnstar.com/piled/dasblinken"
)

type RainbowChaseEffect struct {
	opts ChaseEffectOpts
	ws   WSEngine
	EffectState

	q    float64
	o    float64
	skip int
}

type ChaseEffectOpts struct {
	Base EffectsOpts `json:"-"`

	Speed float64 `json:"name"`
}

func NewRainbowChaseEffect(opts ChaseEffectOpts) *RainbowChaseEffect {
	effect := RainbowChaseEffect{}
	effect.opts = opts
	return &effect
}

func (e *RainbowChaseEffect) Engine() WSEngine {
	return e.ws
}

func (e *RainbowChaseEffect) Start() error {
	return e.StartEffect(e.opts.Base, e)
}

func (e *RainbowChaseEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *RainbowChaseEffect) Opts() EffectsOpts {
	return e.opts.Base
}

func (e *RainbowChaseEffect) SetStripConfig(s StripConfig) {
	e.opts.Base.StripConfig = s
}

func (e *RainbowChaseEffect) Run(engine WSEngine) {
	e.ws = engine

	buffer := make([]RGB, e.opts.Base.Len())
	for e.Running.Load() == true {
		RenderFrame(e.opts.Base.FrameTime, func() {
			e.animate(buffer)
			RenderBuffer(e, buffer)
		})
	}
	e.Done <- true
}

func (e *RainbowChaseEffect) animate(buffer []RGB) {
	ledCount := e.opts.Base.Len()
	e.q = e.q + e.opts.Speed
	if e.q > 3 {
		e.q = 0
	} else if e.q < 0 {
		e.q = 3
	}

	e.o = e.o - 0.01
	if e.o < 0 {
		e.o = 1.0
	}

	ClearBuffer(buffer)

	for i := 0; i < ledCount-3; i = i + 3 {
		w := (e.o + float64(i+int(e.q))/float64(ledCount))
		if w >= 1 {
			w = w - 1
		}
		c := RainbowPalette(w)
		buffer[i+int(e.q)] = c
	}
}
