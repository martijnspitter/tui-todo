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
}

type FilterMode int

const (
	StatusFilter FilterMode = iota
	AllFilter
	TagFilter
	NameDescFilter
)

func NewTuiService() *TuiService {
	return &TuiService{
		KeyMap:      keys.DefaultKeyMap(),
		CurrentView: ListView,
		FilterState: FilterState{
			Mode:            StatusFilter,
			Status:          models.Doing,
			IncludeArchived: false,
		},
	}
}

func (t *TuiService) SwitchPane(key string) {
	switch key {
	case "1":
		t.FilterState.Mode = StatusFilter
		t.FilterState.Status = models.Open
	case "2":
		t.FilterState.Mode = StatusFilter
		t.FilterState.Status = models.Doing
	case "3":
		t.FilterState.Mode = StatusFilter
		t.FilterState.Status = models.Done
	case "4":
		t.FilterState.Mode = AllFilter
	}
}

func (t *TuiService) ToggleShowConfirmQuit() {
	t.ShowConfirmQuit = !t.ShowConfirmQuit
}

func (t *TuiService) ActivateNameFilter() {
	t.FilterState.Mode = NameDescFilter
}

func (t *TuiService) RemoveNameFilter() {
	t.FilterState.Mode = StatusFilter
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

func (t *TuiService) FilterByTag(tag string) {
	t.FilterState.Mode = TagFilter
	t.FilterState.Tag = tag
}

func (t *TuiService) ToggleArchivedInAllView() {
	if t.FilterState.Mode == AllFilter {
		t.FilterState.IncludeArchived = !t.FilterState.IncludeArchived
	}
}
