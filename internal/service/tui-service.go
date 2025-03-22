package service

import "github.com/martijnspitter/tui-todo/internal/keys"

type TuiService struct {
	KeyMap keys.KeyMap
}

func NewTuiService() *TuiService {
	return &TuiService{
		KeyMap: keys.DefaultKeyMap(),
	}
}
