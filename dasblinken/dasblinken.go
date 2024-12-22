package dasblinken

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"sync/atomic"
	"time"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Dasblinken struct {
	active  Effect
	effects map[string]Effect
}

func NewDasblinken() *Dasblinken {
	fmt.Println("Starting dasblinken! n -> next, s -> stop, q -> quit")
	return &Dasblinken{
		effects: make(map[string]Effect),
	}
}

type wsEngine interface {
	Init() error
	Render() error
	Wait() error

	Fini()
	Leds(channel int) []uint32
}

type EffectsOpts struct {
	Name       string
	Pin        int
	Channel    int
	LedCount   int
	Brightness int
	FrameTime  time.Duration
}

type Effect interface {
	Start() error
	Stop()

	run(wsEngine)
	engine() wsEngine
	Opts() EffectsOpts
}

type EffectControl struct {
	running atomic.Bool

	// After setting active to false, wait for this to be closed before
	// calling Fini()
	done chan bool
}

func (ec *EffectControl) stopEffect(engine wsEngine) {
	if ec.running.Load() == false {
		return
	}

	ec.running.Store(false)
	select {
	case <-ec.done:
		engine.Fini()
	}
}

func (ec *EffectControl) startEffect(opts EffectsOpts, e Effect) error {
	if ec.running.Load() {
		return nil
	}

	device, err := device(opts)
	if err != nil {
		return err
	}

	err = device.Init()
	if err != nil {
		fmt.Println("Device could not be initialized")
		return err
	}

	ec.done = make(chan bool, 1)
	ec.running.Store(true)

	go e.run(device)

	return nil
}

func device(opts EffectsOpts) (wsEngine, error) {
	fmt.Println("Configuring Blinken Device at Pin:", opts.Pin, "Channel:", opts.Channel, "LedCount:", opts.LedCount, "Brightness:", opts.Brightness)

	opt := ws281x.DefaultOptions
	opt.Channels[opts.Channel].Brightness = opts.Brightness
	opt.Channels[opts.Channel].LedCount = opts.LedCount
	opt.Channels[opts.Channel].GpioPin = opts.Pin

	return ws281x.MakeWS2811(&opt)
}

func (dbl *Dasblinken) ActiveEffect() Effect {
	return dbl.active
}

func (dbl *Dasblinken) Stop() {
	if dbl.active != nil {
		dbl.active.Stop()
	}
	dbl.active = nil
	fmt.Println("Dasblinken is kaput")
}

func (dbl *Dasblinken) RegisterEffect(effect Effect) {
	dbl.effects[effect.Opts().Name] = effect
}

const (
	defaultChan       = 0
	defaultPin        = 21
	defaultLen        = 144
	defaultBrightness = 128
)

func stripOptsDefString(name string) EffectsOpts {
	// 60 fps
	frameTime := 1000000000 / 60
	return EffectsOpts{
		name,
		21,
		0,
		defaultLen,
		128,
		time.Duration(frameTime),
	}
}

var defaultOpts = stripOptsDefString("default")

func (dbl *Dasblinken) RegisterTestEffects() {

	race1 := NewRaceEffect(RaceEffectOpts{stripOptsDefString(fmt.Sprintf("Single Race")), 18, true})
	dbl.RegisterEffect(race1)

	race2 := NewRaceEffect(RaceEffectOpts{stripOptsDefString(fmt.Sprintf("Double Race")), 18, false})
	dbl.RegisterEffect(race2)

	chase := NewRainbowChaseEffect(ChaseEffectOpts{stripOptsDefString("Rainbow Chase"), 0.25})
	dbl.RegisterEffect(chase)

	fire := NewFireEffect(FireEffectOpts{stripOptsDefString("Fire"), 0.3, 0.02, false, heatColor})
	dbl.RegisterEffect(fire)

	fire2 := NewFireEffect(FireEffectOpts{stripOptsDefString("Fire 2"), 0.4, 0.03, false, heatColor})
	dbl.RegisterEffect(fire2)

	fire3 := NewFireEffect(FireEffectOpts{stripOptsDefString("Double Fire"), 0.3, 0.04, true, heatColor})
	dbl.RegisterEffect(fire3)

	fire4 := NewFireEffect(FireEffectOpts{stripOptsDefString("Double Fire 2"), 0.4, 0.05, true, heatColor})
	dbl.RegisterEffect(fire4)

	fire5 := NewFireEffect(FireEffectOpts{stripOptsDefString("Cold Fire"), 0.4, 0.04, false, coldColor})
	dbl.RegisterEffect(fire5)

	fire6 := NewFireEffect(FireEffectOpts{stripOptsDefString("Double Cold Fire"), 0.3, 0.04, true, coldColor})
	dbl.RegisterEffect(fire6)

	heavySnow := NewSnowEffect(SnowEffectOpts{stripOptsDefString("Heavy Snow"), 0.995, 0.3})
	dbl.RegisterEffect(heavySnow)

	lightSnow := NewSnowEffect(SnowEffectOpts{stripOptsDefString("Light Snow"), 0.995, 0.1})
	dbl.RegisterEffect(lightSnow)

	static := NewStaticEffect(StaticEffectOpts{stripOptsDefString("Static")})
	dbl.RegisterEffect(static)
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
	effectsSlice := make([]Effect, 0, len(dbl.effects))
	for _, effect := range dbl.effects {
		effectsSlice = append(effectsSlice, effect)
	}
	sort.Slice(effectsSlice, func(i, j int) bool {
		return effectsSlice[i].Opts().Name < effectsSlice[j].Opts().Name
	})

	return effectsSlice
}

func (dbl *Dasblinken) SwitchToEffect(name string) error {
	next := dbl.effects[name]
	if nil == next {
		fmt.Printf("Effect %s not found\n", name)
		return fmt.Errorf("Effect not found")
	}
	dbl.Stop()
	dbl.active = next
	dbl.active.Start()
	fmt.Printf("Switched to effect %v\n", name)
	return nil
}
