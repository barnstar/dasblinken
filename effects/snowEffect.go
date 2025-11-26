package effects

import (
	"math/rand/v2"

	. "barnstar.com/dasblinken"
)

type SnowEffect struct {
	opts SnowEffectOpts
	ws   WSEngine
	EffectState
}

type SnowEffectOpts struct {
	base        EffectsOpts
	MeltSpeed   float64
	FlakeChance float64
}

type SnowConfig struct {
	Name       string  `json:"name"`
	Topology   string  `json:"topology"`
	Dampening  float64 `json:"dampening"`
	Snowflakes float64 `json:"snowflakes"`
}

func NewSnowEffect(config SnowConfig, stripConfig StripConfig) *SnowEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))
	opts := SnowEffectOpts{
		base:        baseOpts,
		MeltSpeed:   config.Dampening,
		FlakeChance: config.Snowflakes,
	}
	effect := SnowEffect{}
	effect.opts = opts
	return &effect
}

func (e *SnowEffect) Engine() WSEngine {
	return e.ws
}

func (e *SnowEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *SnowEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *SnowEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *SnowEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *SnowEffect) animate(sprites []LinearSprite) {
	letItSnow := rand.Float64() < e.opts.FlakeChance
	if letItSnow {
		pos := rand.Int() % len(sprites)
		s := &sprites[pos]
		s.Lum = 1.0
	}

	for i, _ := range sprites {
		s := &sprites[i]
		s.Lum = s.Lum * e.opts.MeltSpeed
	}
}

func (e *SnowEffect) render(sprites []LinearSprite) {
	e.ws.Wait()
	for _, s := range sprites {
		WrappingOverlay(e, &s)
	}
	e.ws.Render()
}

func makeFlake(location int) LinearSprite {
	return LinearSprite{
		X:    float32(location),
		V:    0,
		Lum:  0.0,
		Data: SpriteData{RGB{R: 1.0, G: 1.0, B: 1.0}},
	}
}

func (e *SnowEffect) Run(engine WSEngine) {
	e.ws = engine

	flakes := make([]LinearSprite, e.opts.base.Len())
	for i := 0; i < e.opts.base.Len(); i++ {
		flakes[i] = makeFlake(i)
	}

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(flakes)
			e.render(flakes)
		})
	}
	e.Done <- true
}
