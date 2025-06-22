package noteDetail

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	// m.ti.Width = width - 4
	m.vp.Height = height - 1
	m.vp.Width = width - 6

}

func (m *Model) UpdateDisplay(content string) {
	rendered, _ := m.renderer.Render(content)
	m.vp.SetContent(rendered)

}
