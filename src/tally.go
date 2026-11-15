package main

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

type TallyApp struct {
	selectedIndexX int
	selectedIndexY int
	currentCount   int
	boxWidth       int
	boxHeight      int
}

func (a *TallyApp) boxLine(selected bool, line int, label string) string {
	accent := ""
	accentBg := ""
	reset := ""
	if selected {
		accent = Red
		accentBg = RedBg
		reset = Reset
	}

	switch {
	case line == 0:
		return accent + "┌" + strings.Repeat("─", a.boxWidth-2) + "┐ " + reset
	case line == a.boxHeight-1:
		return accent + "└" + strings.Repeat("─", a.boxWidth-2) + "┘ " + reset
	case line == a.boxHeight/2:
		padding := a.boxWidth - 4 - runewidth.StringWidth(label)
		if padding < 0 {
			padding = 0
		}

		return accent + "│" + accentBg + strings.Repeat(" ", (a.boxWidth-len(label))/2) + White + label + strings.Repeat(" ", (a.boxWidth-4-len(label))/2) + reset + accent + "│ " + reset
	default:
		return accent + "│" + accentBg + strings.Repeat(" ", a.boxWidth-2) + reset + accent + "│ " + reset
	}
}

func (a *TallyApp) Name() string { return "Tally" }

func (a *TallyApp) Start() error {
	a.boxWidth = 41
	a.boxHeight = 5
	return nil
}

func (a *TallyApp) Stop() error { return nil }

func (a *TallyApp) HandleEvent(event Event) (Action, error) {
	switch event.Key {
	case KeyLeft:
		a.selectedIndexX = 0
	case KeyRight:
		a.selectedIndexX = 1
	case KeyUp:
		a.selectedIndexY = 0
	case KeyDown:
		a.selectedIndexY = 1
	case KeyOK:
		if a.selectedIndexY == 0 {
			if a.selectedIndexX == 0 {
				a.currentCount--
			} else {
				a.currentCount++
			}
		} else {
			a.currentCount = 0
		}
	}

	return Action{Kind: ActionNone}, nil
}

func (a *TallyApp) Render() error {
	var frame []string

	frame = append(frame, getTitle("Tally", Red))

	for i := 0; i < 3; i++ {
		var row strings.Builder

		row.WriteString(strings.Repeat(" ", 5))

		if i == 1 {
			accentBg := ""
			reset := ""

			if a.selectedIndexX == 0 && a.selectedIndexY == 0 {
				accentBg = RedBg
				reset = Reset
			}

			row.WriteString(accentBg)
			row.WriteString(" - ")
			row.WriteString(reset)
			row.WriteRune(' ')
		} else {
			row.WriteString("    ")
		}

		paddedStr := fmt.Sprintf("%09d", a.currentCount)

		row.WriteString(Grey)

		for j := 0; j < len(paddedStr); j++ {
			currentChar := string(paddedStr[j])

			if currentChar != "0" {
				row.WriteString(Reset)
			}

			if j == len(paddedStr)-1 {
				row.WriteString(Reset)
			}

			row.WriteString(string(artRows(currentChar, numberArt, 4)[i]))

			if j == len(paddedStr)-1 {
				row.WriteString(Grey)
			}
		}

		row.WriteString(Reset)

		if i == 1 {
			accentBg := ""
			reset := ""

			if a.selectedIndexX == 1 && a.selectedIndexY == 0 {
				accentBg = RedBg
				reset = Reset
			}

			row.WriteString(accentBg)
			row.WriteString(" + ")
			row.WriteString(reset)
			row.WriteRune(' ')

		} else {
			row.WriteString("    ")
		}

		frame = append(frame, row.String())
	}

	for line := 0; line < a.boxHeight; line++ {
		var row strings.Builder
		selected := a.selectedIndexY == 1
		row.WriteString(strings.Repeat(" ", 5))
		row.WriteString(a.boxLine(selected, line, "Reset"))
		frame = append(frame, row.String())
	}

	presentFrame(frame)

	return nil
}
