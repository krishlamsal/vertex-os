package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
)

type StopwatchApp struct {
	mu sync.Mutex

	selectedIndexX  int
	stopwatchActive bool
	boxWidth        int
	boxHeight       int
	startTime       time.Time
	elapsed         time.Duration
	laps            []time.Time

	stop      chan struct{}
	rendering bool
}

func (a *StopwatchApp) boxLine(selected bool, line int, label string, accentColor int) string {
	accent := ""
	accentBg := ""
	reset := ""
	if selected {
		accent = Colors[accentColor*2]
		accentBg = Colors[accentColor*2+1]
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

func (a *StopwatchApp) Name() string { return "Stopwatch" }

func (a *StopwatchApp) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.stop != nil {
		close(a.stop)
	}

	a.startTime = time.Now()
	a.boxWidth = 27
	a.boxHeight = 5
	a.stop = make(chan struct{})
	a.rendering = false
	return nil
}

func (a *StopwatchApp) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.stop != nil {
		close(a.stop)
		a.stop = nil
	}
	a.stopwatchActive = false
	return nil
}

func (a *StopwatchApp) HandleEvent(event Event) (Action, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	switch event.Key {
	case KeyLeft:
		a.selectedIndexX = 0
	case KeyRight:
		a.selectedIndexX = 1
	case KeyOK:
		if a.selectedIndexX == 0 {
			if a.stopwatchActive {
				a.elapsed += time.Since(a.startTime)
				a.stopwatchActive = false
			} else {
				a.startTime = time.Now()
				a.stopwatchActive = true
			}
		} else {
			if !a.stopwatchActive {
				a.elapsed = 0
				a.laps = []time.Time{}
			} else {
				a.laps = append(a.laps, time.Time{}.Add(time.Since(a.startTime)).Add(a.elapsed))
			}
		}
	}

	return Action{Kind: ActionNone}, nil
}

func (a *StopwatchApp) Render() error {
	a.mu.Lock()
	if a.stop == nil || a.rendering {
		a.mu.Unlock()
		return nil
	}
	a.rendering = true
	stopCh := a.stop
	a.mu.Unlock()

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		defer func() {
			a.mu.Lock()
			if a.stop == stopCh {
				a.rendering = false
			}
			a.mu.Unlock()
		}()

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				a.mu.Lock()
				if a.stop != stopCh {
					a.mu.Unlock()
					return
				}

				var frame []string

				frame = append(frame, getTitle("Stopwatch", Blue))

				drawLine := func(location string) {
					if location == "top" {
						frame = append(frame, strings.Repeat(" ", 3)+"┌"+strings.Repeat("─", 45)+"┐")
					} else {
						frame = append(frame, strings.Repeat(" ", 3)+"└"+strings.Repeat("─", 45)+"┘")
					}
				}

				drawLine("top")

				d := a.elapsed
				if a.stopwatchActive {
					d += time.Since(a.startTime)
				}

				hours := int(d / time.Hour)
				minutes := int(d/time.Minute) % 60
				seconds := int(d/time.Second) % 60
				milliseconds := int(d/time.Millisecond) % 1000

				formatted := fmt.Sprintf("%d:%02d:%02d.%03d",
					hours, minutes, seconds, milliseconds)

				for i := 0; i < 3; i++ {
					var row strings.Builder

					row.WriteString(strings.Repeat(" ", 3))
					row.WriteString("│ ")
					for j := 0; j < len(formatted); j++ {
						row.WriteString(string(artRows(string(formatted[j]), numberArt, 4)[i]))
					}
					row.WriteString("│")

					frame = append(frame, row.String())
				}
				drawLine("")

				labels := []string{"Start", "Reset"}
				accents := []int{5, 3}

				if a.stopwatchActive {
					labels[0] = "Stop "
					labels[1] = "Lap  "
					accents[0] = 0
				}

				for line := 0; line < a.boxHeight; line++ {
					var row strings.Builder
					for x := 0; x < 2; x++ {
						selected := x == a.selectedIndexX
						row.WriteString(a.boxLine(selected, line, labels[x], accents[x]))
					}
					frame = append(frame, row.String())
				}

				drawLine("top")

				if len(a.laps) == 0 {
					text := "No laps..."
					frame = append(frame, strings.Repeat(" ", 3)+"│"+Blue+text+Reset+strings.Repeat(" ", 45-runewidth.StringWidth(text))+"│")
				}

				for i, v := range a.laps {
					if i == len(a.laps)-14 {
						text := "..."
						frame = append(frame, strings.Repeat(" ", 3)+"│"+Grey+text+Reset+strings.Repeat(" ", 45-runewidth.StringWidth(text))+"│")
					}
					if i > len(a.laps)-14 {
						text := Blue + "Lap " + strconv.Itoa(i+1) + Grey + strings.Repeat(" ", 5-len(strconv.Itoa(i+1))) + "·  " + Cyan + v.Format("15:04:05") + Reset + v.Format(".000")
						cleanText := "Lap " + strconv.Itoa(i+1) + strings.Repeat(" ", 5-len(strconv.Itoa(i+1))) + "·  " + v.Format("15:04:05") + v.Format(".000")
						frame = append(frame, strings.Repeat(" ", 3)+"│"+text+strings.Repeat(" ", 45-runewidth.StringWidth(cleanText))+"│")
					}
				}

				drawLine("")

				a.mu.Unlock()

				presentFrame(frame)
			}

		}
	}()

	return nil
}
