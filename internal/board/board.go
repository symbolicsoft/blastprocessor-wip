/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package board

import (
	"math/rand"
	"time"
)

type BoardStruct struct {
	Size     [2]int
	Circuits [][]Circuit
	Chips    [][]Chip
}

var Board = func(size [2]int) BoardStruct {
	board := BoardStruct{}
	board.Size = size
	rand.Seed(time.Now().UnixMicro())
	board.Reset()
	return board
}([2]int{8, 9})

func (board *BoardStruct) Reset() {
	board.Circuits = make([][]Circuit, board.Size[0])
	board.Chips = make([][]Chip, board.Size[0])
	for row := 0; row < board.Size[0]; row++ {
		colCircuits := make([]Circuit, board.Size[1])
		colChips := make([]Chip, board.Size[1])
		for col := 0; col < board.Size[1]; col++ {
			colCircuits[col] = board.GetRandomCircuit()
			colChips[col] = Chip{ChipKindEmpty, ChipColorA, 0}
		}
		board.Circuits[row] = colCircuits
		board.Chips[row] = colChips
	}
}

func (board *BoardStruct) Randomize() {
	board.PlaceChip(Chip{ChipKindBatt, ChipColorB, 0}, [2]int{0, 0})
	for row := 0; row < board.Size[0]; row++ {
		for col := 0; col < board.Size[1]; col++ {
			if rand.Intn(2) == 0 {
				continue
			}
			randomChip := Chip{ChipKind(rand.Intn(ChipKindCount) + 1), ChipColor(rand.Intn(ChipColorCount)), 0}
			for randomChip.Kind == ChipKindWild {
				randomChip = Chip{ChipKind(rand.Intn(ChipKindCount) + 1), ChipColor(rand.Intn(ChipColorCount)), 0}
			}
			board.PlaceChip(randomChip, [2]int{row, col})
		}
	}
}

func (board *BoardStruct) CodHasChip(cod [2]int) bool {
	switch board.Chips[cod[0]][cod[1]].Kind {
	case ChipKindEmpty:
		return false
	default:
		return true
	}
}

func (board *BoardStruct) AdjacentToChip(cod [2]int) bool {
	if (cod[0] > board.Size[0]) || (cod[1] > board.Size[1]) {
		return false
	}
	adjacentToChip := false
	checkDirs := [][2]int{}
	if cod[0] > 0 {
		checkDirs = append(checkDirs, [2]int{cod[0] - 1, cod[1]})
	}
	if cod[0] < board.Size[0]-1 {
		checkDirs = append(checkDirs, [2]int{cod[0] + 1, cod[1]})
	}
	if cod[1] > 0 {
		checkDirs = append(checkDirs, [2]int{cod[0], cod[1] - 1})
	}
	if cod[1] < board.Size[1]-1 {
		checkDirs = append(checkDirs, [2]int{cod[0], cod[1] + 1})
	}
	for _, dir := range checkDirs {
		if board.CodHasChip(dir) {
			adjacentToChip = true
		}
	}
	return adjacentToChip
}

func (board *BoardStruct) ValidChip(chip Chip, cod [2]int) bool {
	if (cod[0] > board.Size[0]) || (cod[1] > board.Size[1]) {
		return false
	}
	checkDirs := [][2]int{}
	switch chip.Kind {
	case ChipKindEmpty:
		return false
	case ChipKindWild:
		if board.NoChips() {
			return true
		}
	case ChipKindBomb:
		return board.Chips[cod[0]][cod[1]].Kind != ChipKindEmpty
	}
	switch board.CodHasChip(cod) {
	case false:
		if cod[0] > 0 {
			checkDirs = append(checkDirs, [2]int{cod[0] - 1, cod[1]})
		}
		if cod[0] < board.Size[0]-1 {
			checkDirs = append(checkDirs, [2]int{cod[0] + 1, cod[1]})
		}
		if cod[1] > 0 {
			checkDirs = append(checkDirs, [2]int{cod[0], cod[1] - 1})
		}
		if cod[1] < board.Size[1]-1 {
			checkDirs = append(checkDirs, [2]int{cod[0], cod[1] + 1})
		}
		for _, dir := range checkDirs {
			if !chip.Compatible(board.Chips[dir[0]][dir[1]]) {
				return false
			}
		}
		return board.AdjacentToChip(cod)
	default:
		return false
	}
}

