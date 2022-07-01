package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyID = "id"

	numCols = 5
	idWidth = 5

	colWidth = 3
)

type Model struct {
	ready       bool
	table       table.Model
	filterInput textinput.Model
}

func colKey(colNum int) string {
	return fmt.Sprintf("%d", colNum)
}

func genRow(id int) table.Row {
	data := table.RowData{
		columnKeyID: fmt.Sprintf("ID %s", fmt.Sprintf("%d", id)),
	}

	for i := 0; i < numCols; i++ {
		data[colKey(i)] = colWidth
	}

	return table.NewRow(data)
}

func NewModel() Model {
	return Model{}
}

func (m *Model) initTable(width int, height int) {
	rows := []table.Row{}

	// make 2x the rows than the available window height
	for i := 0; i < (height * 2); i++ {
		rows = append(rows, genRow(i))
	}

	cols := []table.Column{
		table.NewColumn(columnKeyID, "ID", idWidth).WithFiltered(true),
	}

	for i := 0; i < numCols; i++ {
		cols = append(cols, table.NewColumn(colKey(i), colKey(i+1), colWidth))
	}

	filterInput := textinput.New()
	filterInput.Prompt = "/ "
	filterInput.Placeholder = "(filter)"
	m.filterInput = filterInput

	m.table = table.New(cols).
		WithRows(rows).
		WithMaxTotalWidth(width).
		WithMaxTotalHeight(height).
		WithFilterInput(&filterInput).
		WithDynamicFooter(updateFooter).
		Filtered(true).
		Focused(true)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			// init the table now that we know the size of the screen
			// We add -2 here, to allow space for the title line and filter
			m.initTable(msg.Width, msg.Height-2)
			m.ready = true
		} else {
			// update the table if the window changes size
			m.table.SetMaxTotalWidth(msg.Width)
			m.table.SetMaxTotalHeight(msg.Height - 2)
		}

	case tea.KeyMsg:
		// event to filter
		if m.filterInput.Focused() {
			switch msg.String() {
			case "enter":
				m.filterInput.Blur()
			case "esc":
				m.filterInput.Reset()
			}
			m.filterInput, _ = m.filterInput.Update(msg)
			m.table = m.table.WithFilterInput(&m.filterInput)

			return m, tea.Batch(cmds...)
		} else {
			switch msg.String() {
			case "/":
				m.filterInput.Focus()
			case "ctrl+c", "esc", "q":
				cmds = append(cmds, tea.Quit)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func updateFooter(m table.Model) string {

	// this adjusts the footer during filtering to show the total row count
	// in the data set, vs the currently displayed amount post-filtering
	if m.GetCanFilter() && m.GetCurrentFilter() != "" {
		return fmt.Sprintf("%d/%d rows (filtered)", m.TotalRows(), m.GetTotalRowsUnfiltered())
	}

	// if we're not actively filtering, just return the total row count
	return fmt.Sprintf("%d rows", m.GetTotalRowsUnfiltered())
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString("Vertical Scroll: ↑/l, ↓/r to scroll, / to filter\n")
	body.WriteString(m.filterInput.View() + "\n")
	body.WriteString(m.table.View())

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
