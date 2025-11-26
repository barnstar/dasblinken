package effects

import (
	"math"
	"math/rand/v2"

	. "barnstar.com/dasblinken"
)

type BallsConfig struct {
	Name     string `json:"name"`
	Topology string `json:"topology"`
	NumBalls int    `json:"numBalls"`
	TrailLen int    `json:"trailLen"`
	Palette  string `json:"palette"`
}

type BallsEffect struct {
	opts BallsEffectOpts
	ws   WSEngine
	EffectState
}

type BallsEffectOpts struct {
	base EffectsOpts

	ballCount    int
	trailLen     int
	palletteFunc func(float64) RGB
}

func NewBallsEffect(config BallsConfig, stripConfig StripConfig) *BallsEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))
	opts := BallsEffectOpts{
		base:         baseOpts,
		ballCount:    config.NumBalls,
		trailLen:     config.TrailLen,
		palletteFunc: getPalette(config.Palette),
	}
	effect := BallsEffect{}
	effect.opts = opts
	return &effect
}

func (e *BallsEffect) Engine() WSEngine {
	return e.ws
}

func (e *BallsEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *BallsEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *BallsEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *BallsEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *BallsEffect) animate(buffer *LedMatrix, balls []ball) {
	ClearBuffer(buffer.Leds)
	for i := range balls {
		balls[i].drawBall(e, buffer)
		e.move(&balls[i])
	}
}

type ball struct {
	x, y   []float64
	dx, dy float64
	color  float64
}

func (e *BallsEffect) randomBall() ball {
	x := rand.Float64() * float64(e.opts.base.Width)
	y := rand.Float64() * float64(e.opts.base.Height)
	trailLen := e.opts.trailLen

	// Initialize slices with trailLen elements of the same value
	xSlice := make([]float64, trailLen+1)
	ySlice := make([]float64, trailLen+1)
	for i := range xSlice {
		xSlice[i] = x
		ySlice[i] = y
	}
	return ball{
		x:     xSlice,
		y:     ySlice,
		dx:    .2 * (rand.Float64()*2 - 1),
		dy:    .2 * (rand.Float64()*2 - 1),
		color: rand.Float64(),
	}
}

func (b *ball) drawBall(e *BallsEffect, buffer *LedMatrix) {
	f := 1.0
	for i := e.opts.trailLen; i >= 0; i-- {
		x := int(math.Round(b.x[i]))
		y := int(math.Round(b.y[i]))
		if x < 0 || x >= buffer.Width || y < 0 || y >= buffer.Height {
			return
		}
		col := e.opts.palletteFunc(b.color)
		buffer.SetPixel(x, y, col, f)
		f -= 0.02
	}
}

func (e *BallsEffect) move(b *ball) {
	for i := e.opts.trailLen - 1; i > 0; i-- {
		b.x[i] = b.x[i-1]
		b.y[i] = b.y[i-1]
	}

	b.x[0] += b.dx
	b.y[0] += b.dy
	if b.x[0] < 0 || b.x[0] >= float64(e.opts.base.Width) {
		b.dx = -b.dx
	}
	if b.y[0] < 0 || b.y[0] >= float64(e.opts.base.Height) {
		b.dy = -b.dy
	}
}

func (e *BallsEffect) Run(engine WSEngine) {
	e.ws = engine

	buffer := LedMatrix{
		Leds:   make([]RGB, e.opts.base.Len()),
		Lum:    make([]float64, e.opts.base.Len()),
		Width:  e.opts.base.Width,
		Height: e.opts.base.Height}

	balls := make([]ball, e.opts.ballCount)
	for i := range balls {
		balls[i] = e.randomBall()
	}

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(&buffer, balls)
			RenderBuffer(e, buffer.Leds)
		})
	}
	e.Done <- true
}
