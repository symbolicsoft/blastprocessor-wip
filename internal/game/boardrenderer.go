/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"strconv"
	"time"

	"blastprocessor.app/internal/board"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BoardRenderer struct {
	TileSize         int
	RenderedImage    *ebiten.Image
	VectorBuffer     *ebiten.Image
	CircuitImages    map[board.CircuitKind]*ebiten.Image
	CorruptionImages map[board.CircuitCorruption]*ebiten.Image
	ChipImages       map[board.Chip]*ebiten.Image
	BracketImages    []*ebiten.Image
	SolderSparks     []SolderSpark
	BoardOffset      [2]float64
	ShakeOffsets     [][2]float64
	ChipOffsets      map[[2]int]BoardRendererChipOffset
	LineClears       []LineClear
	SquareClears     []SquareClear
	SquareHighlight  SquareHighlight
}

type BoardRendererChipOffset struct {
	OffsetX   float64
	OffsetY   float64
	ModifierY float64
	Speed     float64
	Rotation  float64
}

var boardRenderer = func(tileSize int) BoardRenderer {
	boardRenderer := BoardRenderer{}
	boardRenderer.BoardOffset = [2]float64{GameWidth - float64(tileSize*board.Board.Size[1]), 0}
	boardRenderer.TileSize = tileSize
	boardRenderer.RenderedImage = ebiten.NewImage(
		boardRenderer.TileSize*board.Board.Size[1],
		boardRenderer.TileSize*board.Board.Size[0])
	boardRenderer.VectorBuffer = ebiten.NewImage(
		boardRenderer.TileSize*board.Board.Size[1],
		boardRenderer.TileSize*board.Board.Size[0])
	boardRenderer.BracketImages = make([]*ebiten.Image, 8)
	boardRenderer.CircuitImages = map[board.CircuitKind]*ebiten.Image{}
	boardRenderer.CorruptionImages = map[board.CircuitCorruption]*ebiten.Image{}
	boardRenderer.ChipImages = map[board.Chip]*ebiten.Image{}
	boardRenderer.ChipOffsets = map[[2]int]BoardRendererChipOffset{}
	for bracketImage := 0; bracketImage < len(boardRenderer.BracketImages); bracketImage++ {
		boardRenderer.BracketImages[bracketImage] = ebiten.NewImageFromImage(
			assetsGFX["bracket"].(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(40*bracketImage, 0, 40*(bracketImage+1), 40)))
	}
	for kind := 0; kind < board.CircuitKindCount; kind++ {
		boardRenderer.CircuitImages[board.CircuitKind(kind)] = ebiten.NewImageFromImage(
			assetsGFX["circuits"].(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(boardRenderer.TileSize*kind, 0, boardRenderer.TileSize*(kind+1), boardRenderer.TileSize)))
	}
	for corruption := 0; corruption < board.CircuitCorruptionCount; corruption++ {
		switch board.CircuitCorruption(corruption) {
		case board.CircuitCorruptionNone:
		default:
			boardRenderer.CorruptionImages[board.CircuitCorruption(corruption)] = ebiten.NewImageFromImage(
				assetsGFX["corruption"+strconv.Itoa(corruption-1)])
		}
	}
	for kind := 0; kind <= board.ChipKindCount; kind++ {
		switch board.ChipKind(kind) {
		case board.ChipKindEmpty:
			continue
		case board.ChipKindWild:
			for frame := 0; frame < board.ChipMaxFrame[board.ChipKind(kind)]; frame++ {
				boardRenderer.ChipImages[board.Chip{Kind: board.ChipKind(kind), Color: board.ChipColor(board.ChipColorA), Frame: frame}] = ebiten.NewImageFromImage(
					assetsGFX["wild"].(interface {
						SubImage(r image.Rectangle) image.Image
					}).SubImage(image.Rect(boardRenderer.TileSize*frame, 0, boardRenderer.TileSize*(frame+1), boardRenderer.TileSize)))
			}
		case board.ChipKindBomb:
			for frame := 0; frame < board.ChipMaxFrame[board.ChipKind(kind)]; frame++ {
				boardRenderer.ChipImages[board.Chip{Kind: board.ChipKind(kind), Color: board.ChipColor(board.ChipColorA), Frame: frame}] = ebiten.NewImageFromImage(
					assetsGFX["bomb"].(interface {
						SubImage(r image.Rectangle) image.Image
					}).SubImage(image.Rect(boardRenderer.TileSize*frame, 0, boardRenderer.TileSize*(frame+1), boardRenderer.TileSize)))
			}
		default:
			for color := 0; color < board.ChipColorCount; color++ {
				boardRenderer.ChipImages[board.Chip{Kind: board.ChipKind(kind), Color: board.ChipColor(color), Frame: 0}] = ebiten.NewImageFromImage(
					assetsGFX["chips"].(interface {
						SubImage(r image.Rectangle) image.Image
					}).SubImage(image.Rect(boardRenderer.TileSize*(kind-3), boardRenderer.TileSize*color, boardRenderer.TileSize*(kind-2), boardRenderer.TileSize*(color+1))))
			}
		}
	}
	boardRenderer.VectorBuffer.Fill(color.White)
	return boardRenderer
}(45)

