package service_test

import (
	"testing"

	"github.com/martijnspitter/tui-todo/internal/service"
)

// Test NewTuiService for proper initialization
func TestNewTuiService(t *testing.T) {
	svc := service.NewTuiService()

	// Check default state
	if svc.CurrentView != service.TodayPane {
		t.Errorf("Expected default view to be DoingPane, got %v", svc.CurrentView)
	}

	if svc.FilterState.IsFilterActive {
		t.Error("Expected filter to be inactive by default")
	}

	if svc.FilterState.IncludeArchived {
		t.Error("Expected archived items to be excluded by default")
	}

	if svc.FilterState.FilterMode != service.FilterByTitle {
		t.Errorf("Expected default filter mode to be FilterByTitle, got %v", svc.FilterState.FilterMode)
	}

	if svc.ShowConfirmQuit {
		t.Error("Expected ShowConfirmQuit to be false by default")
	}

}

// Test SwitchPane functionality
func TestSwitchPane(t *testing.T) {
	testCases := []struct {
		name         string
		key          string
		expectedView service.ViewType
	}{
		{
			name:         "Switch to Today pane",
			key:          "1",
			expectedView: service.TodayPane,
		},
		{
			name:         "Switch to Open pane",
			key:          "2",
			expectedView: service.OpenPane,
		},
		{
			name:         "Switch to Doing pane",
			key:          "3",
			expectedView: service.DoingPane,
		},
		{
			name:         "Switch to Done pane",
			key:          "4",
			expectedView: service.DonePane,
		},
		{
			name:         "Switch to All pane",
			key:          "5",
			expectedView: service.AllPane,
		},
		{
			name:         "Invalid key doesn't change view",
			key:          "invalid",
			expectedView: service.OpenPane, // Should remain unchanged
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewTuiService()

			// For the last test case, we need to ensure the view is OpenPane
			// before we test that an invalid key doesn't change it
			if tc.name == "Invalid key doesn't change view" {
				svc.CurrentView = service.OpenPane
			}

			svc.SwitchPane(tc.key)

			if svc.CurrentView != tc.expectedView {
				t.Errorf("Expected view %v, got %v", tc.expectedView, svc.CurrentView)
			}
		})
	}
}

// Test filter activation methods
func TestFilterActivation(t *testing.T) {
	t.Run("Activate tag filter", func(t *testing.T) {
		svc := service.NewTuiService()
		svc.ActivateTagFilter()

		if !svc.FilterState.IsFilterActive {
			t.Error("Expected filter to be active")
		}
		if svc.FilterState.FilterMode != service.FilterByTag {
			t.Errorf("Expected filter mode to be FilterByTag, got %v", svc.FilterState.FilterMode)
		}
	})

	t.Run("Activate title filter", func(t *testing.T) {
		svc := service.NewTuiService()
		svc.ActivateTitleFilter()

		if !svc.FilterState.IsFilterActive {
			t.Error("Expected filter to be active")
		}
		if svc.FilterState.FilterMode != service.FilterByTitle {
			t.Errorf("Expected filter mode to be FilterByTitle, got %v", svc.FilterState.FilterMode)
		}
	})
}

// Test filter status check methods
func TestFilterStatusChecks(t *testing.T) {
	testCases := []struct {
		name              string
		setupFunc         func(*service.TuiService)
		expectTagActive   bool
		expectTitleActive bool
	}{
		{
			name: "No filter active",
			setupFunc: func(s *service.TuiService) {
				s.FilterState.IsFilterActive = false
			},
			expectTagActive:   false,
			expectTitleActive: false,
		},
		{
			name: "Tag filter active",
			setupFunc: func(s *service.TuiService) {
				s.FilterState.IsFilterActive = true
				s.FilterState.FilterMode = service.FilterByTag
			},
			expectTagActive:   true,
			expectTitleActive: false,
		},
		{
			name: "Title filter active",
			setupFunc: func(s *service.TuiService) {
				s.FilterState.IsFilterActive = true
				s.FilterState.FilterMode = service.FilterByTitle
			},
			expectTagActive:   false,
			expectTitleActive: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewTuiService()
			tc.setupFunc(svc)

			if svc.IsTagFilterActive() != tc.expectTagActive {
				t.Errorf("Expected IsTagFilterActive() to be %v, got %v",
					tc.expectTagActive, svc.IsTagFilterActive())
			}

			if svc.IsTitleFilterActive() != tc.expectTitleActive {
				t.Errorf("Expected IsTitleFilterActive() to be %v, got %v",
					tc.expectTitleActive, svc.IsTitleFilterActive())
			}
		})
	}
}

