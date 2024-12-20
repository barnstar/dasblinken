package dasblinken

import (
	"math/rand/v2"
)

type WipeEffect struct {
	opts WipeEffectOpts
	ws   wsEngine
	EffectControl
}

type WipeEffectOpts struct {
	base        EffectsOpts
	spriteCount int
}

func NewWipeEffect(opts WipeEffectOpts) *WipeEffect {
	effect := WipeEffect{}
	effect.opts = opts
	return &effect
}

func (e *WipeEffect) engine() wsEngine {
	return e.ws
}

func (e *WipeEffect) Start() error {
	return e.startEffect(e.opts.base, e)
}

func (e *WipeEffect) Stop() {
	e.stopEffect(e.engine())
}

func (e *WipeEffect) Opts() EffectsOpts {
	return e.opts.base
}

var (
	blobs = []spriteData{
		spriteData{{0.5, 0.0, 0.0}, {0.7, 0.0, 0.0}, {1.0, 0.0, 0.0}, {0.7, 0.0, 0.0}, {0.5, 0.0, 0.0}},
		spriteData{{0.0, 0.4, 0.0}, {0.0, 0.6, 0.0}, {0.0, 1.0, 0.0}, {0.0, 0.6, 0.0}, {0.0, 0.4, 0.0}},
	}
)

func makeSprite(speed float32) sprite {
	i := rand.Int() % 2
	return sprite{
		0,
		speed,
		blobs[i],
	}
}

func (e *WipeEffect) render(sprites []sprite) {
	clear(e)
	for _, s := range sprites {
		wrappingOverlay(e, &s)
	}
	e.ws.Render()
}

func (e *WipeEffect) animate(sprites []sprite) {
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

func (e *WipeEffect) run(engine wsEngine) {
	e.ws = engine

	blobs := make([]sprite, 0)
	for range e.opts.spriteCount {
		v := rand.Float32()*2.0 - 1.0
		s := makeSprite(v)
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
