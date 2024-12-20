package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dasblinken "barnstar.com/piled/dasblinken"
	server "barnstar.com/piled/server"
)

func main() {
	das := dasblinken.NewDasblinken()
	das.RegisterTestEffects()

	reader := bufio.NewReaderSize(os.Stdin, 1)
	defer func() {
		das.Stop()
	}()

	input := make(chan rune)

	go readKey(reader, input)
	go handleKeyInput(input, das)
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
	}
	go s.RunServer()
}

func handleKeyInput(input chan rune, das *dasblinken.Dasblinken) {
	for {
		select {
		case i := <-input:
			switch i {
			case 'q':
				p, err := os.FindProcess(os.Getpid())
				if err != nil {
					fmt.Printf("Error finding process: %s\n", err)
					continue
				}
				p.Signal(syscall.SIGINT)
			case 's':
				das.Stop()
			case 'n':
				das.SwitchToEffect(0)
			}
		}
	}
}

func readKey(reader *bufio.Reader, input chan rune) {
	for {
		char, _, _ := reader.ReadRune()
		input <- char
	}
}
