package effects

import (
	"strconv"
	"time"

	. "barnstar.com/dasblinken"
)

// This is optimized for an 8x32 matrix
type ClockEffect struct {
	opts ClockEffectOpts
	ws   WSEngine
	EffectState

	offset float64
	x      float64
}

type ClockConfig struct {
	Name     string `json:"name"`
	Topology string `json:"topology"`
}

type ClockEffectOpts struct {
	base EffectsOpts
}

func NewClockEffect(config ClockConfig, stripConfig StripConfig) *ClockEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))
	opts := ClockEffectOpts{
		base: baseOpts,
	}
	effect := ClockEffect{}
	effect.opts = opts
	return &effect
}

func (e *ClockEffect) Engine() WSEngine {
	return e.ws
}

func (e *ClockEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *ClockEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *ClockEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *ClockEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *ClockEffect) animate(buffer *LedMatrix, fc int) {
	t := time.Now()
	hh := t.Format("03") // Get the current hour as a string
	mm := t.Format("04") // Get the current minute as a string
	ss := t.Format("05") // Get the current second as a string

	ClearBuffer(buffer.Leds)
	x := -2.0
	l := 1.0
	if fc%60 < 30 {
		l = 0.3
	}
	cf := func(float64, float64) RGB {
		return RGB{R: 0.7, G: 0.0, B: 0.0}
	}
	for _, n := range hh {
		c, err := strconv.Atoi(string(n))
		if err != nil || c == 0 {
			x += 5
			continue
		}
		buffer.DrawChar(ClockFont, x, 0, byte(c), cf, 1.0)
		x += 5
	}
	buffer.DrawChar(ClockFont, x, 0, 10, cf, l)
	x += 2
	for _, n := range mm {
		c, _ := strconv.Atoi(string(n))
		buffer.DrawChar(ClockFont, x, 0, byte(c), cf, 1.0)
		x += 5
	}
	buffer.DrawChar(ClockFont, x, 0, 10, cf, l)
	x += 2
	for _, n := range ss {
		c, _ := strconv.Atoi(string(n))
		buffer.DrawChar(ClockFont, x, 0, byte(c), cf, 1.0)
		x += 5
	}
	buffer.ApplyLuminosity()
}

func (e *ClockEffect) Run(engine WSEngine) {
	e.ws = engine
	e.x = float64(e.opts.base.Width)
	e.offset = 0.0
	fc := 0

	buffer := NewLedMatrix(e.opts.base.Width, e.opts.base.Height)

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			e.animate(buffer, fc)
			RenderBuffer(e, buffer.Leds)
		})
		fc += 1
	}
	Clear(e)
	e.Done <- true
}
