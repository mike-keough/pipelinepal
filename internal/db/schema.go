package db

import "time"

type Stage struct {
	ID    int64
	Name  string
	Sort  int
	Color string // not used yet, but future-friendly
}

type Lead struct {
	ID            int64
	FullName      string
	Phone         string
	Email         string
	LeadType      string // "buyer" | "seller" | "other"
	Source        string
	StageID       int64
	StageName     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastContacted *time.Time
}

type Note struct {
	ID        int64
	LeadID    int64
	Body      string
	CreatedAt time.Time
}

type Task struct {
	ID          int64
	LeadID      int64
	LeadName    string // convenient for “All Tasks” view
	Title       string
	DueDate     *time.Time
	Status      string // open|done
	CreatedAt   time.Time
	CompletedAt *time.Time
}
