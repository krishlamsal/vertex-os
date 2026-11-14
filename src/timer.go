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

//go:embed timerSound.wav
var timerSounds embed.FS

var (
	timerSoundOnce sync.Once
	timerSound     *beep.Buffer
	timerSoundErr  error
)

type TimerApp struct {
	mu sync.Mutex

	selectedIndexX int
	selectedIndexY int
	hours          int
	minutes        int
	seconds        int
	timerActive    bool

	duration  time.Duration
	remaining time.Duration
	startTime time.Time

	stop      chan struct{}
	rendering bool

	boxWidth  int
	boxHeight int
}

func loadTimerSound() (*beep.Buffer, error) {
	timerSoundOnce.Do(func() {
		f, err := timerSounds.Open("timerSound.wav")
		if err != nil {
			timerSoundErr = err
			return
		}
		defer f.Close()

		streamer, format, err := wav.Decode(f)
		if err != nil {
			timerSoundErr = err
			return
		}
		defer streamer.Close()

		if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
			timerSoundErr = err
			return
		}

		timerSound = beep.NewBuffer(format)
		timerSound.Append(streamer)
	})

	return timerSound, timerSoundErr
}

func playTimerSound() {
	if timerSound == nil {
		return
	}

	speaker.Play(timerSound.Streamer(0, timerSound.Len()))
}

func (a *TimerApp) configuredDuration() time.Duration {
	return time.Duration(a.hours)*time.Hour +
		time.Duration(a.minutes)*time.Minute +
		time.Duration(a.seconds)*time.Second
}

func (a *TimerApp) boxLine(selected bool, line int, label string, accentColor int) string {
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

func (a *TimerApp) Name() string { return "Timer" }

func (a *TimerApp) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, err := loadTimerSound(); err != nil {
		return err
	}

	a.boxWidth = 41
	a.boxHeight = 5
	a.stop = make(chan struct{})
	a.rendering = false
	a.timerActive = false
	a.duration = a.configuredDuration()
	a.remaining = a.duration

	return nil
}

func (a *TimerApp) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.stop != nil {
		close(a.stop)
		a.stop = nil
	}
	a.timerActive = false

	return nil
}

func (a *TimerApp) HandleEvent(event Event) (Action, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	switch event.Key {
	case KeyOK:
		if !a.timerActive {
			switch a.selectedIndexY {
			case 0:
				switch a.selectedIndexX {
				case 0:
					if a.hours < 99 {
						a.hours++
					}
				case 1:
					if a.minutes < 59 {
						a.minutes++
					}
				case 2:
					if a.seconds < 59 {
						a.seconds++
					}
				}
			case 1:
				switch a.selectedIndexX {
				case 0:
					if a.hours > 0 {
						a.hours--
					}
				case 1:
					if a.minutes > 0 {
						a.minutes--
					}
				case 2:
					if a.seconds > 0 {
						a.seconds--
					}
				}
			}

			if a.selectedIndexY == 0 || a.selectedIndexY == 1 {
				a.duration = a.configuredDuration()
				a.remaining = a.duration
			}
		}
		if a.selectedIndexY == 2 {
			if a.timerActive {
				a.timerActive = false
				a.duration = a.configuredDuration()
				a.remaining = a.duration
			} else {
				a.duration = a.configuredDuration()
				if a.duration > 0 {
					a.remaining = a.duration
					a.startTime = time.Now()
					a.timerActive = true
				}
			}
		}
	case KeyLeft:
		if a.selectedIndexX > 0 {
			a.selectedIndexX--
		}
	case KeyRight:
		if a.selectedIndexX < 2 {
			a.selectedIndexX++
		}
	case KeyUp:
		if a.selectedIndexY > 0 {
			a.selectedIndexY--
		}
	case KeyDown:
		if a.selectedIndexY < 2 {
			a.selectedIndexY++
		}
	}

	return Action{Kind: ActionNone}, nil
}

func (a *TimerApp) Render() error {
	a.mu.Lock()
	if a.stop == nil || a.rendering {
		a.mu.Unlock()
		return nil
	}
	a.rendering = true
	stopCh := a.stop
	a.mu.Unlock()

	renderFrame := func() {
		a.mu.Lock()
		if a.stop != stopCh {
			a.mu.Unlock()
			return
		}

		remaining := a.remaining
		if a.timerActive {
			remaining = a.duration - time.Since(a.startTime)
			if remaining <= 0 {
				remaining = 0
				a.remaining = 0
				a.timerActive = false
				playTimerSound()
			}
		}

		hours := int(remaining / time.Hour)
		minutes := int(remaining/time.Minute) % 60
		seconds := int(remaining/time.Second) % 60
		timeRemaining := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

		var frame []string

		frame = append(frame, getTitle("Timer", Blue))

		var addTime strings.Builder

		if !a.timerActive {
			addTime.WriteString(strings.Repeat(" ", 10))

			for i := 0; i < 3; i++ {
				accentBg := ""
				reset := ""

				if a.selectedIndexX == i && a.selectedIndexY == 0 {
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
		} else {
			frame = append(frame, "\n")
		}

		for i := 0; i < 3; i++ {
			var row strings.Builder

			row.WriteString(strings.Repeat(" ", 12))
			for j := 0; j < len(timeRemaining); j++ {
				accent := ""
				reset := ""
				if a.timerActive && remaining.Seconds() < 10 {
					accent = Red
					reset = Reset
				}
				row.WriteString(accent)
				row.WriteString(string(artRows(string(timeRemaining[j]), numberArt, 4)[i]))
				row.WriteString(reset)
			}
			frame = append(frame, row.String())
		}

		var subtractTime strings.Builder

		if !a.timerActive {
			subtractTime.WriteString("\n")

			subtractTime.WriteString(strings.Repeat(" ", 10))

			for i := 0; i < 3; i++ {
				accentBg := ""
				reset := ""

				if a.selectedIndexX == i && a.selectedIndexY == 1 {
					accentBg = RedBg
					reset = Reset
				}

				subtractTime.WriteString("   ")
				subtractTime.WriteString(accentBg)
				subtractTime.WriteString("  -  ")
				subtractTime.WriteString(reset)
				subtractTime.WriteString("    ")
			}

			frame = append(frame, subtractTime.String())
		} else {
			frame = append(frame, "\n")
		}

		for line := 0; line < a.boxHeight; line++ {
			var row strings.Builder
			row.WriteString(strings.Repeat(" ", 7))
			label := "Start"
			accent := 1
			if a.timerActive {
				label = "Reset"
				accent = 0
			}
			row.WriteString(a.boxLine(a.selectedIndexY == 2, line, label, accent))
			frame = append(frame, row.String())
		}

		a.mu.Unlock()

		presentFrame(frame)
	}

	renderFrame()

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
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
				a.mu.Unlock()
				renderFrame()
			}
		}
	}()

	return nil
}
