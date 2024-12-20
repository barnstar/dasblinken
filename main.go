package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"

	"barnstar.com/piled/dasblinken"
)

func main() {
	fmt.Println("Starting dasblinken! n -> next, s -> stop, q -> quit")
	reader := bufio.NewReader(os.Stdin)
	var effect dasblinken.Effect
	defer func() {
		if effect != nil {
			effect.Stop()
		}
	}()
	input := make(chan rune)
	go readKey(reader, input)

inputLoop:
	for {
		select {
		case i := <-input:
			switch i {
			case 'q':
				break inputLoop
			case 's':
				if effect != nil {
					effect.Stop()
				}
				effect = nil
			case 'n':
				effect = dasblinken.RunWipeEffect(21, int(rand.Float32()*8)+4)
			}
		}
	}

	fmt.Println("Bye!!")
}

func readKey(reader *bufio.Reader, input chan rune) {
	for {
		char, _, _ := reader.ReadRune()
		input <- char
	}
}
