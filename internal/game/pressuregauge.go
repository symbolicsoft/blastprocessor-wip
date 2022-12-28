/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type PressureGauge struct {
	X       int
	Y       int
	Level   float32
	Stopped bool
	Image   *ebiten.Image
}

var pressureGauge = func() PressureGauge {
	pressureGauge := PressureGauge{
		X:       77,
		Y:       100,
		Level:   0,
		Stopped: false,
		Image:   ebiten.NewImageFromImage(assetsGFX["pressureGauge"]),
	}
	return pressureGauge
}()

func (pressureGauge *PressureGauge) Update(tic int) {
	if pressureGauge.Stopped {
		return
	}
	if pressureGauge.Level < 100 {
		pressureGauge.Level += 0.1
	}
	if pressureGauge.Level >= 100 {
		if !pressureGauge.Stopped {
			session.GameOver()
		}
	}
}

func (PressureGauge *PressureGauge) ChipPlaced() {
	pressureGauge.Level -= 20
	if pressureGauge.Level < 0 {
		pressureGauge.Level = 0
	}
}

func (PressureGauge *PressureGauge) ChipMisplaced() {
	pressureGauge.Level += 1
	if pressureGauge.Level > 100 {
		pressureGauge.Level = 100
	}
}

func (PressureGauge *PressureGauge) LineCleared() {
	pressureGauge.Level = 0
}

func (pressureGauge *PressureGauge) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(pressureGauge.X), float64(pressureGauge.Y))
	op.GeoM.Scale(GameScale, GameScale)
	screen.DrawImage(pressureGauge.Image.SubImage(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(pressureGauge.Level * 1.28), Y: 16},
	}).(*ebiten.Image), op)
}
