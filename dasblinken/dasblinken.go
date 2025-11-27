package dasblinken

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"sync"
	"time"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Channel int

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
// It holds all the active effects and the configuration for the strips.
type Dasblinken struct {
	mu      sync.Mutex
	active  map[Channel]Effect
	effects map[string]Effect
	strips  map[Channel]StripConfig
}

// NewDasblinken creates a new Dasblinken instance.
func NewDasblinken() *Dasblinken {
	fmt.Println("Starting dasblinken!")
	return &Dasblinken{
		effects: make(map[string]Effect),
		active:  make(map[Channel]Effect),
		strips:  make(map[Channel]StripConfig),
	}
}

// AddStrip adds a new strip configuration to the Dasblinken instance.
// Each channel can have a separate led configuration.
func (dbl *Dasblinken) AddStrip(c StripConfig) {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	dbl.strips[c.Channel] = c
}

func (dbl *Dasblinken) Config(channel Channel) (StripConfig, bool) {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	config, ok := dbl.strips[channel]
	return config, ok
}

func Device(opts EffectsOpts) (WSEngine, error) {
	fmt.Println("Configuring Blinken Device at Pin:", opts.Pin, "Channel:", opts.Channel, "LedCount:", opts.Len(), "Brightness:", opts.Brightness)

	opt := ws281x.DefaultOptions
	opt.Channels[opts.Channel].Brightness = opts.Brightness

	opt.Channels[opts.Channel].GpioPin = opts.Pin
	opt.Channels[opts.Channel].LedCount = opts.Len()

	return ws281x.MakeWS2811(&opt)
}

func (dbl *Dasblinken) ActiveEffect(channel Channel) Effect {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	return dbl.active[channel]
}

func (dbl *Dasblinken) StopAll() {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	for _, effect := range dbl.active {
		effect.Stop()
	}
	dbl.active = make(map[Channel]Effect)
	fmt.Printf("Dasblinken is kaput")
}

func (dbl *Dasblinken) Stop(channel Channel) {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	if dbl.active[channel] != nil {
		dbl.active[channel].Stop()
	}
	dbl.active[channel] = nil
	fmt.Printf("Dasblinken Channel %v is kaput", channel)
}

func (dbl *Dasblinken) RegisterEffect(effect Effect) {
	stripType := Any
	requires := effect.Opts().Requires

	if effect.Opts().Height > 1 {
		stripType = Matrix
	} else {
		stripType = Linear
	}
	if requires != Any && requires != stripType {
		return
	}

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

func (dbl *Dasblinken) RandomEffect(channel Channel) {
	effects := dbl.Effects()
	if len(effects) > 0 {
		randomIndex := rand.Int() % len(effects)
		effect := effects[randomIndex]
		dbl.SwitchToEffect(effect.Opts().Name, channel)
	}
}

func (dbl *Dasblinken) Effects() []Effect {
	dbl.mu.Lock()
	defer dbl.mu.Unlock()
	effectsSlice := make([]Effect, 0, len(dbl.effects))
	for _, effect := range dbl.effects {
		effectsSlice = append(effectsSlice, effect)
	}
	sort.Slice(effectsSlice, func(i, j int) bool {
		return effectsSlice[i].Opts().Name < effectsSlice[j].Opts().Name
	})

	return effectsSlice
}

func (dbl *Dasblinken) SwitchToEffect(name string, channel Channel) error {
	dbl.mu.Lock()
	next := dbl.effects[name]
	if nil == next {
		return fmt.Errorf("Effect %s not found\n", name)
	}

	config, ok := dbl.strips[channel]
	if !ok {
		return fmt.Errorf("No config for channel %d not found\n", channel)
	}
	dbl.mu.Unlock()

	dbl.Stop(channel)
	dbl.active[channel] = next
	next.SetStripConfig(config)
	next.Start()
	fmt.Printf("Switched to effect %v\n", name)
	return nil
}

func (dbl *Dasblinken) StartStreaming() {

}
