package dasblinken

import "math/rand/v2"

type FireEffect struct {
	opts FireEffectOpts
	ws   wsEngine
	EffectControl
}

type FireEffectOpts struct {
	base     EffectsOpts
	sparking float64
	cooling  float64
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

// Subs two floats and clamps the result at 0.0
func fsub_norm(a, b float64) float64 {
	diff := a - b
	diff = max(diff, 0.0)
	return diff
}

// Adds two floats and clamps the result at 1.0
func fadd_norm(a, b float64) float64 {
	sum := a + b
	sum = min(sum, 1.0)
	return sum
}

func (e *FireEffect) animate(buffer []rgb, heat []float64) {
	ledCount := e.opts.base.LedCount

	// Step 1.  Cool down every cell a little
	for i := 0; i < ledCount; i++ {
		ce := rand.Float64() * e.opts.cooling
		heat[i] = fsub_norm(heat[i], ce)
	}

	// Step 2.  Heat from each cell drifts 'up' and diffuses a little
	for k := ledCount - 1; k > 1; k-- {
		heat[k] = (heat[k-1] + heat[k-2] + heat[k-2]) / 3
	}

	// Step 3.  Randomly ignite new 'sparks' of heat near the bottom
	if rand.Float64() < e.opts.sparking {
		y := rand.Int() % 7
		h := rand.Float64()*0.4 + 0.6
		heat[y] = fadd_norm(heat[y], h)
	}

	// Step 4.  Map from heat cells to LED colors
	for j := 0; j < ledCount; j++ {
		buffer[j] = heatColor(heat[j])
	}
}

func (e *FireEffect) render(buffer []rgb) {
	ledCount := e.opts.base.LedCount

	e.ws.Wait()
	for j := 0; j < ledCount; j++ {
		e.ws.Leds(e.Opts().Channel)[j] = buffer[j].toHex(1.0)
	}
	e.ws.Render()
}

func (e *FireEffect) run(engine wsEngine) {
	e.ws = engine

	buffer := make([]rgb, e.opts.base.LedCount)
	heat := make([]float64, e.opts.base.LedCount)

	for e.running.Load() == true {
		doFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer, heat)
			e.render(buffer)
		})
	}
	e.done <- true
}

func heatColor(temperature float64) rgb {
	var heatcolor rgb

	// now figure out which third of the spectrum we're in:
	if temperature > 0.66 {
		// we're in the hottest third
		heatcolor.r = 1.0                           // full red
		heatcolor.g = 1.0 - (temperature-0.66)/0.33 // full green
		heatcolor.b = (temperature - 0.66) / 0.33   // ramp up blue

	} else if temperature > 0.33 && temperature <= 0.66 {
		// we're in the middle third
		heatcolor.r = 1.0                         // full red
		heatcolor.g = (temperature - 0.33) / 0.33 // ramp up green
		heatcolor.b = 0                           // no blue

	} else {
		// we're in the coolest third
		heatcolor.r = temperature / 0.33 // ramp up red
		heatcolor.g = 0                  // no green
		heatcolor.b = 0                  // no blue
	}

	return heatcolor
}
