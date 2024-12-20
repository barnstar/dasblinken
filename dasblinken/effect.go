package dasblinken

import (
	"fmt"
	"sync/atomic"
	"time"
)

// EffectsOpts holds the basic required configuration for an effect.
type EffectsOpts struct {
	Name string
	StripConfig
	FrameTime time.Duration
}

// An Effect is an interface for a light effect.
type Effect interface {
	Start() error
	Stop()

	Run(WSEngine)
	Engine() WSEngine
	Opts() EffectsOpts
	SetStripConfig(StripConfig)
}

// Effects have an associated state.  The effect will run until 
// it detets Runningis false, at which point it should send true
// to the Done channel.  The engine will clean up.
type EffectState struct {
	Running atomic.Bool

	// After setting active to false, wait for this to be closed before
	// calling Fini()
	Done chan bool
}

// Stops the effect safely, waiting for any animations to complete.
// Cleans up the engine.  The engine should match the provided effect
// state.
func (ec *EffectState) StopEffect(engine WSEngine) {
	if ec.Running.Load() == false {
		return
	}

	ec.Running.Store(false)
	select {
	case <-ec.Done:
		engine.Fini()
	}
}

func (ec *EffectState) StartEffect(opts EffectsOpts, e Effect) error {
	if ec.Running.Load() {
		return nil
	}

	device, err := Device(opts)
	if err != nil {
		return err
	}

	err = device.Init()
	if err != nil {
		fmt.Println("Device could not be initialized")
		return err
	}

	ec.Done = make(chan bool, 1)
	ec.Running.Store(true)

	go e.Run(device)

	return nil
}
