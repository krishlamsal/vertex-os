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

func remove[T any](array []T, index int) []T {
	return append(array[:index], array[index+1:]...)
}

func weatherGraph(weatherReport string) string {
	graph := strings.Trim(weatherReport, "\n")

	splitGraph := strings.Split(graph, "\n")

	for i := 0; i < 5; i++ {
		splitGraph = remove(splitGraph, len(splitGraph)-1) // Info
	}

	splitGraph = remove(splitGraph, 35) // Lines under moon phase
	splitGraph = remove(splitGraph, 34) // Moon phases
	splitGraph = remove(splitGraph, 33) // Blank line
	splitGraph = remove(splitGraph, 30) // Blank line
	splitGraph = remove(splitGraph, 28) // Blank line

	splitGraph = remove(splitGraph, 5) // Blank line
	splitGraph = remove(splitGraph, 4) // Separator lines
	splitGraph = remove(splitGraph, 2) // Blank line
	splitGraph = remove(splitGraph, 1) // Blank line
	splitGraph = remove(splitGraph, 0) // Weather Report: Location

	// Cap off top
	r := []rune(splitGraph[0])
	r[0] = '┬'
	r[len(r)-1] = '┬'
	splitGraph[0] = string(r)

	graph = strings.Join(splitGraph, "\n")

	return graph
}

func weatherInfo(weatherReport string) string {
	splitReport := strings.Split(weatherReport, "\n")
	return strings.Join(splitReport[len(splitReport)-6:], "\n"+Grey+strings.Repeat("─", 54)+Reset+"\n")
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
