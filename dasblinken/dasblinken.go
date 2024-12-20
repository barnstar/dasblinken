package dasblinken

import (
	"fmt"
	"sync/atomic"
	"time"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

type Dasblinken struct {
	active  Effect
	effects []Effect
}

type wsEngine interface {
	Init() error
	Render() error
	Wait() error

	Fini()
	Leds(channel int) []uint32
}

type EffectsOpts struct {
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
	fmt.Println("Configuring Device at Pin:", opts.Pin, "Channel:", opts.Channel, "LedCount:", opts.LedCount, "Brightness:", opts.Brightness)

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
}

func (dbl *Dasblinken) RegisterEffect(effect Effect) int {
	dbl.effects = append(dbl.effects, effect)
	return len(dbl.effects)
}

const (
	defaultChan       = 0
	defaultPin        = 21
	defaultLen        = 64
	defaultBrightness = 128
)

func (dbl *Dasblinken) RegisterTestEffects() {
	for i := 0; i < 10; i++ {
		effectOpts := EffectsOpts{
			defaultPin,
			defaultChan,
			defaultLen,
			defaultBrightness,
			time.Duration(10000000),
		}
		effect := NewWipeEffect(WipeEffectOpts{effectOpts, i + 4})
		dbl.RegisterEffect(effect)
	}
}

func (dbl *Dasblinken) SwitchToEffect(index int) {
	if index < 0 || index > len(dbl.effects) {
		return
	}
	dbl.Stop()
	dbl.active = dbl.effects[index]
	dbl.active.Start()
	fmt.Printf("Switched to effect %d\n", index)
}
