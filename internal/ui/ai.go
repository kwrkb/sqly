package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kwrkb/asql/internal/ai"
	"github.com/kwrkb/asql/internal/db"
)

const aiRequestTimeout = 30 * time.Second

func (m model) updateAI(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.aiLoading {
		if msg.Type == tea.KeyEsc {
			if m.queryCancel != nil {
				m.queryCancel()
				m.queryCancel = nil
			}
			m.aiLoading = false
			m.mode = normalMode
			m.setStatus("Cancelled", false)
			return m, nil
		}
		return m, nil
	}
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = normalMode
		m.aiError = ""
		m.setStatus("Normal mode", false)
		return m, nil
	case tea.KeyEnter:
		prompt := strings.TrimSpace(m.aiInput.Value())
		if prompt == "" {
			return m, nil
		}
		if m.queryCancel != nil {
			m.queryCancel()
		}
		ctx, cancel := context.WithCancel(context.Background())
		m.querySeq++
		m.queryCancel = cancel
		m.aiLoading = true
		m.aiError = ""
		return m, tea.Batch(m.aiSpinner.Tick, generateSQLCmd(ctx, m.aiClient, m.db, prompt, m.querySeq))
	}
	var cmd tea.Cmd
	m.aiInput, cmd = m.aiInput.Update(msg)
	return m, cmd
}

func generateSQLCmd(parent context.Context, client *ai.Client, adapter db.DBAdapter, prompt string, seq uint64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(parent, aiRequestTimeout)
		defer cancel()

		schema, err := adapter.Schema(ctx)
		if err != nil {
			return aiResponseMsg{seq: seq, err: fmt.Errorf("fetching schema: %w", err)}
		}

		sql, err := client.GenerateSQL(ctx, adapter.Type(), schema, prompt)
		return aiResponseMsg{seq: seq, sql: sql, err: err}
	}
}

func (m model) renderWithAIOverlay(background string) string {
	modalWidth := min(m.width-4, 60)
	if modalWidth < 20 {
		modalWidth = 20
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		MarginBottom(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(1, 2).
		Width(modalWidth).
		Background(panelBackground)

	var content string
	if m.aiLoading {
		content = titleStyle.Render("AI Generating SQL...") + "\n" + m.aiSpinner.View() + " Thinking..."
	} else {
		content = titleStyle.Render("AI Assistant (Text-to-SQL)") + "\n" + m.aiInput.View()
		if m.aiError != "" {
			errStyle := lipgloss.NewStyle().Foreground(errorColor).MarginTop(1)
			content += "\n" + errStyle.Render(m.aiError)
		}
	}

	modal := boxStyle.Render(content)

	bgH := lipgloss.Height(background)

	return lipgloss.Place(m.width, bgH, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceBackground(appBackground))
}
