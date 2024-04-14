package ui

type ActiveTab uint

const (
	TabChannels ActiveTab = iota
	TabChannel
)

type Mode uint

const (
	ModeNormal Mode = iota
	ModeInsert
)

type UIState struct {
	ActiveTab ActiveTab
	Mode      Mode
}

func NewUIState() *UIState {
	return &UIState{TabChannels, ModeNormal}
}
