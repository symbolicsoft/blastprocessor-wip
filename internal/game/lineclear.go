/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"image"

	"blastprocessor.app/internal/board"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type LineClear struct {
	IsRow    bool
	Cod      int
	State    int
	SubState float64
	Opacity  float64
	OnClear  func(int)
}

func (lineClear *LineClear) Update(tic int) {
	switch lineClear.State {
	case 0:
		if lineClear.SubState >= 100 {
			lineClear.OnClear(lineClear.Cod)
			lineClear.State = 1
			lineClear.SubState = 0
		} else {
			lineClear.SubState += 6
			lineClear.Opacity = lineClear.SubState / 100
		}
	case 1:
		if lineClear.SubState >= 100 {
			lineClear.State = 2
			lineClear.SubState = 0
		} else {
			lineClear.SubState += 6
			lineClear.Opacity = 1
		}
	case 2:
		lineClear.SubState += 6
		lineClear.Opacity = 1 - (lineClear.SubState / 100)
	}
}

func (lineClear *LineClear) Draw(screen *ebiten.Image, vectorBuffer *ebiten.Image) {
	var path vector.Path
	switch lineClear.IsRow {
	case true:
		path.MoveTo(0, float32(boardRenderer.TileSize)*float32(lineClear.Cod))
		path.LineTo(float32(boardRenderer.TileSize*board.Board.Size[1]), float32(boardRenderer.TileSize)*float32(lineClear.Cod))
		path.LineTo(float32(boardRenderer.TileSize*board.Board.Size[1]), float32(boardRenderer.TileSize)*float32(lineClear.Cod+1))
		path.LineTo(0, float32(boardRenderer.TileSize)*float32(lineClear.Cod+1))
		path.LineTo(0, float32(boardRenderer.TileSize)*float32(lineClear.Cod))
	default:
		path.MoveTo(float32(boardRenderer.TileSize)*float32(lineClear.Cod), 0)
		path.LineTo(float32(boardRenderer.TileSize)*float32(lineClear.Cod), float32(boardRenderer.TileSize*board.Board.Size[0]))
		path.LineTo(float32(boardRenderer.TileSize)*float32(lineClear.Cod+1), float32(boardRenderer.TileSize*board.Board.Size[0]))
		path.LineTo(float32(boardRenderer.TileSize)*float32(lineClear.Cod+1), 0)
		path.LineTo(float32(boardRenderer.TileSize)*float32(lineClear.Cod), 0)
	}
	opVector := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.FillAll,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].ColorR = float32(0x00) / float32(0xff)
		vs[i].ColorG = float32(0xcd) / float32(0xff)
		vs[i].ColorB = float32(0x7e) / float32(0xff)
		vs[i].ColorA = float32(0xff) * float32(lineClear.Opacity) / float32(0xff)
	}
	screen.DrawTriangles(vs, is, vectorBuffer.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), opVector)
}
