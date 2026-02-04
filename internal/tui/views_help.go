package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) viewHelp() string {
	lines := []string{
		m.s.Header.Render("Help"),
		"",
		m.s.Header.Render("Pipeline view"),
		"- arrows / h j k l: navigate",
		"- enter: open lead",
		"- n: new lead",
		"- t: tasks",
		"- H / L: move lead left/right (between stages)",
		"",
		m.s.Header.Render("Lead detail"),
		"- a: add note",
		"- esc: back",
		"",
		m.s.Header.Render("Global"),
		"- t: tasks",
		"- tab: switch Pipeline/Leads",
		"- ?: help",
		"- q: quit",
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
