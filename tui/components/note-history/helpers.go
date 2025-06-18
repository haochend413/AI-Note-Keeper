package noteHistory

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/haochend413/mantis/defs"
)

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	// m.ti.Width = width - 4
	m.tb.SetHeight(height - 1)
	m.tb.SetWidth(width - 1)
}

func (m *Model) UpdateDisplay(data defs.DB_Data) {
	var rows []table.Row
	for _, note := range data.NoteData {
		rows = append(rows, table.Row{
			note.CreatedAt.Format("06-01-02 15:04"), // assuming CreateTime is a string
			strconv.FormatUint(uint64(note.ID), 10), // assuming ID is a string
			note.Content,                            // assuming Content is a string
		})
	}
	m.tb.SetRows(rows)
}

func (m *Model) GetCurrentRowData() table.Row {
	return m.tb.SelectedRow()
}
