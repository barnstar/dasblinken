package dasblinken

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type WipeEffect struct {
	opts WipeEffectOpts
	ws   wsEngine
	EffectControl
}

type WipeEffectOpts struct {
	EffectsOpts
	spriteCount int
	frameTime   time.Duration
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
	return e.startEffect(e.opts.EffectsOpts, e)
}

func (e *WipeEffect) Stop() {
	e.stopEffect(e.engine())
}

var (
	blobs = []spriteData{
		spriteData{{0.5, 0.0, 0.0}, {0.7, 0.0, 0.0}, {1.0, 0.0, 0.0}, {0.7, 0.0, 0.0}, {0.5, 0.0, 0.0}},
		spriteData{{0.0, 0.4, 0.0}, {0.0, 0.6, 0.0}, {0.0, 1.0, 0.0}, {0.0, 0.6, 0.0}, {0.0, 0.4, 0.0}},
	}
)

func (e *WipeEffect) render(sprites []sprite) {
	clear(e)
	for _, s := range sprites {
		wrappingOverlay(e, &s)
	}
	e.ws.Render()
}

func (e *WipeEffect) animate(sprites []sprite) {
	for i, _ := range sprites {
		b := &sprites[i]
		b.x = b.x + b.v
		if b.x < 0 {
			b.x = float32(e.opts.LedCount) + b.x
		} else if b.x > float32(e.opts.LedCount) {
			b.x = b.x - float32(e.opts.LedCount)
		}
	}
}

func (e *WipeEffect) run(engine wsEngine) {
	e.ws = engine

	blobs := make([]sprite, 0)
	for range e.opts.spriteCount {
		v := rand.Float32()*2.0 - 1.0
		blobs = append(blobs, makeSprite(v))
	}

	fmt.Println("Running Wipe Effect")
	for e.running.Load() == true {
		doFrame(e.opts.frameTime, func() {
			e.animate(blobs)
			e.render(blobs)
		})
	}
	fmt.Println("Stopping Wipe Effect")
	e.done <- true
}
