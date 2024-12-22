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
	base          EffectsOpts `json:"-"`
	SpriteCount   int         `json:"spriteCount"`
	Bidirectional bool        `json:"bidirectional"`
	Speed         float64     `json:"speed"`
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

func (e *RaceEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func makeBlob(speed float32) sprite {
	c := rainbowPalette(rand.Float64())
	var data spriteData

	if speed < 0 {
		data = spriteData{c, c.Faded(8.), c.Faded(.6), c.Faded(.4), c.Faded(.2)}
	} else if speed >= 0 {
		data = spriteData{c.Faded(.2), c.Faded(.4), c.Faded(.6), c.Faded(.8), c}
	}

	return sprite{
		0,
		speed,
		1.0,
		data,
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
	for range e.opts.SpriteCount {
		var v float32
		if e.opts.Bidirectional {
			v = rand.Float32()*float32(e.opts.Speed)*2 - float32(e.opts.Speed)
		} else {
			v = rand.Float32() * float32(e.opts.Speed)
		}
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
