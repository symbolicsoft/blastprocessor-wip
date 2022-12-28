/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package board

type Chip struct {
	Kind  ChipKind
	Color ChipColor
	Frame int
}

type ChipKind int
type ChipColor int

const ChipKindCount = 12
const ChipColorCount = 8

var ChipMaxFrame = map[ChipKind]int{
	ChipKindEmpty: 1,
	ChipKindWild:  3,
	ChipKindBomb:  3,
	ChipKindCap:   1,
	ChipKindMem:   1,
	ChipKindRes:   1,
	ChipKindBatt:  1,
	ChipKindSDC:   1,
	ChipKindHeat:  1,
	ChipKindFan:   1,
	ChipKindCore:  1,
	ChipKindPS2:   1,
	ChipKindWiFi:  1,
}

const (
	ChipKindEmpty ChipKind = iota
	ChipKindWild  ChipKind = iota
	ChipKindBomb  ChipKind = iota
	ChipKindCap   ChipKind = iota
	ChipKindMem   ChipKind = iota
	ChipKindRes   ChipKind = iota
	ChipKindBatt  ChipKind = iota
	ChipKindSDC   ChipKind = iota
	ChipKindHeat  ChipKind = iota
	ChipKindFan   ChipKind = iota
	ChipKindCore  ChipKind = iota
	ChipKindPS2   ChipKind = iota
	ChipKindWiFi  ChipKind = iota
)

const (
	ChipColorA ChipColor = iota
	ChipColorB ChipColor = iota
	ChipColorC ChipColor = iota
	ChipColorD ChipColor = iota
	ChipColorE ChipColor = iota
	ChipColorF ChipColor = iota
	ChipColorG ChipColor = iota
	ChipColorH ChipColor = iota
)

func (chip Chip) Identical(chipB Chip) bool {
	return (chip.Kind == chipB.Kind) &&
		(chip.Color == chipB.Color)
}

func (chip Chip) Special() bool {
	switch chip.Kind {
	case ChipKindWild, ChipKindBomb:
		return true
	default:
		return false
	}
}

func (chip Chip) Compatible(chipB Chip) bool {
	return (chip.Kind == ChipKindEmpty) ||
		(chipB.Kind == ChipKindEmpty) ||
		(chip.Kind == ChipKindWild) ||
		(chipB.Kind == ChipKindWild) ||
		(chip.Kind == chipB.Kind) ||
		(chip.Color == chipB.Color)
}

func (chip Chip) Clone() Chip {
	return Chip{
		Kind:  chip.Kind,
		Color: chip.Color,
		Frame: chip.Frame,
	}
}
