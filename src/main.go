package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

const (
	Reset = "\033[0m"

	Red   = "\033[91m"
	RedBg = "\033[101m"

	Blue   = "\033[94m"
	BlueBg = "\033[104m"

	Purple   = "\033[95m"
	PurpleBg = "\033[105m"

	Cyan   = "\033[96m"
	CyanBg = "\033[106m"

	Yellow   = "\033[93m"
	YellowBg = "\033[103m"

	Green   = "\033[92m"
	GreenBg = "\033[102m"

	White = "\033[37m"

	Grey   = "\033[90m"
	GreyBg = "\033[100m"

	Bold = "\033[1m"
)

var Colors []string = []string{Red, RedBg, Blue, BlueBg, Purple, PurpleBg, Cyan, CyanBg, Yellow, YellowBg, Green, GreenBg}

type Key int

const (
	KeyUnknown Key = iota

	KeyUp
	KeyDown
	KeyLeft
	KeyRight

	KeyPageLeft
	KeyPageRight

	KeyOK
	KeyBack
	KeyHome
)

type Event struct {
	Key Key
}

type ActionKind int

const (
	ActionNone ActionKind = iota
	ActionLaunch
	ActionHome
)

type Action struct {
	Kind    ActionKind
	AppName string
}

type App interface {
	Name() string
	Start() error
	HandleEvent(Event) (Action, error)
	Stop() error
	Render() error
}

type Manager struct {
	apps   map[string]App
	home   App
	active App
}

func NewManager(home App) *Manager {
	m := &Manager{
		apps: make(map[string]App),
		home: home,
	}
	m.Register(home)
	return m
}

func (m *Manager) Register(app App) {
	m.apps[app.Name()] = app
}

func (m *Manager) Start() error {
	m.active = m.home
	if err := m.active.Start(); err != nil {
		return err
	}
	return m.active.Render()
}

func (m *Manager) Dispatch(event Event) error {
	if m.active == nil {
		return fmt.Errorf("manager has not been started")
	}

	if event.Key == KeyHome {
		if err := m.GoHome(); err != nil {
			return err
		}
		return m.active.Render()
	}

	action, err := m.active.HandleEvent(event)
	if err != nil {
		return err
	}

	switch action.Kind {
	case ActionNone:
		// No app switch; still redraw in case the app's state changed.
	case ActionHome:
		err = m.GoHome()
	case ActionLaunch:
		err = m.Launch(action.AppName)
	default:
		return fmt.Errorf("unknown action: %d", action.Kind)
	}

	if err != nil {
		return err
	}

	return m.active.Render()
}

func (m *Manager) Launch(name string) error {
	next, ok := m.apps[name]
	if !ok {
		return fmt.Errorf("app not found: %q", name)
	}
	if next == m.active {
		return nil
	}

	if err := m.active.Stop(); err != nil {
		return err
	}
	if err := next.Start(); err != nil {
		// Try to restore the previous app if starting the new one fails.
		_ = m.active.Start()
		return err
	}

	m.active = next
	return nil
}

func (m *Manager) GoHome() error {
	if m.active == m.home {
		return nil // Home cannot be exited.
	}

	if err := m.active.Stop(); err != nil {
		return err
	}
	if err := m.home.Start(); err != nil {
		return err
	}

	m.active = m.home
	return nil
}

func presentFrame(lines []string) {
	var frame strings.Builder
	frame.Grow(2048)

	// Clear the whole screen before drawing so that a frame smaller than the
	// previous app's frame does not leave remnants of it behind. Reset first so
	// the erase uses the default background rather than a leftover colour.
	frame.WriteString(Reset)
	frame.WriteString("\033[2J")
	frame.WriteString("\033[H")

	for i, line := range lines {
		if i > 0 {
			frame.WriteString("\r\n")
		}
		frame.WriteString(line)
		frame.WriteString("\033[K")
	}

	frame.WriteString("\033[J")
	os.Stdout.WriteString(frame.String())
}

func eventFromRune(char rune) (Event, error) {
	var key Key

	switch char {
	// Up/Down/Left/Right
	case '2':
		key = KeyUp
	case '4':
		key = KeyLeft
	case '6':
		key = KeyRight
	case '8':
		key = KeyDown

	// Paging
	case '7':
		key = KeyPageLeft
	case '9':
		key = KeyPageRight

	// OK
	case '5':
		key = KeyOK

	// System
	case '#':
		key = KeyBack
	case '*':
		key = KeyHome

	// Unknown
	default:
		key = KeyUnknown
	}

	return Event{Key: key}, nil
}

func getTitle(titleText string, color string) string {
	linesHalf := strings.Repeat("─", (54-len(titleText)-2)/2)
	return color + linesHalf + " " + titleText + " " + linesHalf + Reset
}

func main() {
	home := &HomeApp{appNames: []string{"Calculator", "Notes", "Calendar", "Camera", "Clock", "Stopwatch", "Timer", "Weather", "Browser", "Phone", "Phonebook", "Messages", "Files", "Text Editor", "Image Viewer", "Audio Player", "Settings", "Help", "System Monitor", "Terminal"},
		appLogos:   []string{"\uf1ec", "\uf4f6", "\uf455", "\uf030", "\uf43a", "\uf520", "\U000f051f", "\U000f0590", "\U000f059f", "\uf095", "\uee21", "\uf27a", "\uf07b", "\uf15c", "\uf03e", "\U000f04c3", "\uf013", "\uf420", "\uf21e", "\uf489"},
		pageTitles: []string{"Utilities", "Time and Weather", "Communication", "File Management", "Advanced"}, maxPage: 4}
	manager := NewManager(home)

	// Utilities
	notes := &NotesApp{}
	manager.Register(notes)
	calculator := &CalculatorApp{}
	manager.Register(calculator)

	// Time and Weather
	clock := &ClockApp{}
	manager.Register(clock)
	stopwatch := &StopwatchApp{}
	manager.Register(stopwatch)
	timer := &TimerApp{}
	manager.Register(timer)

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	if err := manager.Start(); err != nil {
		panic(err)
	}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		char, _, err := keyboard.GetKey()
		if err != nil {
			continue
		}
		event, err := eventFromRune(char)
		if err != nil || event.Key == KeyUnknown {
			continue
		}

		if err := manager.Dispatch(event); err != nil {
			panic(err)
		}
	}
}