// Test ToggleShowConfirmQuit
func TestToggleShowConfirmQuit(t *testing.T) {
	svc := service.NewTuiService()

	// Initially false
	if svc.ShowConfirmQuit {
		t.Error("Expected ShowConfirmQuit to be false initially")
	}

	// Toggle to true
	svc.ToggleShowConfirmQuit()
	if !svc.ShowConfirmQuit {
		t.Error("Expected ShowConfirmQuit to be true after toggle")
	}

	// Toggle back to false
	svc.ToggleShowConfirmQuit()
	if svc.ShowConfirmQuit {
		t.Error("Expected ShowConfirmQuit to be false after second toggle")
	}
}

// Test RemoveNameFilter
func TestRemoveNameFilter(t *testing.T) {
	svc := service.NewTuiService()

	// Setup initial state
	svc.FilterState.IsFilterActive = true
	svc.FilterState.FilterMode = service.FilterByTag

	// Remove filter
	svc.RemoveNameFilter()

	if svc.FilterState.IsFilterActive {
		t.Error("Expected filter to be inactive after removal")
	}

	if svc.FilterState.FilterMode != service.FilterByTitle {
		t.Errorf("Expected filter mode to reset to FilterByTitle, got %v",
			svc.FilterState.FilterMode)
	}
}

// Test view switching methods
func TestViewSwitchingMethods(t *testing.T) {
	testCases := []struct {
		name         string
		switchFunc   func(*service.TuiService)
		expectedView service.ViewType
	}{
		{
			name:         "Switch to list view",
			switchFunc:   func(s *service.TuiService) { s.SwitchToListView() },
			expectedView: service.OpenPane,
		},
		{
			name:         "Switch to edit todo view",
			switchFunc:   func(s *service.TuiService) { s.SwitchToEditTodoView() },
			expectedView: service.AddEditModal,
		},
		{
			name:         "Switch to confirm delete view",
			switchFunc:   func(s *service.TuiService) { s.SwitchToConfirmDeleteView() },
			expectedView: service.ConfirmDeleteModal,
		},
		{
			name:         "Switch to Update modal view",
			switchFunc:   func(s *service.TuiService) { s.SwitchToUpdateModalView() },
			expectedView: service.UpdateModal,
		},
		{
			name:         "Switch to About modal view",
			switchFunc:   func(s *service.TuiService) { s.SwitchToAboutModalView() },
			expectedView: service.AboutModal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewTuiService()
			tc.switchFunc(svc)

			if svc.CurrentView != tc.expectedView {
				t.Errorf("Expected view %v, got %v", tc.expectedView, svc.CurrentView)
			}
		})
	}
}

// Test ShouldShowModal
func TestShouldShowModal(t *testing.T) {
	testCases := []struct {
		name        string
		view        service.ViewType
		expectModal bool
	}{
		{
			name:        "Open pane is not modal",
			view:        service.OpenPane,
			expectModal: false,
		},
		{
			name:        "Doing pane is not modal",
			view:        service.DoingPane,
			expectModal: false,
		},
		{
			name:        "Done pane is not modal",
			view:        service.DonePane,
			expectModal: false,
		},
		{
			name:        "All pane is not modal",
			view:        service.AllPane,
			expectModal: false,
		},
		{
			name:        "Add/Edit modal is modal",
			view:        service.AddEditModal,
			expectModal: true,
		},
		{
			name:        "Confirm delete modal is modal",
			view:        service.ConfirmDeleteModal,
			expectModal: true,
		},
		{
			name:        "UpdateModal is modal",
			view:        service.UpdateModal,
			expectModal: true,
		},
		{
			name:        "About Modal is modal",
			view:        service.AboutModal,
			expectModal: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewTuiService()
			svc.CurrentView = tc.view

			if svc.ShouldShowModal() != tc.expectModal {
				t.Errorf("Expected ShouldShowModal() to be %v for view %v, got %v",
					tc.expectModal, tc.view, svc.ShouldShowModal())
			}
		})
	}
}

