//go:build js

/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package main

import (
	"regexp"
	"syscall/js"

	"blastprocessor.app/internal/game"
)

func main() {
	if DRM() {
		js.Global().Get("document").Call("getElementById", "loading").Call("remove")
		var g game.Game
		g.Construct(BUILDNUM, DEBUG, true)
	} else {
		js.Global().Get("window").Call("alert", "DRM check failed")
	}
}

func DRM() bool {
	validHref := `^https:\/\/(www\.)?blastprocessor\.app\/wasm\/?$`
	locationHref := js.Global().Get("location").Get("href").String()
	matchHref, err := regexp.MatchString(validHref, locationHref)
	return matchHref && (err == nil)
}
