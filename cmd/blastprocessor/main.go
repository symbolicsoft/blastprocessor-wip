//go:build !js

//go:generate goversioninfo -64=true ../../build/windows/versioninfo.json

/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package main

import (
	"os"
	"runtime/pprof"

	"blastprocessor.app/internal/game"
)

func main() {
	if DEBUG {
		pprofFile, _ := os.Create("cpu.pprof")
		pprof.StartCPUProfile(pprofFile)
		defer pprof.StopCPUProfile()
	}
	var g game.Game
	g.Construct(BUILDNUM, DEBUG, true)
}
