package effects

import (
	"math"
	"math/rand/v2"

	. "barnstar.com/dasblinken"
)

type FlickerEffect struct {
	opts FlickerEffectOpts
	ws   WSEngine
	EffectState
}

type FlickerEffectOpts struct {
	base         EffectsOpts
	NumGroups    int
	MinRadius    int
	MaxRadius    int
	FlickerSpeed float64
	FlickerAmt   float64
	PalletteFunc func(float64) RGB
}

type FlickerConfig struct {
	Name         string  `json:"name"`
	Topology     string  `json:"topology"`
	NumGroups    int     `json:"numGroups"`
	MinRadius    int     `json:"minRadius"`
	MaxRadius    int     `json:"maxRadius"`
	FlickerSpeed float64 `json:"flickerSpeed"`
	FlickerAmt   float64 `json:"flickerAmt"`
	Palette      string  `json:"palette"`
}

type flickerGroup struct {
	center      int
	radius      int
	intensity   float64
	flickerAmt  float64
	speed       float64
	phase       float64
	lifetime    float64 // Total lifetime of this group (5-10 seconds)
	age         float64 // Current age in seconds
	fadeInTime  float64 // Time to fade in (1 second)
	fadeOutTime float64 // Time to fade out (1 second)
}

func NewFlickerEffect(config FlickerConfig, stripConfig StripConfig) *FlickerEffect {
	baseOpts := StripOptsDefString(config.Name, stripConfig, getTopology(config.Topology))

	numGroups := config.NumGroups
	if numGroups == 0 {
		numGroups = 5
	}
	minRadius := config.MinRadius
	if minRadius == 0 {
		minRadius = 15
	}
	maxRadius := config.MaxRadius
	if maxRadius == 0 {
		maxRadius = 30
	}
	flickerSpeed := config.FlickerSpeed
	if flickerSpeed == 0 {
		flickerSpeed = 2.0
	}
	flickerAmt := config.FlickerAmt
	if flickerAmt == 0 {
		flickerAmt = 0.5
	}

	opts := FlickerEffectOpts{
		base:         baseOpts,
		NumGroups:    numGroups,
		MinRadius:    minRadius,
		MaxRadius:    maxRadius,
		FlickerSpeed: flickerSpeed,
		FlickerAmt:   flickerAmt,
		PalletteFunc: getPalette(config.Palette),
	}

	effect := FlickerEffect{}
	effect.opts = opts
	return &effect
}

func (e *FlickerEffect) Engine() WSEngine {
	return e.ws
}

func (e *FlickerEffect) Start() error {
	return e.StartEffect(e.opts.base, e)
}

func (e *FlickerEffect) Stop() {
	Clear(e)
	e.ws.Render()
	e.ws.Wait()
	e.StopEffect(e.Engine())
}

func (e *FlickerEffect) Opts() EffectsOpts {
	return e.opts.base
}

func (e *FlickerEffect) SetStripConfig(s StripConfig) {
	e.opts.base.StripConfig = s
}

func (e *FlickerEffect) initializeGroup(ledCount int) flickerGroup {
	return flickerGroup{
		center:      rand.IntN(ledCount),
		radius:      rand.IntN(e.opts.MaxRadius-e.opts.MinRadius) + e.opts.MinRadius,
		intensity:   rand.Float64(),
		flickerAmt:  rand.Float64()*0.4 + 0.3,
		speed:       rand.Float64()*e.opts.FlickerSpeed + 0.5,
		phase:       rand.Float64() * math.Pi * 2,
		lifetime:    rand.Float64()*10.0 + 5.0, // 5-15 seconds
		age:         0,
		fadeInTime:  2.0, // 2 second fade in
		fadeOutTime: 2.0, // 2 second fade out
	}
}

func (e *FlickerEffect) Run(engine WSEngine) {
	e.ws = engine

	ledCount := e.opts.base.Len()
	buffer := make([]RGB, ledCount)

	// Initialize flicker groups
	groups := make([]flickerGroup, e.opts.NumGroups)
	for i := range groups {
		groups[i] = e.initializeGroup(ledCount)
		// Stagger initial ages so they don't all restart at once
		groups[i].age = rand.Float64() * groups[i].lifetime
	}

	frameCount := 0.0

	for e.Running.Load() == true {
		RenderFrame(e.opts.base.FrameTime, func() {
			// Clear buffers
			ClearBuffer(buffer)

			deltaTime := e.opts.base.FrameTime.Seconds()
			frameCount += deltaTime

			// Calculate contribution from each group
			for g := range groups {
				group := &groups[g]

				// Update group age
				group.age += deltaTime

				// Check if group lifetime expired, re-initialize if so
				if group.age >= group.lifetime {
					groups[g] = e.initializeGroup(ledCount)
					group = &groups[g]
				}

				// Calculate envelope (fade in/out)
				envelope := 1.0
				if group.age < group.fadeInTime {
					// Fade in
					envelope = group.age / group.fadeInTime
				} else if group.age > group.lifetime-group.fadeOutTime {
					// Fade out
					timeUntilEnd := group.lifetime - group.age
					envelope = timeUntilEnd / group.fadeOutTime
				}

				// Update flicker using sine wave plus noise
				flicker := math.Sin(frameCount*group.speed+group.phase) * group.flickerAmt
				flicker += (rand.Float64() - 0.5) * 0.2
				brightness := 1.0 + flicker

				// Apply to pixels within radius
				for i := 0; i < ledCount; i++ {
					distance := math.Abs(float64(i - group.center))

					// Calculate fade based on distance from center
					fade := 1.0 - (distance / float64(group.radius))
					if fade < 0 {
						fade = 0
					}

					// Apply smooth falloff (cosine fade)
					if fade > 0 {
						fade = (math.Cos((1-fade)*math.Pi) + 1) / 2
					}

					// Apply envelope to intensity
					intensity := fade * brightness * group.intensity * envelope

					// Get normalized color values
					var r, g, b float64
					if e.opts.PalletteFunc != nil {
						color := e.opts.PalletteFunc(intensity)
						r = float64(color.R) / 255.0
						g = float64(color.G) / 255.0
						b = float64(color.B) / 255.0
					} else {
						// Default flame colors: orange/red/yellow (normalized)
						r = intensity
						g = intensity * 0.6
						b = intensity * 0.12
					}

					// Additive blending with normalized floats
					buffer[i].R = math.Min(1.0, buffer[i].R+r)
					buffer[i].G = math.Min(1.0, buffer[i].G+g)
					buffer[i].B = math.Min(1.0, buffer[i].B+b)
				}
			}
			RenderBuffer(e, buffer)
		})
	}
	Clear(e)
	e.Done <- true
}
