package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "i":
		m.mode = insertMode
		m.textarea.Focus()
		m.setStatus("Insert mode", false)
	case "t":
		if m.width >= minWidthForSidebar {
			m.sidebarOpen = true
			m.mode = sidebarMode
			m.textarea.Blur()
			m.sidebarCursor = 0
			m.setStatus("Sidebar", false)
			m.resize()
		} else {
			m.setStatus("Terminal too narrow for sidebar", true)
		}
	case "e":
		if len(m.lastResult.Columns) == 0 {
			m.setStatus("No query results to export", true)
		} else {
			m.mode = exportMode
			m.exportCursor = 0
			m.setStatus("Export mode", false)
		}
	case "ctrl+k":
		if m.aiEnabled {
			m.mode = aiMode
			m.aiInput.Reset()
			m.aiInput.Focus()
			m.aiError = ""
			m.aiLoading = false
			m.setStatus("AI mode", false)
			return m, textinput.Blink
		}
		m.setStatus("AI not configured", true)
	case "j":
		m.table.MoveDown(1)
	case "k":
		m.table.MoveUp(1)
	case "h", "left":
		if m.colCursor > 0 {
			m.colCursor--
		}
	case "l", "right":
		if len(m.lastResult.Columns) > 0 && m.colCursor < len(m.lastResult.Columns)-1 {
			m.colCursor++
		}
	case "s":
		if len(m.lastResult.Columns) > 0 {
			m.toggleSort()
		}
	}
	m.syncViewport()
	return m, nil
}
