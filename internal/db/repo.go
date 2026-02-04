package db

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type Repo struct {
	db *DB
}

func NewRepo(db *DB) *Repo { return &Repo{db: db} }

// -------- Stages --------

func (r *Repo) ListStages(ctx context.Context) ([]Stage, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, sort, color FROM stages ORDER BY sort ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Stage
	for rows.Next() {
		var s Stage
		if err := rows.Scan(&s.ID, &s.Name, &s.Sort, &s.Color); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// -------- Leads --------

func (r *Repo) ListLeads(ctx context.Context, q string) ([]Lead, error) {
	q = strings.TrimSpace(q)
	var rows *sql.Rows
	var err error

	if q == "" {
		rows, err = r.db.QueryContext(ctx, `
SELECT l.id, l.full_name, l.phone, l.email, l.lead_type, l.source,
       l.stage_id, s.name,
       l.created_at, l.updated_at, l.last_contacted
FROM leads l
JOIN stages s ON s.id = l.stage_id
ORDER BY l.updated_at DESC, l.id DESC
`)
	} else {
		like := "%" + q + "%"
		rows, err = r.db.QueryContext(ctx, `
SELECT l.id, l.full_name, l.phone, l.email, l.lead_type, l.source,
       l.stage_id, s.name,
       l.created_at, l.updated_at, l.last_contacted
FROM leads l
JOIN stages s ON s.id = l.stage_id
WHERE l.full_name LIKE ? OR l.phone LIKE ? OR l.email LIKE ? OR l.source LIKE ?
ORDER BY l.updated_at DESC, l.id DESC
`, like, like, like, like)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Lead
	for rows.Next() {
		var l Lead
		var created, updated string
		var last sql.NullString
		if err := rows.Scan(
			&l.ID, &l.FullName, &l.Phone, &l.Email, &l.LeadType, &l.Source,
			&l.StageID, &l.StageName,
			&created, &updated, &last,
		); err != nil {
			return nil, err
		}
		l.CreatedAt = mustParseTime(created)
		l.UpdatedAt = mustParseTime(updated)
		if last.Valid && last.String != "" {
			t := mustParseTime(last.String)
			l.LastContacted = &t
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *Repo) ListLeadsByStage(ctx context.Context) (map[int64][]Lead, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT l.id, l.full_name, l.phone, l.email, l.lead_type, l.source,
       l.stage_id, s.name,
       l.created_at, l.updated_at, l.last_contacted
FROM leads l
JOIN stages s ON s.id = l.stage_id
ORDER BY s.sort ASC, l.updated_at DESC, l.id DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[int64][]Lead)
	for rows.Next() {
		var l Lead
		var created, updated string
		var last sql.NullString
		if err := rows.Scan(
			&l.ID, &l.FullName, &l.Phone, &l.Email, &l.LeadType, &l.Source,
			&l.StageID, &l.StageName,
			&created, &updated, &last,
		); err != nil {
			return nil, err
		}
		l.CreatedAt = mustParseTime(created)
		l.UpdatedAt = mustParseTime(updated)
		if last.Valid && last.String != "" {
			t := mustParseTime(last.String)
			l.LastContacted = &t
		}
		out[l.StageID] = append(out[l.StageID], l)
	}
	return out, rows.Err()
}

func (r *Repo) CreateLead(ctx context.Context, fullName, phone, email, leadType, source string, stageID int64) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO leads(full_name, phone, email, lead_type, source, stage_id)
VALUES (?, ?, ?, ?, ?, ?)
`, fullName, phone, email, leadType, source, stageID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repo) MoveLeadStage(ctx context.Context, leadID, newStageID int64) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE leads
SET stage_id = ?, updated_at = datetime('now')
WHERE id = ?
`, newStageID, leadID)
	return err
}

func (r *Repo) GetLead(ctx context.Context, id int64) (Lead, error) {
	var l Lead
	var created, updated string
	var last sql.NullString
	err := r.db.QueryRowContext(ctx, `
SELECT l.id, l.full_name, l.phone, l.email, l.lead_type, l.source,
       l.stage_id, s.name,
       l.created_at, l.updated_at, l.last_contacted
FROM leads l
JOIN stages s ON s.id = l.stage_id
WHERE l.id = ?
`, id).Scan(
		&l.ID, &l.FullName, &l.Phone, &l.Email, &l.LeadType, &l.Source,
		&l.StageID, &l.StageName,
		&created, &updated, &last,
	)
	if err != nil {
		return Lead{}, err
	}
	l.CreatedAt = mustParseTime(created)
	l.UpdatedAt = mustParseTime(updated)
	if last.Valid && last.String != "" {
		t := mustParseTime(last.String)
		l.LastContacted = &t
	}
	return l, nil
}

// -------- Notes --------

func (r *Repo) ListNotes(ctx context.Context, leadID int64) ([]Note, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, lead_id, body, created_at
FROM notes
WHERE lead_id = ?
ORDER BY created_at DESC, id DESC
`, leadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Note
	for rows.Next() {
		var n Note
		var created string
		if err := rows.Scan(&n.ID, &n.LeadID, &n.Body, &created); err != nil {
			return nil, err
		}
		n.CreatedAt = mustParseTime(created)
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *Repo) AddNote(ctx context.Context, leadID int64, body string) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO notes(lead_id, body) VALUES (?, ?)
`, leadID, body)
	if err != nil {
		return 0, err
	}
	_, _ = r.db.ExecContext(ctx, `UPDATE leads SET updated_at = datetime('now') WHERE id = ?`, leadID)
	return res.LastInsertId()
}

// -------- Tasks --------

func (r *Repo) ListOpenTasks(ctx context.Context) ([]Task, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT t.id, t.lead_id, l.full_name, t.title, t.due_date, t.status, t.created_at, t.completed_at
FROM tasks t
JOIN leads l ON l.id = t.lead_id
WHERE t.status = 'open'
ORDER BY
  CASE WHEN t.due_date IS NULL OR t.due_date = '' THEN 1 ELSE 0 END,
  t.due_date ASC,
  t.id DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		var due sql.NullString
		var created string
		var completed sql.NullString
		if err := rows.Scan(&t.ID, &t.LeadID, &t.LeadName, &t.Title, &due, &t.Status, &created, &completed); err != nil {
			return nil, err
		}
		t.CreatedAt = mustParseTime(created)
		if due.Valid && due.String != "" {
			dd := mustParseDate(due.String)
			t.DueDate = &dd
		}
		if completed.Valid && completed.String != "" {
			ct := mustParseTime(completed.String)
			t.CompletedAt = &ct
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) ListTasksForLead(ctx context.Context, leadID int64) ([]Task, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT t.id, t.lead_id, '' as lead_name, t.title, t.due_date, t.status, t.created_at, t.completed_at
FROM tasks t
WHERE t.lead_id = ?
ORDER BY
  CASE WHEN t.status = 'open' THEN 0 ELSE 1 END,
  CASE WHEN t.due_date IS NULL OR t.due_date = '' THEN 1 ELSE 0 END,
  t.due_date ASC,
  t.id DESC
`, leadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		var due sql.NullString
		var created string
		var completed sql.NullString
		if err := rows.Scan(&t.ID, &t.LeadID, &t.LeadName, &t.Title, &due, &t.Status, &created, &completed); err != nil {
			return nil, err
		}
		t.CreatedAt = mustParseTime(created)
		if due.Valid && due.String != "" {
			dd := mustParseDate(due.String)
			t.DueDate = &dd
		}
		if completed.Valid && completed.String != "" {
			ct := mustParseTime(completed.String)
			t.CompletedAt = &ct
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) CreateTask(ctx context.Context, leadID int64, title string, due *time.Time) (int64, error) {
	var dueStr any = nil
	if due != nil {
		dueStr = due.Format("2006-01-02")
	}
	res, err := r.db.ExecContext(ctx, `
INSERT INTO tasks(lead_id, title, due_date)
VALUES (?, ?, ?)
`, leadID, title, dueStr)
	if err != nil {
		return 0, err
	}
	_, _ = r.db.ExecContext(ctx, `UPDATE leads SET updated_at = datetime('now') WHERE id = ?`, leadID)
	return res.LastInsertId()
}

func (r *Repo) CompleteTask(ctx context.Context, taskID int64) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE tasks
SET status = 'done',
    completed_at = datetime('now')
WHERE id = ?
`, taskID)
	return err
}

// -------- Helpers --------

func mustParseTime(s string) time.Time {
	// SQLite datetime('now') => "YYYY-MM-DD HH:MM:SS"
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		return t
	}
	// fallback: RFC3339
	t2, err2 := time.Parse(time.RFC3339, s)
	if err2 == nil {
		return t2
	}
	return time.Time{}
}

func mustParseDate(s string) time.Time {
	// YYYY-MM-DD
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		return t
	}
	return time.Time{}
}
