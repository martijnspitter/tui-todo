package service

import "github.com/martijnspitter/tui-todo/internal/keys"

type SelectedPane int

const (
	Main SelectedPane = iota
	New
)

type TuiService struct {
	KeyMap       keys.KeyMap
	SelectedPane SelectedPane
}

func NewTuiService() *TuiService {
	return &TuiService{
		KeyMap: keys.DefaultKeyMap(),
	}
}
