package app

import (
	"context"

	"github.com/mike-keough/pipelinepal/internal/db"
	"github.com/mike-keough/pipelinepal/internal/tui"
)

type App struct {
	DB   *db.DB
	Repo *db.Repo
}

func New(dbPath string) (*App, error) {
	d, err := db.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &App{
		DB:   d,
		Repo: db.NewRepo(d),
	}, nil
}

func (a *App) Close() error { return a.DB.Close() }

func (a *App) Bootstrap(ctx context.Context) error {
	if err := a.DB.Migrate(ctx); err != nil {
		return err
	}
	return nil
}

func (a *App) Model() tui.Model {
	return tui.New(a.Repo)
}
