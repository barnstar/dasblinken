package effects

import (
	"encoding/json"
	"os"

	. "barnstar.com/dasblinken"
)

// Most effets were scaled around a 144 1m led strip
// So some of the properties are scaled as a proportion of the
var sf = 0.2

type RegisterFn func(effect Effect)

type EffectsConfig struct {
	Balls      []BallsConfig      `json:"balls"`
	Race       []RaceConfig       `json:"race"`
	Wave       []WaveConfig       `json:"wave"`
	Chase      []ChaseConfig      `json:"chase"`
	Fire       []FireConfig       `json:"fire"`
	FireMatrix []FireMatrixConfig `json:"fireMatrix"`
	Snow       []SnowConfig       `json:"snow"`
	Solid      []SolidConfig      `json:"solid"`
	TextScroll []TextScrollConfig `json:"textScroll"`
	FontTest   []FontTestConfig   `json:"fontTest"`
	Static     []StaticConfig     `json:"static"`
	Clock      []ClockConfig      `json:"clock"`
	Flicker    []FlickerConfig    `json:"flicker"`
}

func getPalette(name string) func(float64) RGB {
	switch name {
	case "rainbow":
		return RainbowPalette
	case "heat":
		return HeatPalette
	case "cold":
		return ColdPalette
	case "green":
		return GreenFire
	case "ice":
		return IcePalette
	case "festive":
		return FestivePalette
	case "fullfire":
		return FullFirePalette
	default:
		return nil
	}
}

func getTopology(name string) Topology {
	switch name {
	case "linear":
		return Linear
	case "matrix":
		return Matrix
	case "any":
		return Any
	default:
		return Any
	}
}

func getColorTransform(name string) ColorTransform {
	switch name {
	case "rotate":
		return Rotate
	case "random":
		return RandomColor
	default:
		return Rotate
	}
}

func LoadEffectsFromFile(filename string, registerFn RegisterFn, config StripConfig) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var effectsConfig EffectsConfig
	if err := json.Unmarshal(data, &effectsConfig); err != nil {
		return err
	}

	registerEffects(registerFn, config, effectsConfig)
	return nil
}

func registerEffects(registerFn RegisterFn, stripCfg StripConfig, ec EffectsConfig) {
	// Register Balls effects
	for _, cfg := range ec.Balls {
		registerFn(NewBallsEffect(cfg, stripCfg))
	}

	// Register Race effects
	for _, cfg := range ec.Race {
		registerFn(NewRaceEffect(cfg, stripCfg))
	}

	// Register Wave effects
	for _, cfg := range ec.Wave {
		registerFn(NewWaveEffect(cfg, stripCfg))
	}

	// Register Chase effects
	for _, cfg := range ec.Chase {
		registerFn(NewChaseEffect(cfg, stripCfg))
	}

	// Register Fire effects
	for _, cfg := range ec.Fire {
		registerFn(NewFireEffect(cfg, stripCfg))
	}

	// Register FireMatrix effects
	for _, cfg := range ec.FireMatrix {
		registerFn(NewFireMatrixEffect(cfg, stripCfg))
	}

	// Register Snow effects
	for _, cfg := range ec.Snow {
		registerFn(NewSnowEffect(cfg, stripCfg))
	}

	// Register Solid effects
	for _, cfg := range ec.Solid {
		registerFn(NewSolidEffect(cfg, stripCfg))
	}

	// Register TextScroll effects
	for _, cfg := range ec.TextScroll {
		registerFn(NewTextScrollEffect(cfg, stripCfg))
	}

	// Register FontTest effects
	for _, cfg := range ec.FontTest {
		registerFn(NewFontTestEffect(cfg, stripCfg))
	}

	// Register Static effects
	for _, cfg := range ec.Static {
		registerFn(NewStaticEffect(cfg, stripCfg))
	}

	// Register Clock effects
	for _, cfg := range ec.Clock {
		registerFn(NewClockEffect(cfg, stripCfg))
	}

	// Register Flicker effects
	for _, cfg := range ec.Flicker {
		registerFn(NewFlickerEffect(cfg, stripCfg))
	}
}

func RegisterDefaultEffects(f RegisterFn, config StripConfig) {
	// Try to load from file first, fall back to hardcoded defaults
	err := LoadEffectsFromFile("effects.json", f, config)
	if err != nil {
		// Fall back to hardcoded defaults
		registerDefaultEffectsHardcoded(f, config)
	}
}

func registerDefaultEffectsHardcoded(f RegisterFn, config StripConfig) {
	// ... existing hardcoded effects as fallback ...
}
