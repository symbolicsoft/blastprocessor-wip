/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Keyboard struct {
	Speed int
}

var keyboard = Keyboard{
	Speed: 6,
}

func (keyboard *Keyboard) Update(tic int, gameState GameState) {
	switch gameState {
	case GameStateInitializing:
	case GameStateTitleScreen:
		switch {
		case ebiten.IsKeyPressed(ebiten.KeyEnter):
			titleScreen.StartGame()
		}
	case GameStateInGame:
		switch {
		case ebiten.IsKeyPressed(ebiten.KeyW),
			ebiten.IsKeyPressed(ebiten.KeyUp):
			GameCurrentInputType = GameInputTypeKeyboard
			mouse.Y -= keyboard.Speed
		case ebiten.IsKeyPressed(ebiten.KeyD),
			ebiten.IsKeyPressed(ebiten.KeyRight):
			GameCurrentInputType = GameInputTypeKeyboard
			mouse.X += keyboard.Speed
		case ebiten.IsKeyPressed(ebiten.KeyS),
			ebiten.IsKeyPressed(ebiten.KeyDown):
			GameCurrentInputType = GameInputTypeKeyboard
			mouse.Y += keyboard.Speed
		case ebiten.IsKeyPressed(ebiten.KeyA),
			ebiten.IsKeyPressed(ebiten.KeyLeft):
			GameCurrentInputType = GameInputTypeKeyboard
			mouse.X -= keyboard.Speed
		case inpututil.IsKeyJustPressed(ebiten.KeySpace),
			inpututil.IsKeyJustPressed(ebiten.KeyEnter):
			GameCurrentInputType = GameInputTypeKeyboard
			mouse.LeftClick()
		case inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft),
			inpututil.IsKeyJustPressed(ebiten.KeyShiftRight):
			GameCurrentInputType = GameInputTypeKeyboard
			mouse.RightClick()
		}
	}
}
