package dasblinken

type RainbowChaseEffect struct {
	opts ChaseEffectOpts
	ws   wsEngine
	EffectControl

	q    float64
	o    float64
	skip int
}

type ChaseEffectOpts struct {
	Base EffectsOpts `json:"-"`

	Speed float64 `json:"name"`
}

func NewRainbowChaseEffect(opts ChaseEffectOpts) *RainbowChaseEffect {
	effect := RainbowChaseEffect{}
	effect.opts = opts
	return &effect
}

func (e *RainbowChaseEffect) engine() wsEngine {
	return e.ws
}

func (e *RainbowChaseEffect) Start() error {
	return e.startEffect(e.opts.Base, e)
}

func (e *RainbowChaseEffect) Stop() {
	clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.stopEffect(e.engine())
}

func (e *RainbowChaseEffect) Opts() EffectsOpts {
	return e.opts.Base
}

func (e *RainbowChaseEffect) SetStripConfig(s StripConfig) {
	e.opts.Base.StripConfig = s
}

func (e *RainbowChaseEffect) run(engine wsEngine) {
	e.ws = engine

	buffer := make([]rgb, e.opts.Base.LedCount)
	for e.running.Load() == true {
		doFrame(e.opts.Base.FrameTime, func() {
			e.animate(buffer)
			renderBuffer(e, buffer)
		})
	}
	e.done <- true
}

func (e *RainbowChaseEffect) animate(buffer []rgb) {
	ledCount := e.opts.Base.LedCount
	e.q = e.q + e.opts.Speed
	if e.q > 3 {
		e.q = 0
	} else if e.q < 0 {
		e.q = 3
	}

	e.o = e.o - 0.01
	if e.o < 0 {
		e.o = 1.0
	}

	clearBuffer(buffer)

	for i := 0; i < ledCount-3; i = i + 3 {
		w := (e.o + float64(i+int(e.q))/float64(ledCount))
		if w >= 1 {
			w = w - 1
		}
		c := rainbowPalette(w)
		buffer[i+int(e.q)] = c
	}
}
