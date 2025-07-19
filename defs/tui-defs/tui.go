package tui_defs

type ViewType string

type AppStatus struct {
	CurrentView     ViewType
	LastRowSelected []string
}
