/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type SquareHighlight struct {
	Cod       [2]int
	Frame     int
	HueOffset int
	Opacity   float64
}

func (squareHighlight *SquareHighlight) Show(cod [2]int, hueOffset int) {
	squareHighlight.Cod = cod
	squareHighlight.HueOffset = hueOffset
	squareHighlight.Opacity = 1
}

func (squareHighlight *SquareHighlight) Hide() {
	squareHighlight.Cod = [2]int{-1, -1}
	squareHighlight.Opacity = 0
}

func (squareHighlight *SquareHighlight) Update(tic int) {
	if squareHighlight.Opacity > 0 {
		if tic%6 == 0 {
			squareHighlight.Frame = (squareHighlight.Frame + 1) % 8
		}
	}
}

func (squareHighlight *SquareHighlight) Draw(screen *ebiten.Image) {
	if squareHighlight.Opacity > 0 {
		op := &colorm.DrawImageOptions{}
		co := colorm.ColorM{}
		op.GeoM.Translate(float64(squareHighlight.Cod[1]*boardRenderer.TileSize)+2, float64(squareHighlight.Cod[0]*boardRenderer.TileSize)+2)
		co.ChangeHSV(float64(squareHighlight.HueOffset%360)*2*math.Pi/360, 2, 1)
		colorm.DrawImage(screen, boardRenderer.BracketImages[squareHighlight.Frame], co, op)
	}
}
