package effects

import (
	. "barnstar.com/dasblinken"
)

// Most effets were scaled around a 144 1m led strip
// So some of the properties are scaled as a proportion of the
var sf = 0.2

func RegisterDefaultEffects(dbl *Dasblinken, channel Channel) {

	config, ok := dbl.Config(channel)

	//Scaling factor
	if !ok {
		panic("No default strip configuration")
	}

	balls10 := NewBallsEffect(
		BallsEffectOpts{
			StripOptsDefString("10 Balls", config),
			10,
			30,
			RainbowPalette,
		})
	dbl.RegisterEffect(balls10)

	balls20 := NewBallsEffect(
		BallsEffectOpts{
			StripOptsDefString("20 Balls", config),
			20,
			20,
			RainbowPalette,
		})
	dbl.RegisterEffect(balls20)

	race1 := NewRaceEffect(
		RaceEffectOpts{
			StripOptsDefString("Single Race", config),
			18,
			false,
			4,
		})
	dbl.RegisterEffect(race1)

	race2 := NewRaceEffect(
		RaceEffectOpts{StripOptsDefString("Double Race", config),
			18,
			true,
			4,
		})
	dbl.RegisterEffect(race2)

	wave := NewWaveEffect(
		WaveEffectOpts{StripOptsDefString("Wave Chase", config),
			RainbowPalette,
		})
	dbl.RegisterEffect(wave)

	chase := NewRainbowChaseEffect(
		ChaseEffectOpts{StripOptsDefString("Rainbow Chase", config),
			0.25,
		})
	dbl.RegisterEffect(chase)

	fire := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Fire", config),
			0.3 * sf,
			0.02 / sf,
			false,
			HeatPalette,
		})
	dbl.RegisterEffect(fire)

	mfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{StripOptsDefString("Fire Matrix", config),
			0.7 * sf,
			0.02 / sf,
			HeatPalette,
		})
	dbl.RegisterEffect(mfire)

	gfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{StripOptsDefString("Fire Matrix (Green)", config),
			0.7 * sf,
			0.02 / sf,
			GreenFire,
		})
	dbl.RegisterEffect(gfire)

	cfire := NewFireMatrixEffect(
		FireMatrixEffectOpts{StripOptsDefString("Fire Matrix (Blue)", config),
			0.7 * sf,
			0.02 / sf,
			ColdPalette,
		})
	dbl.RegisterEffect(cfire)

	fire2 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Fire 2", config),
			0.4 * sf,
			0.03 / sf,
			false,
			HeatPalette,
		})
	dbl.RegisterEffect(fire2)

	fire3 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Double Fire", config),
			0.3 * sf,
			0.04 / sf,
			true,
			HeatPalette,
		})
	dbl.RegisterEffect(fire3)

	fire4 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Double Fire 2", config),
			0.4 * sf,
			0.05 / sf,
			true,
			HeatPalette,
		})
	dbl.RegisterEffect(fire4)

	fire5 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Cold Fire", config),
			0.4 * sf,
			0.04 / sf,
			false,
			ColdPalette,
		})
	dbl.RegisterEffect(fire5)

	fire6 := NewFireEffect(
		FireEffectOpts{StripOptsDefString("Double Cold Fire", config),
			0.3 * sf,
			0.04 / sf,
			true,
			ColdPalette,
		})
	dbl.RegisterEffect(fire6)

	heavySnow := NewSnowEffect(
		SnowEffectOpts{StripOptsDefString("Heavy Snow", config),
			0.995,
			0.3 * sf,
		})
	dbl.RegisterEffect(heavySnow)

	lightSnow := NewSnowEffect(
		SnowEffectOpts{StripOptsDefString("Light Snow", config),
			0.995,
			0.1 * sf,
		})
	dbl.RegisterEffect(lightSnow)

	rotation := NewSolidEffect(
		SolidEffectOpts{StripOptsDefString("Rotating Rainbow", config),
			240,
			RainbowPalette,
			RandomColor,
		})
	dbl.RegisterEffect(rotation)

	rotation2 := NewSolidEffect(
		SolidEffectOpts{StripOptsDefString("Rotating Heat", config),
			4,
			RainbowPalette,
			Rotate,
		})
	dbl.RegisterEffect(rotation2)

	marquee := NewTextScrollEffect(
		TextScrollEffectOpts{StripOptsDefString("Marquee", config),
			"Hello World!",
			RainbowPalette,
		})
	dbl.RegisterEffect(marquee)

	font := NewFontTestEffect(
		FontTestEffectOpts{StripOptsDefString("Font Test", config),
			RainbowPalette,
		})
	dbl.RegisterEffect(font)

	static := NewStaticEffect(
		StaticEffectOpts{StripOptsDefString("Static", config)})
	dbl.RegisterEffect(static)

	clock := NewClockEffect(
		ClockEffectOpts{StripOptsDefString("Clock", config)})
	dbl.RegisterEffect(clock)
}
