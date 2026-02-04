package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mike-keough/pipelinepal/internal/db"
)

type taskItem struct {
	T db.Task
}

func (i taskItem) Title() string {
	t := i.T
	due := "No due date"
	if t.DueDate != nil && !t.DueDate.IsZero() {
		due = t.DueDate.Format("2006-01-02")
	}
	return fmt.Sprintf("%s • %s", t.Title, due)
}

func (i taskItem) Description() string {
	t := i.T
	return fmt.Sprintf("Lead: %s", t.LeadName)
}

func (i taskItem) FilterValue() string { return "" }

type tasksState struct {
	list   list.Model
	items  []db.Task
	loaded bool
}

func newTasksState() tasksState {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Title = "Open Tasks"
	return tasksState{list: l}
}

func (m Model) updateTasks(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.view = ViewPipeline
		return m, m.cmdLoadPipeline()

	case key.Matches(msg, m.keys.Enter):
		if it, ok := m.tasks.list.SelectedItem().(taskItem); ok {
			return m, m.cmdLoadLeadDetail(it.T.LeadID)
		}
		return m, nil

	case key.Matches(msg, m.keys.Complete):
		if it, ok := m.tasks.list.SelectedItem().(taskItem); ok {
			cmd := func() tea.Msg {
				if err := m.repo.CompleteTask(m.ctx, it.T.ID); err != nil {
					return errMsg{err}
				}
				return statusMsg("Task completed.")
			}
			return m, tea.Batch(cmd, m.cmdLoadTasks())
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.tasks.list, cmd = m.tasks.list.Update(msg)
	return m, cmd
}

func (m Model) viewTasks() string {
	if !m.tasks.loaded {
		return "Loading tasks…"
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		m.s.Header.Render("Open Tasks"),
		m.s.Subtle.Render("enter: open lead • c: complete • esc: back"),
		"",
	)

	return header + m.tasks.list.View()
}
