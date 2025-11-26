package effects

import (
	"fmt"
	"math/rand/v2"

	. "barnstar.com/dasblinken"
)

type FireMatrixEffect struct {
	opts FireMatrixEffectOpts
	ws   WSEngine
	EffectState
}

type FireMatrixEffectOpts struct {
	base         EffectsOpts
	Sparking     float64
	Cooling      float64
	palletteFunc func(float64) RGB
}

type FireMatrixConfig struct {
	Name     string  `json:"name"`
	Topology string  `json:"topology"`
	Sparking float64 `json:"sparking"`
	Cooling  float64 `json:"cooling"`
	Palette  string  `json:"palette"`
}

func NewFireMatrixEffect(config FireMatrixConfig, stripConfig StripConfig) *FireMatrixEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))
	opts := FireMatrixEffectOpts{
		base:         baseOpts,
		Sparking:     config.Sparking,
		Cooling:      config.Cooling,
		palletteFunc: getPalette(config.Palette),
	}
	effect := FireMatrixEffect{}
	effect.opts = opts
	return &effect
}

func (e *FireMatrixEffect) Engine() WSEngine {
	return e.ws
}

func (e *FireMatrixEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *FireMatrixEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *FireMatrixEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *FireMatrixEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

type FireMatrixLayout struct {
	start int
	end   int
}

func (e *FireMatrixEffect) animate(buffer *LedMatrix, heat [][]float64) {

	for row := 0; row < buffer.Height; row++ {
		// Step 1.  Cool down every cell a little
		cf := 1.0
		heat := heat[row]
		for i := 0; i < buffer.Width; i++ {
			ce := rand.Float64() * e.opts.Cooling * cf
			cf += 0.03
			heat[i] = Fsub_norm(heat[i], ce)
		}

		for k := buffer.Width - 1; k > 1; k-- {
			if k > 2 {
				heat[k] = (heat[k-1] + heat[k-2] + heat[k-3]) / 3
			} else if k == 2 {
				heat[k] = (heat[k-1] + heat[k-2]) / 2
			} else if k == 1 {
				heat[k] = (heat[k-1])
			}
		}

		if rand.Float64() < e.opts.Sparking {
			y := rand.Int() % (buffer.Width / 16)
			h := rand.Float64()*0.5 + 0.5
			heat[y] = Fadd_norm(heat[y], h)
		}

		for col := 0; col < buffer.Width; col++ {
			h := heat[col]
			if col == 0 && h < 0.2 {
				h = max(h+0.2, 1.0)
			}
			c := e.opts.palletteFunc(h)
			buffer.SetPixel(col, row, c, 1.0)
		}
	}
}

func (e *FireMatrixEffect) Run(engine WSEngine) {
	e.ws = engine
	fmt.Printf("New FireMatrixEffect with width %d, height %d\n", e.opts.base.Width, e.opts.base.Height)

	buffer := NewLedMatrix(e.opts.base.Width, e.opts.base.Height)

	heat := make([]([]float64), e.opts.base.Height)
	for i := range heat {
		heat[i] = make([]float64, e.opts.base.Width)
	}

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer, heat)
			RenderBuffer(e, buffer.Leds)
		})
	}
	e.Done <- true
}
