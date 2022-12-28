/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type SoundPlayer struct {
	Context *audio.Context
	Sounds  map[string]*audio.Player
}

var soundPlayer = func() SoundPlayer {
	soundPlayer := SoundPlayer{
		Context: audio.NewContext(48000),
		Sounds:  map[string]*audio.Player{},
	}
	for sndName, sndDecoded := range assetsSND {
		sndPlayer, _ := soundPlayer.Context.NewPlayer(sndDecoded)
		if sndName == "music" || sndName == "title" {
			sndPlayer, _ = soundPlayer.Context.NewPlayer(audio.NewInfiniteLoopWithIntro(sndDecoded, 0, sndDecoded.Length()))
		}
		soundPlayer.Sounds[sndName] = sndPlayer
		soundPlayer.Sounds[sndName].Rewind()
	}
	return soundPlayer
}()

func (soundplayer *SoundPlayer) AdjustVolume(sndName string) {
	soundPlayer.Sounds[sndName].SetVolume(map[string]float64{
		"bomb0":     0.75,
		"bomb1":     0.75,
		"discard0":  1,
		"discard1":  1,
		"lineClear": 0.5,
		"music":     1,
		"title":     1,
		"popup":     0.5,
		"thunk":     1,
		"wild":      0.25,
	}[sndName])
}

func (soundPlayer *SoundPlayer) Play(sndName string) {
	soundPlayer.AdjustVolume(sndName)
	soundPlayer.Sounds[sndName].Rewind()
	soundPlayer.Sounds[sndName].Play()
}

func (soundPlayer *SoundPlayer) Pause(sndName string) {
	soundPlayer.Sounds[sndName].Pause()
}
