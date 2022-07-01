package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) moveHighlightUp() {
	m.rowCursorIndex--
	var wrapped bool
	var lastRowIndex int
	if m.rowCursorIndex < 0 {
		lastRowIndex = m.TotalRows() - 1
		m.rowCursorIndex = lastRowIndex
		wrapped = true
	}

	if m.pageSize > 0 {
		m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	}

	currentWindowHeight := min(m.maxVisibleRows, m.TotalRows())
	if m.maxTotalHeight > 0 {
		m.verticalScrollWindowYOffset--
		if m.verticalScrollWindowYOffset < 0 {
			if wrapped {
				m.verticalScrollWindowStart = lastRowIndex - currentWindowHeight + 1
				m.verticalScrollWindowEnd = lastRowIndex
				m.verticalScrollWindowYOffset = currentWindowHeight - 1
			} else {
				m.verticalScrollWindowStart--
				m.verticalScrollWindowEnd--
				m.verticalScrollWindowYOffset = 0
			}
		}
	}
}

func (m *Model) moveHighlightDown() {
	m.rowCursorIndex++
	var wrapped bool

	if m.rowCursorIndex >= m.TotalRows() {
		m.rowCursorIndex = 0
		wrapped = true
	}

	if m.pageSize > 0 {
		m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	}
	currentWindowHeight := min(m.maxVisibleRows, m.TotalRows())
	if m.maxTotalHeight > 0 {
		m.verticalScrollWindowYOffset++
		if m.verticalScrollWindowYOffset > (currentWindowHeight - 1) {
			if wrapped {
				m.verticalScrollWindowStart = 0
				m.verticalScrollWindowEnd = currentWindowHeight - 1
				m.verticalScrollWindowYOffset = 0
			} else {
				m.verticalScrollWindowStart++
				m.verticalScrollWindowEnd++
				m.verticalScrollWindowYOffset = currentWindowHeight - 1
			}
		}
	}
}

func (m *Model) toggleSelect() {
	if !m.selectableRows || len(m.GetVisibleRows()) == 0 {
		return
	}

	rows := make([]Row, len(m.GetVisibleRows()))
	copy(rows, m.GetVisibleRows())

	currentSelectedState := rows[m.rowCursorIndex].selected

	rows[m.rowCursorIndex].selected = !currentSelectedState

	m.rows = rows

	m.appendUserEvent(UserEventRowSelectToggled{
		RowIndex:   m.rowCursorIndex,
		IsSelected: !currentSelectedState,
	})
}

func (m Model) updateFilterTextInput(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var reset bool
	var blurred bool
	prevTotalRows := m.TotalRows()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.FilterBlur) {
			m.filterTextInput.Blur()
			blurred = true
		}
		if key.Matches(msg, m.keyMap.FilterClear) {
			reset = true
		}
	}
	m.filterTextInput, cmd = m.filterTextInput.Update(msg)
	if len(m.rows) > m.TotalRows() || prevTotalRows != m.TotalRows() || reset || blurred {
		// filter did something
		if m.maxTotalHeight > 0 {
			// vertical scrolling mode
			m.rowCursorIndex = 0
			m.recalculateHeight()
		}
	}
	if m.pageSize > 0 {
		m.pageFirst()
	}

	return m, cmd
}

// nolint: cyclop // This is a series of Matches tests with minimal logic
func (m *Model) handleKeypress(msg tea.KeyMsg) {
	previousRowIndex := m.rowCursorIndex

	if key.Matches(msg, m.keyMap.RowDown) {
		m.moveHighlightDown()
	}

	if key.Matches(msg, m.keyMap.RowUp) {
		m.moveHighlightUp()
	}

	if key.Matches(msg, m.keyMap.RowSelectToggle) {
		m.toggleSelect()
	}

	if key.Matches(msg, m.keyMap.PageDown) {
		m.pageDown()
	}

	if key.Matches(msg, m.keyMap.PageUp) {
		m.pageUp()
	}

	if key.Matches(msg, m.keyMap.PageFirst) {
		m.pageFirst()
	}

	if key.Matches(msg, m.keyMap.PageLast) {
		m.pageLast()
	}

	if key.Matches(msg, m.keyMap.Filter) {
		if m.filtered {
			m.filterTextInput.Focus()
		}
	}

	if key.Matches(msg, m.keyMap.FilterClear) {
		if m.filtered {
			m.filterTextInput.Reset()
		}
	}

	if key.Matches(msg, m.keyMap.ScrollRight) {
		m.scrollRight()
	}

	if key.Matches(msg, m.keyMap.ScrollLeft) {
		m.scrollLeft()
	}

	if m.rowCursorIndex != previousRowIndex {
		m.appendUserEvent(UserEventHighlightedIndexChanged{
			PreviousRowIndex: previousRowIndex,
			SelectedRowIndex: m.rowCursorIndex,
		})
	}
}

// Update responds to input from the user or other messages from Bubble Tea.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	m.clearUserEvents()

	if !m.focused {
		return m, nil
	}

	if m.filterTextInput.Focused() {
		return m.updateFilterTextInput(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.handleKeypress(msg)
	}

	return m, nil
}
