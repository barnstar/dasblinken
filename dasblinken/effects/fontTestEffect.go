package effects

import (
	. "barnstar.com/piled/dasblinken"
)

type FontTestEffect struct {
	opts FontTestEffectOpts
	ws   WSEngine
	EffectState

	offset float64
	x      float64
	s      string
}

type FontTestEffectOpts struct {
	base EffectsOpts

	palletteFunc func(float64) RGB `json:"-"`
}

func NewFontTestEffect(opts FontTestEffectOpts) *FontTestEffect {
	effect := FontTestEffect{}
	effect.opts = opts
	return &effect
}

func (e *FontTestEffect) Engine() WSEngine {
	return e.ws
}

func (e *FontTestEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *FontTestEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *FontTestEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *FontTestEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *FontTestEffect) animate(buffer *LedMatrix) {
	s := e.s
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

func (e *FontTestEffect) Run(engine WSEngine) {
	e.ws = engine
	e.x = float64(e.opts.base.Width)
	e.offset = 0.0
	e.s = ""
	for i := 0; i < 256; i++ {
		e.s += string(i)
	}

	buffer := NewLedMatrix(e.opts.base.Width, e.opts.base.Height)

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer)
			RenderBuffer(e, buffer.Leds)
		})
	}
	e.Done <- true
}
