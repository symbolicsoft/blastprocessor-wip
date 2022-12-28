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
	"math/rand"
	"strconv"

	"blastprocessor.app/internal/board"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type SideBar struct {
	Fonts          map[string]font.Face
	Images         map[string]*ebiten.Image
	FaceFrames     map[string][]*ebiten.Image
	CurrentFrame   map[string]int
	VectorBuffer   *ebiten.Image
	TargetDialog   string
	RenderedDialog string
	DialogBracket  string
}

var sideBar = func() SideBar {
	var sideBar = SideBar{
		Fonts: map[string]font.Face{},
	}
	sideBar.Fonts["text"], _ = opentype.NewFace(assetsFNT["chevyRayOeuf"], &opentype.FaceOptions{
		Size:    13 * 1,
		DPI:     72 * GameScale,
		Hinting: font.HintingFull,
	})
	sideBar.Fonts["digits"], _ = opentype.NewFace(assetsFNT["bittyPix"], &opentype.FaceOptions{
		Size:    9,
		DPI:     72 * GameScale,
		Hinting: font.HintingNone,
	})
	sideBar.Fonts["dialog"], _ = opentype.NewFace(assetsFNT["futilePro"], &opentype.FaceOptions{
		Size:    16,
		DPI:     72 * GameScale,
		Hinting: font.HintingFull,
	})
	sideBar.Images = map[string]*ebiten.Image{}
	sideBar.CurrentFrame = map[string]int{}
	imageNames := []string{
		"sideBar", "hardDisk", "hardDiskFrantic",
		"psu", "psuFrantic", "boards",
		"boardsFrantic", "pipe", "pipeFrantic",
	}
	for _, imageName := range imageNames {
		sideBar.Images[imageName] = ebiten.NewImageFromImage(assetsGFX[imageName])
		sideBar.CurrentFrame[imageName] = 0
	}
	sideBar.FaceFrames = map[string][]*ebiten.Image{}
	sideBar.FaceFrames["drDokkan"] = make([]*ebiten.Image, 33)
	for frame := 0; frame < 33; frame++ {
		sideBar.FaceFrames["drDokkan"][frame] = ebiten.NewImageFromImage(assetsGFX["drDokkan"].(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(71*frame, 0, 71*(frame+1), 71)))
	}
	sideBar.CurrentFrame["drDokkan"] = 0
	sideBar.FaceFrames["kubrik"] = make([]*ebiten.Image, 28)
	for frame := 0; frame < 28; frame++ {
		sideBar.FaceFrames["kubrik"][frame] = ebiten.NewImageFromImage(assetsGFX["kubrik"].(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(71*frame, 0, 71*(frame+1), 71)))
	}
	sideBar.CurrentFrame["kubrik"] = 0
	sideBar.VectorBuffer = ebiten.NewImage(3, 3)
	sideBar.VectorBuffer.Fill(color.White)
	return sideBar
}()

func (sideBar *SideBar) Update(tic int) {
	if tic%60 == 0 {
		session.Time++
	}
	if tic%6 == 0 {
		sideBar.UpdateHoldAndNextChips(tic)
		sideBar.UpdateFrames(tic)
	}
	if tic%4 == 0 {
		sideBar.UpdateDialog(tic)
	}
}

func (sideBar *SideBar) UpdateHoldAndNextChips(tic int) {
	if board.ChipMaxFrame[session.HoldChip.Kind] > 1 {
		session.HoldChip.Frame = (session.HoldChip.Frame + 1) %
			board.ChipMaxFrame[session.HoldChip.Kind]
	}
	if board.ChipMaxFrame[session.NextChip.Kind] > 1 {
		session.NextChip.Frame = (session.NextChip.Frame + 1) %
			board.ChipMaxFrame[session.NextChip.Kind]
	}
}

func (sideBar *SideBar) UpdateFrames(tic int) {
	sideBar.CurrentFrame["hardDisk"] = (sideBar.CurrentFrame["hardDisk"] + 1) % 2
	sideBar.CurrentFrame["hardDiskFrantic"] = (sideBar.CurrentFrame["hardDiskFrantic"] + 1) % 6
	sideBar.CurrentFrame["psu"] = (sideBar.CurrentFrame["psu"] + 1) % 6
	sideBar.CurrentFrame["psuFrantic"] = (sideBar.CurrentFrame["psuFrantic"] + 1) % 12
	sideBar.CurrentFrame["boards"] = (sideBar.CurrentFrame["boards"] + 1) % 8
	sideBar.CurrentFrame["boardsFrantic"] = (sideBar.CurrentFrame["boardsFrantic"] + 1) % 12
	sideBar.CurrentFrame["pipe"] = (sideBar.CurrentFrame["pipe"] + 1) % 4
	sideBar.CurrentFrame["pipeFrantic"] = (sideBar.CurrentFrame["pipeFrantic"] + 1) % 12
	sideBar.CurrentFrame["drDokkan"] = (sideBar.CurrentFrame["drDokkan"] + 1) % 33
	sideBar.CurrentFrame["kubrik"] = (sideBar.CurrentFrame["kubrik"] + 1) % 28
}

