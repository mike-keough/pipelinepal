package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mike-keough/pipelinepal/internal/db"
)

type newLeadForm struct {
	step    int
	stageID int64

	name   textinput.Model
	phone  textinput.Model
	email  textinput.Model
	ltype  textinput.Model
	source textinput.Model

	// reused by Leads view as a quick store (keeps things simple for now)
	leads []db.Lead // overwritten in code? we keep separate below
}

func newNewLeadForm() newLeadForm {
	mk := func(ph string, w int) textinput.Model {
		ti := textinput.New()
		ti.Placeholder = ph
		ti.Width = w
		return ti
	}

	f := newLeadForm{
		name:   mk("Full name", 40),
		phone:  mk("Phone", 24),
		email:  mk("Email", 40),
		ltype:  mk("Lead type: buyer/seller/other", 30),
		source: mk("Source: Zillow, referral, sign call…", 40),
	}
	f.ltype.SetValue("buyer")
	return f
}

func (f *newLeadForm) reset() {
	f.step = 0
	f.name.SetValue("")
	f.phone.SetValue("")
	f.email.SetValue("")
	f.ltype.SetValue("buyer")
	f.source.SetValue("")
	f.name.Focus()
}

func (m Model) updateNewLead(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	fields := []*textinput.Model{&m.newLead.name, &m.newLead.phone, &m.newLead.email, &m.newLead.ltype, &m.newLead.source}

	// escape cancels
	if key.Matches(msg, m.keys.Back) {

		m.view = ViewPipeline
		return m, nil
	}

	switch msg.String() {
	case "enter":
		// next step or save
		if m.newLead.step < len(fields)-1 {
			fields[m.newLead.step].Blur()
			m.newLead.step++
			fields[m.newLead.step].Focus()
			return m, nil
		}

		// save
		fullName := strings.TrimSpace(m.newLead.name.Value())
		if fullName == "" {
			m.err = errString("name is required")
			return m, nil
		}

		phone := strings.TrimSpace(m.newLead.phone.Value())
		email := strings.TrimSpace(m.newLead.email.Value())
		leadType := strings.ToLower(strings.TrimSpace(m.newLead.ltype.Value()))
		if leadType == "" {
			leadType = "buyer"
		}
		source := strings.TrimSpace(m.newLead.source.Value())

		stageID := m.newLead.stageID
		if stageID == 0 && len(m.pipe.Stages) > 0 {
			stageID = m.pipe.Stages[0].ID
		}

		cmd := func() tea.Msg {
			if _, err := m.repo.CreateLead(m.ctx, fullName, phone, email, leadType, source, stageID); err != nil {
				return errMsg{err}
			}
			return statusMsg("Lead created.")
		}

		m.view = ViewPipeline
		return m, tea.Batch(cmd, m.cmdLoadPipeline())
	}

	// text input update
	var c tea.Cmd
	cur := fields[m.newLead.step]
	*cur, c = cur.Update(msg)
	return m, c
}

type errString string

func (e errString) Error() string { return string(e) }

func (m Model) viewNewLead() string {
	lines := []string{
		m.s.Header.Render("New Lead"),
		"",
		m.s.Subtle.Render("enter: next/save • esc: cancel"),
		"",
	}

	labels := []string{"Name", "Phone", "Email", "Type", "Source"}
	fields := []string{
		m.newLead.name.View(),
		m.newLead.phone.View(),
		m.newLead.email.View(),
		m.newLead.ltype.View(),
		m.newLead.source.View(),
	}

	for i := range fields {
		box := m.s.Border.Render(fields[i])
		if m.newLead.step == i {
			box = m.s.BorderFocus.Render(fields[i])
		}

		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(8).Render(labels[i]+":"),
			box,
		))
		lines = append(lines, "")
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
