package dasblinken

import (
	"fmt"
	"sync/atomic"
	"time"

	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

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
}

type Effect interface {
	Start() error
	Stop()

	run(wsEngine)
	engine() wsEngine
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

var activeEffect Effect

func RunWipeEffect(pin int, count int) Effect {
	if activeEffect != nil {
		activeEffect.Stop()
	}
	effectOpts := EffectsOpts{pin, 0, 64, 128}
	activeEffect = NewWipeEffect(WipeEffectOpts{effectOpts, count, time.Duration(10000000)})
	activeEffect.Start()
	return activeEffect
}
