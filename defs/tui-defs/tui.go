package tui_defs

type ViewType string

type AppStatus struct {
	//data
	CurrentView ViewType
	CurrentID   int
}
