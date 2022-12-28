/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package blastprocessormobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"

	"blastprocessor.app/internal/game"
)

func init() {
	var g game.Game
	g.Construct(BUILDNUM, DEBUG, false)
	mobile.SetGame(&g)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() bool {
	return true
}