func (sideBar *SideBar) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(GameScale, GameScale)
	screen.DrawImage(sideBar.Images["sideBar"], op)
	sideBar.DrawHardDisk(screen)
	sideBar.DrawPSU(screen)
	sideBar.DrawBoards(screen)
	sideBar.DrawPipe(screen)
	sideBar.DrawTimeAndScore(screen)
	sideBar.DrawHoldAndNextChips(screen)
	sideBar.DrawDialog(screen)
	op.GeoM.Translate(3*GameScale, 86*GameScale)
	screen.DrawImage(sideBar.FaceFrames["kubrik"][sideBar.CurrentFrame["kubrik"]], op)
	op.GeoM.Translate(0, 166*GameScale)
	screen.DrawImage(sideBar.FaceFrames["drDokkan"][sideBar.CurrentFrame["drDokkan"]], op)
	sideBar.DrawFloppyBlink(screen)
}

func (sideBar *SideBar) DrawTimeAndScore(screen *ebiten.Image) {
	timeMin := strconv.Itoa(session.Time / 60)
	timeSec := strconv.Itoa(session.Time % 60)
	if len(timeMin) > 2 {
		timeMin = "99"
	}
	for len(timeMin) < 2 {
		timeMin = "0" + timeMin
	}
	for len(timeSec) < 2 {
		timeSec = "0" + timeSec
	}
	timeString := timeMin + ":" + timeSec
	scoreString := strconv.Itoa(session.Score)
	for len(scoreString) < 6 {
		scoreString = "0" + scoreString
	}
	textX := int(92 * GameScale)
	textY := int(150 * GameScale)
	text.Draw(screen, "88888", sideBar.Fonts["digits"],
		textX+int(4*GameScale), textY,
		color.RGBA{0x30, 0x30, 0x30, 0xff})
	text.Draw(screen, timeString, sideBar.Fonts["digits"],
		textX+int(4*GameScale), textY,
		color.RGBA{0xe6, 0x00, 0x51, 0xff})
	text.Draw(screen, "888888", sideBar.Fonts["digits"],
		textX+int(65*GameScale), textY,
		color.RGBA{0x30, 0x30, 0x30, 0xff})
	text.Draw(screen, scoreString, sideBar.Fonts["digits"],
		textX+int(65*GameScale), textY,
		color.RGBA{0xe6, 0x00, 0x51, 0xff})
}

func (sideBar *SideBar) DrawHoldAndNextChips(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(84, 190)
	op.GeoM.Scale(GameScale, GameScale)
	switch session.HoldChip.Kind {
	case board.ChipKindEmpty:
	default:
		screen.DrawImage(boardRenderer.ChipImages[session.HoldChip], op)
	}
	op.GeoM.Translate(66*GameScale, 0)
	screen.DrawImage(boardRenderer.ChipImages[session.NextChip], op)
}

func (sideBar *SideBar) DrawHardDisk(screen *ebiten.Image) {
	hardDiskSpriteType := "hardDisk"
	if pressureGauge.Level >= 60 {
		hardDiskSpriteType = "hardDiskFrantic"
	}
	op := &ebiten.DrawImageOptions{}
	hardDiskSprite := sideBar.Images[hardDiskSpriteType].SubImage(image.Rectangle{
		Min: image.Point{X: 51 * sideBar.CurrentFrame[hardDiskSpriteType], Y: 0},
		Max: image.Point{X: 51 * (sideBar.CurrentFrame[hardDiskSpriteType] + 1), Y: 78},
	}).(*ebiten.Image)
	op.GeoM.Scale(GameScale, GameScale)
	screen.DrawImage(hardDiskSprite, op)
}

