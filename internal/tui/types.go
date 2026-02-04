package tui

import "github.com/mike-keough/pipelinepal/internal/db"

type View int

const (
	ViewPipeline View = iota
	ViewLeads
	ViewLeadDetail
	ViewNewLead
	ViewTasks
	ViewHelp
)

type PipelineState struct {
	Stages     []db.Stage
	ByStage    map[int64][]db.Lead
	StageIndex int
	LeadIndex  int
}

type LeadDetailState struct {
	LeadID    int64
	Lead      db.Lead
	Tasks     []db.Task
	TaskIndex int
	Notes     []db.Note
}
