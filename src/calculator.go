package main

import (
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

type CalculatorApp struct {
	selectedIndexX   int
	selectedIndexY   int
	currentOperation string
	currentTotal     float64
	currentTotal2    float64
	total2Active     bool
	inputBuffer      string
	boxWidth         int
	boxHeight        int
}

func (a *CalculatorApp) boxLine(selected bool, line int, label string) string {
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
		return accent + "│" + accentBg + strings.Repeat(" ", (a.boxWidth-2)/2) + White + label + strings.Repeat(" ", (a.boxWidth-2)/2) + reset + accent + "│ " + reset
	default:
		return accent + "│" + accentBg + strings.Repeat(" ", a.boxWidth-2) + reset + accent + "│ " + reset
	}
}

func (a *CalculatorApp) activeValue() *float64 {
	if a.total2Active {
		return &a.currentTotal2
	}
	return &a.currentTotal
}

func (a *CalculatorApp) commitBuffer() error {
	if a.inputBuffer == "" {
		return nil
	}

	value, err := strconv.ParseFloat(a.inputBuffer, 64)
	if err != nil {
		return err
	}

	*a.activeValue() = value
	return nil
}

func decideOperation(a *CalculatorApp, x float64) error {
	a.inputBuffer += strconv.FormatFloat(x, 'f', -1, 64)
	return a.commitBuffer()
}

func enterDecimal(a *CalculatorApp) error {
	if strings.Contains(a.inputBuffer, ".") {
		return nil
	}
	if a.inputBuffer == "" {
		a.inputBuffer = "0"
	}

	a.inputBuffer += "."
	return a.commitBuffer()
}

func backspace(a *CalculatorApp) error {
	if a.inputBuffer == "" {
		a.inputBuffer = strconv.FormatFloat(*a.activeValue(), 'f', -1, 64)
	}

	if len(a.inputBuffer) <= 1 {
		a.inputBuffer = ""
		*a.activeValue() = 0
		return nil
	}

	a.inputBuffer = a.inputBuffer[:len(a.inputBuffer)-1]
	if a.inputBuffer == "" || a.inputBuffer == "-" {
		a.inputBuffer = ""
		*a.activeValue() = 0
		return nil
	}

	return a.commitBuffer()
}

func negate(a *CalculatorApp) error {
	if a.inputBuffer == "" {
		a.inputBuffer = strconv.FormatFloat(*a.activeValue(), 'f', -1, 64)
	}

	if strings.HasPrefix(a.inputBuffer, "-") {
		a.inputBuffer = a.inputBuffer[1:]
	} else {
		a.inputBuffer = "-" + a.inputBuffer
	}

	if a.inputBuffer == "" || a.inputBuffer == "-" {
		a.inputBuffer = ""
		*a.activeValue() = 0
		return nil
	}

	return a.commitBuffer()
}

func percent(a *CalculatorApp) error {
	*a.activeValue() /= 100
	a.inputBuffer = strconv.FormatFloat(*a.activeValue(), 'f', -1, 64)
	return nil
}

func applyOperation(a *CalculatorApp) {
	switch a.currentOperation {
	case "+":
		a.currentTotal += a.currentTotal2
	case "-":
		a.currentTotal -= a.currentTotal2
	case "*":
		a.currentTotal *= a.currentTotal2
	case "/":
		a.currentTotal /= a.currentTotal2
	}

	a.currentOperation = ""
	a.currentTotal2 = 0
}

func setOperation(a *CalculatorApp, operation string) {
	if a.total2Active && a.currentOperation != "" {
		applyOperation(a)
	}

	a.currentOperation = operation
	a.total2Active = true
	a.inputBuffer = ""
}

func (a *CalculatorApp) Name() string { return "Calculator" }

func (a *CalculatorApp) Start() error {
	a.boxWidth = 13
	a.boxHeight = 5
	return nil
}

func (a *CalculatorApp) Stop() error { return nil }

func (a *CalculatorApp) HandleEvent(event Event) (Action, error) {
	var err error = nil

	switch event.Key {
	case KeyUp:
		if a.selectedIndexY > 0 {
			a.selectedIndexY--
		}
	case KeyDown:
		if a.selectedIndexY < 4 {
			a.selectedIndexY++
		}
	case KeyLeft:
		if a.selectedIndexX > 0 {
			a.selectedIndexX--
		}
	case KeyRight:
		if a.selectedIndexX < 3 {
			a.selectedIndexX++
		}
	case KeyOK:
		switch a.selectedIndexY {
		case 0:
			switch a.selectedIndexX {
			case 0:
				a.currentTotal = 0
				a.currentTotal2 = 0
				a.total2Active = false
				a.currentOperation = ""
				a.inputBuffer = ""
			case 1:
				err = backspace(a)
			case 2:
				err = negate(a)
			case 3:
				setOperation(a, "/")
			}
		case 1:
			switch a.selectedIndexX {
			case 0:
				decideOperation(a, 7)
			case 1:
				decideOperation(a, 8)
			case 2:
				decideOperation(a, 9)
			case 3:
				setOperation(a, "*")
			}
		case 2:
			switch a.selectedIndexX {
			case 0:
				decideOperation(a, 4)
			case 1:
				decideOperation(a, 5)
			case 2:
				decideOperation(a, 6)
			case 3:
				setOperation(a, "-")
			}
		case 3:
			switch a.selectedIndexX {
			case 0:
				decideOperation(a, 1)
			case 1:
				decideOperation(a, 2)
			case 2:
				decideOperation(a, 3)
			case 3:
				setOperation(a, "+")
			}
		case 4:
			switch a.selectedIndexX {
			case 0:
				decideOperation(a, 0)
			case 1:
				err = percent(a)
			case 2:
				err = enterDecimal(a)
			case 3:
				applyOperation(a)
				a.total2Active = false
				a.inputBuffer = ""
			}
		}
	}

	return Action{Kind: ActionNone}, err
}

func (a *CalculatorApp) Render() error {
	// Layout:
	// C ⌫  ±  ÷
	// 7  8  9  ×
	// 4  5  6  -
	// 1  2  3  +
	// 0  %  .  =

	lines := strings.Repeat("─", 54)
	currentTotalStr := strconv.FormatFloat(a.currentTotal, 'f', -1, 64)
	currentTotalStr2 := strconv.FormatFloat(a.currentTotal2, 'f', -1, 64)

	var frame []string
	frame = append(frame, getTitle("Calculator", Red))

	display := a.inputBuffer
	if display == "" {
		if a.total2Active {
			display = currentTotalStr2
		} else {
			display = currentTotalStr
		}
	}

	padding := 54 - runewidth.StringWidth(display)
	if padding < 0 {
		padding = 0
	}
	frame = append(frame, Blue+strings.Repeat(" ", padding)+display+Reset)
	frame = append(frame, Grey+lines+Reset)

	labels := [20]string{"C", "⌫", "±", "÷", "7", "8", "9", "×", "4", "5", "6", "-", "1", "2", "3", "+", "0", "%", ".", "="}

	for y := 0; y < 5; y++ {
		for line := 0; line < a.boxHeight; line++ {
			var row strings.Builder
			for x := 0; x < 4; x++ {
				selected := x == a.selectedIndexX && y == a.selectedIndexY
				row.WriteString(a.boxLine(selected, line, labels[x+y*4]))
			}
			frame = append(frame, row.String())
		}
	}

	presentFrame(frame)

	return nil
}
