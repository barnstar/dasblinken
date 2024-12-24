package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dasblinken "barnstar.com/piled/dasblinken"
	server "barnstar.com/piled/server"
)

func main() {
	das := dasblinken.NewDasblinken()
	//Channel 0, pin 21, 144 LEDs, 128 brightness
	das.AddStrip(0, 21, 32, 8, 128)
	das.RegisterTestEffects()

	defer func() {
		das.StopAll()
	}()

	runServer(das)

	exitOnSignal()
}

func exitOnSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Printf("Received signal: %s. Shutting down...\n", sig)
}

func runServer(das *dasblinken.Dasblinken) {
	s := &server.LedControlServer{
		EffectHandler: das.SwitchToEffect,
		StopHandler:   das.Stop,
		EffectFetcher: das.Effects,
	}
	go s.RunServer()
}
