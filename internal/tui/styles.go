package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	App    lipgloss.Style
	Header lipgloss.Style
	Subtle lipgloss.Style

	Border      lipgloss.Style
	BorderFocus lipgloss.Style

	Col      lipgloss.Style
	ColSel   lipgloss.Style
	ColTitle lipgloss.Style

	Card    lipgloss.Style
	CardSel lipgloss.Style

	Badge lipgloss.Style
	Error lipgloss.Style
}

func makeStyles() styles {
	return styles{
		// Root container (NO background here)
		App: lipgloss.NewStyle().
			Padding(1, 2),

		// Headers (royal blue)
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#1E3A8A", Dark: "#93C5FD"}),

		// Subtle text
		Subtle: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}),

		// Generic bordered panel
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#CBD5E1", Dark: "#334155"}).
			Background(lipgloss.AdaptiveColor{Light: "#F3F4F6", Dark: "#0F172A"}).
			Padding(0, 1),

		// Focused input border (royal blue)
		BorderFocus: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#1E3A8A", Dark: "#93C5FD"}).
			Background(lipgloss.AdaptiveColor{Light: "#F3F4F6", Dark: "#0F172A"}).
			Padding(0, 1),

		// Pipeline columns
		Col: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#CBD5E1", Dark: "#334155"}).
			Background(lipgloss.AdaptiveColor{Light: "#F3F4F6", Dark: "#0F172A"}).
			Padding(0, 1),

		// Selected column
		ColSel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#1E3A8A", Dark: "#93C5FD"}).
			Background(lipgloss.AdaptiveColor{Light: "#F3F4F6", Dark: "#0F172A"}).
			Padding(0, 1),

		// Column title
		ColTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#1E3A8A", Dark: "#93C5FD"}).
			MarginBottom(1),

		// Lead cards
		Card: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#CBD5E1", Dark: "#334155"}).
			Background(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#020617"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#374151", Dark: "#D1D5DB"}).
			Padding(0, 1).
			MarginBottom(1),

		// Selected lead card
		CardSel: lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#1E3A8A", Dark: "#93C5FD"}).
			Background(lipgloss.AdaptiveColor{Light: "#DBEAFE", Dark: "#1E293B"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#374151", Dark: "#D1D5DB"}).
			Padding(0, 1).
			MarginBottom(1).
			Bold(true),

		// Count badges
		Badge: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1E3A8A", Dark: "#93C5FD"}).
			Background(lipgloss.AdaptiveColor{Light: "#DBEAFE", Dark: "#1E293B"}).
			Padding(0, 1).
			MarginLeft(1).
			Bold(true),

		// Error messages
		Error: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#F87171"}),
	}
}
