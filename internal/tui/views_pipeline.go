package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mike-keough/pipelinepal/internal/db"
)

func (m Model) updatePipeline(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.NewLead):
		m.view = ViewNewLead
		m.newLead.reset()
		if len(m.pipe.Stages) > 0 {
			m.newLead.stageID = m.pipe.Stages[m.pipe.StageIndex].ID
		}
		return m, nil

	case key.Matches(msg, m.keys.Left):
		m.pipe.StageIndex = clamp(m.pipe.StageIndex-1, 0, len(m.pipe.Stages)-1)
		m.pipe.LeadIndex = clamp(m.pipe.LeadIndex, 0, m.maxLeadIndexInStage())
		return m, nil

	case key.Matches(msg, m.keys.Right):
		m.pipe.StageIndex = clamp(m.pipe.StageIndex+1, 0, len(m.pipe.Stages)-1)
		m.pipe.LeadIndex = clamp(m.pipe.LeadIndex, 0, m.maxLeadIndexInStage())
		return m, nil

	case key.Matches(msg, m.keys.Up):
		m.pipe.LeadIndex = clamp(m.pipe.LeadIndex-1, 0, m.maxLeadIndexInStage())
		return m, nil

	case key.Matches(msg, m.keys.Down):
		m.pipe.LeadIndex = clamp(m.pipe.LeadIndex+1, 0, m.maxLeadIndexInStage())
		return m, nil

	case key.Matches(msg, m.keys.Enter):
		lead, ok := m.selectedLead()
		if !ok {
			return m, nil
		}
		return m, m.cmdLoadLeadDetail(lead.ID)

	case key.Matches(msg, m.keys.MoveL):
		return m.moveSelectedLead(-1)

	case key.Matches(msg, m.keys.MoveR):
		return m.moveSelectedLead(+1)
	}

	return m, nil
}

func (m Model) viewPipeline() string {
	if len(m.pipe.Stages) == 0 {
		return "No stages found."
	}

	cols := make([]string, 0, len(m.pipe.Stages))
	colW := m.columnWidth()

	for i, st := range m.pipe.Stages {
		leads := m.pipe.ByStage[st.ID]

		title := m.s.ColTitle.Render(st.Name) + m.s.Badge.Render(fmt.Sprintf("%d", len(leads)))

		var cards []string
		for j, ld := range leads {
			line := fmtLeadLine(
				ellipsize(ld.FullName, colW-6),
				ld.LeadType,
				ellipsize(ld.Source, 12),
			)

			cardStyle := m.s.Card
			if i == m.pipe.StageIndex && j == m.pipe.LeadIndex {
				cardStyle = m.s.CardSel
			}
			cards = append(cards, cardStyle.Width(colW-4).Render(line))
		}

		if len(cards) == 0 {
			cards = []string{m.s.Subtle.Render("(empty)")}
		}

		colStyle := m.s.Col
		if i == m.pipe.StageIndex {
			colStyle = m.s.ColSel
		}

		col := colStyle.Width(colW).Render(
			title + "\n" + lipgloss.JoinVertical(lipgloss.Left, cards...),
		)
		cols = append(cols, col)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)

	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}

func (m Model) maxLeadIndexInStage() int {
	if len(m.pipe.Stages) == 0 {
		return 0
	}
	stage := m.pipe.Stages[m.pipe.StageIndex]
	leads := m.pipe.ByStage[stage.ID]
	if len(leads) == 0 {
		return 0
	}
	return len(leads) - 1
}

func (m Model) selectedLead() (db.Lead, bool) {
	if len(m.pipe.Stages) == 0 {
		return db.Lead{}, false
	}
	stage := m.pipe.Stages[m.pipe.StageIndex]
	leads := m.pipe.ByStage[stage.ID]
	if len(leads) == 0 {
		return db.Lead{}, false
	}
	if m.pipe.LeadIndex < 0 || m.pipe.LeadIndex >= len(leads) {
		return db.Lead{}, false
	}
	return leads[m.pipe.LeadIndex], true
}

func (m Model) moveSelectedLead(dir int) (tea.Model, tea.Cmd) {
	lead, ok := m.selectedLead()
	if !ok {
		return m, nil
	}
	newStageIdx := m.pipe.StageIndex + dir
	if newStageIdx < 0 || newStageIdx >= len(m.pipe.Stages) {
		return m, nil
	}
	newStageID := m.pipe.Stages[newStageIdx].ID

	m.pending = pendingSelection{
		leadID:  lead.ID,
		stageID: newStageID,
		active:  true,
	}

	cmd := func() tea.Msg {
		if err := m.repo.MoveLeadStage(m.ctx, lead.ID, newStageID); err != nil {
			return errMsg{err}
		}
		return statusMsg("Moved lead.")
	}
	return m, tea.Batch(cmd, m.cmdLoadPipeline())
}

func (m Model) columnWidth() int {
	// Reserve some breathing room for outer padding
	usable := m.w - 6
	if usable < 40 {
		return 24
	}

	n := len(m.pipe.Stages)
	if n <= 0 {
		return 26
	}

	// Try to fit all columns; clamp so it stays pretty
	w := usable / n
	if w < 22 {
		w = 22
	}
	if w > 36 {
		w = 36
	}
	return w
}
