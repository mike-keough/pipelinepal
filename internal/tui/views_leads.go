package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mike-keough/pipelinepal/internal/db"
)

type leadItem db.Lead

func (i leadItem) Title() string { return string(i.FullName) }
func (i leadItem) Description() string {
	return fmt.Sprintf("%s • %s", strings.ToUpper(i.LeadType), i.StageName)
}
func (i leadItem) FilterValue() string {
	return fmt.Sprintf("%s %s %s %s", i.FullName, i.Phone, i.Email, i.Source)
}

type leadsState struct {
	search textinput.Model
	list   list.Model
	items  []db.Lead
	loaded bool
}

func newLeadsState() leadsState {
	ti := textinput.New()
	ti.Placeholder = "Search leads (name, phone, email, source)…"
	ti.Width = 52

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Title = "Leads"

	return leadsState{search: ti, list: l}
}

func (m Model) updateLeads(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.leads.search.Blur()
		m.view = ViewPipeline
		return m, m.cmdLoadPipeline()
	}

	// Jump into search
	if msg.String() == "/" {
		m.leads.search.Focus()
		return m, nil
	}

	// If search is focused, typing should go there (not the list)
	if m.leads.search.Focused() {
		switch msg.String() {
		case "enter":
			q := strings.TrimSpace(m.leads.search.Value())
			m.leads.search.Blur() // after applying search, go to list navigation
			return m, m.cmdLoadLeads(q)
		case "esc":
			m.leads.search.Blur()
			return m, nil
		}

		var cmd tea.Cmd
		m.leads.search, cmd = m.leads.search.Update(msg)
		return m, cmd
	}

	// List navigation mode
	switch {
	case key.Matches(msg, m.keys.Enter):
		if it, ok := m.leads.list.SelectedItem().(leadItem); ok {
			ld := db.Lead(it)
			return m, m.cmdLoadLeadDetail(ld.ID)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.leads.list, cmd = m.leads.list.Update(msg)
	return m, cmd
}

func (m Model) viewLeads() string {
	if !m.leads.loaded {
		return "Loading leads…"
	}

	// Focused search box looks different
	box := m.s.Border.Render(m.leads.search.View())
	if m.leads.search.Focused() {
		box = m.s.BorderFocus.Render(m.leads.search.View())
	}

	top := lipgloss.JoinVertical(lipgloss.Left,
		m.s.Header.Render("Leads"),
		m.s.Subtle.Render("Type / to search • Enter applies search • Enter on a lead opens details • esc back"),
		"",
		box,
		"",
	)

	return top + m.leads.list.View()
}
