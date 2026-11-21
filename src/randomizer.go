package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

type RandomizerApp struct {
	min            int
	max            int
	rnd            int
	selectedIndexX int
	selectedIndexY int
	boxWidth       int
	boxHeight      int
}

func (a *RandomizerApp) boxLine(selected bool, line int, label string) string {
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

func (a *RandomizerApp) Name() string { return "Randomizer" }

func (a *RandomizerApp) Start() error {
	if a.min >= a.max {
		a.max = 10
	}
	a.boxWidth = 51
	a.boxHeight = 5
	return nil
}

func (a *RandomizerApp) Stop() error { return nil }

func (a *RandomizerApp) HandleEvent(event Event) (Action, error) {
	switch event.Key {
	case KeyOK:
		switch a.selectedIndexY {
		case 0:
			switch a.selectedIndexX {
			case 0:
				if a.min+10 < a.max {
					a.min += 10
				}
			case 1:
				a.max += 10
			}
		case 1:
			switch a.selectedIndexX {
			case 0:
				if a.min+1 < a.max {
					a.min++
				}
			case 1:
				a.max++
			}
		case 2:
			switch a.selectedIndexX {
			case 0:
				if a.min-1 >= 0 {
					a.min--
				}
			case 1:
				if a.max-1 > a.min {
					a.max--
				}
			}
		case 3:
			switch a.selectedIndexX {
			case 0:
				if a.min-10 >= 0 {
					a.min -= 10
				}
			case 1:
				if a.max-10 > a.min {
					a.max -= 10
				}
			}
		}
		if a.selectedIndexY == 4 {
			a.rnd = rand.Intn(a.max-a.min) + a.min
		}
	case KeyLeft:
		if a.selectedIndexX > 0 {
			a.selectedIndexX--
		}
	case KeyRight:
		if a.selectedIndexX < 1 {
			a.selectedIndexX++
		}
	case KeyUp:
		if a.selectedIndexY > 0 {
			a.selectedIndexY--
		}
	case KeyDown:
		if a.selectedIndexY < 4 {
			a.selectedIndexY++
		}
	}

	return Action{Kind: ActionNone}, nil
}

func (a *RandomizerApp) Render() error {
	var frame []string

	frame = append(frame, getTitle("Randomizer", Red))

	margin := 15

	// ADD 10
	var addTime10 strings.Builder
	addTime10.WriteString(strings.Repeat(" ", margin))

	for i := 0; i < 2; i++ {
		accentBg := ""
		reset := ""

		if a.selectedIndexX == i && a.selectedIndexY == 0 {
			accentBg = RedBg
			reset = Reset
		}

		addTime10.WriteString("  ")
		addTime10.WriteString(accentBg)
		addTime10.WriteString("  +10  ")
		addTime10.WriteString(reset)
		addTime10.WriteString("   ")
	}
	addTime10.WriteString("\n")
	frame = append(frame, addTime10.String())

	// ADD 1
	var addTime strings.Builder
	addTime.WriteString(strings.Repeat(" ", margin))

	for i := 0; i < 2; i++ {
		accentBg := ""
		reset := ""

		if a.selectedIndexX == i && a.selectedIndexY == 1 {
			accentBg = RedBg
			reset = Reset
		}

		addTime.WriteString("   ")
		addTime.WriteString(accentBg)
		addTime.WriteString("  +  ")
		addTime.WriteString(reset)
		addTime.WriteString("    ")
	}
	addTime.WriteString("\n")
	frame = append(frame, addTime.String())

	// MIN and MAX
	frame = append(frame, strings.Repeat(" ", 20)+Bold+Red+strconv.Itoa(a.min)+strings.Repeat(" ", 11-len(strconv.Itoa(a.min))+1)+strconv.Itoa(a.max)+Reset)
	frame = append(frame, "")

	// SUB 1
	var subTime strings.Builder
	subTime.WriteString(strings.Repeat(" ", margin))

	for i := 0; i < 2; i++ {
		accentBg := ""
		reset := ""

		if a.selectedIndexX == i && a.selectedIndexY == 2 {
			accentBg = RedBg
			reset = Reset
		}

		subTime.WriteString("   ")
		subTime.WriteString(accentBg)
		subTime.WriteString("  -  ")
		subTime.WriteString(reset)
		subTime.WriteString("    ")
	}
	subTime.WriteString("\n")
	frame = append(frame, subTime.String())

	// SUB 10
	var subTime10 strings.Builder
	subTime10.WriteString(strings.Repeat(" ", margin))

	for i := 0; i < 2; i++ {
		accentBg := ""
		reset := ""

		if a.selectedIndexX == i && a.selectedIndexY == 3 {
			accentBg = RedBg
			reset = Reset
		}

		subTime10.WriteString("  ")
		subTime10.WriteString(accentBg)
		subTime10.WriteString("  -10  ")
		subTime10.WriteString(reset)
		subTime10.WriteString("   ")
	}
	subTime10.WriteString("\n")
	frame = append(frame, subTime10.String())

	for i := 0; i < 3; i++ {
		var row strings.Builder

		row.WriteString(strings.Repeat(" ", 9))

		paddedStr := fmt.Sprintf("%09d", a.rnd)

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

		frame = append(frame, row.String())
	}

	for line := 0; line < a.boxHeight; line++ {
		var row strings.Builder
		selected := a.selectedIndexY == 4
		row.WriteString(a.boxLine(selected, line, "Randomize"))
		frame = append(frame, row.String())
	}

	presentFrame(frame)
	return nil
}
