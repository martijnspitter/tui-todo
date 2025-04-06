package service_test

import (
	"testing"

	"github.com/martijnspitter/tui-todo/internal/service"
)

// Test NewTuiService for proper initialization
func TestNewTuiService(t *testing.T) {
	svc := service.NewTuiService()

	// Check default state
	if svc.CurrentView != service.OpenPane {
		t.Errorf("Expected default view to be OpenPane, got %v", svc.CurrentView)
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
			name:         "Switch to Open pane",
			key:          "1",
			expectedView: service.OpenPane,
		},
		{
			name:         "Switch to Doing pane",
			key:          "2",
			expectedView: service.DoingPane,
		},
		{
			name:         "Switch to Done pane",
			key:          "3",
			expectedView: service.DonePane,
		},
		{
			name:         "Switch to All pane",
			key:          "4",
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
