package main

import (
	"embed"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"github.com/mattn/go-runewidth"
)

//go:embed metronomeTick.wav
var metronomeSounds embed.FS

var (
	metronomeSoundOnce sync.Once
	metronomeSound     *beep.Buffer
	metronomeSoundErr  error
)

func loadMetronomeSound() (*beep.Buffer, error) {
	metronomeSoundOnce.Do(func() {
		f, err := metronomeSounds.Open("metronomeTick.wav")
		if err != nil {
			metronomeSoundErr = err
			return
		}
		defer f.Close()

		streamer, format, err := wav.Decode(f)
		if err != nil {
			metronomeSoundErr = err
			return
		}
		defer streamer.Close()

		if err := initSpeaker(format); err != nil {
			metronomeSoundErr = err
			return
		}

		metronomeSound = beep.NewBuffer(format)
		metronomeSound.Append(streamer)
	})

	return metronomeSound, metronomeSoundErr
}

func playMetronomeTick(buf *beep.Buffer) {
	if buf == nil {
		return
	}
	speaker.Play(buf.Streamer(0, buf.Len()))
}

type MetronomeApp struct {
	mu sync.Mutex

	selectedIndexX  int
	selectedIndexY  int
	bpm             int
	metronomeActive bool
	boxWidth        int
	boxHeight       int

	tickStop  chan struct{}
	tickSound *beep.Buffer
}

func (a *MetronomeApp) startTickingLocked() {
	if a.tickStop != nil || a.tickSound == nil {
		return
	}

	stopCh := make(chan struct{})
	a.tickStop = stopCh
	buf := a.tickSound
	playMetronomeTick(buf)

	go func() {
		for {
			a.mu.Lock()
			if a.tickStop != stopCh || !a.metronomeActive {
				a.mu.Unlock()
				return
			}
			bpm := a.bpm
			a.mu.Unlock()

			timer := time.NewTimer(time.Minute / time.Duration(bpm))
			select {
			case <-stopCh:
				timer.Stop()
				return
			case <-timer.C:
			}

			a.mu.Lock()
			if a.tickStop != stopCh || !a.metronomeActive {
				a.mu.Unlock()
				return
			}
			a.mu.Unlock()
			playMetronomeTick(buf)
		}
	}()
}

func (a *MetronomeApp) stopTickingLocked() {
	if a.tickStop != nil {
		close(a.tickStop)
		a.tickStop = nil
	}
}

func (a *MetronomeApp) boxLine(selected bool, line int, label string, accentColor int) string {
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

func (a *MetronomeApp) Name() string { return "Metronome" }

func (a *MetronomeApp) Start() error {
	buf, err := loadMetronomeSound()
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.boxWidth = 21
	a.boxHeight = 5

	if a.bpm < 40 {
		a.bpm = 120
	}

	a.tickStop = nil
	a.tickSound = buf
	a.metronomeActive = false

	return nil
}

func (a *MetronomeApp) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.stopTickingLocked()
	a.metronomeActive = false

	return nil
}

func (a *MetronomeApp) HandleEvent(event Event) (Action, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

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
				if a.bpm > 40 {
					a.bpm--
				}
			} else {
				if a.bpm < 600 {
					a.bpm++
				}
			}
		} else {
			if a.selectedIndexX == 0 {
				if a.metronomeActive {
					a.metronomeActive = false
					a.stopTickingLocked()
				} else {
					a.metronomeActive = true
					a.startTickingLocked()
				}
			} else {
				a.bpm = 120
				a.metronomeActive = false
				a.stopTickingLocked()
			}
		}
	}

	return Action{Kind: ActionNone}, nil
}

func (a *MetronomeApp) Render() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var frame []string
	frame = append(frame, getTitle("Metronome", Red))

	for i := 0; i < 3; i++ {
		var row strings.Builder

		row.WriteString(strings.Repeat(" ", 17))

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

		paddedStr := fmt.Sprintf("%03d", a.bpm)

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

	frame = append(frame, strings.Repeat(" ", 25)+Bold+Red+"BPM"+Reset)

	labels := []string{"Start", "Reset"}
	accents := []int{5, 3}

	if a.metronomeActive {
		labels[0] = "Stop "
		accents[0] = 0
	}

	for line := 0; line < a.boxHeight; line++ {
		var row strings.Builder
		row.WriteString(strings.Repeat(" ", 5))
		for x := 0; x < 2; x++ {
			selected := x == a.selectedIndexX && a.selectedIndexY == 1
			row.WriteString(a.boxLine(selected, line, labels[x], accents[x]))
		}
		frame = append(frame, row.String())
	}

	presentFrame(frame)
	return nil
}