// Test ToggleArchivedInAllView
func TestToggleArchivedInAllView(t *testing.T) {
	testCases := []struct {
		name            string
		initialView     service.ViewType
		initialArchived bool
		expectArchived  bool
	}{
		{
			name:            "Toggle in All view",
			initialView:     service.AllPane,
			initialArchived: false,
			expectArchived:  true,
		},
		{
			name:            "Toggle in All view when already true",
			initialView:     service.AllPane,
			initialArchived: true,
			expectArchived:  false,
		},
		{
			name:            "No toggle in Open view",
			initialView:     service.OpenPane,
			initialArchived: false,
			expectArchived:  false,
		},
		{
			name:            "No toggle in Doing view",
			initialView:     service.DoingPane,
			initialArchived: false,
			expectArchived:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewTuiService()
			svc.CurrentView = tc.initialView
			svc.FilterState.IncludeArchived = tc.initialArchived

			svc.ToggleArchivedInAllView()

			if svc.FilterState.IncludeArchived != tc.expectArchived {
				t.Errorf("Expected IncludeArchived to be %v after toggle, got %v",
					tc.expectArchived, svc.FilterState.IncludeArchived)
			}
		})
	}
}

func TestSwitchToListView(t *testing.T) {
	testCases := []struct {
		name         string
		initialView  service.ViewType
		prevView     service.ViewType
		expectedView service.ViewType
	}{
		{
			name:         "Return to DoingPane when it was previous tab",
			initialView:  service.AddEditModal,
			prevView:     service.DoingPane,
			expectedView: service.DoingPane,
		},
		{
			name:         "Return to OpenPane when it was previous tab",
			initialView:  service.ConfirmDeleteModal,
			prevView:     service.OpenPane,
			expectedView: service.OpenPane,
		},
		{
			name:         "Return to DonePane when it was previous tab",
			initialView:  service.AddEditModal,
			prevView:     service.DonePane,
			expectedView: service.DonePane,
		},
		{
			name:         "Return to AllPane when it was previous tab",
			initialView:  service.ConfirmDeleteModal,
			prevView:     service.AllPane,
			expectedView: service.AllPane,
		},
		{
			name:         "Default to OpenPane when previous view was not a tab",
			initialView:  service.AddEditModal,
			prevView:     service.ConfirmDeleteModal, // Non-tab previous view
			expectedView: service.OpenPane,
		},
		{
			name:         "Default to OpenPane when previous view is invalid",
			initialView:  service.AddEditModal,
			prevView:     99, // Invalid view type
			expectedView: service.OpenPane,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewTuiService()
			svc.CurrentView = tc.initialView
			svc.PrevView = tc.prevView

			svc.SwitchToListView()

			if svc.CurrentView != tc.expectedView {
				t.Errorf("Expected CurrentView to be %v after SwitchToListView(), got %v",
					tc.expectedView, svc.CurrentView)
			}
		})
	}
}

func TestSwitchToUpdateModalView(t *testing.T) {
	svc := service.NewTuiService()
	// Switch to update modal view
	svc.SwitchToUpdateModalView()

	// Check if CurrentView is now UpdateModal
	if svc.CurrentView != service.UpdateModal {
		t.Errorf("Expected CurrentView to be UpdateModal, got %v", svc.CurrentView)
	}
}

func TestModalPreservesPrevView(t *testing.T) {
	testCases := []struct {
		name             string
		initialView      service.ViewType
		expectedPrevView service.ViewType
	}{
		{
			name:             "Preserve OpenPane as previous view",
			initialView:      service.OpenPane,
			expectedPrevView: service.OpenPane,
		},
		{
			name:             "Preserve DoingPane as previous view",
			initialView:      service.DoingPane,
			expectedPrevView: service.DoingPane,
		},
		{
			name:             "Preserve DonePane as previous view",
			initialView:      service.DonePane,
			expectedPrevView: service.DonePane,
		},
		{
			name:             "Preserve AllPane as previous view",
			initialView:      service.AllPane,
			expectedPrevView: service.AllPane,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewTuiService()
			svc.CurrentView = tc.initialView

			// Switch to update modal view
			svc.SwitchToEditTodoView()

			if svc.PrevView != tc.expectedPrevView {
				t.Errorf("Expected PrevView to be %v after switching to modal, got %v",
					tc.expectedPrevView, svc.PrevView)
			}
		})
	}
}