func (boardRenderer *BoardRenderer) Update(tic int) {
	for i := 0; i < len(boardRenderer.SolderSparks); i++ {
		boardRenderer.SolderSparks[i].Update(tic)
		if boardRenderer.SolderSparks[i].Opacity <= 0 {
			boardRenderer.SolderSparks = append(boardRenderer.SolderSparks[:i], boardRenderer.SolderSparks[i+1:]...)
		}
	}
	if tic%6 == 0 {
		boardRenderer.UpdateChipFrames(tic)
	}
	for i := 0; i < len(boardRenderer.LineClears); i++ {
		boardRenderer.LineClears[i].Update(tic)
	}
	for i := 0; i < len(boardRenderer.SquareClears); i++ {
		boardRenderer.SquareClears[i].Update(tic)
	}
	boardRenderer.SquareHighlight.Update(tic)
	if pressureGauge.Level > 70 && pressureGauge.Level < 100 {
		if tic%4 == 0 {
			boardRenderer.AddShake(float64(pressureGauge.Level/40), int(pressureGauge.Level/20))
		}
	}
	if len(boardRenderer.ChipOffsets) > 0 {
		boardRenderer.UpdateExplosion()
	}
}

func (boardRenderer *BoardRenderer) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if len(boardRenderer.ShakeOffsets) > 0 {
		op.GeoM.Translate(boardRenderer.BoardOffset[0]+boardRenderer.ShakeOffsets[0][0], boardRenderer.BoardOffset[1]+boardRenderer.ShakeOffsets[0][1])
		boardRenderer.ShakeOffsets = boardRenderer.ShakeOffsets[1:]
	} else {
		op.GeoM.Translate(boardRenderer.BoardOffset[0], boardRenderer.BoardOffset[1])
	}
	boardRenderer.Render()
	boardRenderer.SquareHighlight.Draw(boardRenderer.RenderedImage)
	for i := 0; i < len(boardRenderer.LineClears); i++ {
		boardRenderer.LineClears[i].Draw(boardRenderer.RenderedImage, boardRenderer.VectorBuffer)
		switch boardRenderer.LineClears[i].State {
		case 2:
			if boardRenderer.LineClears[i].SubState >= 100 {
				boardRenderer.LineClears = boardRenderer.LineClears[1:]
			}
		}
	}
	for i := 0; i < len(boardRenderer.SquareClears); i++ {
		boardRenderer.SquareClears[i].Draw(boardRenderer.RenderedImage, boardRenderer.VectorBuffer)
		switch boardRenderer.SquareClears[i].State {
		case 2:
			if boardRenderer.SquareClears[i].SubState >= 100 {
				boardRenderer.SquareClears = boardRenderer.SquareClears[1:]
			}
		}
	}
	op.GeoM.Scale(GameScale, GameScale)
	screen.DrawImage(boardRenderer.RenderedImage, op)
	for i := range boardRenderer.SolderSparks {
		boardRenderer.SolderSparks[i].Draw(screen)
	}
}

