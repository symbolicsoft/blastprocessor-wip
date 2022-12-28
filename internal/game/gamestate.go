/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

type GameState int

const (
	GameStateInitializing GameState = iota
	GameStateTitleScreen  GameState = iota
	GameStateInGame       GameState = iota
)

func (g *Game) GameStateNext() {
	switch g.GameState {
	case GameStateInitializing:
		g.GameState = GameStateInitializingFunc()
	case GameStateTitleScreen:
		g.GameState = GameStateTitleScreenFunc()
	case GameStateInGame:
		g.GameState = GameStateInGameFunc()
	}
}

func GameStateInitializingFunc() GameState {
	return GameStateTitleScreen
}

func GameStateTitleScreenFunc() GameState {
	switch titleScreen.GameStart {
	case true:
		soundPlayer.Pause("title")
		soundPlayer.Play("music")
		return GameStateInGame
	default:
		return GameStateTitleScreen
	}
}

func GameStateInGameFunc() GameState {
	return GameStateInGame
}
