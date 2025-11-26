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

type RaceConfig struct {
	Name      string  `json:"name"`
	Topology  string  `json:"topology"`
	Length    int     `json:"length"`
	Mirrored  bool    `json:"mirrored"`
	NumRacers float64 `json:"numRacers"`
}

type RaceEffectOpts struct {
	base          EffectsOpts
	SpriteCount   int
	Bidirectional bool
	Speed         float64
}

func NewRaceEffect(config RaceConfig, stripConfig StripConfig) *RaceEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))
	opts := RaceEffectOpts{
		base:          baseOpts,
		SpriteCount:   int(config.NumRacers),
		Bidirectional: config.Mirrored,
		Speed:         float64(config.Length),
	}
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
		X:    0,
		V:    speed,
		Lum:  1.0,
		Data: data,
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
