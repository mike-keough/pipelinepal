package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mike-keough/pipelinepal/internal/app"
)

func main() {
	ctx := context.Background()

	// Store DB in a local app folder (works on Linux/macOS/Windows)
	base := defaultDataDir()
	if err := os.MkdirAll(base, 0o755); err != nil {
		fmt.Println("failed creating data dir:", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(base, "pipelinepal.sqlite")

	a, err := app.New(dbPath)
	if err != nil {
		fmt.Println("init error:", err)
		os.Exit(1)
	}
	defer a.Close()

	if err := a.Bootstrap(ctx); err != nil {
		fmt.Println("bootstrap error:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(a.Model(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("run error:", err)
		os.Exit(1)
	}
}

func defaultDataDir() string {
	// Simple + predictable: ./data if present, else user home
	if _, err := os.Stat("data"); err == nil {
		return "data"
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "data"
	}
	return filepath.Join(home, ".pipelinepal")
}
