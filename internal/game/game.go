/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"image"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	StaticOverlay      []*ebiten.Image
	StaticOverlayFrame int
	Debug              bool
	Tic                int
	GameState          GameState
}

type GameInputType int

const (
	GameInputTypeMouse    GameInputType = iota
	GameInputTypeKeyboard GameInputType = iota
	GameInputTypeGamepad  GameInputType = iota
)

var (
	GameZoomLevel        float64       = 2
	GameWidth            float64       = 640
	GameHeight           float64       = 360
	GameTPS              int           = 60
	GameScale            float64       = math.Round(ebiten.DeviceScaleFactor())
	GameCurrentInputType GameInputType = GameInputTypeMouse
)

func (g *Game) Construct(buildNum int, debug bool, runGame bool) {
	ebiten.SetWindowTitle("Blast Processor")
	ebiten.SetWindowIcon([]image.Image{assetsGFX["icon"]})
	ebiten.SetWindowSize(
		int(GameWidth*GameZoomLevel),
		int(GameHeight*GameZoomLevel),
	)
	ebiten.SetTPS(GameTPS)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetRunnableOnUnfocused(true)
	g.StaticOverlay = make([]*ebiten.Image, 6)
	g.GameState = GameStateInitializing
	for i := 0; i < len(g.StaticOverlay); i++ {
		g.StaticOverlay[i] = ebiten.NewImageFromImage(assetsGFX["static"+strconv.Itoa(i)])
	}
	if runGame {
		ebiten.RunGame(g)
	}
}

func (g *Game) Update() error {
	switch g.GameState {
	case GameStateInitializing:
		session.NewGame(1, 10, 6)
		soundPlayer.Play("title")
	case GameStateTitleScreen:
		titleScreen.Update(g.Tic)
		keyboard.Update(g.Tic, g.GameState)
	case GameStateInGame:
		if g.Tic%4 == 0 {
			g.StaticOverlayFrame = (g.StaticOverlayFrame + 1) % len(g.StaticOverlay)
		}
		boardRenderer.Update(g.Tic)
		pressureGauge.Update(g.Tic)
		sideBar.Update(g.Tic)
		mouse.Update(g.Tic)
		keyboard.Update(g.Tic, g.GameState)
		gamepad.Update(g.Tic)
	}
	g.GameStateNext()
	g.Tic++
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.GameState {
	case GameStateInitializing:
	case GameStateTitleScreen:
		titleScreen.Draw(screen)
	case GameStateInGame:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(GameScale, GameScale)
		sideBar.Draw(screen)
		boardRenderer.Draw(screen)
		pressureGauge.Draw(screen)
		mouse.Draw(screen)
		op = &ebiten.DrawImageOptions{
			Blend: ebiten.BlendLighter,
		}
		op.ColorScale.ScaleAlpha(0.02)
		op.GeoM.Scale(GameScale, GameScale)
		screen.DrawImage(g.StaticOverlay[g.StaticOverlayFrame], op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(GameWidth * GameScale), int(GameHeight * GameScale)
}
