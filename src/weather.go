package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type WeatherApp struct {
	view string
}

type location struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	City    string `json:"city"`
	Region  string `json:"regionName"`
	Country string `json:"country"`
}

func isWeatherInfoLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	for _, prefix := range []string{"Weather:", "Timezone:", "Now", "Zenith:", "Location:"} {
		if strings.Contains(trimmed, prefix) {
			return true
		}
	}
	return false
}

func isMoonLine(line string) bool {
	for _, r := range line {
		switch r {
		case '🌑', '🌒', '🌓', '🌔', '🌕', '🌖', '🌗', '🌘':
			return true
		}
	}
	return false
}

func isMoonUnderlineLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	runes := []rune(trimmed)
	if len(runes) < 2 || runes[0] != '│' || runes[len(runes)-1] != '│' {
		return false
	}
	hasDash := false
	for _, r := range runes[1 : len(runes)-1] {
		if r == '─' {
			hasDash = true
			continue
		}
		if r != ' ' {
			return false
		}
	}
	return hasDash
}

func isBottomBorderLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	runes := []rune(trimmed)
	if len(runes) < 2 || runes[0] != '└' || runes[len(runes)-1] != '┘' {
		return false
	}
	for _, r := range runes[1 : len(runes)-1] {
		if r != '─' {
			return false
		}
	}
	return true
}

func isBorderOnlyLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 2 {
		return false
	}
	runes := []rune(trimmed)
	if runes[0] != '│' || runes[len(runes)-1] != '│' {
		return false
	}
	for _, r := range runes[1 : len(runes)-1] {
		if r != ' ' {
			return false
		}
	}
	return true
}

func isTopSeparatorLine(line string) bool {
	if !strings.Contains(line, "╷") {
		return false
	}
	trimmed := strings.TrimSpace(line)
	runes := []rune(trimmed)
	if len(runes) < 2 || runes[0] != '│' || runes[len(runes)-1] != '│' {
		return false
	}
	for _, r := range runes[1 : len(runes)-1] {
		if r != ' ' && r != '╷' {
			return false
		}
	}
	return true
}

func weatherGraph(weatherReport string) string {
	lines := strings.Split(strings.Trim(weatherReport, "\n"), "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if strings.Contains(line, "Weather report:") {
			continue
		}
		if isWeatherInfoLine(line) {
			continue
		}
		if isTopSeparatorLine(line) {
			continue
		}
		if isBorderOnlyLine(line) {
			continue
		}
		if isMoonLine(line) || isMoonUnderlineLine(line) {
			continue
		}
		filtered = append(filtered, line)
	}

	if len(filtered) == 0 {
		return ""
	}

	r := []rune(filtered[0])
	if len(r) >= 2 {
		r[0] = '┬'
		r[len(r)-1] = '┬'
		filtered[0] = string(r)
	}

	return strings.Join(filtered, "\n")
}

func weatherInfo(weatherReport string) string {
	lines := strings.Split(strings.Trim(weatherReport, "\n"), "\n")
	info := make([]string, 0, 6)
	for _, line := range lines {
		if isWeatherInfoLine(line) {
			info = append(info, line)
		}
	}
	if len(info) == 0 {
		if len(lines) > 5 {
			info = lines[len(lines)-5:]
		} else {
			info = lines
		}
	}
	return strings.Join(info, "\n"+Grey+strings.Repeat("─", 54)+Reset+"\n")
}

func (a *WeatherApp) Name() string { return "Weather" }

func (a *WeatherApp) Start() error {
	a.view = "Graph"
	return nil
}

func (a *WeatherApp) Stop() error { return nil }

func (a *WeatherApp) HandleEvent(event Event) (Action, error) {
	switch event.Key {
	case KeyPageLeft:
		a.view = "Graph"
	case KeyPageRight:
		a.view = "Info "
	}

	return Action{Kind: ActionNone}, nil
}

func (a *WeatherApp) Render() error {
	resp, err := http.Get("http://ip-api.com/json/?fields=status,message,city,regionName,country")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var loc location
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return err
	}
	if loc.Status != "success" {
		return fmt.Errorf("location lookup failed: %s", loc.Message)
	}

	var frame []string

	pageDots := "1  •"
	if a.view == "Info " {
		pageDots = "•  2"
	}

	frame = append(frame, getTitle("Weather  ·  "+pageDots+"  ·  "+loc.City, Blue))

	req, err := http.NewRequest("GET",
		"https://v2d.wttr.in/"+url.PathEscape(loc.City), nil)
	req.Header.Set("User-Agent", "curl")

	wttrResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer wttrResp.Body.Close()

	wttrBytes, err := io.ReadAll(wttrResp.Body)
	if err != nil {
		return err
	}
	wttrStr := string(wttrBytes)

	if a.view == "Graph" {
		frame = append(frame, weatherGraph(wttrStr))
	} else {
		frame = append(frame, weatherInfo(wttrStr))
	}

	presentFrame(frame)
	return nil
}
