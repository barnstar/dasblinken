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
	fmt.Println("Starting dasblinken! n -> next, s -> stop, q -> quit")

	das := dasblinken.Dasblinken{}
	das.RegisterTestEffects()

	reader := bufio.NewReader(os.Stdin)
	defer func() {
		das.Stop()
	}()
	input := make(chan rune)

	go readKey(reader, input)
	go handleKeyInput(input, &das)

	// Run the server in a separate goroutine
	go func() {
		s := &server.LedControlServer{}
		s.EffectHandler = func(index int) {
			das.SwitchToEffect(index)
		}
		s.RunServer()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Printf("Received signal: %s. Shutting down...\n", sig)
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
