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
	authkey := flag.String("authkey", "", "Tailscale auth key")
	configFile := flag.String("config", "config.json", "Config file path")
	effectsDef := flag.String("effects", "effects.json", "Effects definition file")

	flag.Parse()

	das := dasblinken.NewDasblinken()

	config, err := dasblinken.LoadStripConfig(*configFile)
	if err != nil {
		fmt.Printf("Unable to load config from %s: %s\n", *configFile, err)
		return
	}
	fmt.Printf("Loaded strip config: %+v\n", config)
	das.SetStrip(config)
	effects.LoadEffectsFromFile(*effectsDef, das.RegisterEffect, config)

	if *authkey == "" {
		keybytes, err := os.ReadFile("authkey.txt")
		if err == nil {
			*authkey = string(keybytes)
		}
	}

	if *authkey == "" {
		fmt.Println("No Tailscale auth key found or provided, exiting (use --authkey or authkey.json).")
		return
	}

	if config.Hostname == "" {
		fmt.Println("No hostname specified in config, exiting.")
		return
	}

	defer func() {
		das.StopAll()
	}()

	fmt.Printf("Starting server with configuration %+v\n", config)

	s := &LedControlServer{
		Das:         das,
		Hostname:    config.Hostname,
		AuthKey:     *authkey,
		ConfigFile:  *configFile,
		EffectsFile: *effectsDef,
	}
	go s.RunServer()

	exitOnSignal()
}

func exitOnSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Printf("Received signal: %s. Shutting down...\n", sig)
}
