package effects

import (
	. "barnstar.com/dasblinken"
)

type ChaseEffect struct {
	opts ChaseEffectOpts
	ws   WSEngine
	EffectState

	q    float64
	o    float64
	skip int
}

type ChaseConfig struct {
	Name     string  `json:"name"`
	Topology string  `json:"topology"`
	Speed    float64 `json:"speed"`
	Palette  string  `json:"palette"`
}

type ChaseEffectOpts struct {
	Base EffectsOpts

	Speed   float64
	Palette func(float64) RGB
}

func NewChaseEffect(config ChaseConfig, stripConfig StripConfig) *ChaseEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))
	opts := ChaseEffectOpts{
		Base:    baseOpts,
		Speed:   config.Speed,
		Palette: getPalette(config.Palette),
	}
	effect := ChaseEffect{}
	effect.opts = opts
	return &effect
}

func (e *ChaseEffect) Engine() WSEngine {
	return e.ws
}

func (e *ChaseEffect) Start() error {
	return e.StartEffect(e.opts.Base, e)
}

func (e *ChaseEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *ChaseEffect) Opts() EffectsOpts {
	return e.opts.Base
}

func (e *ChaseEffect) SetStripConfig(s StripConfig) {
	e.opts.Base.StripConfig = s
}

func (e *ChaseEffect) Run(engine WSEngine) {
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

func (e *ChaseEffect) animate(buffer []RGB) {
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
		c := e.opts.Palette(w)
		buffer[i+int(e.q)] = c
	}
}
