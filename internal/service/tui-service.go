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
	UpdateModal
	AboutModal
)

type TuiService struct {
	KeyMap          keys.KeyMap
	CurrentView     ViewType
	PrevView        ViewType
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
	if t.isPrevViewATab() {
		t.CurrentView = t.PrevView
	} else {
		t.CurrentView = OpenPane
	}
}

func (t *TuiService) SwitchToEditTodoView() {
	t.PrevView = t.CurrentView
	t.CurrentView = AddEditModal
}

func (t *TuiService) SwitchToConfirmDeleteView() {
	t.PrevView = t.CurrentView
	t.CurrentView = ConfirmDeleteModal
}

func (t *TuiService) ShouldShowModal() bool {
	return (t.CurrentView == AddEditModal ||
		t.CurrentView == ConfirmDeleteModal ||
		t.CurrentView == UpdateModal ||
		t.CurrentView == AboutModal)
}

func (t *TuiService) ToggleArchivedInAllView() {
	if t.CurrentView == AllPane {
		t.FilterState.IncludeArchived = !t.FilterState.IncludeArchived
	}
}

func (t *TuiService) isPrevViewATab() bool {
	return t.PrevView == OpenPane || t.PrevView == DoingPane || t.PrevView == DonePane || t.PrevView == AllPane
}

var (
	minWidthTitle      = 50.0
	maxWidthTitleRatio = 0.2

	minWidthDesc      = 50.0
	maxWidthDescRatio = 0.4
)

func (t *TuiService) DetermineMaxWidthsForTodo(screenWidth, requiredItemsWidth, dueDateWidth int) (titleWidth, desciptionWidth, leftWidth, remainderWidth int) {
	availableW := float64(screenWidth - requiredItemsWidth)

	titleW := max(availableW*maxWidthTitleRatio, minWidthTitle)
	descriptionW := max(availableW*maxWidthDescRatio, minWidthDesc)
	leftW := titleW + descriptionW
	remainderW := availableW - leftW - float64(dueDateWidth)

	return int(titleW), int(descriptionW), int(leftW), int(remainderW)
}

func (t *TuiService) SwitchToUpdateModalView() {
	t.CurrentView = UpdateModal
}

func (t *TuiService) SwitchToAboutModalView() {
	t.PrevView = t.CurrentView
	t.CurrentView = AboutModal
}