func (sideBar *SideBar) DrawPSU(screen *ebiten.Image) {
	psuSpriteType := "psu"
	if pressureGauge.Level >= 60 {
		psuSpriteType = "psuFrantic"
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(154, 0)
	psuSprite := sideBar.Images[psuSpriteType].SubImage(image.Rectangle{
		Min: image.Point{X: 51 * sideBar.CurrentFrame[psuSpriteType], Y: 0},
		Max: image.Point{X: 51 * (sideBar.CurrentFrame[psuSpriteType] + 1), Y: 78},
	}).(*ebiten.Image)
	op.GeoM.Scale(GameScale, GameScale)
	screen.DrawImage(psuSprite, op)
}

func (sideBar *SideBar) DrawFloppyBlink(screen *ebiten.Image) {
	if (5 + rand.Intn(5)) > 8 {
		return
	}
	var path vector.Path
	blinkX := float32(128)
	blinkY := float32(148)
	blinkW := float32(6)
	blinkH := float32(4)
	path.MoveTo(blinkX, blinkY)
	path.LineTo(blinkX+blinkW, blinkY)
	path.LineTo(blinkX+blinkW, blinkY+blinkH)
	path.LineTo(blinkX, blinkY+blinkH)
	path.LineTo(blinkX, blinkY)
	opVector := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.FillAll,
	}
	colorR := 0x00
	colorG := 0xff
	if pressureGauge.Level > 45 {
		colorR = 0xed
		colorG = 0x9d
	}
	if pressureGauge.Level > 60 {
		colorR = 0xff
		colorG = 0x00
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].ColorR = float32(colorR) / float32(0xff)
		vs[i].ColorG = float32(colorG) / float32(0xff)
		vs[i].ColorB = float32(0x00) / float32(0xff)
		vs[i].ColorA = float32(0xff) / float32(0xff)
	}
	screen.DrawTriangles(vs, is, sideBar.VectorBuffer.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), opVector)
}

func (sideBar *SideBar) DrawBoards(screen *ebiten.Image) {
	boardsSpriteType := "boards"
	if pressureGauge.Level >= 60 {
		boardsSpriteType = "boardsFrantic"
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 169)
	boardsSprite := sideBar.Images[boardsSpriteType].SubImage(image.Rectangle{
		Min: image.Point{X: 37 * sideBar.CurrentFrame[boardsSpriteType], Y: 0},
		Max: image.Point{X: 37 * (sideBar.CurrentFrame[boardsSpriteType] + 1), Y: 76},
	}).(*ebiten.Image)
	op.GeoM.Scale(GameScale, GameScale)
	screen.DrawImage(boardsSprite, op)
}

func (sideBar *SideBar) DrawPipe(screen *ebiten.Image) {
	pipeSpriteType := "pipe"
	if pressureGauge.Level >= 60 {
		pipeSpriteType = "pipeFrantic"
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(40, 169)
	pipeSprite := sideBar.Images[pipeSpriteType].SubImage(image.Rectangle{
		Min: image.Point{X: 37 * sideBar.CurrentFrame[pipeSpriteType], Y: 0},
		Max: image.Point{X: 37 * (sideBar.CurrentFrame[pipeSpriteType] + 1), Y: 76},
	}).(*ebiten.Image)
	op.GeoM.Scale(GameScale, GameScale)
	screen.DrawImage(pipeSprite, op)
}

func (sideBar *SideBar) SetDialog(dialogItem DialogItem) {
	sideBar.TargetDialog = dialogItem.Text
	sideBar.RenderedDialog = ""
}

func (sideBar *SideBar) UpdateDialog(tic int) {
	if len(sideBar.TargetDialog) > len(sideBar.RenderedDialog) {
		sideBar.RenderedDialog = sideBar.TargetDialog[:len(sideBar.RenderedDialog)+1]
		sideBar.DialogBracket = "□"
	} else if tic%12 == 0 {
		switch len(sideBar.DialogBracket) {
		case 0:
			sideBar.DialogBracket = "□"
		default:
			sideBar.DialogBracket = ""
		}
	}
}

func (sideBar *SideBar) DrawDialog(screen *ebiten.Image) {
	text.Draw(screen, sideBar.RenderedDialog+sideBar.DialogBracket, sideBar.Fonts["dialog"],
		int(86*GameScale), int(270*GameScale),
		color.RGBA{0x30, 0x30, 0x30, 0xff})
	text.Draw(screen, sideBar.RenderedDialog+sideBar.DialogBracket, sideBar.Fonts["dialog"],
		int(84*GameScale), int(268*GameScale),
		color.RGBA{0x30, 0x30, 0x30, 0xff})
	text.Draw(screen, sideBar.RenderedDialog+sideBar.DialogBracket, sideBar.Fonts["dialog"],
		int(85*GameScale), int(269*GameScale),
		color.RGBA{0xdd, 0xdd, 0xdd, 0xff})
}
