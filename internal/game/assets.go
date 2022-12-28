/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"bytes"
	"embed"
	"image"
	"path"
	"strings"

	// Required
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

//go:embed assets
var assets embed.FS

var assetsGFX = func() map[string]image.Image {
	gfx := map[string]image.Image{}
	gfxPath := path.Join("assets", "gfx")
	gfxDir, _ := assets.ReadDir(gfxPath)
	for _, gfxFile := range gfxDir {
		gfxBytes, _ := assets.ReadFile(path.Join(gfxPath, gfxFile.Name()))
		gfxName := strings.TrimSuffix(gfxFile.Name(), path.Ext(gfxFile.Name()))
		gfx[gfxName], _, _ = image.Decode(bytes.NewReader(gfxBytes))
	}
	return gfx
}()

var assetsSND = func() map[string]*vorbis.Stream {
	snd := map[string]*vorbis.Stream{}
	sndPath := path.Join("assets", "snd")
	sndDir, _ := assets.ReadDir(sndPath)
	for _, sndFile := range sndDir {
		sndBytes, _ := assets.ReadFile(path.Join(sndPath, sndFile.Name()))
		sndName := strings.TrimSuffix(sndFile.Name(), path.Ext(sndFile.Name()))
		snd[sndName], _ = vorbis.DecodeWithoutResampling(bytes.NewReader(sndBytes))
	}
	return snd
}()

var assetsFNT = func() map[string]*sfnt.Font {
	fnt := map[string]*sfnt.Font{}
	fntPath := path.Join("assets", "fnt")
	fntDir, _ := assets.ReadDir(fntPath)
	for _, fntFile := range fntDir {
		fntName := strings.TrimSuffix(fntFile.Name(), path.Ext(fntFile.Name()))
		fntBytes, _ := assets.ReadFile(path.Join(fntPath, fntFile.Name()))
		fnt[fntName], _ = opentype.Parse(fntBytes)
	}
	return fnt
}()

var assetsTXT = func() map[string][]byte {
	txt := map[string][]byte{}
	txtPath := path.Join("assets", "txt")
	txtDir, _ := assets.ReadDir(txtPath)
	for _, txtFile := range txtDir {
		txtName := strings.TrimSuffix(txtFile.Name(), path.Ext(txtFile.Name()))
		txtBytes, _ := assets.ReadFile(path.Join(txtPath, txtFile.Name()))
		txt[txtName] = txtBytes
	}
	return txt
}()
