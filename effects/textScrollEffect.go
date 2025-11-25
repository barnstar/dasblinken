package effects

import (
	. "barnstar.com/dasblinken"
)

type TextScrollEffect struct {
	opts TextScrollEffectOpts
	ws   WSEngine
	EffectState

	offset float64
	x      float64
}

type TextScrollEffectOpts struct {
	base EffectsOpts

	text         string
	palletteFunc func(float64) RGB `json:"-"`
}

func NewTextScrollEffect(opts TextScrollEffectOpts) *TextScrollEffect {
	effect := TextScrollEffect{}
	effect.opts = opts
	return &effect
}

func (e *TextScrollEffect) Engine() WSEngine {
	return e.ws
}

func (e *TextScrollEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *TextScrollEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *TextScrollEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *TextScrollEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *TextScrollEffect) animate(buffer *LedMatrix) {
	s := e.opts.text
	e.offset += 0.01
	if e.offset > 1.0 {
		e.offset = 0.01
	}
	e.x -= 0.3
	if e.x < -StringWidth(s) {
		e.x = float64(e.opts.base.Width)
	}
	ClearBuffer(buffer.Leds)
	cf := func(x float64, y float64) RGB {
		c := max(0.0, min(1.0, x/float64(e.opts.base.Width))) + (1.0 - e.offset)
		if c > 1.0 {
			c -= 1.0
		}
		return e.opts.palletteFunc(c)
	}
	buffer.DrawString(e.x, 0, s, cf, 1.0)
}

func (e *TextScrollEffect) Run(engine WSEngine) {
	e.ws = engine
	e.x = float64(e.opts.base.Width)
	e.offset = 0.0

	buffer := NewLedMatrix(e.opts.base.Width, e.opts.base.Height)

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer)
			RenderBuffer(e, buffer.Leds)
		})
	}
	e.Done <- true
}
