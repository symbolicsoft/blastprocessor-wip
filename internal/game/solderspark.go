/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SolderSpark struct {
	X         float64
	Y         float64
	OffsetX   float64
	OffsetY   float64
	ModifierY float64
	Opacity   float64
	Speed     float64
	Size      int
	Color     color.RGBA
}

func (solderSpark *SolderSpark) GenSparkPair(sparkColor color.RGBA, colorMutationFactor int, cod [2]int) SolderSpark {
	mutatedColor := color.RGBA{
		R: sparkColor.R - uint8(colorMutationFactor/2) + uint8(rand.Intn(colorMutationFactor)),
		G: sparkColor.G - uint8(colorMutationFactor/2) + uint8(rand.Intn(colorMutationFactor)),
		B: sparkColor.B - uint8(colorMutationFactor/2) + uint8(rand.Intn(colorMutationFactor)),
		A: sparkColor.A,
	}
	solderSpark.X = (boardRenderer.BoardOffset[0] + float64(boardRenderer.TileSize)*float64(cod[1]) + (float64(boardRenderer.TileSize) / 4)) * GameScale
	solderSpark.Y = (boardRenderer.BoardOffset[1] + float64(boardRenderer.TileSize)*float64(cod[0]) + (float64(boardRenderer.TileSize) / 1.5) - rand.Float64()*10) * GameScale
	solderSpark.OffsetX = (-5 + rand.Float64()*10)
	solderSpark.ModifierY = math.Pow((rand.Float64()*50), 2) * math.Pow((rand.Float64()*10), 2)
	solderSpark.OffsetY = (math.Pow(solderSpark.OffsetX, 2)) / solderSpark.ModifierY
	solderSpark.Opacity = 1
	solderSpark.Speed = -2 + rand.Float64()
	solderSpark.Size = 1 * int(GameScale)
	solderSpark.Color = mutatedColor
	sparkB := SolderSpark{
		X:         (boardRenderer.BoardOffset[0] + float64(boardRenderer.TileSize)*float64(cod[1]) + (float64(boardRenderer.TileSize) / 1.25)) * GameScale,
		Y:         (boardRenderer.BoardOffset[1] + solderSpark.Y),
		OffsetX:   solderSpark.OffsetX,
		OffsetY:   0,
		ModifierY: solderSpark.ModifierY,
		Opacity:   1,
		Speed:     -solderSpark.Speed,
		Size:      solderSpark.Size,
		Color:     mutatedColor,
	}
	sparkB.OffsetY = (math.Pow(sparkB.OffsetX, 2)) / sparkB.ModifierY
	return sparkB
}

func (solderSpark *SolderSpark) Update(tic int) {
	solderSpark.OffsetX += solderSpark.Speed
	solderSpark.Opacity -= 0.02
	solderSpark.OffsetY = math.Pow(solderSpark.OffsetX, 4) / solderSpark.ModifierY
}

func (solderSpark *SolderSpark) Draw(screen *ebiten.Image) {
	var path vector.Path
	path.MoveTo(float32(solderSpark.X+solderSpark.OffsetX), float32(solderSpark.Y+solderSpark.OffsetY))
	path.LineTo(float32(solderSpark.X+solderSpark.OffsetX)+float32(solderSpark.Size), float32(solderSpark.Y+solderSpark.OffsetY))
	path.LineTo(float32(solderSpark.X+solderSpark.OffsetX)+float32(solderSpark.Size), float32(solderSpark.Y+solderSpark.OffsetY)+float32(solderSpark.Size))
	path.LineTo(float32(solderSpark.X+solderSpark.OffsetX), float32(solderSpark.Y+solderSpark.OffsetY)+float32(solderSpark.Size))
	path.LineTo(float32(solderSpark.X+solderSpark.OffsetX), float32(solderSpark.Y+solderSpark.OffsetY))
	opVector := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.FillAll,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].ColorR = float32(solderSpark.Color.R) / float32(0xff)
		vs[i].ColorG = float32(solderSpark.Color.G) / float32(0xff)
		vs[i].ColorB = float32(solderSpark.Color.B) / float32(0xff)
		vs[i].ColorA = float32(solderSpark.Color.A) * float32(solderSpark.Opacity) / float32(0xff)
	}
	screen.DrawTriangles(vs, is, boardRenderer.VectorBuffer.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), opVector)
}
