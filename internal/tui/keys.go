package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit key.Binding

	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding

	Enter key.Binding
	Back  key.Binding

	NewLead key.Binding
	MoveL   key.Binding
	MoveR   key.Binding

	Notes key.Binding
	Help  key.Binding

	Tab key.Binding

	TasksView key.Binding
	FollowUp  key.Binding
	Complete  key.Binding
}

func keys() keyMap {
	return keyMap{
		Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),

		Left:  key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
		Right: key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
		Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),

		Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open/select")),
		Back:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),

		NewLead: key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new lead")),
		MoveL:   key.NewBinding(key.WithKeys("H"), key.WithHelp("H", "move lead left")),
		MoveR:   key.NewBinding(key.WithKeys("L"), key.WithHelp("L", "move lead right")),

		Notes:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add note")),
		Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Tab:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch view")),
		TasksView: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "tasks")),
		FollowUp:  key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "new follow-up")),
		Complete:  key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "complete task")),
	}
}
