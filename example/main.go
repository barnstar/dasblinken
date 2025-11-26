package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dasblinken "barnstar.com/dasblinken"
	effects "barnstar.com/effects"
)

var ()

func main() {
	tsname := flag.String("tsname", "dasblinken", "Tailscale hostname")
	authkey := flag.String("authkey", "", "Tailscale auth key")
	channel := flag.Int("c", 0, "Channel number")
	pin := flag.Int("p", 21, "Pin number")
	width := flag.Int("w", 32, "Width (length) of LED matrix")
	height := flag.Int("h", 8, "Height of LED matrix")
	fps := flag.Int("f", 60, "FPS")
	brightness := flag.Int("brightness", 128, "Brightness level")
	configFile := flag.String("config", "", "Config File")
	effectsDef := flag.String("effects", "effects.json", "Effects definition file")

	flag.Parse()

	// Create a new Dasblinken instance
	das := dasblinken.NewDasblinken()

	if *configFile != "" {
		sc, err := dasblinken.LoadStripConfig(*configFile)
		if err != nil {
			fmt.Printf("Error loading config file: %s\n", err)
			return
		}
		das.AddStrip(sc)
	} else {
		sc := dasblinken.NewStripConfig(*pin, *channel, *width, *height, *brightness, *fps)
		das.AddStrip(sc)
	}
	ch := dasblinken.Channel(*channel)
	config, ok := das.Config(ch)
	if !ok {
		fmt.Printf("No default strip configuration")
		os.Exit(1)
	}

	effects.LoadEffectsFromFile(*effectsDef, das.RegisterEffect, config)

	defer func() {
		das.StopAll()
	}()

	s := &LedControlServer{
		EffectHandler: das.SwitchToEffect,
		StopHandler:   das.Stop,
		EffectFetcher: das.Effects,
		Hostname:      *tsname,
		AuthKey:       *authkey,
	}
	go s.RunServer()

	// Wait for a signal to exit
	exitOnSignal()
}

func exitOnSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Printf("Received signal: %s. Shutting down...\n", sig)
}
