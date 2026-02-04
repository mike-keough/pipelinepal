package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type addNoteForm struct {
	active bool
	input  textinput.Model
}

func newAddNoteForm() addNoteForm {
	ti := textinput.New()
	ti.Placeholder = "Write a note and press enter…"
	ti.CharLimit = 500
	ti.Width = 60
	return addNoteForm{input: ti}
}

func (f *addNoteForm) open() {
	f.active = true
	f.input.SetValue("")
	f.input.Focus()
}
func (f *addNoteForm) close() {
	f.active = false
	f.input.Blur()
}

type addTaskForm struct {
	active bool
	step   int // 0=title, 1=due
	title  textinput.Model
	due    textinput.Model
}

func newAddTaskForm() addTaskForm {
	t := textinput.New()
	t.Placeholder = "Follow-up task (e.g. Call about pre-approval)…"
	t.Width = 60

	d := textinput.New()
	d.Placeholder = "Due date (YYYY-MM-DD) optional"
	d.Width = 30

	return addTaskForm{title: t, due: d}
}

func (f *addTaskForm) open() {
	f.active = true
	f.step = 0
	f.title.SetValue("")
	f.due.SetValue("")
	f.title.Focus()
	f.due.Blur()
}
func (f *addTaskForm) close() {
	f.active = false
	f.title.Blur()
	f.due.Blur()
}

func parseOptionalDue(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (m Model) updateLeadDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// add note mode
	if m.addNote.active {
		switch msg.String() {
		case "esc":
			m.addNote.close()
			return m, nil
		case "enter":
			body := strings.TrimSpace(m.addNote.input.Value())
			if body == "" {
				m.addNote.close()
				return m, nil
			}
			leadID := m.dtl.LeadID
			m.addNote.close()

			cmd := func() tea.Msg {
				if _, err := m.repo.AddNote(m.ctx, leadID, body); err != nil {
					return errMsg{err}
				}
				return statusMsg("Note added.")
			}
			return m, tea.Batch(cmd, m.cmdLoadLeadDetail(leadID), m.cmdLoadPipeline(), m.cmdLoadTasks())
		}

		var c tea.Cmd
		m.addNote.input, c = m.addNote.input.Update(msg)
		return m, c
	}

	// add task mode
	if m.addTask.active {
		switch msg.String() {
		case "esc":
			m.addTask.close()
			return m, nil
		case "enter":
			if m.addTask.step == 0 {
				m.addTask.title.Blur()
				m.addTask.step = 1
				m.addTask.due.Focus()
				return m, nil
			}

			title := strings.TrimSpace(m.addTask.title.Value())
			if title == "" {
				m.addTask.close()
				return m, nil
			}

			due, err := parseOptionalDue(m.addTask.due.Value())
			if err != nil {
				m.err = err
				return m, nil
			}

			leadID := m.dtl.LeadID
			m.addTask.close()

			cmd := func() tea.Msg {
				if _, err := m.repo.CreateTask(m.ctx, leadID, title, due); err != nil {
					return errMsg{err}
				}
				return statusMsg("Follow-up created.")
			}
			return m, tea.Batch(cmd, m.cmdLoadLeadDetail(leadID), m.cmdLoadPipeline(), m.cmdLoadTasks())
		}

		var c tea.Cmd
		if m.addTask.step == 0 {
			m.addTask.title, c = m.addTask.title.Update(msg)
		} else {
			m.addTask.due, c = m.addTask.due.Update(msg)
		}
		return m, c
	}

	// normal mode
	switch {
	case key.Matches(msg, m.keys.Back):
		m.view = ViewPipeline
		return m, m.cmdLoadPipeline()

	case key.Matches(msg, m.keys.Notes):
		m.addNote.open()
		return m, nil

	case key.Matches(msg, m.keys.FollowUp):
		m.addTask.open()
		return m, nil

	case key.Matches(msg, m.keys.Up):
		if len(m.dtl.Tasks) > 0 {
			m.dtl.TaskIndex = clamp(m.dtl.TaskIndex-1, 0, len(m.dtl.Tasks)-1)
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if len(m.dtl.Tasks) > 0 {
			m.dtl.TaskIndex = clamp(m.dtl.TaskIndex+1, 0, len(m.dtl.Tasks)-1)
		}
		return m, nil

	case key.Matches(msg, m.keys.Complete):
		if len(m.dtl.Tasks) == 0 {
			return m, nil
		}
		t := m.dtl.Tasks[m.dtl.TaskIndex]
		cmd := func() tea.Msg {
			if err := m.repo.CompleteTask(m.ctx, t.ID); err != nil {
				return errMsg{err}
			}
			return statusMsg("Task completed.")
		}
		return m, tea.Batch(cmd, m.cmdLoadLeadDetail(m.dtl.LeadID), m.cmdLoadTasks())
	}

	return m, nil
}

func (m Model) viewLeadDetail() string {
	l := m.dtl.Lead

	lines := []string{
		m.s.Header.Render("Lead Detail"),
		"",
		fmt.Sprintf("%s %s", m.s.Badge.Render(strings.ToUpper(l.LeadType)), m.s.Header.Render(l.FullName)),
		m.s.Subtle.Render(fmt.Sprintf("Stage: %s • Source: %s", l.StageName, emptyDash(l.Source))),
		m.s.Subtle.Render(fmt.Sprintf("Phone: %s • Email: %s", emptyDash(l.Phone), emptyDash(l.Email))),
		m.s.Subtle.Render(fmt.Sprintf("Updated: %s", l.UpdatedAt.Format("2006-01-02 15:04"))),
		"",
		m.s.Header.Render("Follow-ups (tasks)"),
		m.s.Subtle.Render("f: new follow-up • c: complete selected • j/k: select"),
		"",
	}

	if m.addTask.active {
		lines = append(lines, m.s.BorderFocus.Render(m.addTask.title.View()))
		lines = append(lines, m.s.BorderFocus.Render(m.addTask.due.View()))
		lines = append(lines, m.s.Subtle.Render("enter: next/save • esc: cancel"))
		lines = append(lines, "")
	}

	if len(m.dtl.Tasks) == 0 {
		lines = append(lines, m.s.Subtle.Render("(no follow-ups yet)"))
	} else {
		for i, t := range m.dtl.Tasks {
			due := "—"
			if t.DueDate != nil && !t.DueDate.IsZero() {
				due = t.DueDate.Format("2006-01-02")
			}
			row := fmt.Sprintf("%s  [%s]  %s", ellipsize(t.Title, 52), due, strings.ToUpper(t.Status))
			if i == m.dtl.TaskIndex {
				lines = append(lines, m.s.CardSel.Render(row))
			} else {
				lines = append(lines, m.s.Card.Render(row))
			}
		}
	}

	lines = append(lines, "", m.s.Header.Render("Notes"))

	if m.addNote.active {
		lines = append(lines, m.s.BorderFocus.Render(m.addNote.input.View()))
		lines = append(lines, m.s.Subtle.Render("enter: save • esc: cancel"))
		lines = append(lines, "")
	}

	if len(m.dtl.Notes) == 0 {
		lines = append(lines, m.s.Subtle.Render("(no notes yet)"))
	} else {
		for _, n := range m.dtl.Notes {
			lines = append(lines,
				m.s.Border.Render(fmt.Sprintf("%s\n%s",
					m.s.Subtle.Render(n.CreatedAt.Format("2006-01-02 15:04")),
					ellipsize(n.Body, 400),
				)),
			)
		}
	}

	lines = append(lines, "", m.s.Subtle.Render("a: add note • esc: back • q: quit"))
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func emptyDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}
