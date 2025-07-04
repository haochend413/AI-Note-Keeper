package noteHistory

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/haochend413/mantis/defs"
	tui_defs "github.com/haochend413/mantis/defs/tui-defs"
)

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	// m.ti.Width = width - 4
	m.tb.SetHeight(height - 1)
	m.tb.SetWidth(width - 1)
}

func TopicsToString(pt []*defs.Topic) string {
	output := ""
	for i, topic := range pt {
		s := topic.Topic
		output += s
		if i != len(pt)-1 {
			output += " | "
		}
	}
	return output
}

func ContextFiltering(context tui_defs.Context, row []*defs.Note) []*defs.Note {
	now := time.Now()
	switch context {
	case tui_defs.Default:
		return row
	case tui_defs.Day:
		var result []*defs.Note
		for _, note := range row {
			if note.CreatedAt.Year() == now.Year() &&
				note.CreatedAt.Month() == now.Month() &&
				note.CreatedAt.Day() == now.Day() {
				result = append(result, note)
			}
			if note.ID == 0 {
				result = append(result, note)
			}
		}
		return result
	case tui_defs.Week:
		var result []*defs.Note
		for _, note := range row {
			y1, w1 := note.CreatedAt.ISOWeek()
			y2, w2 := now.ISOWeek()
			if y1 == y2 && w1 == w2 {
				result = append(result, note)
			}
			if note.ID == 0 {
				result = append(result, note)
			}
		}
		return result
	case tui_defs.Month:
		var result []*defs.Note
		for _, note := range row {
			if note.CreatedAt.Year() == now.Year() &&
				note.CreatedAt.Month() == now.Month() {
				result = append(result, note)
			}
			if note.ID == 0 {
				result = append(result, note)
			}
		}
		return result
	default:
		return row
	}
}

// This should depend on the context
func (m *Model) UpdateDisplay(data defs.DB_Data) {
	var rows []table.Row
	//filter by context
	filtered_notes := ContextFiltering(m.context, data.NoteData)
	for _, note := range filtered_notes {
		// transform topics into an array of strings;
		// This function will require further customizing;
		topics := TopicsToString(note.Topics)
		rows = append(rows, table.Row{
			note.CreatedAt.Format("06-01-02 15:04"), // assuming CreateTime is a string
			strconv.FormatUint(uint64(note.ID), 10), // assuming ID is a string
			note.Content,                            // assuming Content is a string
			topics,
		})
	}
	//Filter based on context
	m.tb.SetRows(rows)
}

func (m *Model) GetCurrentRowData() table.Row {
	return m.tb.SelectedRow()
}
