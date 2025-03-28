package service

import (
	"github.com/martijnspitter/tui-todo/internal/keys"
	"github.com/martijnspitter/tui-todo/internal/models"
)

type ViewType int

const (
	ListView ViewType = iota
	NewView
	EditView
	ConfirmDelete
)

type TuiService struct {
	KeyMap      keys.KeyMap
	CurrentView ViewType
	FilterState FilterState
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

func (t *TuiService) SelectFilter(key string) {
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

func (t *TuiService) SwitchToNewTodoView() {
	t.CurrentView = NewView
}

func (t *TuiService) SwitchToListView() {
	t.CurrentView = ListView
}

func (t *TuiService) SwitchToEditTodoView() {
	t.CurrentView = EditView
}

func (t *TuiService) SwitchToConfirmDeleteView() {
	t.CurrentView = ConfirmDelete
}

func (t *TuiService) ShouldShowModal() bool {
	return t.CurrentView == EditView || t.CurrentView == ConfirmDelete
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
