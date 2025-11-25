package effects

import (
	"math/rand/v2"

	. "barnstar.com/dasblinken"
)

type RaceEffect struct {
	opts RaceEffectOpts
	ws   WSEngine
	EffectState
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

func (e *RaceEffect) Engine() WSEngine {
	return e.ws
}

func (e *RaceEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *RaceEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *RaceEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *RaceEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func makeBlob(speed float32) LinearSprite {
	c := RainbowPalette(rand.Float64())
	var data SpriteData

	if speed < 0 {
		data = SpriteData{c, c.Faded(8.), c.Faded(.6), c.Faded(.4), c.Faded(.2)}
	} else if speed >= 0 {
		data = SpriteData{c.Faded(.2), c.Faded(.4), c.Faded(.6), c.Faded(.8), c}
	}

	return LinearSprite{
		0,
		speed,
		1.0,
		data,
	}
}

func (e *RaceEffect) render(sprites []LinearSprite) {
	e.ws.Wait()
	Clear(e)
	for _, s := range sprites {
		WrappingOverlay(e, &s)
	}
	e.ws.Render()
}

func (e *RaceEffect) animate(sprites []LinearSprite) {
	for i, _ := range sprites {
		s := &sprites[i]
		s.X = s.X + s.V
		if s.X < 0 {
			s.X = float32(e.Opts().Len()) + s.X
		} else if s.X > float32(e.Opts().Len()) {
			s.X = s.X - float32(e.Opts().Len())
		}
	}
}

func (e *RaceEffect) Run(engine WSEngine) {
	e.ws = engine

	blobs := make([]LinearSprite, 0)
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

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(blobs)
			e.render(blobs)
		})
	}
	e.Done <- true
}
