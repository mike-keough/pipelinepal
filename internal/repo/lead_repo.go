package repo

import (
	"database/sql"
	"time"

	"github.com/mike-keough/pipelinepal/internal/models"
)

type LeadRepo struct {
	db *sql.DB
}

func NewLeadRepo(db *sql.DB) *LeadRepo {
	return &LeadRepo{db: db}
}

func (r *LeadRepo) Add(l *models.Lead) (int64, error) {
	now := time.Now().UTC()
	l.CreatedAt = now

	res, err := r.db.Exec(`
		INSERT INTO leads (name, phone, email, source, kind, status, notes, created_at, next_follow_up_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		l.Name, l.Phone, l.Email, l.Source, l.Kind, l.Status, l.Notes,
		l.CreatedAt.Format(time.RFC3339),
		nullTime(l.NextFollowUpAt),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *LeadRepo) List(kind string) ([]models.Lead, error) {
	rows, err := r.db.Query(`
		SELECT id, name, phone, email, source, kind, status, notes, created_at, last_contact_at, next_follow_up_at
		FROM leads
		WHERE (? = '' OR kind = ?)
		ORDER BY datetime(created_at) DESC
	`, kind, kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Lead
	for rows.Next() {
		var l models.Lead
		var created, last, next sql.NullString
		if err := rows.Scan(&l.ID, &l.Name, &l.Phone, &l.Email, &l.Source, &l.Kind, &l.Status, &l.Notes, &created, &last, &next); err != nil {
			return nil, err
		}
		if created.Valid {
			t, _ := time.Parse(time.RFC3339, created.String)
			l.CreatedAt = t
		}
		if last.Valid {
			t, _ := time.Parse(time.RFC3339, last.String)
			l.LastContactAt = &t
		}
		if next.Valid {
			t, _ := time.Parse(time.RFC3339, next.String)
			l.NextFollowUpAt = &t
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func nullTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}
