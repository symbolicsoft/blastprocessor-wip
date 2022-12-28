/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SquareClear struct {
	Cod      [2]int
	State    int
	SubState float64
	Opacity  float64
	Color    color.RGBA
	OnClear  func([2]int)
}

func (squareClear *SquareClear) Update(tic int) {
	switch squareClear.State {
	case 0:
		if squareClear.SubState >= 100 {
			squareClear.OnClear(squareClear.Cod)
			squareClear.State = 1
			squareClear.SubState = 0
		} else {
			squareClear.SubState += 12
			squareClear.Opacity = squareClear.SubState / 100
		}
	case 1:
		if squareClear.SubState >= 100 {
			squareClear.State = 2
			squareClear.SubState = 0
		} else {
			squareClear.SubState += 12
			squareClear.Opacity = 1
		}
	case 2:
		squareClear.SubState += 12
		squareClear.Opacity = 1 - (squareClear.SubState / 100)
	}
}

func (squareClear *SquareClear) Draw(screen *ebiten.Image, vectorBuffer *ebiten.Image) {
	var path vector.Path
	path.MoveTo(float32(boardRenderer.TileSize)*float32(squareClear.Cod[1]), float32(boardRenderer.TileSize)*float32(squareClear.Cod[0]))
	path.LineTo(float32(boardRenderer.TileSize)*float32(squareClear.Cod[1]+1), float32(boardRenderer.TileSize)*float32(squareClear.Cod[0]))
	path.LineTo(float32(boardRenderer.TileSize)*float32(squareClear.Cod[1]+1), float32(boardRenderer.TileSize)*float32(squareClear.Cod[0]+1))
	path.LineTo(float32(boardRenderer.TileSize)*float32(squareClear.Cod[1]), float32(boardRenderer.TileSize)*float32(squareClear.Cod[0]+1))
	path.LineTo(float32(boardRenderer.TileSize)*float32(squareClear.Cod[1]), float32(boardRenderer.TileSize)*float32(squareClear.Cod[0]))
	opVector := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.FillAll,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].ColorR = float32(squareClear.Color.R) / float32(0xff)
		vs[i].ColorG = float32(squareClear.Color.G) / float32(0xff)
		vs[i].ColorB = float32(squareClear.Color.B) / float32(0xff)
		vs[i].ColorA = float32(squareClear.Color.A) * float32(squareClear.Opacity) / float32(0xff)
	}
	screen.DrawTriangles(vs, is, vectorBuffer.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), opVector)
}
