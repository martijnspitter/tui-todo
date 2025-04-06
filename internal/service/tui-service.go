package service

import (
	"github.com/martijnspitter/tui-todo/internal/keys"
	"github.com/martijnspitter/tui-todo/internal/models"
)

type ViewType int

const (
	ListView ViewType = iota
	AddEditView
	ConfirmDelete
)

type TuiService struct {
	KeyMap          keys.KeyMap
	CurrentView     ViewType
	FilterState     FilterState
	ShowConfirmQuit bool
}

type FilterState struct {
	Mode            FilterMode
	Status          models.Status
	Tag             string
	IncludeArchived bool
	FilterMode      FilterInputMode
}

type FilterMode int

const (
	StatusPanes FilterMode = iota
	AllPane
	Filtering
)

type FilterInputMode int

const (
	FilterByTitle FilterInputMode = iota
	FilterByTag
)

func NewTuiService() *TuiService {
	return &TuiService{
		KeyMap:      keys.DefaultKeyMap(),
		CurrentView: ListView,
		FilterState: FilterState{
			Mode:            StatusPanes,
			Status:          models.Doing,
			IncludeArchived: false,
			FilterMode:      FilterByTitle,
		},
	}
}

func (t *TuiService) SwitchPane(key string) {
	switch key {
	case "1":
		t.FilterState.Mode = StatusPanes
		t.FilterState.Status = models.Open
	case "2":
		t.FilterState.Mode = StatusPanes
		t.FilterState.Status = models.Doing
	case "3":
		t.FilterState.Mode = StatusPanes
		t.FilterState.Status = models.Done
	case "4":
		t.FilterState.Mode = AllPane
	}
}

func (t *TuiService) ActivateTagFilter() {
	t.FilterState.FilterMode = FilterByTag
}

func (t *TuiService) ActivateTitleFilter() {
	t.FilterState.FilterMode = FilterByTitle
}

func (t *TuiService) IsTagFilterActive() bool {
	return t.FilterState.Mode == AllPane &&
		t.FilterState.FilterMode == FilterByTag
}

func (t *TuiService) IsTitleFilterActive() bool {
	return t.FilterState.FilterMode == FilterByTitle
}

func (t *TuiService) ToggleShowConfirmQuit() {
	t.ShowConfirmQuit = !t.ShowConfirmQuit
}

func (t *TuiService) ActivateNameFilter() {
	t.FilterState.Mode = Filtering
}

func (t *TuiService) RemoveNameFilter() {
	t.FilterState.FilterMode = FilterByTitle
}

func (t *TuiService) SwitchToListView() {
	t.CurrentView = ListView
}

func (t *TuiService) SwitchToEditTodoView() {
	t.CurrentView = AddEditView
}

func (t *TuiService) SwitchToConfirmDeleteView() {
	t.CurrentView = ConfirmDelete
}

func (t *TuiService) ShouldShowModal() bool {
	return t.CurrentView == AddEditView || t.CurrentView == ConfirmDelete
}

func (t *TuiService) ToggleArchivedInAllView() {
	if t.FilterState.Mode == AllPane {
		t.FilterState.IncludeArchived = !t.FilterState.IncludeArchived
	}
}
