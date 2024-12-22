package dasblinken

import "math/rand/v2"

type FireEffect struct {
	opts FireEffectOpts
	ws   wsEngine
	EffectControl
}

type FireEffectOpts struct {
	base         EffectsOpts       `json:"-"`
	Sparking     float64           `json:"sparking"`
	Cooling      float64           `json:"cooling"`
	DoubleEnded  bool              `json:"doubleEnded"`
	palletteFunc func(float64) rgb `json:"-"`
}

func NewFireEffect(opts FireEffectOpts) *FireEffect {
	effect := FireEffect{}
	effect.opts = opts
	return &effect
}

func (e *FireEffect) engine() wsEngine {
	return e.ws
}

func (e *FireEffect) Start() error {
	return e.startEffect(e.opts.base, e)
}

func (e *FireEffect) Stop() {
	clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.stopEffect(e.engine())
}

func (e *FireEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *FireEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

type fireLayout struct {
	start int
	end   int
}

func (e *FireEffect) animate(buffer []rgb, heat []float64, reverse bool, limit int) {
	ledCount := e.opts.base.LedCount

	// Step 1.  Cool down every cell a little
	for i := 0; i < ledCount; i++ {
		ce := rand.Float64() * e.opts.Cooling
		heat[i] = fsub_norm(heat[i], ce)
	}

	for k := ledCount - 1; k > 1; k-- {
		heat[k] = (heat[k-1] + heat[k-2] + heat[k-2]) / 3
	}

	if rand.Float64() < e.opts.Sparking {
		y := rand.Int() % ledCount / 10
		h := rand.Float64()*0.4 + 0.6
		heat[y] = fadd_norm(heat[y], h)
	}

	if !reverse {
		for j := 0; j < limit; j++ {
			buffer[j] = e.opts.palletteFunc(heat[j])
		}
	} else {
		for j := ledCount - 1; j >= limit; j-- {
			buffer[j] = e.opts.palletteFunc(heat[ledCount-j-1])
		}
	}
}

func (e *FireEffect) run(engine wsEngine) {
	e.ws = engine

	ledCount := e.opts.base.LedCount

	buffer := make([]rgb, ledCount)
	heat1 := make([]float64, ledCount)
	heat2 := make([]float64, ledCount)

	for e.running.Load() == true {
		doFrame(e.opts.base.FrameTime, func() {
			if e.opts.DoubleEnded {
				e.animate(buffer, heat1, false, ledCount/2)
				e.animate(buffer, heat2, true, ledCount/2)
			} else {
				e.animate(buffer, heat1, false, ledCount)
			}
			renderBuffer(e, buffer)
		})
	}
	e.done <- true
}