func (board *BoardStruct) PlaceChip(chip Chip, cod [2]int) {
	switch board.Circuits[cod[0]][cod[1]].Filled {
	case false:
		board.Circuits[cod[0]][cod[1]].Filled = true
	}
	switch chip.Kind {
	case ChipKindBomb:
		board.Chips[cod[0]][cod[1]] = Chip{
			ChipKindEmpty, ChipColorA,
			rand.Intn(ChipMaxFrame[chip.Kind]),
		}
	default:
		board.Chips[cod[0]][cod[1]] = Chip{
			chip.Kind, chip.Color,
			rand.Intn(ChipMaxFrame[chip.Kind]),
		}
	}
}

func (board *BoardStruct) GetClearable() ([]int, []int) {
	rows := []int{}
	cols := []int{}
	for row := 0; row < board.Size[0]; row++ {
		shouldClearRow := true
		for col := 0; col < board.Size[1]; col++ {
			if !board.CodHasChip([2]int{row, col}) {
				shouldClearRow = false
			}
		}
		if shouldClearRow {
			rows = append(rows, row)
		}
	}
	for col := 0; col < board.Size[1]; col++ {
		shouldClearCol := true
		for row := 0; row < board.Size[0]; row++ {
			if !board.CodHasChip([2]int{row, col}) {
				shouldClearCol = false
			}
		}
		if shouldClearCol {
			cols = append(cols, col)
		}
	}
	return rows, cols
}

func (board *BoardStruct) ClearChip(cod [2]int) {
	board.Chips[cod[0]][cod[1]] = Chip{ChipKindEmpty, ChipColorA, 0}
}

func (board *BoardStruct) ClearRow(row int) {
	for col := 0; col < board.Size[1]; col++ {
		board.ClearChip([2]int{row, col})
	}
}

func (board *BoardStruct) ClearCol(col int) {
	for row := 0; row < board.Size[0]; row++ {
		board.ClearChip([2]int{row, col})
	}
}

func (board *BoardStruct) CircuitComplete() bool {
	for row := 0; row < board.Size[0]; row++ {
		for col := 0; col < board.Size[1]; col++ {
			switch board.Circuits[row][col].Filled {
			case false:
				return false
			}
		}
	}
	return true
}

func (board *BoardStruct) NoChips() bool {
	for row := 0; row < board.Size[0]; row++ {
		for col := 0; col < board.Size[1]; col++ {
			if board.CodHasChip([2]int{row, col}) {
				return false
			}
		}
	}
	return true
}

func (board *BoardStruct) ValidChipCods(chip Chip) [][2]int {
	validCods := [][2]int{}
	for row := 0; row < board.Size[0]; row++ {
		for col := 0; col < board.Size[1]; col++ {
			if board.ValidChip(chip, [2]int{row, col}) {
				validCods = append(validCods, [2]int{row, col})
			}
		}
	}
	return validCods
}

func (board *BoardStruct) GetRandomCircuit() Circuit {
	return Circuit{
		Kind:       CircuitKind(rand.Intn(CircuitKindCount)),
		Filled:     false,
		Corruption: CircuitCorruptionNone,
	}
}

func (board *BoardStruct) GenRandomChip(chipKindMax int, chipColorMax int) Chip {
	randomChip := Chip{
		Kind:  ChipKind(rand.Intn(chipKindMax+2) + 1),
		Color: ChipColor(rand.Intn(chipColorMax)),
	}
	switch randomChip.Kind {
	case ChipKindWild:
		if rand.Intn(2) == 0 {
			randomChip.Color = ChipColorA
		} else {
			return board.GenRandomChip(chipKindMax, chipColorMax)
		}
	case ChipKindBomb:
		if rand.Intn(2) == 0 {
			randomChip.Color = ChipColorA
		} else {
			return board.GenRandomChip(chipKindMax, chipColorMax)
		}
	}
	return randomChip
}

func (board *BoardStruct) CoordinatesAreAllFilledCircuits(cods [][2]int) bool {
	for _, cod := range cods {
		if !board.Circuits[cod[0]][cod[1]].Filled {
			return false
		}
	}
	return true
}

func (board *BoardStruct) GetAllChipCods() [][2]int {
	cods := [][2]int{}
	for row := 0; row < board.Size[0]; row++ {
		for col := 0; col < board.Size[1]; col++ {
			if board.CodHasChip([2]int{row, col}) {
				cods = append(cods, [2]int{row, col})
			}
		}
	}
	return cods
}
