package main

import (
	"strconv"
	"strings"
)

type FlashlightApp struct {
	selectedMode int
}

func printLight(char string, color string, frame []string) string {
	var light strings.Builder
	for i := 0; i < 23; /*3 lines needed for mode display*/ i++ {
		if i != 0 {
			light.WriteString("\n")
		}
		light.WriteString(color)
		light.WriteString(strings.Repeat(char, 54))
		light.WriteString(Reset)
	}
	return light.String()
}

func (a *FlashlightApp) Name() string { return "Flashlight" }

func (a *FlashlightApp) Start() error {
	a.selectedMode = 2
	return nil
}

func (a *FlashlightApp) Stop() error { return nil }

func (a *FlashlightApp) HandleEvent(event Event) (Action, error) {
	switch event.Key {
	case KeyLeft:
		if a.selectedMode > 0 {
			a.selectedMode--
		}
	case KeyRight:
		if a.selectedMode < 2 {
			a.selectedMode++
		}
	}

	return Action{Kind: ActionNone}, nil
}

func (a *FlashlightApp) Render() error {
	var frame []string
	var char string
	var color string

	if a.selectedMode > 0 {
		char = "█"
		if a.selectedMode == 1 {
			color = Grey
		} else {
			color = White
		}
	} else {
		char = " "
		color = Reset
	}

	frame = append(frame, printLight(char, color, frame))

	frame = append(frame, Red+" ┌─────────┐"+Reset)

	var indicator strings.Builder
	for i := 0; i < 3; i++ {
		if i == 0 {
			indicator.WriteString(Red)
			indicator.WriteString("├┤")
			indicator.WriteString(Reset)
		}
		if a.selectedMode == i {
			indicator.WriteString(Red)
			indicator.WriteString("█")
			indicator.WriteString(Reset)
			indicator.WriteString(RedBg)
			indicator.WriteString(strconv.Itoa(i))
			indicator.WriteString(Red)
			indicator.WriteString("█")
			indicator.WriteString(Reset)
		} else {
			indicator.WriteString(Grey)
			indicator.WriteString("███")
			indicator.WriteString(Reset)
		}

		if i == 2 {
			indicator.WriteString(Red)
			indicator.WriteString("├┤")
			indicator.WriteString(Reset)
		}
	}
	frame = append(frame, indicator.String())

	frame = append(frame, Red+" └─────────┘"+Reset)

	presentFrame(frame)
	return nil
}