func (boardRenderer *BoardRenderer) Render() {
	boardRenderer.RenderedImage.Fill(color.White)
	var pathA vector.Path
	var pathB vector.Path
	pathAWidth := 1
	pathBWidth := 8
	for row := 0; row < board.Board.Size[0]; row++ {
		for col := 0; col < board.Board.Size[1]; col++ {
			switch board.Board.Circuits[row][col].Filled {
			case true:
				if row > 0 && !board.Board.Circuits[row-1][col].Filled {
					pathA.MoveTo(float32(boardRenderer.TileSize*col), float32(boardRenderer.TileSize*row))
					pathA.LineTo(float32(boardRenderer.TileSize*(col+1)), float32(boardRenderer.TileSize*row))
					pathB.MoveTo(float32(boardRenderer.TileSize*col), float32(boardRenderer.TileSize*row-pathAWidth))
					pathB.LineTo(float32(boardRenderer.TileSize*(col+1)-(pathBWidth/2)), float32(boardRenderer.TileSize*row-pathAWidth))
				}
				if row < board.Board.Size[0]-1 && !board.Board.Circuits[row+1][col].Filled {
					pathA.MoveTo(float32(boardRenderer.TileSize*col), float32(boardRenderer.TileSize*(row+1))+1)
					pathA.LineTo(float32(boardRenderer.TileSize*(col+1)), float32(boardRenderer.TileSize*(row+1))+1)
				}
				if col > 0 && !board.Board.Circuits[row][col-1].Filled {
					pathA.MoveTo(float32(boardRenderer.TileSize*col), float32(boardRenderer.TileSize*row))
					pathA.LineTo(float32(boardRenderer.TileSize*col), float32(boardRenderer.TileSize*(row+1)))
					pathB.MoveTo(float32(boardRenderer.TileSize*col-pathAWidth), float32(boardRenderer.TileSize*row))
					pathB.LineTo(float32(boardRenderer.TileSize*col-pathAWidth), float32(boardRenderer.TileSize*(row+1)-(pathBWidth/2)))
				}
				if col < board.Board.Size[1]-1 && !board.Board.Circuits[row][col+1].Filled {
					pathA.MoveTo(float32(boardRenderer.TileSize*(col+1))+1, float32(boardRenderer.TileSize*row))
					pathA.LineTo(float32(boardRenderer.TileSize*(col+1))+1, float32(boardRenderer.TileSize*(row+1)))
				}
			}
		}
	}
	opVector := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.FillAll,
	}
	vsA, isA := pathA.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    float32(pathAWidth),
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	})
	vsB, isB := pathB.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    float32(pathBWidth),
		LineCap:  vector.LineCapSquare,
		LineJoin: vector.LineJoinMiter,
	})
	for i := range vsA {
		vsA[i].ColorR = float32(0x00) / float32(0xff)
		vsA[i].ColorG = float32(0x60) / float32(0xff)
		vsA[i].ColorB = float32(0x40) / float32(0xff)
		vsA[i].ColorA = float32(0xff) / float32(0xff)
	}
	for i := range vsB {
		vsB[i].ColorR = float32(0x00) / float32(0xff)
		vsB[i].ColorG = float32(0x00) / float32(0xff)
		vsB[i].ColorB = float32(0x00) / float32(0xff)
		vsB[i].ColorA = float32(0x50) / float32(0xff)
	}
	boardRenderer.RenderedImage.DrawTriangles(vsB, isB, boardRenderer.VectorBuffer.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), opVector)
	for row := 0; row < board.Board.Size[0]; row++ {
		for col := 0; col < board.Board.Size[1]; col++ {
			var op *colorm.DrawImageOptions
			co := colorm.ColorM{}
			switch board.Board.Circuits[row][col].Filled {
			case false:
				op = &colorm.DrawImageOptions{
					Blend: ebiten.Blend{
						BlendFactorSourceRGB:        ebiten.BlendFactorOne,
						BlendFactorSourceAlpha:      ebiten.BlendFactorOne,
						BlendFactorDestinationRGB:   ebiten.BlendFactorOne,
						BlendFactorDestinationAlpha: ebiten.BlendFactorOne,
						BlendOperationRGB:           ebiten.BlendOperationReverseSubtract,
						BlendOperationAlpha:         ebiten.BlendOperationAdd,
					},
				}
				op.GeoM.Translate(float64(boardRenderer.TileSize*col), float64(boardRenderer.TileSize*row))
				co.ChangeHSV(float64(3.8+(pressureGauge.Level/50)), 1, float64(0.9+(pressureGauge.Level/400)))
			default:
				op = &colorm.DrawImageOptions{}
				op.GeoM.Translate(float64(boardRenderer.TileSize*col), float64(boardRenderer.TileSize*row))
			}
			colorm.DrawImage(boardRenderer.RenderedImage, boardRenderer.CircuitImages[board.Board.Circuits[row][col].Kind], co, op)
			switch board.Board.Circuits[row][col].Corruption {
			case board.CircuitCorruptionNone:
			case board.CircuitCorruptionPopup:
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(boardRenderer.TileSize*col), float64(boardRenderer.TileSize*row))
				boardRenderer.RenderedImage.DrawImage(boardRenderer.CorruptionImages[board.Board.Circuits[row][col].Corruption], op)
			}
		}
	}
	boardRenderer.RenderedImage.DrawTriangles(vsA, isA, boardRenderer.VectorBuffer.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), opVector)
	for row := 0; row < board.Board.Size[0]; row++ {
		for col := 0; col < board.Board.Size[1]; col++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(boardRenderer.TileSize*col), float64(boardRenderer.TileSize*row))
			chipCod := [2]int{row, col}
			if board.Board.CodHasChip(chipCod) {
				scale := math.Abs(boardRenderer.ChipOffsets[chipCod].OffsetY) / 5
				op.GeoM.Scale(1+scale, 1+scale)
				op.GeoM.Translate(
					boardRenderer.ChipOffsets[chipCod].OffsetX-(scale*5*float64(boardRenderer.TileSize)),
					boardRenderer.ChipOffsets[chipCod].OffsetY-(scale*5*float64(boardRenderer.TileSize)),
				)
				op.GeoM.Rotate(boardRenderer.ChipOffsets[chipCod].Rotation)
				boardRenderer.RenderedImage.DrawImage(boardRenderer.ChipImages[board.Board.Chips[row][col]], op)
			}
		}
	}
}

