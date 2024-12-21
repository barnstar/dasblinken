package dasblinken

import "math/rand/v2"

type SnowEffect struct {
	opts SnowEffectOpts
	ws   wsEngine
	EffectControl
}

type SnowEffectOpts struct {
	base        EffectsOpts
	meltSpeed   float64
	flakeChance float64
}

func NewSnowEffect(opts SnowEffectOpts) *SnowEffect {
	effect := SnowEffect{}
	effect.opts = opts
	return &effect
}

func (e *SnowEffect) engine() wsEngine {
	return e.ws
}

func (e *SnowEffect) Start() error {
	return e.startEffect(e.opts.base, e)
}

func (e *SnowEffect) Stop() {
	clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.stopEffect(e.engine())
}

func (e *SnowEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *SnowEffect) animate(sprites []sprite) {
	letItSnow := rand.Float64() < e.opts.flakeChance
	if letItSnow {
		pos := rand.Int() % len(sprites)
		s := &sprites[pos]
		s.lum = 1.0
	}

	for i, _ := range sprites {
		s := &sprites[i]
		s.lum = s.lum * e.opts.meltSpeed
	}
}

func (e *SnowEffect) render(sprites []sprite) {
	e.ws.Wait()
	for _, s := range sprites {
		wrappingOverlay(e, &s)
	}
	e.ws.Render()
}

func makeFlake(location int) sprite {
	return sprite{
		float32(location),
		0,
		0.0,
		spriteData{{1.0, 1.0, 1.0}},
	}
}

func (e *SnowEffect) run(engine wsEngine) {
	e.ws = engine

	flakes := make([]sprite, e.opts.base.LedCount)
	for i := 0; i < e.opts.base.LedCount; i++ {
		flakes[i] = makeFlake(i)
	}

	for e.running.Load() == true {
		doFrame(e.opts.base.FrameTime, func() {
			e.animate(flakes)
			e.render(flakes)
		})
	}
	e.done <- true
}
