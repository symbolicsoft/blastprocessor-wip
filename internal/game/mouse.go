/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"image"
	"math"
	"math/rand"

	"blastprocessor.app/internal/board"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Mouse struct {
	X           int
	Y           int
	LeftClicked bool
	Images      map[string]*ebiten.Image
	Sparkles    []MouseSparkle
}

type MouseSparkle struct {
	X       float64
	Y       float64
	OffsetX float64
	OffsetY float64
	Scale   float64
	Hue     int
	Opacity float64
}

var mouse = func() Mouse {
	mouse := Mouse{}
	mouse.LeftClicked = false
	mouse.Images = map[string]*ebiten.Image{}
	mouse.Images["point"] = ebiten.NewImageFromImage(
		assetsGFX["mouse"].(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(0, 0, 45, 45)))
	mouse.Images["click"] = ebiten.NewImageFromImage(
		assetsGFX["mouse"].(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(45, 0, 90, 45)))
	mouse.Images["open"] = ebiten.NewImageFromImage(
		assetsGFX["mouse"].(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(90, 0, 135, 45)))
	mouse.Images["grab"] = ebiten.NewImageFromImage(
		assetsGFX["mouse"].(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(135, 0, 180, 45)))
	mouse.Images["sparkle"] = ebiten.NewImageFromImage(assetsGFX["sparkle"])
	return mouse
}()

func (mouse *Mouse) Update(tic int) {
	switch GameCurrentInputType {
	case GameInputTypeMouse:
		mouse.X, mouse.Y = ebiten.CursorPosition()
	}
	mouse.UpdateSparkles(tic)
	if tic%16 == 0 {
		switch session.CurrentChip.Kind {
		case board.ChipKindBomb:
		default:
			// mouse.AddSparkle()
		}
	}
	if tic%6 == 0 {
		mouse.UpdateChipFrames(tic)
	}
	mouse.UpdateSquareHighlight()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		GameCurrentInputType = GameInputTypeMouse
		mouse.LeftClick()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mouse.LeftClicked = false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		GameCurrentInputType = GameInputTypeMouse
		mouse.RightClick()
	}
}

func (mouse *Mouse) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(mouse.X), float64(mouse.Y))
	if session.Active &&
		(session.CurrentChip.Kind != board.ChipKindEmpty) &&
		(mouse.X > int(boardRenderer.BoardOffset[0]-float64(boardRenderer.TileSize)/4)) {
		op.GeoM.Translate(-15-3, -15-3)
		op.ColorScale.Scale(0, 0, 0, 0.4)
		op.GeoM.Scale(GameScale, GameScale)
		screen.DrawImage(boardRenderer.ChipImages[session.CurrentChip], op)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(mouse.X), float64(mouse.Y))
		op.GeoM.Translate(-15, -15)
		op.GeoM.Scale(GameScale, GameScale)
		screen.DrawImage(boardRenderer.ChipImages[session.CurrentChip], op)
		mouse.DrawSparkles(screen)
	}
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(mouse.X)-3, float64(mouse.Y)-3)
	op.ColorScale.Scale(0, 0, 0, 0.4)
	op.GeoM.Scale(GameScale, GameScale)
	if mouse.X <= int(boardRenderer.BoardOffset[0]-float64(boardRenderer.TileSize)/4) {
		switch mouse.LeftClicked {
		case true:
			screen.DrawImage(mouse.Images["click"], op)
		default:
			screen.DrawImage(mouse.Images["point"], op)
		}
	} else {
		switch mouse.LeftClicked {
		case true:
			screen.DrawImage(mouse.Images["open"], op)
		default:
			screen.DrawImage(mouse.Images["grab"], op)
		}
	}
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(mouse.X), float64(mouse.Y))
	op.GeoM.Scale(GameScale, GameScale)
	if mouse.X <= int(boardRenderer.BoardOffset[0]-float64(boardRenderer.TileSize)/4) {
		switch mouse.LeftClicked {
		case true:
			screen.DrawImage(mouse.Images["click"], op)
		default:
			screen.DrawImage(mouse.Images["point"], op)
		}
	} else {
		if mouse.LeftClicked {
			screen.DrawImage(mouse.Images["open"], op)
		} else {
			screen.DrawImage(mouse.Images["grab"], op)
		}
	}
}

