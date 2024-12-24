package dasblinken

import (
	"fmt"
	"math/rand/v2"
)

type FireMatrixEffect struct {
	opts FireMatrixEffectOpts
	ws   wsEngine
	EffectControl
}

type FireMatrixEffectOpts struct {
	base         EffectsOpts       `json:"-"`
	Sparking     float64           `json:"sparking"`
	Cooling      float64           `json:"cooling"`
	palletteFunc func(float64) rgb `json:"-"`
}

func NewFireMatrixEffect(opts FireMatrixEffectOpts) *FireMatrixEffect {
	effect := FireMatrixEffect{}
	effect.opts = opts
	return &effect
}

func (e *FireMatrixEffect) engine() wsEngine {
	return e.ws
}

func (e *FireMatrixEffect) Start() error {
	return e.startEffect(e.opts.base, e)
}

func (e *FireMatrixEffect) Stop() {
	clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.stopEffect(e.engine())
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

	for row := 0; row < buffer.height; row++ {
		// Step 1.  Cool down every cell a little
		heat := heat[row]
		for i := 0; i < buffer.width; i++ {
			ce := rand.Float64() * e.opts.Cooling
			heat[i] = fsub_norm(heat[i], ce)
		}

		for k := buffer.width - 1; k > 1; k-- {
			if k > 2 {
				heat[k] = (heat[k-1] + heat[k-2] + heat[k-3]) / 3
			} else if k == 2 {
				heat[k] = (heat[k-1] + heat[k-2]) / 2
			} else if k == 1 {
				heat[k] = (heat[k-1])
			}
		}

		if rand.Float64() < e.opts.Sparking {
			y := rand.Int() % buffer.width / 10
			if heat[0] < 0.1 {
				y = 0
			}
			h := rand.Float64()*0.4 + 0.6
			heat[y] = fadd_norm(heat[y], h)
		}

		for col := 0; col < buffer.width; col++ {
			c := e.opts.palletteFunc(heat[col])
			buffer.setPixel(col, row, c, 1.0)
		}
	}
}

func (e *FireMatrixEffect) run(engine wsEngine) {
	e.ws = engine

	ledCount := e.opts.base.LedCount
	fmt.Printf("New FireMatrixEffect with width %d, height %d\n", e.opts.base.Width, e.opts.base.Height)

	buffer := LedMatrix{make([]rgb, ledCount), e.opts.base.Width, e.opts.base.Height}
	heat := make([]([]float64), e.opts.base.Height)
	for i := range heat {
		heat[i] = make([]float64, e.opts.base.Width)
	}

	for e.running.Load() == true {
		doFrame(e.opts.base.FrameTime, func() {
			e.animate(&buffer, heat)
			renderBuffer(e, buffer.leds)
		})
	}
	e.done <- true
}
