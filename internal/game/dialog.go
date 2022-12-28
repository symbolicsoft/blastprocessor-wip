/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

package game

type DialogItem struct {
	Face string
	Text string
}

var DialogDrDokkanGameStart = []DialogItem{
	{
		Face: "",
		Text: "Still blowin' up\nboards after\nall these\nyears?!",
	},
}

var DialogDrDokkanCompleted = []DialogItem{
	{
		Face: "",
		Text: "Completed",
	},
}

var DialogDrDokkanGameOver = []DialogItem{
	{
		Face: "",
		Text: "Still blowin' up\nboards after\nall these\nyears?!",
	},
}
