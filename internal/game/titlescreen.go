/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import "github.com/hajimehoshi/ebiten/v2"

type TitleScreen struct {
	Background *ebiten.Image
	GameStart  bool
}

var titleScreen TitleScreen = func() TitleScreen {
	return TitleScreen{
		Background: ebiten.NewImageFromImage(assetsGFX["titleScreen"]),
	}
}()

func (titleScreen *TitleScreen) Update(tic int) {

}

func (titleScreen *TitleScreen) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(GameScale, GameScale)
	screen.DrawImage(titleScreen.Background, op)
}

func (titleScreen *TitleScreen) StartGame() {
	titleScreen.GameStart = true
}
