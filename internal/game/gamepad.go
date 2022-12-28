/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Gamepad struct {
	IDsBuf         []ebiten.GamepadID
	IDs            map[ebiten.GamepadID]struct{}
	PressedButtons map[ebiten.GamepadID][]string
	Speed          int
}

var gamepad = func() Gamepad {
	gameControllerDB := assetsTXT["gamecontrollerdb"]
	ebiten.UpdateStandardGamepadLayoutMappings(string(gameControllerDB))
	return Gamepad{
		IDsBuf:         []ebiten.GamepadID{},
		IDs:            map[ebiten.GamepadID]struct{}{},
		PressedButtons: map[ebiten.GamepadID][]string{},
		Speed:          6,
	}
}()

func (gamepad *Gamepad) Scan() {
	gamepad.IDsBuf = inpututil.AppendJustConnectedGamepadIDs(gamepad.IDsBuf[:0])
	for _, id := range gamepad.IDsBuf {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			gamepad.IDs[id] = struct{}{}
		}
	}
	for id := range gamepad.IDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			delete(gamepad.IDs, id)
		}
	}
}

func (gamepad *Gamepad) Update(tic int) {
	gamepad.Scan()
	gamepad.PressedButtons = map[ebiten.GamepadID][]string{}
	for id := range gamepad.IDs {
		switch {
		case ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) <= -0.9:
			GameCurrentInputType = GameInputTypeGamepad
			mouse.Y -= gamepad.Speed
		case ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) >= 0.9:
			GameCurrentInputType = GameInputTypeGamepad
			mouse.X += gamepad.Speed
		case ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical) >= 0.9:
			GameCurrentInputType = GameInputTypeGamepad
			mouse.Y += gamepad.Speed
		case ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal) <= -0.9:
			GameCurrentInputType = GameInputTypeGamepad
			mouse.X -= gamepad.Speed
		}
		switch {
		case inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightRight):
			GameCurrentInputType = GameInputTypeGamepad
			mouse.LeftClick()
		case inpututil.IsStandardGamepadButtonJustPressed(id, ebiten.StandardGamepadButtonRightTop):
			GameCurrentInputType = GameInputTypeGamepad
			mouse.RightClick()
		}
	}
}
