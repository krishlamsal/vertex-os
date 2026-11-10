package main

type NotesApp struct {
}

func (a *NotesApp) Name() string { return "Notes" }

func (a *NotesApp) Start() error { return nil }

func (a *NotesApp) Stop() error { return nil }

func (a *NotesApp) HandleEvent(event Event) (Action, error) {
	switch event.Key {

	}

	return Action{Kind: ActionNone}, nil
}

func (a *NotesApp) Render() error {
	presentFrame([]string{"N"})
	return nil
}
