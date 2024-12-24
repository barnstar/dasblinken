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
	active  map[int]Effect
	effects map[string]Effect
	strips  map[int]StripConfig
}

func NewDasblinken() *Dasblinken {
	fmt.Println("Starting dasblinken! n -> next, s -> stop, q -> quit")
	return &Dasblinken{
		effects: make(map[string]Effect),
		active:  make(map[int]Effect),
		strips:  make(map[int]StripConfig),
	}
}

func (dbl *Dasblinken) AddStrip(
	channel int,
	pin int,
	width int,
	height int,
	brightness int) {
	dbl.strips[channel] = NewStripConfig(pin, channel, width, height, brightness)
}

type wsEngine interface {
	Init() error
	Render() error
	Wait() error

	Fini()
	Leds(channel int) []uint32
}

type StripConfig struct {
	Pin        int
	Channel    int
	LedCount   int
	Brightness int
	Width      int
	Height     int
}

func NewStripConfig(pin, channel, width, height, brightness int) StripConfig {
	return StripConfig{pin, channel, width * height, brightness, width, height}
}

type EffectsOpts struct {
	Name string
	StripConfig
	FrameTime time.Duration
}

type Effect interface {
	Start() error
	Stop()

	run(wsEngine)
	engine() wsEngine
	Opts() EffectsOpts
	SetStripConfig(StripConfig)
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

func (dbl *Dasblinken) ActiveEffect(channel int) Effect {
	return dbl.active[channel]
}

func (dbl *Dasblinken) StopAll() {
	for _, effect := range dbl.active {
		effect.Stop()
	}
	dbl.active = make(map[int]Effect)
	fmt.Printf("Dasblinken is kaput")

}

func (dbl *Dasblinken) Stop(channel int) {

	if dbl.active[channel] != nil {
		dbl.active[channel].Stop()
	}
	dbl.active[channel] = nil
	fmt.Printf("Dasblinken Channel %v is kaput", channel)
}

func (dbl *Dasblinken) RegisterEffect(effect Effect) {
	dbl.effects[effect.Opts().Name] = effect
}

const (
	defaultChan       = 0
	defaultPin        = 21
	defaultLen        = 144
	defaultBrightness = 128
	defaultfps        = 40
)

var stringLen = 32

func stripOptsDefString(name string, config StripConfig) EffectsOpts {
	// 60 fps
	frameTime := 1000000000 / defaultfps
	return EffectsOpts{
		name,
		config,
		time.Duration(frameTime),
	}
}

func (dbl *Dasblinken) RegisterTestEffects() {

	//Scaling factor
	sf := float64(stringLen) / float64(defaultLen)
	config, ok := dbl.strips[defaultChan]
	if !ok {
		panic("No default strip configuration")
	}

	balls10 := NewBallsEffect(
		BallsEffectOpts{
			stripOptsDefString("10 Balls", config),
			10,
			30,
			rainbowPalette,
		})
	dbl.RegisterEffect(balls10)

	balls20 := NewBallsEffect(
		BallsEffectOpts{
			stripOptsDefString("20 Balls", config),
			20,
			20,
			rainbowPalette,
		})
	dbl.RegisterEffect(balls20)

	race1 := NewRaceEffect(
		RaceEffectOpts{
			stripOptsDefString("Single Race", config),
			18,
			false,
			4,
		})
	dbl.RegisterEffect(race1)

	race2 := NewRaceEffect(
		RaceEffectOpts{stripOptsDefString("Double Race", config),
			18,
			true,
			4,
		})
	dbl.RegisterEffect(race2)

	chase := NewRainbowChaseEffect(
		ChaseEffectOpts{stripOptsDefString("Rainbow Chase", config),
			0.25,
		})
	dbl.RegisterEffect(chase)

	fire := NewFireEffect(
		FireEffectOpts{stripOptsDefString("Fire", config),
			0.3 * sf,
			0.02 / sf,
			false,
			heatPalette,
		})
	dbl.RegisterEffect(fire)

	mfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{stripOptsDefString("Fire Matrix", config),
			0.5 * sf,
			0.025 / sf,
			heatPalette,
		})
	dbl.RegisterEffect(mfire)

	gfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{stripOptsDefString("Fire Matrix (Greed)", config),
			0.5 * sf,
			0.025 / sf,
			greenFire,
		})
	dbl.RegisterEffect(gfire)

	cfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{stripOptsDefString("Fire Matrix (Blue)", config),
			0.5 * sf,
			0.025 / sf,
			coldPalette,
		})
	dbl.RegisterEffect(cfire)

	fire2 := NewFireEffect(
		FireEffectOpts{stripOptsDefString("Fire 2", config),
			0.4 * sf,
			0.03 / sf,
			false,
			heatPalette,
		})
	dbl.RegisterEffect(fire2)

	fire3 := NewFireEffect(
		FireEffectOpts{stripOptsDefString("Double Fire", config),
			0.3 * sf,
			0.04 / sf,
			true,
			heatPalette,
		})
	dbl.RegisterEffect(fire3)

	fire4 := NewFireEffect(
		FireEffectOpts{stripOptsDefString("Double Fire 2", config),
			0.4 * sf,
			0.05 / sf,
			true,
			heatPalette,
		})
	dbl.RegisterEffect(fire4)

	fire5 := NewFireEffect(
		FireEffectOpts{stripOptsDefString("Cold Fire", config),
			0.4 * sf,
			0.04 / sf,
			false,
			coldPalette,
		})
	dbl.RegisterEffect(fire5)

	fire6 := NewFireEffect(
		FireEffectOpts{stripOptsDefString("Double Cold Fire", config),
			0.3 * sf,
			0.04 / sf,
			true,
			coldPalette,
		})
	dbl.RegisterEffect(fire6)

	heavySnow := NewSnowEffect(
		SnowEffectOpts{stripOptsDefString("Heavy Snow", config),
			0.995,
			0.3 * sf,
		})
	dbl.RegisterEffect(heavySnow)

	lightSnow := NewSnowEffect(
		SnowEffectOpts{stripOptsDefString("Light Snow", config),
			0.995,
			0.1 * sf,
		})
	dbl.RegisterEffect(lightSnow)

	rotation := NewSolidEffect(
		SolidEffectOpts{stripOptsDefString("Rotating Rainbow", config),
			240,
			rainbowPalette,
			randomColor,
		})
	dbl.RegisterEffect(rotation)

	rotation2 := NewSolidEffect(
		SolidEffectOpts{stripOptsDefString("Rotating Heat", config),
			4,
			rainbowPalette,
			rotate,
		})
	dbl.RegisterEffect(rotation2)

	static := NewStaticEffect(
		StaticEffectOpts{stripOptsDefString("Static", config)})
	dbl.RegisterEffect(static)
}

func (dbl *Dasblinken) RandomEffect(channel int) {
	effects := dbl.Effects()
	if len(effects) > 0 {
		randomIndex := rand.Int() % len(effects)
		effect := effects[randomIndex]
		dbl.SwitchToEffect(effect.Opts().Name, channel)
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

func (dbl *Dasblinken) SwitchToEffect(name string, channel int) error {
	next := dbl.effects[name]
	if nil == next {
		return fmt.Errorf("Effect %s not found\n", name)
	}

	config, ok := dbl.strips[channel]
	if !ok {
		return fmt.Errorf("No config for channel %d not found\n", channel)
	}

	dbl.Stop(channel)
	dbl.active[channel] = next
	next.SetStripConfig(config)
	next.Start()
	fmt.Printf("Switched to effect %v\n", name)
	return nil
}

func (dbl *Dasblinken) StartStreaming() {

}
