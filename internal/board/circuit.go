/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package board

type Circuit struct {
	Kind       CircuitKind
	Filled     bool
	Corruption CircuitCorruption
}

type CircuitKind int

const CircuitKindCount = 32

type CircuitCorruption int

const CircuitCorruptionCount = 2

const (
	CircuitCorruptionNone  CircuitCorruption = iota
	CircuitCorruptionPopup CircuitCorruption = iota
)