func (mouse *Mouse) BoardCoordinates() [2]int {
	row := (mouse.Y - int(boardRenderer.BoardOffset[1])) / boardRenderer.TileSize
	col := (mouse.X - int(boardRenderer.BoardOffset[0])) / boardRenderer.TileSize
	if (row > board.Board.Size[0]-1) || (row < 0) {
		row = -1
	}
	if (col > board.Board.Size[1]-1) || (col < 0) {
		col = -1
	}
	return [2]int{row, col}
}

func (mouse *Mouse) LeftClick() {
	mouse.LeftClicked = true
	mouseCod := mouse.BoardCoordinates()
	if (mouseCod[0] >= 0) && (mouseCod[1] >= 0) {
		session.PlaceChip(mouseCod)
	}
}

func (mouse *Mouse) RightClick() {
	mouseCod := mouse.BoardCoordinates()
	if (mouseCod[0] >= 0) && (mouseCod[1] >= 0) {
		session.ShuffleHold()
	}
}

func (mouse *Mouse) AddSparkle() {
	sparkle := MouseSparkle{
		X: float64(mouse.X) + 3,
		Y: float64(mouse.Y) + 10,
		OffsetX: func() float64 {
			spread := rand.Intn(12)
			if spread == 0 {
				return -15 + (rand.Float64() * 30)
			} else if spread >= 1 && spread <= 3 {
				return -10 + (rand.Float64() * 20)
			} else {
				return -8 + (rand.Float64() * 16)
			}
		}(),
		OffsetY: rand.Float64() * 5,
		Scale:   ((rand.Float64() * 4) + 4) / 10 / GameScale,
		Hue:     rand.Intn(360),
		Opacity: 1,
	}
	mouse.Sparkles = append(mouse.Sparkles, sparkle)
}

func (mouse *Mouse) UpdateSparkles(tic int) {
	for i := 0; i < len(mouse.Sparkles); i++ {
		dx := mouse.Sparkles[i].X - float64(mouse.X)
		dy := mouse.Sparkles[i].Y - float64(mouse.Y)
		mouse.Sparkles[i].Opacity -= ((dx * dx) + (dy * dy)) / 200000
		mouse.Sparkles[i].Opacity -= 0.0025
		mouse.Sparkles[i].OffsetY += 0.1
		mouse.Sparkles[i].Hue = (mouse.Sparkles[i].Hue + 2) % 360
		if mouse.Sparkles[i].Opacity <= 0 {
			mouse.Sparkles = append(mouse.Sparkles[:i], mouse.Sparkles[i+1:]...)
		}
	}
}

func (mouse *Mouse) DrawSparkles(screen *ebiten.Image) {
	for _, sparkle := range mouse.Sparkles {
		op := &colorm.DrawImageOptions{}
		co := colorm.ColorM{}
		op.GeoM.Scale(sparkle.Scale, sparkle.Scale)
		co.ChangeHSV(float64(sparkle.Hue%360)*2*math.Pi/360, math.Pi*2, math.Pi*2)
		co.Scale(1, 1, 1, sparkle.Opacity)
		op.GeoM.Translate(float64(sparkle.X)+sparkle.OffsetX, float64(sparkle.Y)+sparkle.OffsetY)
		op.GeoM.Scale(GameScale, GameScale)
		colorm.DrawImage(screen, mouse.Images["sparkle"], co, op)
	}
}

func (mouse *Mouse) UpdateChipFrames(tic int) {
	if board.ChipMaxFrame[session.CurrentChip.Kind] > 1 {
		session.CurrentChip.Frame = (session.CurrentChip.Frame + 1) %
			board.ChipMaxFrame[session.CurrentChip.Kind]
	}
}

func (mouse *Mouse) UpdateSquareHighlight() {
	mouseCod := mouse.BoardCoordinates()
	if (mouseCod[0] < 0) || (mouseCod[1] < 0) {
		return
	}
	if boardRenderer.SquareHighlight.Cod[0] != mouseCod[0] ||
		boardRenderer.SquareHighlight.Cod[1] != mouseCod[1] {
		boardRenderer.SquareHighlight.Hide()
		if session.Complete {
			return
		}
		if board.Board.AdjacentToChip(mouseCod) {
			switch session.CurrentChip.Kind {
			case board.ChipKindBomb:
				if board.Board.CodHasChip(mouseCod) {
					boardRenderer.SquareHighlight.Show(mouseCod, 50)
				}
			default:
				if !board.Board.CodHasChip(mouseCod) {
					boardRenderer.SquareHighlight.Show(mouseCod, 240)
				}
			}
		}
	}
}
