package tui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mike-keough/pipelinepal/internal/db"
)

type Model struct {
	repo *db.Repo
	ctx  context.Context

	w, h int

	view View
	keys keyMap
	s    styles

	pipe PipelineState
	dtl  LeadDetailState

	// new lead form
	newLead newLeadForm

	// add note form
	addNote addNoteForm

	leads leadsState
	tasks tasksState

	addTask addTaskForm

	pending pendingSelection

	status string
	err    error
}

func New(repo *db.Repo) Model {
	m := Model{
		repo:    repo,
		ctx:     context.Background(),
		view:    ViewPipeline,
		keys:    keys(),
		s:       makeStyles(),
		newLead: newNewLeadForm(),
		addNote: newAddNoteForm(),
		leads:   newLeadsState(),
		tasks:   newTasksState(),
		addTask: newAddTaskForm(),
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.cmdLoadPipeline(),
		m.cmdLoadLeads(""),
		m.cmdLoadTasks(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		// keep lists sized
		if m.leads.list.Width() == 0 {
			m.leads.list.SetSize(m.w-4, m.h-10)
		} else {
			m.leads.list.SetSize(m.w-4, m.h-10)
		}
		if m.tasks.list.Width() == 0 {
			m.tasks.list.SetSize(m.w-4, m.h-8)
		} else {
			m.tasks.list.SetSize(m.w-4, m.h-8)
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case statusMsg:
		m.status = string(msg)
		return m, nil

	case pipelineLoadedMsg:
		m.pipe.Stages = msg.stages
		m.pipe.ByStage = msg.byStage

		// If we just moved a lead, keep selection on that same lead in the new stage.
		if m.pending.active {
			// set stage index to the stage we moved into
			for i, st := range m.pipe.Stages {
				if st.ID == m.pending.stageID {
					m.pipe.StageIndex = i
					break
				}
			}

			// set lead index to the moved lead inside that stage
			leads := m.pipe.ByStage[m.pending.stageID]
			for i, ld := range leads {
				if ld.ID == m.pending.leadID {
					m.pipe.LeadIndex = i
					break
				}
			}

			m.pending.active = false
		}

		// normal clamps
		m.pipe.StageIndex = clamp(m.pipe.StageIndex, 0, len(m.pipe.Stages)-1)
		m.pipe.LeadIndex = clamp(m.pipe.LeadIndex, 0, m.maxLeadIndexInStage())
		m.err = nil
		return m, nil

	case leadDetailLoadedMsg:
		m.dtl = msg.detail
		m.dtl.TaskIndex = clamp(m.dtl.TaskIndex, 0, len(m.dtl.Tasks)-1)
		m.view = ViewLeadDetail
		m.err = nil
		return m, nil

	case leadsLoadedMsg:
		m.leads.items = msg.leads
		items := make([]list.Item, 0, len(msg.leads))
		for _, l := range msg.leads {
			items = append(items, leadItem(l))
		}
		m.leads.list.SetItems(items)
		m.leads.loaded = true
		return m, nil

	case tasksLoadedMsg:
		m.tasks.items = msg.tasks
		items := make([]list.Item, 0, len(msg.tasks))
		for _, t := range msg.tasks {
			items = append(items, taskItem{T: t})
		}
		m.tasks.list.SetItems(items)
		m.tasks.loaded = true
		return m, nil

	case tea.KeyMsg:

		// Global keys (ONLY when not typing)
		if !m.isTyping() {
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.TasksView):
				m.view = ViewTasks
				return m, m.cmdLoadTasks()

			case key.Matches(msg, m.keys.Help):
				if m.view == ViewHelp {
					m.view = ViewPipeline
					return m, nil
				}
				m.view = ViewHelp
				return m, nil

			case key.Matches(msg, m.keys.Tab):
				if m.view == ViewPipeline {
					m.view = ViewLeads
					m.leads.search.Focus()
					return m, m.cmdLoadLeads(strings.TrimSpace(m.leads.search.Value()))
				}
				if m.view == ViewLeads {
					m.leads.search.Blur()
					m.view = ViewPipeline
					return m, m.cmdLoadPipeline()
				}
			}
		}

		// View-specific handling (always allowed)
		switch m.view {
		case ViewPipeline:
			return m.updatePipeline(msg)
		case ViewLeads:
			return m.updateLeads(msg)
		case ViewLeadDetail:
			return m.updateLeadDetail(msg)
		case ViewNewLead:
			return m.updateNewLead(msg)
		case ViewTasks:
			return m.updateTasks(msg)
		case ViewHelp:
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.w == 0 {
		return "Loading…"
	}

	header := m.s.Header.Render("PipelinePal") + "  " + m.s.Subtle.Render("tab: leads • t: tasks • n: new lead • q: quit • ?: help")

	var body string
	switch m.view {
	case ViewPipeline:
		body = m.viewPipeline()
	case ViewTasks:
		body = m.viewTasks()

	case ViewLeads:
		body = m.viewLeads()
	case ViewLeadDetail:
		body = m.viewLeadDetail()
	case ViewNewLead:
		body = m.viewNewLead()
	case ViewHelp:
		body = m.viewHelp()
	}

	foot := ""
	if m.err != nil {
		foot = "\n" + m.s.Error.Render("Error: "+m.err.Error())
	} else if m.status != "" {
		foot = "\n" + m.s.Subtle.Render(m.status)
	}

	return m.s.App.Render(header + "\n\n" + body + foot)
}

func (m Model) isTyping() bool {
	return m.newLead.name.Focused() ||
		m.newLead.phone.Focused() ||
		m.newLead.email.Focused() ||
		m.newLead.ltype.Focused() ||
		m.newLead.source.Focused() ||
		m.addNote.active ||
		m.addTask.active ||
		m.leads.search.Focused()
}

// ---------- Commands + messages ----------

type errMsg struct{ err error }
type statusMsg string

type pipelineLoadedMsg struct {
	stages  []db.Stage
	byStage map[int64][]db.Lead
}

type leadDetailLoadedMsg struct {
	detail LeadDetailState
}

type leadsLoadedMsg struct {
	leads []db.Lead
}

type tasksLoadedMsg struct {
	tasks []db.Task
}

type pendingSelection struct {
	leadID  int64
	stageID int64
	active  bool
}

func (m Model) cmdLoadPipeline() tea.Cmd {
	return func() tea.Msg {
		stages, err := m.repo.ListStages(m.ctx)
		if err != nil {
			return errMsg{err}
		}
		byStage, err := m.repo.ListLeadsByStage(m.ctx)
		if err != nil {
			return errMsg{err}
		}
		return pipelineLoadedMsg{stages: stages, byStage: byStage}
	}
}

func (m Model) cmdLoadTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.repo.ListOpenTasks(m.ctx)
		if err != nil {
			return errMsg{err}
		}
		return tasksLoadedMsg{tasks: tasks}
	}
}

func (m Model) cmdLoadLeadDetail(id int64) tea.Cmd {
	return func() tea.Msg {
		lead, err := m.repo.GetLead(m.ctx, id)
		if err != nil {
			return errMsg{err}
		}
		tasks, err := m.repo.ListTasksForLead(m.ctx, id)
		if err != nil {
			return errMsg{err}
		}
		notes, err := m.repo.ListNotes(m.ctx, id)
		if err != nil {
			return errMsg{err}
		}
		return leadDetailLoadedMsg{detail: LeadDetailState{
			LeadID: id,
			Lead:   lead,
			Tasks:  tasks,
			Notes:  notes,
		}}
	}
}

func (m Model) cmdLoadLeads(q string) tea.Cmd {
	return func() tea.Msg {
		leads, err := m.repo.ListLeads(m.ctx, q)
		if err != nil {
			return errMsg{err}
		}
		return leadsLoadedMsg{leads: leads}
	}
}
