package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/kwrkb/asql/internal/db"
	"github.com/kwrkb/asql/internal/db/dbutil"
)

const minWidthForCompare = 80

// focusedTableStyles returns table styles with the focused (selected row) highlight.
func focusedTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(panelBorder).
		BorderBottom(true).
		Bold(true).
		Foreground(accentColor)
	s.Selected = s.Selected.
		Foreground(textColor).
		Background(lipgloss.Color("#1D4ED8")).
		Bold(false)
	s.Cell = s.Cell.Foreground(textColor)
	return s
}

// unfocusedTableStyles returns table styles without the selected row highlight.
func unfocusedTableStyles() table.Styles {
	s := focusedTableStyles()
	s.Selected = s.Selected.
		Background(lipgloss.Color("#334155")).
		Foreground(mutedTextColor)
	return s
}

// pinnedPane holds a snapshot of a query result for side-by-side comparison.
type pinnedPane struct {
	result      db.QueryResult
	connName    string
	table       table.Model
	displayRows []table.Row
	colWidths   []int
	colCursor   int
	colOffset   int
	sortCol     int
	sortDir     sortOrder
	lastVisStart  int
	lastVisEnd    int
	viewportDirty bool
}

// pinCurrentResult creates a pinnedPane from the current active result.
func (m *model) pinCurrentResult() *pinnedPane {
	if len(m.lastResult.Columns) == 0 {
		return nil
	}

	// Clone the table model — starts unfocused since right pane gets initial focus
	tbl := table.New(
		table.WithColumns([]table.Column{{Title: "Result", Width: 30}}),
		table.WithRows([]table.Row{{"pinned"}}),
		table.WithFocused(false),
		table.WithHeight(m.table.Height()),
	)
	tbl.SetStyles(unfocusedTableStyles())

	// Copy displayRows
	rows := make([]table.Row, len(m.displayRows))
	copy(rows, m.displayRows)

	// Copy colWidths
	widths := make([]int, len(m.cachedColWidths))
	copy(widths, m.cachedColWidths)

	return &pinnedPane{
		result:      m.lastResult,
		connName:    m.connMgr.ActiveName(),
		table:       tbl,
		displayRows: rows,
		colWidths:   widths,
		colCursor:   m.colCursor,
		colOffset:   m.colOffset,
		sortCol:     m.sortCol,
		sortDir:     m.sortDir,
		viewportDirty: true,
	}
}

// comparePaneWidth returns the width for each pane in side-by-side mode.
func (m *model) comparePaneWidth() int {
	return m.fullContentWidth() / 2
}

// pinnedVisibleColumnRange returns the visible column range for the pinned pane.
func (p *pinnedPane) visibleColumnRange(availWidth int) (int, int) {
	if len(p.colWidths) == 0 {
		return 0, 0
	}
	available := availWidth - 8
	start := p.colOffset
	if start >= len(p.colWidths) {
		start = 0
	}
	sum := 0
	for i := start; i < len(p.colWidths); i++ {
		w := p.colWidths[i] + 1
		if sum+w > available && i > start {
			return start, i
		}
		sum += w
	}
	return start, len(p.colWidths)
}

// adjustColOffset ensures colCursor is within the visible window for pinned pane.
func (p *pinnedPane) adjustColOffset(availWidth int) {
	if p.colCursor < p.colOffset {
		p.colOffset = p.colCursor
	}
	_, visEnd := p.visibleColumnRange(availWidth)
	for p.colCursor >= visEnd && p.colOffset < len(p.colWidths)-1 {
		p.colOffset++
		_, visEnd = p.visibleColumnRange(availWidth)
	}
	p.viewportDirty = true
}

// togglePinnedSort toggles sort on the pinned pane.
func (m *model) togglePinnedSort() {
	p := m.pinned
	if p.colCursor == p.sortCol {
		switch p.sortDir {
		case sortNone:
			p.sortDir = sortAsc
		case sortAsc:
			p.sortDir = sortDesc
		case sortDesc:
			p.sortDir = sortNone
		}
	} else {
		p.sortCol = p.colCursor
		p.sortDir = sortAsc
	}
	// Re-sort and rebuild displayRows
	sorted := sortedRows(p.result.Rows, p.sortCol, p.sortDir)
	p.displayRows = make([]table.Row, 0, len(sorted))
	for _, row := range sorted {
		p.displayRows = append(p.displayRows, table.Row(row))
	}
	if len(p.displayRows) == 0 {
		sentinel := make(table.Row, len(p.result.Columns))
		sentinel[0] = "(no rows)"
		p.displayRows = []table.Row{sentinel}
	}
	p.table.GotoTop()
	p.viewportDirty = true
}

