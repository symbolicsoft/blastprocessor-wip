/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

import (
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"blastprocessor.app/internal/board"
)

type Session struct {
	Active       bool
	Complete     bool
	CurrentChip  board.Chip
	NextChip     board.Chip
	HoldChip     board.Chip
	Level        int
	Score        int
	Time         int
	Turn         int
	ChipKindMax  int
	ChipColorMax int
}

var session = func() Session {
	session := Session{}
	return session
}()

func (session *Session) NewGame(level int, chipKindMax int, chipColorMax int) {
	rand.Seed(time.Now().UnixMicro())
	board.Board.PlaceChip(board.Chip{
		Kind:  board.ChipKindWild,
		Color: board.ChipColorA,
	}, [2]int{3, 4})
	session.Active = true
	session.Complete = false
	session.Level = level
	session.Time = 0
	session.Turn = 1
	session.ChipKindMax = chipKindMax
	session.ChipColorMax = chipColorMax
	session.HoldChip = board.Chip{
		Kind:  board.ChipKindEmpty,
		Color: board.ChipColorA,
	}
	session.GenNextChip()
	session.GenNextChip()
	pressureGauge.Level = 0
	pressureGauge.Stopped = false
	sideBar.SetDialog(DialogDrDokkanGameStart[rand.Intn(len(DialogDrDokkanGameStart))])
}

func (session *Session) GenNextChip() {
	for i := 0; true; i++ {
		randomChip := board.Board.GenRandomChip(session.ChipKindMax, session.ChipColorMax)
		if i == 0 {
			if randomChip.Special() {
				continue
			}
		}
		validCods := board.Board.ValidChipCods(randomChip)
		switch len(validCods) {
		case 0:
			continue
		case 1:
			if !randomChip.Special() {
				continue
			}
		}
		if !randomChip.Special() {
			if i < session.Turn {
				if board.Board.CoordinatesAreAllFilledCircuits(validCods) {
					continue
				}
			}
		}
		session.CurrentChip = session.NextChip.Clone()
		session.NextChip = randomChip
		break
	}
}

func (session *Session) PlaceChip(cod [2]int) {
	if session.Complete {
		return
	}
	if !board.Board.ValidChip(session.CurrentChip, cod) {
		pressureGauge.ChipMisplaced()
		return
	}
	board.Board.PlaceChip(session.CurrentChip, cod)
	session.Score += 10
	session.Turn++
	switch session.CurrentChip.Kind {
	case board.ChipKindWild:
		soundPlayer.Play("wild")
		boardRenderer.AddSolderSparks(50, color.RGBA{0x00, 0x00, 0x00, 0xff}, 255, 200, cod)
		boardRenderer.AddShake(float64(pressureGauge.Level/5), int(pressureGauge.Level/5))
	case board.ChipKindBomb:
		soundPlayer.Play("bomb" + strconv.Itoa(rand.Intn(2)))
		boardRenderer.AddSolderSparks(200, color.RGBA{0xcd, 0x80, 0x80, 0xff}, 100, 400, cod)
		boardRenderer.AddShake(20, 500)
	default:
		soundPlayer.Play("thunk")
		boardRenderer.AddSolderSparks(2*int(pressureGauge.Level), color.RGBA{0xcd, 0xcd, 0x32, 0xff}, 100, 200, cod)
		boardRenderer.AddShake(float64(pressureGauge.Level/5), int(pressureGauge.Level/5))
	}
	session.CurrentChip = board.Chip{
		Kind:  board.ChipKindEmpty,
		Color: board.ChipColorA,
		Frame: 0,
	}
	pressureGauge.ChipPlaced()
	clearableRows, clearableCols := board.Board.GetClearable()
	willClearLines := (len(clearableRows) + len(clearableCols)) > 0
	if willClearLines {
		pressureGauge.LineCleared()
		soundPlayer.Play("lineClear")
	}
	for _, row := range clearableRows {
		boardRenderer.AddLineClear(true, row, func(cod int) {
			board.Board.ClearRow(cod)
			session.Score += 100
			if board.Board.CircuitComplete() {
				session.LevelComplete()
			} else {
				session.GenNextChip()
			}
		})
	}
	for _, col := range clearableCols {
		boardRenderer.AddLineClear(false, col, func(cod int) {
			board.Board.ClearCol(cod)
			session.Score += 100
			if board.Board.CircuitComplete() {
				session.LevelComplete()
			} else {
				session.GenNextChip()
			}
		})
	}
	if !willClearLines {
		if board.Board.CircuitComplete() {
			session.LevelComplete()
		} else {
			session.GenNextChip()
		}
	}
	mouse.Sparkles = []MouseSparkle{}
	boardRenderer.SquareHighlight.Hide()
}

