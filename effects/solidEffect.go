package effects

import (
	"math/rand/v2"

	. "barnstar.com/dasblinken"
)

type SolidEffect struct {
	opts SolidEffectOpts
	ws   WSEngine
	EffectState

	color      float64
	destColour float64
	colorStep  float64
}
type SolidConfig struct {
	Name       string  `json:"name"`
	Topology   string  `json:"topology"`
	FrameDelay float64 `json:"frameDelay"`
	Palette    string  `json:"palette"`
	Mode       string  `json:"mode"`
}

type SolidEffectOpts struct {
	base EffectsOpts

	rotationFrames float64
	palletteFunc   func(float64) RGB
	rotationFunc   func(float64) float64
}

func NewSolidEffect(config SolidConfig, stripConfig StripConfig) *SolidEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))
	opts := SolidEffectOpts{
		base:           baseOpts,
		rotationFrames: config.FrameDelay,
		palletteFunc:   getPalette(config.Palette),
		rotationFunc:   getColorTransform(config.Mode),
	}
	effect := SolidEffect{}
	effect.opts = opts
	return &effect
}

func (e *SolidEffect) Engine() WSEngine {
	return e.ws
}

func (e *SolidEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *SolidEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *SolidEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *SolidEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *SolidEffect) animate(buffer []RGB) {
	e.color = e.color + e.colorStep
	if (e.colorStep > 0 && e.color >= e.destColour) ||
		(e.colorStep < 0 && e.color <= e.destColour) {
		e.pickNewDestColour()
	}

	for i := 0; i < e.opts.base.Len(); i++ {
		buffer[i] = e.opts.palletteFunc(e.color)
	}
}

func (e *SolidEffect) pickNewDestColour() {
	e.destColour = e.opts.rotationFunc(e.color)
	e.colorStep = (e.destColour - e.color) / e.opts.rotationFrames
}

func (e *SolidEffect) Run(engine WSEngine) {
	e.ws = engine

	e.color = rand.Float64()
	buffer := make([]RGB, e.opts.base.Len())
	for i := 0; i < e.opts.base.Len(); i++ {
		buffer[i] = e.opts.palletteFunc(e.color)
	}
	RenderBuffer(e, buffer)
	e.pickNewDestColour()

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer)
			RenderBuffer(e, buffer)
		})
	}
	e.Done <- true
}