func (boardRenderer *BoardRenderer) AddSolderSparks(sparkCount int, sparkColor color.RGBA, colorMutationFactor int, duration int, cod [2]int) {
	for i := 0; i < sparkCount; i++ {
		go func() {
			time.Sleep(time.Duration(int(time.Millisecond) * rand.Intn(duration)))
			sparkA := SolderSpark{}
			sparkB := sparkA.GenSparkPair(sparkColor, colorMutationFactor, cod)
			boardRenderer.SolderSparks = append(boardRenderer.SolderSparks, sparkA)
			boardRenderer.SolderSparks = append(boardRenderer.SolderSparks, sparkB)
		}()
	}
}

func (BoardRenderer *BoardRenderer) UpdateChipFrames(tic int) {
	for row := 0; row < board.Board.Size[0]; row++ {
		for col := 0; col < board.Board.Size[1]; col++ {
			if board.ChipMaxFrame[board.Board.Chips[row][col].Kind] > 1 {
				board.Board.Chips[row][col].Frame = (board.Board.Chips[row][col].Frame + 1) %
					board.ChipMaxFrame[board.Board.Chips[row][col].Kind]
			}
		}
	}
}

func (boardRenderer *BoardRenderer) AddShake(power float64, duration int) {
	boardRenderer.ShakeOffsets = make([][2]float64, duration)
	for d := 0; d < duration; d++ {
		pmd := power - ((power / float64(duration)) * (math.Pow(float64(d), 3) / 5))
		if pmd <= 0 {
			continue
		}
		boardRenderer.ShakeOffsets[d] = [2]float64{
			-(pmd / 2) + (rand.Float64() * pmd),
			-(pmd / 2) + (rand.Float64() * pmd),
		}
	}
}

func (boardRenderer *BoardRenderer) AddLineClear(isRow bool, cod int, onClear func(cod int)) {
	lineClear := LineClear{
		IsRow:    isRow,
		Cod:      cod,
		State:    0,
		SubState: 0,
		OnClear:  onClear,
	}
	boardRenderer.LineClears = append(boardRenderer.LineClears, lineClear)
}

func (boardRenderer *BoardRenderer) AddSquareClear(cod [2]int, circuitColor color.RGBA, onClear func([2]int)) {
	squareClear := SquareClear{
		Cod:      cod,
		State:    0,
		SubState: 0,
		Color:    circuitColor,
		OnClear:  onClear,
	}
	boardRenderer.SquareClears = append(boardRenderer.SquareClears, squareClear)
}

func (boardRenderer *BoardRenderer) CreateExplosion() {
	chipCods := board.Board.GetAllChipCods()
	for _, chipCod := range chipCods {
		if _, ok := boardRenderer.ChipOffsets[chipCod]; !ok {
			boardRenderer.ChipOffsets[chipCod] = BoardRendererChipOffset{
				OffsetX:   0,
				OffsetY:   0,
				ModifierY: math.Pow(25, 2),
				Speed:     1,
				Rotation:  0,
			}
			boardRenderer.AddSolderSparks(50, color.RGBA{0xcd, 0x80, 0x80, 0xff}, 100, 400, chipCod)
		}
	}
	boardRenderer.AddShake(100, 10000)
}

func (boardRenderer *BoardRenderer) UpdateExplosion() {
	for chipCod, chipOffset := range boardRenderer.ChipOffsets {
		if math.Abs(chipOffset.OffsetX) > 500 {
			board.Board.ClearChip(chipCod)
			delete(boardRenderer.ChipOffsets, chipCod)
		}
	}
	for chipCod, chipOffset := range boardRenderer.ChipOffsets {
		if chipOffset.OffsetX == 0 {
			switch rand.Intn(2) {
			case 0:
				chipOffset.OffsetX -= chipOffset.Speed
			case 1:
				chipOffset.OffsetX += chipOffset.Speed
			}
		} else if chipOffset.OffsetX < 0 {
			chipOffset.OffsetX -= chipOffset.Speed
			chipOffset.Rotation -= 0.035 + rand.Float64()*0.02
		} else if chipOffset.OffsetX > 0 {
			chipOffset.OffsetX += chipOffset.Speed
			chipOffset.Rotation += 0.035 + rand.Float64()*0.02
		}
		chipOffset.OffsetY = math.Pow(chipOffset.OffsetX, 2) / chipOffset.ModifierY
		chipOffset.Speed *= 1.15
		boardRenderer.ChipOffsets[chipCod] = chipOffset
	}
}
