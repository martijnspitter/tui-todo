package service

import (
	"github.com/martijnspitter/tui-todo/internal/keys"
)

type ViewType int

const (
	OpenPane ViewType = iota
	DoingPane
	DonePane
	AllPane
	AddEditModal
	ConfirmDeleteModal
)

type TuiService struct {
	KeyMap          keys.KeyMap
	CurrentView     ViewType
	FilterState     FilterState
	ShowConfirmQuit bool
}

type FilterState struct {
	IsFilterActive  bool
	IncludeArchived bool
	FilterMode      FilterInputMode
}

type FilterInputMode int

const (
	FilterByTitle FilterInputMode = iota
	FilterByTag
)

func NewTuiService() *TuiService {
	return &TuiService{
		KeyMap:      keys.DefaultKeyMap(),
		CurrentView: DoingPane,
		FilterState: FilterState{
			IncludeArchived: false,
			IsFilterActive:  false,
			FilterMode:      FilterByTitle,
		},
	}
}

func (t *TuiService) SwitchPane(key string) {
	switch key {
	case "1":
		t.CurrentView = OpenPane
	case "2":
		t.CurrentView = DoingPane
	case "3":
		t.CurrentView = DonePane
	case "4":
		t.CurrentView = AllPane
	}
}

func (t *TuiService) ActivateTagFilter() {
	t.FilterState.IsFilterActive = true
	t.FilterState.FilterMode = FilterByTag
}

func (t *TuiService) ActivateTitleFilter() {
	t.FilterState.IsFilterActive = true
	t.FilterState.FilterMode = FilterByTitle
}

func (t *TuiService) IsTagFilterActive() bool {
	return t.FilterState.IsFilterActive &&
		t.FilterState.FilterMode == FilterByTag
}

func (t *TuiService) IsTitleFilterActive() bool {
	return t.FilterState.IsFilterActive &&
		t.FilterState.FilterMode == FilterByTitle
}

func (t *TuiService) ToggleShowConfirmQuit() {
	t.ShowConfirmQuit = !t.ShowConfirmQuit
}

func (t *TuiService) RemoveNameFilter() {
	t.FilterState.IsFilterActive = false
	t.FilterState.FilterMode = FilterByTitle
}

func (t *TuiService) SwitchToListView() {
	t.CurrentView = OpenPane
}

func (t *TuiService) SwitchToEditTodoView() {
	t.CurrentView = AddEditModal
}

func (t *TuiService) SwitchToConfirmDeleteView() {
	t.CurrentView = ConfirmDeleteModal
}

func (t *TuiService) ShouldShowModal() bool {
	return t.CurrentView == AddEditModal || t.CurrentView == ConfirmDeleteModal
}

func (t *TuiService) ToggleArchivedInAllView() {
	if t.CurrentView == AllPane {
		t.FilterState.IncludeArchived = !t.FilterState.IncludeArchived
	}
}
