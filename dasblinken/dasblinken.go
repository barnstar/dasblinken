package dasblinken

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"sync"
	"time"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Topology int

const (
	Any Topology = iota
	Linear
	Matrix
)

// WSEngine is the interface wrapper for the rpi-ws281x-go library.
type WSEngine interface {
	Init() error
	Render() error
	Wait() error

	Fini()
	Leds(channel int) []uint32
}

// Dasblinken is the main struct for the dasblinken package.
// It holds the active effect and the configuration for the strip.
type Dasblinken struct {
	mu      sync.Mutex
	active  Effect
	effects map[string]Effect
	strip   StripConfig
}

// NewDasblinken creates a new Dasblinken instance.
func NewDasblinken() *Dasblinken {
	fmt.Println("Starting dasblinken!")
	return &Dasblinken{
		effects: make(map[string]Effect),
	}
}

// SetStrip sets the strip configuration.
func (dbl *Dasblinken) SetStrip(c StripConfig) {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	dbl.strip = c
}

func (dbl *Dasblinken) Config() StripConfig {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	return dbl.strip
}

// UpdateConfig updates the strip configuration and stops any active effect
func (dbl *Dasblinken) UpdateConfig(config StripConfig) {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	if dbl.active != nil {
		dbl.active.Stop()
		dbl.active = nil
	}
	dbl.strip = config
}

func Device(opts EffectsOpts) (WSEngine, error) {
	fmt.Println("Configuring Blinken Device at Pin:", opts.Pin, "LedCount:", opts.Len(), "Brightness:", opts.Brightness)

	opt := ws281x.DefaultOptions
	opt.Channels[0].Brightness = opts.Brightness
	opt.Channels[0].GpioPin = opts.Pin
	opt.Channels[0].LedCount = opts.Len()

	return ws281x.MakeWS2811(&opt)
}

func (dbl *Dasblinken) ActiveEffect() Effect {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	return dbl.active
}

func (dbl *Dasblinken) StopAll() {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	if dbl.active != nil {
		dbl.active.Stop()
	}
	dbl.active = nil
	fmt.Println("Dasblinken is kaput")
}

func (dbl *Dasblinken) Stop() {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	if dbl.active != nil {
		dbl.active.Stop()
	}
	dbl.active = nil
	fmt.Println("Dasblinken Stopped")
}

func (dbl *Dasblinken) RegisterEffect(effect Effect) {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	dbl.effects[effect.Opts().Name] = effect
}

func StripOptsDefString(name string, config StripConfig, requires Topology) EffectsOpts {
	frameTime := 1000000000 / config.Fps
	return EffectsOpts{
		name,
		config,
		time.Duration(frameTime),
		requires,
	}
}

func (dbl *Dasblinken) RandomEffect() {
	effects := dbl.Effects()
	if len(effects) > 0 {
		randomIndex := rand.Int() % len(effects)
		effect := effects[randomIndex]
		dbl.SwitchToEffect(effect.Opts().Name)
	}
}

func (dbl *Dasblinken) Effects() []Effect {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()

	config := dbl.strip
	topology := Linear
	if config.Height > 1 {
		topology = Matrix
	}

	effectsSlice := make([]Effect, 0)
	for _, effect := range dbl.effects {
		requires := effect.Opts().Requires
		if requires == Any || requires == topology {
			effectsSlice = append(effectsSlice, effect)
		}
	}
	sort.Slice(effectsSlice, func(i, j int) bool {
		return effectsSlice[i].Opts().Name < effectsSlice[j].Opts().Name
	})

	return effectsSlice
}

func (dbl *Dasblinken) ClearEffects() {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	dbl.effects = make(map[string]Effect)
}

func (dbl *Dasblinken) SwitchToEffect(name string) error {
	dbl.mu.Lock()
	next := dbl.effects[name]
	if nil == next {
		dbl.mu.Unlock()
		return fmt.Errorf("Effect %s not found", name)
	}
	config := dbl.strip
	dbl.mu.Unlock()

	dbl.Stop()
	dbl.mu.Lock()
	dbl.active = next
	dbl.mu.Unlock()
	next.SetStripConfig(config)
	next.Start()
	fmt.Printf("Switched to effect %v\n", name)
	return nil
}

func (dbl *Dasblinken) StartStreaming() {

}
