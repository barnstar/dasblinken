package effects

import (
	"fmt"
	"math"

	. "barnstar.com/dasblinken"
)

type WaveEffect struct {
	opts WaveEffectOpts
	ws   WSEngine
	EffectState

	offset float64
}

type WaveConfig struct {
	Name     string `json:"name"`
	Topology string `json:"topology"`
	Palette  string `json:"palette"`
}
type WaveEffectOpts struct {
	base EffectsOpts

	palletteFunc func(float64) RGB
}

func NewWaveEffect(config WaveConfig, stripConfig StripConfig) *WaveEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))
	opts := WaveEffectOpts{
		base:         baseOpts,
		palletteFunc: getPalette(config.Palette),
	}
	effect := WaveEffect{}
	effect.opts = opts
	return &effect
}

func (e *WaveEffect) Engine() WSEngine {
	return e.ws
}

func (e *WaveEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *WaveEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *WaveEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *WaveEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *WaveEffect) animate(buffer *LedMatrix) {
	e.offset += 0.05
	ClearBuffer(buffer.Leds)
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
		buffer.SetPixel(x, int(math.Round(y-3)), e.opts.palletteFunc(c), 0.3)
		buffer.SetPixel(x, int(math.Round(y-2)), e.opts.palletteFunc(c), 0.5)
		buffer.SetPixel(x, int(math.Round(y-1)), e.opts.palletteFunc(c), 0.7)
		buffer.SetPixel(x, int(math.Round(y)), e.opts.palletteFunc(c), 1.4)
		buffer.SetPixel(x, int(math.Round(y+1)), e.opts.palletteFunc(c), 0.7)
		buffer.SetPixel(x, int(math.Round(y-2)), e.opts.palletteFunc(c), 0.5)
		buffer.SetPixel(x, int(math.Round(y-3)), e.opts.palletteFunc(c), 0.3)
	}
	buffer.ApplyLuminosity()
}

func (e *WaveEffect) Run(engine WSEngine) {
	e.ws = engine

	fmt.Printf("New FireMatrixEffect with width %d, height %d\n", e.opts.base.Width, e.opts.base.Height)

	buffer := NewLedMatrix(e.opts.base.Width, e.opts.base.Height)

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer)
			RenderBuffer(e, buffer.Leds)
		})
	}
	e.Done <- true
}
