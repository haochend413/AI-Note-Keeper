package note

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.ti.Width = width - 4
}
