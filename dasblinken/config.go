package dasblinken

import (
	"encoding/json"
	"os"
)

// StripConfig holds the configuration for a strip of LEDs.
type StripConfig struct {
	Pin        int     `json:"pin"`        // The GPIO Pin (BCM) to which the strip is connected
	Channel    Channel `json:"channel"`    // The channel number for the strip
	Brightness int     `json:"brightness"` // The brightness level of the strip
	Width      int     `json:"width"`      // The width of the strip
	Height     int     `json:"height"`     // The height of the strip
	Fps        int     `json:"fps"`        // The desired framerate
	Invert     bool    `json:"invert"`     // Invert the pwm signal
}

func NewStripConfig(pin, channel, width, height, brightness, fps int) StripConfig {
	return StripConfig{pin, Channel(channel), brightness, width, height, fps, false}
}

func LoadStripConfig(path string) (StripConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return StripConfig{}, err
	}
	var s StripConfig
	err = json.Unmarshal(data, &s)
	if err != nil {
		return StripConfig{}, err
	}
	return s, nil
}

func (s StripConfig) SaveTo(path string) error {
	jsondata, err := json.Marshal(s)
	if err != nil {
		return err
	}
	os.WriteFile(path, jsondata, 0644)
	return nil
}

func (s StripConfig) Len() int {
	return s.Width * s.Height
}
