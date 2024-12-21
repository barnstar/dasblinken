package dasblinken

import (
	"math/rand/v2"
)

type RaceEffect struct {
	opts RaceEffectOpts
	ws   wsEngine
	EffectControl
}

type RaceEffectOpts struct {
	base        EffectsOpts
	spriteCount int
}

func NewRaceEffect(opts RaceEffectOpts) *RaceEffect {
	effect := RaceEffect{}
	effect.opts = opts
	return &effect
}

func (e *RaceEffect) engine() wsEngine {
	return e.ws
}

func (e *RaceEffect) Start() error {
	return e.startEffect(e.opts.base, e)
}

func (e *RaceEffect) Stop() {
	clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.stopEffect(e.engine())
}

func (e *RaceEffect) Opts() EffectsOpts {
	return e.opts.base
}

var (
	blobs = []spriteData{
		spriteData{{0.5, 0.0, 0.0}, {0.7, 0.0, 0.0}, {1.0, 0.0, 0.0}, {0.7, 0.0, 0.0}, {0.5, 0.0, 0.0}},
		spriteData{{0.0, 0.4, 0.0}, {0.0, 0.6, 0.0}, {0.0, 1.0, 0.0}, {0.0, 0.6, 0.0}, {0.0, 0.4, 0.0}},
		spriteData{{0.4, 0.4, 0.0}, {0.6, 0.6, 0.0}, {1.0, 1.0, 0.0}, {0.6, 0.6, 0.0}, {0.4, 0.4, 0.0}},
		spriteData{{0.0, 0.4, 0.4}, {0.0, 0.6, 0.6}, {0.0, 1.0, 1.0}, {0.0, 0.6, 6.0}, {0.0, 0.4, 4.0}},
		spriteData{{0.6, 0.6, 0.6}, {1.0, 1.0, 1.0}, {0.6, 0.6, 6.0}},
	}
)

func makeBlob(speed float32) sprite {
	return sprite{
		0,
		speed,
		1.0,
		blobs[rand.Int()%len(blobs)],
	}
}

func (e *RaceEffect) render(sprites []sprite) {
	e.ws.Wait()
	clear(e)
	for _, s := range sprites {
		wrappingOverlay(e, &s)
	}
	e.ws.Render()
}

func (e *RaceEffect) animate(sprites []sprite) {
	for i, _ := range sprites {
		s := &sprites[i]
		s.x = s.x + s.v
		if s.x < 0 {
			s.x = float32(e.Opts().LedCount) + s.x
		} else if s.x > float32(e.Opts().LedCount) {
			s.x = s.x - float32(e.Opts().LedCount)
		}
	}
}

func (e *RaceEffect) run(engine wsEngine) {
	e.ws = engine

	blobs := make([]sprite, 0)
	for range e.opts.spriteCount {
		v := rand.Float32()*2.0 - 1.0
		s := makeBlob(v)
		blobs = append(blobs, s)
	}

	for e.running.Load() == true {
		doFrame(e.opts.base.FrameTime, func() {
			e.animate(blobs)
			e.render(blobs)
		})
	}
	e.done <- true
}
