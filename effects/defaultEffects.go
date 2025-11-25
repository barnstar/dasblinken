package effects

import (
	. "barnstar.com/dasblinken"
)

// Most effets were scaled around a 144 1m led strip
// So some of the properties are scaled as a proportion of the
var sf = 0.2

type RegisterFn func(effect Effect)

func RegisterDefaultEffects(f RegisterFn, config StripConfig) {

	balls10 := NewBallsEffect(
		BallsEffectOpts{
			StripOptsDefString("10 Balls", config),
			10,
			30,
			RainbowPalette,
		})
	f(balls10)

	balls20 := NewBallsEffect(
		BallsEffectOpts{
			StripOptsDefString("20 Balls", config),
			20,
			20,
			RainbowPalette,
		})
	f(balls20)

	race1 := NewRaceEffect(
		RaceEffectOpts{
			StripOptsDefString("Single Race", config),
			18,
			false,
			4,
		})
	f(race1)

	race2 := NewRaceEffect(
		RaceEffectOpts{StripOptsDefString("Double Race", config),
			18,
			true,
			4,
		})
	f(race2)

	wave := NewWaveEffect(
		WaveEffectOpts{StripOptsDefString("Wave Chase", config),
			RainbowPalette,
		})
	f(wave)

	chase := NewRainbowChaseEffect(
		ChaseEffectOpts{StripOptsDefString("Rainbow Chase", config),
			0.25,
		})
	f(chase)

	fire := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Fire", config),
			0.3 * sf,
			0.02 / sf,
			false,
			HeatPalette,
		})
	f(fire)

	mfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{StripOptsDefString("Fire Matrix", config),
			0.7 * sf,
			0.02 / sf,
			HeatPalette,
		})
	f(mfire)

	gfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{StripOptsDefString("Fire Matrix (Green)", config),
			0.7 * sf,
			0.02 / sf,
			GreenFire,
		})
	f(gfire)

	cfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{StripOptsDefString("Fire Matrix (Blue)", config),
			0.7 * sf,
			0.02 / sf,
			ColdPalette,
		})
	f(cfire)

	fire2 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Fire 2", config),
			0.4 * sf,
			0.03 / sf,
			false,
			HeatPalette,
		})
	f(fire2)

	fire3 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Double Fire", config),
			0.3 * sf,
			0.04 / sf,
			true,
			HeatPalette,
		})
	f(fire3)

	fire4 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Double Fire 2", config),
			0.4 * sf,
			0.05 / sf,
			true,
			HeatPalette,
		})
	f(fire4)

	fire5 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Cold Fire", config),
			0.4 * sf,
			0.04 / sf,
			false,
			ColdPalette,
		})
	f(fire5)

	fire6 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Double Cold Fire", config),
			0.3 * sf,
			0.04 / sf,
			true,
			ColdPalette,
		})
	f(fire6)

	heavySnow := NewSnowEffect(
		SnowEffectOpts{StripOptsDefString("Heavy Snow", config),
			0.995,
			0.3 * sf,
		})
	f(heavySnow)

	lightSnow := NewSnowEffect(
		SnowEffectOpts{StripOptsDefString("Light Snow", config),
			0.995,
			0.1 * sf,
		})
	f(lightSnow)

	rotation := NewSolidEffect(
		SolidEffectOpts{StripOptsDefString("Rotating Rainbow", config),
			240,
			RainbowPalette,
			RandomColor,
		})
	f(rotation)

	rotation2 := NewSolidEffect(
		SolidEffectOpts{StripOptsDefString("Rotating Heat", config),
			4,
			RainbowPalette,
			Rotate,
		})
	f(rotation2)

	marquee := NewTextScrollEffect(
		TextScrollEffectOpts{StripOptsDefString("Marquee", config),
			"Hello World!",
			RainbowPalette,
		})
	f(marquee)

	font := NewFontTestEffect(
		FontTestEffectOpts{StripOptsDefString("Font Test", config),
			RainbowPalette,
		})
	f(font)

	static := NewStaticEffect(
		StaticEffectOpts{StripOptsDefString("Static", config)})
	f(static)

	clock := NewClockEffect(
		ClockEffectOpts{StripOptsDefString("Clock", config)})
	f(clock)
}