func (session *Session) ShuffleHold() {
	if session.Complete {
		return
	}
	switch session.HoldChip.Kind {
	case board.ChipKindEmpty:
		mouse.Sparkles = []MouseSparkle{}
		session.HoldChip = session.CurrentChip.Clone()
		session.GenNextChip()
		// soundPlayer.Play("discard" + strconv.Itoa(rand.Intn(2)))
	default:
		currentChip := session.CurrentChip.Clone()
		session.CurrentChip = session.HoldChip.Clone()
		session.HoldChip = currentChip.Clone()
	}
}

func (session *Session) LevelComplete() {
	if session.Complete {
		return
	}
	session.Complete = true
	pressureGauge.Stopped = true
	go func() {
		chipCods := board.Board.GetAllChipCods()
		rand.Shuffle(len(chipCods), func(i, j int) {
			chipCods[i], chipCods[j] = chipCods[j], chipCods[i]
		})
		for _, chipCod := range chipCods {
			time.Sleep(time.Millisecond * 50)
			if board.Board.CodHasChip(chipCod) {
				boardRenderer.AddSquareClear(chipCod,
					color.RGBA{
						uint8(128 + rand.Intn(127)),
						uint8(128 + rand.Intn(127)),
						uint8(128 + rand.Intn(127)), 0xff},
					func(cod [2]int) {
						board.Board.ClearChip(cod)
						session.Score += 50
					},
				)
			}
		}
		time.Sleep(time.Millisecond * 2000)
		session.ResetBoard(color.RGBA{0x00, 0xcd, 0x7e, 0xff})
	}()
}

func (session *Session) ResetBoard(squareClearColor color.RGBA) {
	session.CurrentChip = board.Chip{
		Kind:  board.ChipKindEmpty,
		Color: board.ChipColorA,
	}
	session.HoldChip = board.Chip{
		Kind:  board.ChipKindEmpty,
		Color: board.ChipColorA,
	}
	for row := 0; row < board.Board.Size[0]; row++ {
		for col := 0; col < board.Board.Size[1]; col++ {
			time.Sleep(time.Millisecond * 20)
			if pressureGauge.Level > 0 {
				pressureGauge.Level -= 1.4
			}
			if pressureGauge.Level < 0 {
				pressureGauge.Level = 0
			}
			squareClearColor.B += 1
			boardRenderer.AddSquareClear([2]int{row, col},
				squareClearColor,
				func(cod [2]int) {
					board.Board.Circuits[cod[0]][cod[1]].Filled = false
					if session.Complete {
						session.Score += 10
					}
					if cod[0] == 3 && cod[1] == 4 {
						board.Board.PlaceChip(board.Chip{
							Kind:  board.ChipKindWild,
							Color: board.ChipColorA,
							Frame: 0,
						}, cod)
					}
				},
			)
		}
	}
	session.NewGame(session.Level+1, session.ChipKindMax, session.ChipColorMax)
}

func (session *Session) GameOver() {
	pressureGauge.Stopped = true
	boardRenderer.CreateExplosion()
	go func() {
		time.Sleep(1000 * time.Millisecond)
		go session.ResetBoard(color.RGBA{0xff, 0x00, 0x00, 0xff})
		sideBar.SetDialog(DialogDrDokkanGameOver[rand.Intn(len(DialogDrDokkanGameOver))])
	}()
}
