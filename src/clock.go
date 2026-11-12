package main

import (
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-runewidth"
)

type ClockApp struct {
	mu   sync.Mutex
	stop chan struct{}
}

func (a *ClockApp) Name() string { return "Clock" }

func (a *ClockApp) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.stop = make(chan struct{})
	return nil
}

func (a *ClockApp) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.stop != nil {
		close(a.stop)
		a.stop = nil
	}
	return nil
}

func (a *ClockApp) HandleEvent(event Event) (Action, error) { return Action{Kind: ActionNone}, nil }

func (a *ClockApp) Render() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	stopCh := a.stop
	if stopCh == nil {
		return nil
	}

	cities := []string{
		"America/New_York",
		"America/Los_Angeles",
		"America/Sao_Paulo",
		"Europe/London",
		"Europe/Paris",
		"Asia/Dubai",
		"Asia/Kolkata",
		"Asia/Singapore",
		"Asia/Tokyo",
		"Australia/Sydney",
	}

	friendlyCities := map[string]string{
		"America/New_York":    "New York",
		"America/Los_Angeles": "Los Angeles",
		"America/Sao_Paulo":   "São Paulo",
		"Europe/London":       "London",
		"Europe/Paris":        "Paris",
		"Asia/Dubai":          "Dubai",
		"Asia/Kolkata":        "Mumbai",
		"Asia/Singapore":      "Singapore",
		"Asia/Tokyo":          "Tokyo",
		"Australia/Sydney":    "Sydney",
	}

	renderFrame := func() {
		var frame []string
		currentTime := time.Now().Format("15:04:05")
		currentDate := time.Now().Format("2006-01-02")

		frame = append(frame, getTitle("Clock", Blue))

		frame = append(frame, strings.Repeat(" ", 12)+"┌─"+Cyan+currentDate+Reset+strings.Repeat("─", 32-len(currentDate))+"┐")

		for i := 0; i < 3; i++ {
			var row strings.Builder

			row.WriteString(strings.Repeat(" ", 12))
			row.WriteString("│ ")
			for j := 0; j < len(currentTime); j++ {
				row.WriteString(string(artRows(string(currentTime[j]), numberArt, 4)[i]))
			}
			row.WriteString("│")

			frame = append(frame, row.String())
		}

		frame = append(frame, strings.Repeat(" ", 12)+"└"+strings.Repeat("─", 33)+"┘")

		var cityList strings.Builder

		for _, city := range cities {
			// Load the location database for city
			loc, err := time.LoadLocation(city)
			if err != nil {
				cityList.WriteString(Red)
				cityList.WriteString("Error loading location")
				cityList.WriteString(Reset)
				cityList.WriteString("\n")
				continue
			}

			cityTime := time.Now().In(loc)

			cityList.WriteString(strings.Repeat(" ", 10))
			cityList.WriteString(Red)
			cityList.WriteString(friendlyCities[city])
			cityList.WriteString(Reset)
			cityList.WriteString(strings.Repeat(" ", 11-runewidth.StringWidth(friendlyCities[city])))
			cityList.WriteString("  ·  ")
			cityList.WriteString(cityTime.Format("2006-01-02"))
			cityList.WriteString(Blue)
			cityList.WriteString(cityTime.Format(" 15:04"))
			cityList.WriteString(Reset)
			cityList.WriteString(cityTime.Format(" (MST)\n"))
		}

		frame = append(frame, cityList.String())

		presentFrame(frame)
	}

	renderFrame()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				renderFrame()
			}
		}
	}()

	return nil
}