// syncPinnedTable rebuilds the pinned pane's table for the given width/height.
func (m *model) syncPinnedTable(paneWidth, paneHeight int) {
	p := m.pinned
	if p == nil || len(p.result.Columns) == 0 {
		return
	}

	p.adjustColOffset(paneWidth)
	visStart, visEnd := p.visibleColumnRange(paneWidth)

	rebuildNeeded := visStart != p.lastVisStart || visEnd != p.lastVisEnd || p.viewportDirty
	if !rebuildNeeded {
		return
	}

	selectedStyle := lipgloss.NewStyle().Reverse(true)
	columns := make([]table.Column, 0, visEnd-visStart)
	for i := visStart; i < visEnd; i++ {
		header := sanitize(p.result.Columns[i])
		if i < len(p.result.ColumnTypes) && p.result.ColumnTypes[i] != "" {
			shortType := dbutil.ShortenTypeName(sanitize(p.result.ColumnTypes[i]))
			header = header + " " + typeStyle.Render(shortType)
		}
		if i == p.sortCol && p.sortDir != sortNone {
			header += sortIndicator(p.sortDir)
		}
		if m.comparePane == 0 && i == p.colCursor {
			header = selectedStyle.Render(header)
		}
		columns = append(columns, table.Column{Title: header, Width: p.colWidths[i]})
	}

	rows := make([]table.Row, 0, len(p.displayRows))
	for _, row := range p.displayRows {
		windowed := make(table.Row, 0, visEnd-visStart)
		for i := visStart; i < visEnd; i++ {
			if i < len(row) {
				windowed = append(windowed, sanitize(row[i]))
			} else {
				windowed = append(windowed, "")
			}
		}
		rows = append(rows, windowed)
	}

	cursor := p.table.Cursor()
	p.table.SetRows([]table.Row{})
	p.table.SetColumns(columns)
	p.table.SetRows(rows)
	p.table.SetCursor(cursor)
	p.table.SetHeight(max(paneHeight-4, 3))
	p.lastVisStart = visStart
	p.lastVisEnd = visEnd
	p.viewportDirty = false
}

// renderCompareView renders side-by-side panels.
func (m *model) renderCompareView() string {
	paneWidth := m.comparePaneWidth()
	paneHeight := m.resultsHeight() - 1 // subtract 1 for label row

	m.syncPinnedTable(paneWidth, paneHeight)

	focusedBorder := lipgloss.Color("#38BDF8") // accentColor
	unfocusedBorder := panelBorder

	// Left pane (pinned)
	leftBorderColor := unfocusedBorder
	if m.comparePane == 0 {
		leftBorderColor = focusedBorder
	}
	leftLabel := sanitize(m.pinned.connName)
	if leftLabel == "" {
		leftLabel = "pinned"
	}
	leftPanel := lipgloss.NewStyle().
		Width(max(paneWidth, 0)).
		Height(max(paneHeight, 0)).
		Background(panelBackground).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(leftBorderColor).
		Padding(0, 1).
		Render(m.pinned.table.View())

	// Right pane (active)
	rightBorderColor := unfocusedBorder
	if m.comparePane == 1 {
		rightBorderColor = focusedBorder
	}
	rightPanel := lipgloss.NewStyle().
		Width(max(paneWidth, 0)).
		Height(max(paneHeight, 0)).
		Background(panelBackground).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(rightBorderColor).
		Padding(0, 1).
		Render(m.table.View())

	// Label above each pane
	leftLabelStr := lipgloss.NewStyle().
		Width(paneWidth).
		Foreground(accentColor).
		Background(appBackground).
		Align(lipgloss.Center).
		Render("[" + leftLabel + "]")
	rightLabel := sanitize(m.connMgr.ActiveName())
	if rightLabel == "" {
		rightLabel = "active"
	}
	rightLabelStr := lipgloss.NewStyle().
		Width(paneWidth).
		Foreground(accentColor).
		Background(appBackground).
		Align(lipgloss.Center).
		Render("[" + rightLabel + "]")

	labels := lipgloss.JoinHorizontal(lipgloss.Top, leftLabelStr, rightLabelStr)
	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	return lipgloss.JoinVertical(lipgloss.Left, labels, panels)
}
