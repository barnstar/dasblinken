module barnstar.com/piled

go 1.23.4

require barnstar.com/dasblinken v0.0.0

require barnstar.com/effects v0.0.0

require (
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rpi-ws281x/rpi-ws281x-go v1.0.10 // indirect
)

replace barnstar.com/dasblinken => ./../dasblinken

replace barnstar.com/effects => ./../effects
