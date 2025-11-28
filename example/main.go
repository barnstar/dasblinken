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
	configFile := flag.String("config", "", "Config file path")
	effectsDef := flag.String("effects", "effects.json", "Effects definition file")

	flag.Parse()

	// Create a new Dasblinken instance
	das := dasblinken.NewDasblinken()

	if *configFile != "" {
		// Load the single strip configuration
		config, err := dasblinken.LoadStripConfig(*configFile)
		if err != nil {
			fmt.Printf("Unable to load config from %s: %s\n", *configFile, err)
			return
		}
		fmt.Printf("Loaded strip config: %+v\n", config)
		das.SetStrip(config)
		// Load effects for the strip
		effects.LoadEffectsFromFile(*effectsDef, das.RegisterEffect, config)
	} else {
		fmt.Printf("No config file specified\n")
		return
	}

	defer func() {
		das.StopAll()
	}()

	s := &LedControlServer{
		Das:         das,
		Hostname:    *tsname,
		AuthKey:     *authkey,
		ConfigFile:  *configFile,
		EffectsFile: *effectsDef,
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
