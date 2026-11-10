package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

type HomeApp struct {
	selectedIndexX int
	selectedIndexY int
	currentPage    int
	maxPage        int
	appNames       []string
	appLogos       []string
	pageTitles     []string
	accent         string
	accentBg       string
	boxWidth       int
	boxHeight      int
}

func randomizeColors(a *HomeApp) {
	oldColor := a.accent

	for oldColor == a.accent {
		colorIdx := rand.Intn(len(Colors)/2) * 2
		a.accent = Colors[colorIdx]
		a.accentBg = Colors[colorIdx+1]
	}
}

func (a *HomeApp) Name() string { return "Home" }

func (a *HomeApp) Start() error {
	// The app boxes are 26 columns wide. Terminal cells are taller than they
	// are wide, so 13 rows makes each box appear roughly square on screen.

	a.boxWidth = 26
	a.boxHeight = 13

	randomizeColors(a)
	return nil
}

func (a *HomeApp) Stop() error { return nil }

func (a *HomeApp) HandleEvent(event Event) (Action, error) {
	oldX, oldY := a.selectedIndexX, a.selectedIndexY
	oldPage := a.currentPage

	switch event.Key {
	case KeyOK:
		if len(a.appNames) == 0 {
			return Action{Kind: ActionNone}, nil
		} else {
			var index int
			if a.selectedIndexX == 0 && a.selectedIndexY == 0 {
				index = 0
			} else if a.selectedIndexX == 1 && a.selectedIndexY == 0 {
				index = 1
			} else if a.selectedIndexX == 0 && a.selectedIndexY == 1 {
				index = 2
			} else if a.selectedIndexX == 1 && a.selectedIndexY == 1 {
				index = 3
			}

			return Action{
				Kind:    ActionLaunch,
				AppName: a.appNames[a.currentPage*4+index],
			}, nil
		}
	case KeyUp:
		a.selectedIndexY = 0
	case KeyDown:
		a.selectedIndexY = 1
	case KeyLeft:
		a.selectedIndexX = 0
	case KeyRight:
		a.selectedIndexX = 1
	case KeyPageLeft:
		if a.currentPage > 0 {
			a.currentPage--
		} else {
			a.currentPage = a.maxPage
		}
	case KeyPageRight:
		if a.currentPage < a.maxPage {
			a.currentPage++
		} else {
			a.currentPage = 0
		}
	}

	if a.selectedIndexX != oldX || a.selectedIndexY != oldY || a.currentPage != oldPage {
		randomizeColors(a)
	}

	return Action{Kind: ActionNone}, nil
}

func (a *HomeApp) Render() error {
	pageCount := (len(a.appNames) + 3) / 4
	if pageCount == 0 {
		presentFrame([]string{"No apps installed.", "", "Use arrows to navigate."})
		return nil
	}

	// Keep the page within the valid range.
	if a.currentPage < 0 {
		a.currentPage = 0
	}
	if a.currentPage >= pageCount {
		a.currentPage = pageCount - 1
	}

	start := a.currentPage * 4

	gridWidth := 52
	pageTitle := a.pageTitles[a.currentPage]
	pageTitleLines := strings.Repeat("─", gridWidth/2-len(pageTitle)/2-1)
	pageTitleColor := Colors[a.currentPage*2]

	var frame []string
	frame = append(frame, pageTitleColor+pageTitleLines+" "+pageTitle+" "+pageTitleLines+Reset)

	for y := 0; y < 2; y++ {
		var logos, labels [2]string
		var art [2][]string

		for x := 0; x < 2; x++ {
			index := start + y*2 + x
			if index < len(a.appNames) {
				labels[x] = a.appNames[index]
				logos[x] = a.appLogos[index]
				art[x] = artRows(labels[x], appArt, 22 /* App Art maximum width*/)
			}
		}

		for line := 0; line < a.boxHeight; line++ {
			var row strings.Builder
			for x := 0; x < 2; x++ {
				selected := x == a.selectedIndexX && y == a.selectedIndexY
				row.WriteString(a.boxLine(selected, line, logos[x], labels[x], art[x]))
			}
			frame = append(frame, row.String())
		}
	}

	// Page Dots
	var pageDots, plainDots string

	for i := 0; i < pageCount; i++ {
		if i == a.currentPage {
			label := "  " + strconv.Itoa(i+1)
			pageDots += Reset + label
			plainDots += label
		} else {
			label := "  •"
			pageDots += Grey + label + Reset
			plainDots += label
		}
	}

	padding := (gridWidth - runewidth.StringWidth(plainDots)) / 2
	if padding < 0 {
		padding = 0
	}

	frame = append(frame, strings.Repeat(" ", padding)+pageDots)

	presentFrame(frame)

	return nil
}

func (a *HomeApp) boxLine(selected bool, line int, logo, label string, art []string) string {
	accent := ""
	reset := ""
	if selected {
		accent = a.accent
		reset = Reset
	}

	switch {
	case line == 0:
		return accent + "┌" + strings.Repeat("─", a.boxWidth-2) + "┐ " + reset
	case line == a.boxHeight-1:
		return accent + "└" + strings.Repeat("─", a.boxWidth-2) + "┘ " + reset
	case line == a.boxHeight-2:
		marker := " "
		if selected {
			marker = ">"
		}

		plain := marker + "  " + logo + "   " + label
		padding := a.boxWidth - 4 - runewidth.StringWidth(plain)
		if padding < 0 {
			padding = 0
		}

		if selected {
			return accent + "│ " + marker + Reset + " " + a.accentBg + " " + logo + "   " + label +
				strings.Repeat(" ", padding) + Reset + accent + " │ " + reset
		}
		return "│ " + marker + "  " + logo + "   " + Grey + label +
			strings.Repeat(" ", padding) + Reset + " │ "
	case line >= 1 && line <= a.boxHeight-3:
		// Rows between the top border and the footer hold the app sprite.
		regionTop := 1
		regionHeight := a.boxHeight - 3
		offset := regionTop + (regionHeight-len(art))/2
		if offset < regionTop {
			offset = regionTop
		}

		if idx := line - offset; idx >= 0 && idx < len(art) {
			color := Reset
			if selected {
				color = a.accent
			}
			return color + "│ " + art[idx] + " │ " + reset
		}

		if selected {
			return accent + "│" + strings.Repeat(" ", a.boxWidth-2) + "│ " + reset
		}
		return "│" + strings.Repeat(" ", a.boxWidth-2) + "│ "
	default:
		if selected {
			return accent + "│" + strings.Repeat(" ", a.boxWidth-2) + "│ " + reset
		}
		return "│" + strings.Repeat(" ", a.boxWidth-2) + "│ "
	}
}
