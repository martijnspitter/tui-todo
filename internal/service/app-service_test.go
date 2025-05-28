package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/repository"
	"github.com/martijnspitter/tui-todo/internal/service"
	"pgregory.net/rapid"
)

// MockTodoRepository implements repository.TodoRepository for testing
type MockTodoRepository struct {
	// These fields track calls to the methods
	CreatedTodos []*models.Todo
	UpdatedTodos []*models.Todo
	DeletedIDs   []int64
	AddedTags    map[int64][]string
	RemovedTags  map[int64][]string

	// Mock data to return
	MockTodos    []*models.Todo
	MockTodo     *models.Todo
	MockError    error
	MockTodoTags []string
	MockTags     []models.Tag
	SearchQuery  string
}

// Implement all repository methods...
func (m *MockTodoRepository) Create(todo *models.Todo) error {
	if m.MockError != nil {
		return m.MockError
	}
	m.CreatedTodos = append(m.CreatedTodos, todo)
	todo.ID = 1 // Simulate auto-increment ID
	return nil
}

func (m *MockTodoRepository) GetByID(id int64) (*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodo, nil
}

func (m *MockTodoRepository) GetAll(filters ...repository.Filter) ([]*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodos, nil
}

func (m *MockTodoRepository) Update(todo *models.Todo) error {
	if m.MockError != nil {
		return m.MockError
	}
	m.UpdatedTodos = append(m.UpdatedTodos, todo)
	return nil
}

func (m *MockTodoRepository) Delete(id int64) error {
	if m.MockError != nil {
		return m.MockError
	}
	m.DeletedIDs = append(m.DeletedIDs, id)
	return nil
}

// Mock the filtered queries
func (m *MockTodoRepository) GetOpen() ([]*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodos, nil
}

func (m *MockTodoRepository) GetActive() ([]*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodos, nil
}

func (m *MockTodoRepository) GetCompleted() ([]*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodos, nil
}

func (m *MockTodoRepository) GetBlocked() ([]*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodos, nil
}

func (m *MockTodoRepository) GetArchived() ([]*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodos, nil
}

func (m *MockTodoRepository) Search(query string) ([]*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	m.SearchQuery = query
	return m.MockTodos, nil
}

// Tag methods
func (m *MockTodoRepository) AddTagToTodo(id int64, tagname string) error {
	if m.MockError != nil {
		return m.MockError
	}
	if m.AddedTags == nil {
		m.AddedTags = make(map[int64][]string)
	}
	m.AddedTags[id] = append(m.AddedTags[id], tagname)
	return nil
}

func (m *MockTodoRepository) RemoveTagFromTodo(id int64, tagname string) error {
	if m.MockError != nil {
		return m.MockError
	}
	if m.RemovedTags == nil {
		m.RemovedTags = make(map[int64][]string)
	}
	m.RemovedTags[id] = append(m.RemovedTags[id], tagname)
	return nil
}

func (m *MockTodoRepository) GetTodoTags(id int64) ([]string, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodoTags, nil
}

func (m *MockTodoRepository) FindTodosByTag(tagName string) ([]*models.Todo, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	return m.MockTodos, nil
}

func (m *MockTodoRepository) GetAllTags() ([]*models.Tag, error) {
	if m.MockError != nil {
		return nil, m.MockError
	}
	if m.MockTags == nil {
		m.MockTags = []models.Tag{
			{Name: "test", ID: 1},
			{Name: "important", ID: 2},
			{Name: "urgent", ID: 3},
		}
	}
	// Convert []models.Tag to []*models.Tag
	tagPointers := make([]*models.Tag, len(m.MockTags))
	for i := range m.MockTags {
		tagPointers[i] = &m.MockTags[i]
	}
	return tagPointers, nil
}

func (m *MockTodoRepository) DeleteTag(id int64) error {
	if m.MockError != nil {
		return m.MockError
	}
	// Simulate tag deletion
	for i, tag := range m.MockTags {
		if tag.ID == id {
			m.MockTags = append(m.MockTags[:i], m.MockTags[i+1:]...)
			return nil
		}
	}
	return errors.New("tag not found")
}

// Helper function to create a test todo
func createTestTodo(id int64) *models.Todo {
	now := time.Now()
	return &models.Todo{
		ID:          id,
		Title:       "Test Todo",
		Description: "Test Description",
		Status:      models.Open,
		Priority:    models.Medium,
		CreatedAt:   now,
		UpdatedAt:   now,
		Tags:        []string{"test"},
	}
}

// Test SaveTodo (which calls either CreateTodo or UpdateTodo)
func TestSaveTodo(t *testing.T) {
	testCases := []struct {
		name         string
		todo         *models.Todo
		tags         []string
		mockError    error
		wantError    bool
		expectCreate bool
		expectUpdate bool
	}{
		{
			name:         "Create new todo",
			todo:         &models.Todo{Title: "New Todo", Description: "New Description", Priority: models.High},
			tags:         []string{"new", "important"},
			mockError:    nil,
			wantError:    false,
			expectCreate: true,
			expectUpdate: false,
		},
		{
			name:         "Update existing todo",
			todo:         &models.Todo{ID: 5, Title: "Update Todo", Description: "Update Description", Priority: models.Low},
			tags:         []string{"update"},
			mockError:    nil,
			wantError:    false,
			expectCreate: false,
			expectUpdate: true,
		},
		{
			name:         "Error during save",
			todo:         &models.Todo{Title: "Error Todo"},
			tags:         []string{},
			mockError:    errors.New("save error"),
			wantError:    true,
			expectCreate: false,
			expectUpdate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			mockRepo := &MockTodoRepository{
				MockError: tc.mockError,
				MockTodo:  tc.todo, // For GetByID during update
			}

			// Create service
			svc := service.NewAppService(mockRepo)

			// Call method
			err := svc.SaveTodo(tc.todo, tc.tags)

			// Check expectations
			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				if tc.expectCreate {
					if len(mockRepo.CreatedTodos) != 1 {
						t.Errorf("Expected todo to be created but it wasn't")
					}
				} else if tc.expectUpdate {
					if len(mockRepo.UpdatedTodos) != 1 {
						t.Errorf("Expected todo to be updated but it wasn't")
					}
				}
			}
		})
	}
}

// Tests for CreateTodo functionality
func TestCreateTodo(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name        string
		title       string
		description string
		priority    models.Priority
		tags        []string
		mockError   error
		wantError   bool
		dueDate     *time.Time
		status      models.Status
	}{
		{
			name:        "Successful creation",
			title:       "Test Todo",
			description: "Test Description",
			priority:    models.Medium,
			tags:        []string{"test", "todo"},
			mockError:   nil,
			wantError:   false,
			dueDate:     &now,
			status:      models.Doing,
		},
		{
			name:        "Repository error",
			title:       "Error Todo",
			description: "Error Description",
			priority:    models.Low,
			tags:        []string{},
			mockError:   errors.New("db error"),
			wantError:   true,
			dueDate:     &now,
			status:      models.Doing,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := &MockTodoRepository{
				MockError: tc.mockError,
			}

			// Create service with mock repo
			svc := service.NewAppService(mockRepo)

			// Call method under test
			err := svc.CreateTodo(tc.title, tc.description, tc.priority, tc.tags, tc.dueDate, tc.status)

			// Check error expectation
			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				// Verify todo was created with correct data
				if len(mockRepo.CreatedTodos) != 1 {
					t.Fatalf("Expected 1 todo to be created, got %d", len(mockRepo.CreatedTodos))
				}

				todo := mockRepo.CreatedTodos[0]
				if todo.Title != tc.title {
					t.Errorf("Expected title %q, got %q", tc.title, todo.Title)
				}
				if todo.Description != tc.description {
					t.Errorf("Expected description %q, got %q", tc.description, todo.Description)
				}
				if todo.Priority != tc.priority {
					t.Errorf("Expected priority %v, got %v", tc.priority, todo.Priority)
				}
				if todo.DueDate == nil || !todo.DueDate.Equal(now) {
					expected := now.Format(time.RFC3339)
					actual := "nil"
					if todo.DueDate != nil {
						actual = todo.DueDate.Format(time.RFC3339)
					}
					t.Errorf("Expected due date %s, got %s", expected, actual)
				}
				if todo.Status != tc.status {
					t.Errorf("Expected status %v, got %v", tc.status, todo.Status)
				}

				// Verify tags were added
				if len(tc.tags) > 0 {
					if mockRepo.AddedTags == nil || len(mockRepo.AddedTags[1]) != len(tc.tags) {
						t.Errorf("Expected %d tags to be added, got %d",
							len(tc.tags),
							len(mockRepo.AddedTags[1]))
					}
				}
			}
		})
	}
}

// Tests for GetAllTodos functionality
func TestGetAllTodos(t *testing.T) {
	now := time.Now()
	mockTodos := []*models.Todo{
		{
			ID:          1,
			Title:       "Todo 1",
			Description: "Description 1",
			Status:      models.Open,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Title:       "Todo 2",
			Description: "Description 2",
			Status:      models.Doing,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	testCases := []struct {
		name      string
		mockTodos []*models.Todo
		mockError error
		wantError bool
		wantCount int
	}{
		{
			name:      "Success with todos",
			mockTodos: mockTodos,
			mockError: nil,
			wantError: false,
			wantCount: 2,
		},
		{
			name:      "Success with empty list",
			mockTodos: []*models.Todo{},
			mockError: nil,
			wantError: false,
			wantCount: 0,
		},
		{
			name:      "Repository error",
			mockTodos: nil,
			mockError: errors.New("db error"),
			wantError: true,
			wantCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			mockRepo := &MockTodoRepository{
				MockTodos: tc.mockTodos,
				MockError: tc.mockError,
			}

			// Create service
			svc := service.NewAppService(mockRepo)

			// Call method
			todos, err := svc.GetAllTodos(false)

			// Check expectations
			if tc.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				if len(todos) != tc.wantCount {
					t.Errorf("Expected %d todos, got %d", tc.wantCount, len(todos))
				}

				if tc.wantCount > 0 {
					// Check if todos match
					for i, todo := range todos {
						if todo.ID != tc.mockTodos[i].ID {
							t.Errorf("Expected todo ID %d, got %d", tc.mockTodos[i].ID, todo.ID)
						}
						if todo.Title != tc.mockTodos[i].Title {
							t.Errorf("Expected todo title %q, got %q", tc.mockTodos[i].Title, todo.Title)
						}
					}
				}
			}
		})
	}
}

// Test GetTodo functionality
func TestGetTodo(t *testing.T) {
	testCases := []struct {
		name      string
		todoID    int64
		mockTodo  *models.Todo
		mockError error
		wantError bool
	}{
		{
			name:      "Success getting todo",
			todoID:    1,
			mockTodo:  createTestTodo(1),
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Todo not found",
			todoID:    99,
			mockTodo:  nil,
			mockError: errors.New("not found"),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			mockRepo := &MockTodoRepository{
				MockTodo:  tc.mockTodo,
				MockError: tc.mockError,
			}

			// Create service
			svc := service.NewAppService(mockRepo)

			// Call method
			todo, err := svc.GetTodo(tc.todoID)

			// Check expectations
			if tc.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				if todo == nil {
					t.Fatal("Expected todo but got nil")
				}

				if todo.ID != tc.mockTodo.ID {
					t.Errorf("Expected todo ID %d, got %d", tc.mockTodo.ID, todo.ID)
				}
				if todo.Title != tc.mockTodo.Title {
					t.Errorf("Expected todo title %q, got %q", tc.mockTodo.Title, todo.Title)
				}
			}
		})
	}
}

// Test UpdateTodo functionality
func TestUpdateTodo(t *testing.T) {
	testCases := []struct {
		name      string
		todo      *models.Todo
		tags      []string
		mockError error
		wantError bool
	}{
		{
			name:      "Success updating todo",
			todo:      createTestTodo(1),
			tags:      []string{"updated", "important"},
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Error updating todo",
			todo:      createTestTodo(2),
			tags:      []string{},
			mockError: errors.New("update error"),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get time before update to check it was updated
			beforeUpdate := tc.todo.UpdatedAt

			// Setup mock
			mockRepo := &MockTodoRepository{
				MockError: tc.mockError,
			}

			// Create service
			svc := service.NewAppService(mockRepo)

			// Call method
			err := svc.UpdateTodo(tc.todo, tc.tags)

			// Check expectations
			if tc.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				// Verify todo was updated
				if len(mockRepo.UpdatedTodos) != 1 {
					t.Errorf("Expected todo to be updated")
				} else {
					updatedTodo := mockRepo.UpdatedTodos[0]

					// Check updated_at timestamp was changed
					if !updatedTodo.UpdatedAt.After(beforeUpdate) {
						t.Errorf("Expected updated_at to be updated")
					}

					// Verify tags were added
					if len(tc.tags) > 0 {
						if mockRepo.AddedTags == nil {
							t.Errorf("Expected tags to be added")
						} else {
							addedTags := mockRepo.AddedTags[tc.todo.ID]
							if len(addedTags) != len(tc.tags) {
								t.Errorf("Expected %d tags, got %d", len(tc.tags), len(addedTags))
							}
						}
					}
				}
			}
		})
	}
}

// Test DeleteTodo functionality
func TestDeleteTodo(t *testing.T) {
	testCases := []struct {
		name      string
		todoID    int64
		mockError error
		wantError bool
	}{
		{
			name:      "Success deleting todo",
			todoID:    1,
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Error deleting todo",
			todoID:    2,
			mockError: errors.New("delete error"),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			mockRepo := &MockTodoRepository{
				MockError: tc.mockError,
			}

			// Create service
			svc := service.NewAppService(mockRepo)

			// Call method
			err := svc.DeleteTodo(tc.todoID)

			// Check expectations
			if tc.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				// Verify todo was deleted
				if len(mockRepo.DeletedIDs) != 1 || mockRepo.DeletedIDs[0] != tc.todoID {
					t.Errorf("Expected todo ID %d to be deleted", tc.todoID)
				}
			}
		})
	}
}

// Test status change methods
func TestStatusChangeMethods(t *testing.T) {
	// Create a struct for all status test cases
	statusTests := []struct {
		methodName     string
		methodFunc     func(*service.AppService, int64) error
		expectedStatus models.Status
	}{
		{
			methodName:     "MarkAsOpen",
			methodFunc:     func(svc *service.AppService, id int64) error { return svc.MarkAsOpen(id) },
			expectedStatus: models.Open,
		},
		{
			methodName:     "MarkAsDoing",
			methodFunc:     func(svc *service.AppService, id int64) error { return svc.MarkAsDoing(id) },
			expectedStatus: models.Doing,
		},
		{
			methodName:     "MarkAsDone",
			methodFunc:     func(svc *service.AppService, id int64) error { return svc.MarkAsDone(id) },
			expectedStatus: models.Done,
		},
	}

	// For each status method, run a set of test cases
	for _, statusTest := range statusTests {
		t.Run(statusTest.methodName, func(t *testing.T) {
			testCases := []struct {
				name      string
				todoID    int64
				mockTodo  *models.Todo
				mockError error
				wantError bool
			}{
				{
					name:      "Success changing status",
					todoID:    1,
					mockTodo:  createTestTodo(1),
					mockError: nil,
					wantError: false,
				},
				{
					name:      "Error fetching todo",
					todoID:    2,
					mockTodo:  nil,
					mockError: errors.New("fetch error"),
					wantError: true,
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Setup mock
					mockRepo := &MockTodoRepository{
						MockTodo:  tc.mockTodo,
						MockError: tc.mockError,
					}

					// Create service
					svc := service.NewAppService(mockRepo)

					// Get time before update to check it was updated
					var beforeUpdate time.Time
					if tc.mockTodo != nil {
						beforeUpdate = tc.mockTodo.UpdatedAt
					}

					// Call the method being tested
					err := statusTest.methodFunc(svc, tc.todoID)

					// Check expectations
					if tc.wantError {
						if err == nil {
							t.Error("Expected error but got nil")
						}
					} else {
						if err != nil {
							t.Errorf("Expected no error but got: %v", err)
						}

						// Verify todo was updated
						if len(mockRepo.UpdatedTodos) != 1 {
							t.Errorf("Expected todo to be updated")
						} else {
							updatedTodo := mockRepo.UpdatedTodos[0]

							// Check status was changed
							if updatedTodo.Status != statusTest.expectedStatus {
								t.Errorf("Expected status %v, got %v",
									statusTest.expectedStatus, updatedTodo.Status)
							}

							// Check updated_at timestamp was changed
							if !updatedTodo.UpdatedAt.After(beforeUpdate) {
								t.Errorf("Expected updated_at to be updated")
							}
						}
					}
				})
			}
		})
	}
}

// Test filtered query methods
func TestFilteredQueryMethods(t *testing.T) {
	// Create a struct for all filtered query test cases
	queryTests := []struct {
		methodName string
		methodFunc func(*service.AppService) ([]*models.Todo, error)
	}{
		{
			methodName: "GetOpenTodos",
			methodFunc: func(svc *service.AppService) ([]*models.Todo, error) { return svc.GetOpenTodos() },
		},
		{
			methodName: "GetActiveTodos",
			methodFunc: func(svc *service.AppService) ([]*models.Todo, error) { return svc.GetActiveTodos() },
		},
		{
			methodName: "GetCompletedTodos",
			methodFunc: func(svc *service.AppService) ([]*models.Todo, error) { return svc.GetCompletedTodos() },
		},
	}

	mockTodos := []*models.Todo{
		createTestTodo(1),
		createTestTodo(2),
	}

	// For each query method, run a set of test cases
	for _, queryTest := range queryTests {
		t.Run(queryTest.methodName, func(t *testing.T) {
			testCases := []struct {
				name      string
				mockTodos []*models.Todo
				mockError error
				wantError bool
				wantCount int
			}{
				{
					name:      "Success with todos",
					mockTodos: mockTodos,
					mockError: nil,
					wantError: false,
					wantCount: 2,
				},
				{
					name:      "Empty result",
					mockTodos: []*models.Todo{},
					mockError: nil,
					wantError: false,
					wantCount: 0,
				},
				{
					name:      "Repository error",
					mockTodos: nil,
					mockError: errors.New("query error"),
					wantError: true,
					wantCount: 0,
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Setup mock
					mockRepo := &MockTodoRepository{
						MockTodos: tc.mockTodos,
						MockError: tc.mockError,
					}

					// Create service
					svc := service.NewAppService(mockRepo)

					// Call the method being tested
					todos, err := queryTest.methodFunc(svc)

					// Check expectations
					if tc.wantError {
						if err == nil {
							t.Error("Expected error but got nil")
						}
					} else {
						if err != nil {
							t.Errorf("Expected no error but got: %v", err)
						}

						if len(todos) != tc.wantCount {
							t.Errorf("Expected %d todos, got %d", tc.wantCount, len(todos))
						}
					}
				})
			}
		})
	}
}

// Test tag methods
func TestTagMethods(t *testing.T) {
	// Test AddTagToTodo
	t.Run("AddTagToTodo", func(t *testing.T) {
		testCases := []struct {
			name      string
			todoID    int64
			tag       string
			mockError error
			wantError bool
		}{
			{
				name:      "Successfully add tag",
				todoID:    1,
				tag:       "important",
				mockError: nil,
				wantError: false,
			},
			{
				name:      "Error adding tag",
				todoID:    2,
				tag:       "error",
				mockError: errors.New("tag error"),
				wantError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				mockRepo := &MockTodoRepository{
					MockError: tc.mockError,
				}

				// Create service
				svc := service.NewAppService(mockRepo)

				// Call method
				err := svc.AddTagToTodo(tc.todoID, tc.tag)

				// Check expectations
				if tc.wantError {
					if err == nil {
						t.Error("Expected error but got nil")
					}
				} else {
					if err != nil {
						t.Errorf("Expected no error but got: %v", err)
					}

					// Verify tag was added
					if mockRepo.AddedTags == nil {
						t.Error("Expected tag to be added")
					} else {
						tags := mockRepo.AddedTags[tc.todoID]
						if len(tags) != 1 || tags[0] != tc.tag {
							t.Errorf("Expected tag %q to be added", tc.tag)
						}
					}
				}
			})
		}
	})

	// Test RemoveTagFromTodo
	t.Run("RemoveTagFromTodo", func(t *testing.T) {
		testCases := []struct {
			name      string
			todoID    int64
			tag       string
			mockError error
			wantError bool
		}{
			{
				name:      "Successfully remove tag",
				todoID:    1,
				tag:       "important",
				mockError: nil,
				wantError: false,
			},
			{
				name:      "Error removing tag",
				todoID:    2,
				tag:       "error",
				mockError: errors.New("tag error"),
				wantError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				mockRepo := &MockTodoRepository{
					MockError: tc.mockError,
				}

				// Create service
				svc := service.NewAppService(mockRepo)

				// Call method
				err := svc.RemoveTagFromTodo(tc.todoID, tc.tag)

				// Check expectations
				if tc.wantError {
					if err == nil {
						t.Error("Expected error but got nil")
					}
				} else {
					if err != nil {
						t.Errorf("Expected no error but got: %v", err)
					}

					// Verify tag was removed
					if mockRepo.RemovedTags == nil {
						t.Error("Expected tag to be removed")
					} else {
						tags := mockRepo.RemovedTags[tc.todoID]
						if len(tags) != 1 || tags[0] != tc.tag {
							t.Errorf("Expected tag %q to be removed", tc.tag)
						}
					}
				}
			})
		}
	})
}

// Test due date methods
func TestDueDateMethods(t *testing.T) {
	// Test SetDueDate
	t.Run("SetDueDate", func(t *testing.T) {
		dueDate := time.Now().Add(24 * time.Hour)

		testCases := []struct {
			name      string
			todoID    int64
			dueDate   time.Time
			mockTodo  *models.Todo
			mockError error
			wantError bool
		}{
			{
				name:      "Successfully set due date",
				todoID:    1,
				dueDate:   dueDate,
				mockTodo:  createTestTodo(1),
				mockError: nil,
				wantError: false,
			},
			{
				name:      "Error setting due date",
				todoID:    2,
				dueDate:   dueDate,
				mockTodo:  nil,
				mockError: errors.New("due date error"),
				wantError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				mockRepo := &MockTodoRepository{
					MockTodo:  tc.mockTodo,
					MockError: tc.mockError,
				}

				// Create service
				svc := service.NewAppService(mockRepo)

				// Call method
				err := svc.SetDueDate(tc.todoID, tc.dueDate)

				// Check expectations
				if tc.wantError {
					if err == nil {
						t.Error("Expected error but got nil")
					}
				} else {
					if err != nil {
						t.Errorf("Expected no error but got: %v", err)
					}

					// Verify due date was set
					if len(mockRepo.UpdatedTodos) != 1 {
						t.Error("Expected todo to be updated")
					} else {
						updatedTodo := mockRepo.UpdatedTodos[0]
						if updatedTodo.DueDate == nil {
							t.Error("Expected due date to be set")
						} else if !updatedTodo.DueDate.Equal(tc.dueDate) {
							t.Errorf("Expected due date %v, got %v", tc.dueDate, *updatedTodo.DueDate)
						}
					}
				}
			})
		}
	})

	// Test ClearDueDate
	t.Run("ClearDueDate", func(t *testing.T) {
		dueDate := time.Now().Add(24 * time.Hour)
		todo := createTestTodo(1)
		todo.DueDate = &dueDate

		testCases := []struct {
			name      string
			todoID    int64
			mockTodo  *models.Todo
			mockError error
			wantError bool
		}{
			{
				name:      "Successfully clear due date",
				todoID:    1,
				mockTodo:  todo,
				mockError: nil,
				wantError: false,
			},
			{
				name:      "Error clearing due date",
				todoID:    2,
				mockTodo:  nil,
				mockError: errors.New("due date error"),
				wantError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup mock
				mockRepo := &MockTodoRepository{
					MockTodo:  tc.mockTodo,
					MockError: tc.mockError,
				}

				// Create service
				svc := service.NewAppService(mockRepo)

				// Call method
				err := svc.ClearDueDate(tc.todoID)

				// Check expectations
				if tc.wantError {
					if err == nil {
						t.Error("Expected error but got nil")
					}
				} else {
					if err != nil {
						t.Errorf("Expected no error but got: %v", err)
					}

					// Verify due date was cleared
					if len(mockRepo.UpdatedTodos) != 1 {
						t.Error("Expected todo to be updated")
					} else {
						updatedTodo := mockRepo.UpdatedTodos[0]
						if updatedTodo.DueDate != nil {
							t.Errorf("Expected due date to be nil, got %v", *updatedTodo.DueDate)
						}
					}
				}
			})
		}
	})
}

// Test SetPriority functionality
func TestSetPriority(t *testing.T) {
	testCases := []struct {
		name      string
		todoID    int64
		priority  models.Priority
		mockTodo  *models.Todo
		mockError error
		wantError bool
	}{
		{
			name:      "Successfully set priority",
			todoID:    1,
			priority:  models.High,
			mockTodo:  createTestTodo(1),
			mockError: nil,
			wantError: false,
		},
		{
			name:      "Error setting priority",
			todoID:    2,
			priority:  models.Low,
			mockTodo:  nil,
			mockError: errors.New("priority error"),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			mockRepo := &MockTodoRepository{
				MockTodo:  tc.mockTodo,
				MockError: tc.mockError,
			}

			// Create service
			svc := service.NewAppService(mockRepo)

			// Call method
			err := svc.SetPriority(tc.todoID, tc.priority)

			// Check expectations
			if tc.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				// Verify priority was set
				if len(mockRepo.UpdatedTodos) != 1 {
					t.Error("Expected todo to be updated")
				} else {
					updatedTodo := mockRepo.UpdatedTodos[0]
					if updatedTodo.Priority != tc.priority {
						t.Errorf("Expected priority %v, got %v", tc.priority, updatedTodo.Priority)
					}
				}
			}
		})
	}
}

// Test sortTodos functionality (which is used by several methods)
func TestSortTodos(t *testing.T) {
	now := time.Now()
	older := now.Add(-1 * time.Hour)

	// Create unsorted todos with various priorities and timestamps
	unsortedTodos := []*models.Todo{
		{ID: 1, Title: "Low Priority Newer", Priority: models.Low, UpdatedAt: now},
		{ID: 2, Title: "High Priority Older", Priority: models.High, UpdatedAt: older},
		{ID: 3, Title: "Medium Priority Older", Priority: models.Medium, UpdatedAt: older},
		{ID: 4, Title: "Medium Priority Newer", Priority: models.Medium, UpdatedAt: now},
		{ID: 5, Title: "High Priority Newer", Priority: models.High, UpdatedAt: now},
	}

	mockRepo := &MockTodoRepository{
		MockTodos: unsortedTodos,
	}

	svc := service.NewAppService(mockRepo)

	// Call one of the methods that uses sortTodos internally
	todos, err := svc.GetAllTodos(false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for i := 0; i < len(todos)-1; i++ {
		current := todos[i]
		next := todos[i+1]

		// If priorities are different, higher priority should come first
		if current.Priority != next.Priority {
			if current.Priority < next.Priority {
				t.Errorf("At position %d, expected higher priority (%d) to come before lower priority (%d)",
					i, next.Priority, current.Priority)
			}
		} else {
			// If priorities are the same, newer update time should come first
			if current.UpdatedAt.Before(next.UpdatedAt) {
				t.Errorf("At position %d, expected newer update time to come before older update time", i)
			}
		}
	}
}

// Test archive methods
func TestArchiveMethods(t *testing.T) {
	// Test structure for both archive and unarchive operations
	archiveTests := []struct {
		methodName  string
		methodFunc  func(*service.AppService, int64) error
		setArchived bool // Expected archive status after operation
	}{
		{
			methodName:  "ArchiveTodo",
			methodFunc:  func(svc *service.AppService, id int64) error { return svc.ArchiveTodo(id) },
			setArchived: true,
		},
		{
			methodName:  "UnarchiveTodo",
			methodFunc:  func(svc *service.AppService, id int64) error { return svc.UnarchiveTodo(id) },
			setArchived: false,
		},
	}

	for _, archiveTest := range archiveTests {
		t.Run(archiveTest.methodName, func(t *testing.T) {
			testCases := []struct {
				name      string
				todoID    int64
				mockTodo  *models.Todo
				mockError error
				wantError bool
			}{
				{
					name:      "Success operation",
					todoID:    1,
					mockTodo:  createTestTodo(1),
					mockError: nil,
					wantError: false,
				},
				{
					name:      "Error fetching todo",
					todoID:    2,
					mockTodo:  nil,
					mockError: errors.New("fetch error"),
					wantError: true,
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// Setup mock
					mockRepo := &MockTodoRepository{
						MockTodo:  tc.mockTodo,
						MockError: tc.mockError,
					}

					// Create service
					svc := service.NewAppService(mockRepo)

					// Get time before update
					var beforeUpdate time.Time
					if tc.mockTodo != nil {
						beforeUpdate = tc.mockTodo.UpdatedAt
					}

					// Call the method being tested
					err := archiveTest.methodFunc(svc, tc.todoID)

					// Check expectations
					if tc.wantError {
						if err == nil {
							t.Error("Expected error but got nil")
						}
					} else {
						if err != nil {
							t.Errorf("Expected no error but got: %v", err)
						}

						// Verify todo was updated
						if len(mockRepo.UpdatedTodos) != 1 {
							t.Errorf("Expected todo to be updated")
						} else {
							updatedTodo := mockRepo.UpdatedTodos[0]

							// Check archived was changed
							if updatedTodo.Archived != archiveTest.setArchived {
								t.Errorf("Expected archived status %v, got %v",
									archiveTest.setArchived, updatedTodo.Archived)
							}

							// Check updated_at timestamp was changed
							if !updatedTodo.UpdatedAt.After(beforeUpdate) {
								t.Errorf("Expected updated_at to be updated")
							}
						}
					}
				})
			}
		})
	}
}

// Test AdvanceStatus functionality
func TestAdvanceStatus(t *testing.T) {
	testCases := []struct {
		name           string
		todoID         int64
		initialStatus  models.Status
		expectedStatus models.Status
		mockError      error
		wantError      bool
	}{
		{
			name:           "Advance from Open to Doing",
			todoID:         1,
			initialStatus:  models.Open,
			expectedStatus: models.Doing,
			mockError:      nil,
			wantError:      false,
		},
		{
			name:           "Advance from Doing to Done",
			todoID:         2,
			initialStatus:  models.Doing,
			expectedStatus: models.Done,
			mockError:      nil,
			wantError:      false,
		},
		{
			name:           "Advance from Done error",
			todoID:         3,
			initialStatus:  models.Done,
			expectedStatus: models.Open,
			mockError:      errors.New("update error"),
			wantError:      true,
		},
		{
			name:           "Error fetching todo",
			todoID:         4,
			initialStatus:  models.Open,
			expectedStatus: models.Doing,
			mockError:      errors.New("fetch error"),
			wantError:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a todo with the initial status
			todo := createTestTodo(tc.todoID)
			todo.Status = tc.initialStatus

			// Setup mock
			mockRepo := &MockTodoRepository{
				MockTodo:  todo,
				MockError: tc.mockError,
			}

			// Create service
			svc := service.NewAppService(mockRepo)

			// Call method
			newStatus, err := svc.AdvanceStatus(tc.todoID)

			// Check expectations
			if tc.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				if newStatus != tc.expectedStatus {
					t.Errorf("Expected status %v, got %v", tc.expectedStatus, newStatus)
				}

				// Verify the todo was updated by checking the underlying method calls
				// For Open->Doing, should have called MarkAsDoing
				// For Doing->Done, should have called MarkAsDone
				// For Done->Open, should have called MarkAsOpen
				if len(mockRepo.UpdatedTodos) == 0 {
					t.Error("Expected todo to be updated")
				} else {
					updatedTodo := mockRepo.UpdatedTodos[0]
					if updatedTodo.Status != tc.expectedStatus {
						t.Errorf("Expected todo status to be %v, got %v",
							tc.expectedStatus, updatedTodo.Status)
					}
				}
			}
		})
	}
}

// Test GetFilteredTodos functionality
func TestGetFilteredTodos(t *testing.T) {
	// First, define ViewType enum to match service implementation
	testCases := []struct {
		name         string
		viewType     service.ViewType
		showArchived bool
		mockTodos    []*models.Todo
		mockError    error
		wantError    bool
		expectedRepo string // Which repository method should be called
	}{
		{
			name:         "Get Open Todos",
			viewType:     service.OpenPane,
			showArchived: false,
			mockTodos:    []*models.Todo{createTestTodo(1), createTestTodo(2)},
			mockError:    nil,
			wantError:    false,
			expectedRepo: "GetOpen",
		},
		{
			name:         "Get Doing Todos",
			viewType:     service.DoingPane,
			showArchived: false,
			mockTodos:    []*models.Todo{createTestTodo(3)},
			mockError:    nil,
			wantError:    false,
			expectedRepo: "GetActive",
		},
		{
			name:         "Get Done Todos",
			viewType:     service.DonePane,
			showArchived: false,
			mockTodos:    []*models.Todo{createTestTodo(4), createTestTodo(5)},
			mockError:    nil,
			wantError:    false,
			expectedRepo: "GetCompleted",
		},
		{
			name:         "Get All Todos without archived",
			viewType:     service.AllPane,
			showArchived: false,
			mockTodos:    []*models.Todo{createTestTodo(6), createTestTodo(7)},
			mockError:    nil,
			wantError:    false,
			expectedRepo: "GetAll",
		},
		{
			name:         "Get All Todos with archived",
			viewType:     service.AllPane,
			showArchived: true,
			mockTodos:    []*models.Todo{createTestTodo(8)},
			mockError:    nil,
			wantError:    false,
			expectedRepo: "GetAll",
		},
		{
			name:         "Error in repository",
			viewType:     service.OpenPane,
			showArchived: false,
			mockTodos:    nil,
			mockError:    errors.New("repository error"),
			wantError:    true,
			expectedRepo: "GetOpen",
		},
		{
			name:         "Invalid view type",
			viewType:     service.ViewType(999), // Invalid value
			showArchived: false,
			mockTodos:    nil,
			mockError:    nil,
			wantError:    true,
			expectedRepo: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := &MockTodoRepository{
				MockTodos: tc.mockTodos,
				MockError: tc.mockError,
			}

			// Create service
			svc := service.NewAppService(mockRepo)

			// Call method
			todos, err := svc.GetFilteredTodos(service.ViewType(tc.viewType), tc.showArchived)

			// Check expectations
			if tc.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				// Verify we got the expected todos
				if len(todos) != len(tc.mockTodos) {
					t.Errorf("Expected %d todos, got %d", len(tc.mockTodos), len(todos))
				}

				// Additional verification could check that the correct repository method was called
				// This would require enhancing the mock to track method calls
			}
		})
	}
}

func TestUpdateInfoMethods(t *testing.T) {
	// Test SetUpdateInfo and GetUpdateInfo
	t.Run("SetUpdateInfo_GetUpdateInfo", func(t *testing.T) {
		// Setup
		mockRepo := &MockTodoRepository{}
		svc := service.NewAppService(mockRepo)

		version := "1.2.3"
		url := "https://example.com/releases"
		notes := "Bug fixes and performance improvements"
		forceUpdate := true
		hasUpdate := true

		// Execute
		svc.SetUpdateInfo(version, url, notes, forceUpdate, hasUpdate)
		updateInfo := svc.GetUpdateInfo()

		// Verify
		if updateInfo == nil {
			t.Fatal("Expected update info to be set, got nil")
		}

		if updateInfo.Version != version {
			t.Errorf("Expected version %s, got %s", version, updateInfo.Version)
		}

		if updateInfo.URL != url {
			t.Errorf("Expected URL %s, got %s", url, updateInfo.URL)
		}

		if updateInfo.Notes != notes {
			t.Errorf("Expected release notes %s, got %s", notes, updateInfo.Notes)
		}

		if updateInfo.ForceUpdate != forceUpdate {
			t.Errorf("Expected forceUpdate %v, got %v", forceUpdate, updateInfo.ForceUpdate)
		}

		if updateInfo.HasUpdate != hasUpdate {
			t.Errorf("Expected hasUpdate %v, got %v", hasUpdate, updateInfo.HasUpdate)
		}

		// CheckedAt should be recent
		timeDiff := time.Since(updateInfo.CheckedAt)
		if timeDiff > 1*time.Second {
			t.Errorf("CheckedAt time should be recent, but was %v ago", timeDiff)
		}
	})

	// Test HasUpdate method
	t.Run("HasUpdate", func(t *testing.T) {
		testCases := []struct {
			name      string
			hasUpdate bool
			setupInfo bool
			expected  bool
		}{
			{
				name:      "Has update true",
				hasUpdate: true,
				setupInfo: true,
				expected:  true,
			},
			{
				name:      "Has update false",
				hasUpdate: false,
				setupInfo: true,
				expected:  false,
			},
			{
				name:      "No update info",
				hasUpdate: false,
				setupInfo: false,
				expected:  false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				mockRepo := &MockTodoRepository{}
				svc := service.NewAppService(mockRepo)

				if tc.setupInfo {
					svc.SetUpdateInfo("1.0.0", "url", "notes", false, tc.hasUpdate)
				}

				// Execute
				result := svc.HasUpdate()

				// Verify
				if result != tc.expected {
					t.Errorf("Expected HasUpdate() to return %v, got %v", tc.expected, result)
				}
			})
		}
	})

	// Test NeedsForceUpdate method
	t.Run("NeedsForceUpdate", func(t *testing.T) {
		testCases := []struct {
			name        string
			forceUpdate bool
			setupInfo   bool
			expected    bool
		}{
			{
				name:        "Force update true",
				forceUpdate: true,
				setupInfo:   true,
				expected:    true,
			},
			{
				name:        "Force update false",
				forceUpdate: false,
				setupInfo:   true,
				expected:    false,
			},
			{
				name:        "No update info",
				forceUpdate: false,
				setupInfo:   false,
				expected:    false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				mockRepo := &MockTodoRepository{}
				svc := service.NewAppService(mockRepo)

				if tc.setupInfo {
					svc.SetUpdateInfo("1.0.0", "url", "notes", tc.forceUpdate, true)
				}

				// Execute
				result := svc.NeedsForceUpdate()

				// Verify
				if result != tc.expected {
					t.Errorf("Expected NeedsForceUpdate() to return %v, got %v", tc.expected, result)
				}
			})
		}
	})

	// Test initialization
	t.Run("UpdateInfo_Initialization", func(t *testing.T) {
		// Setup
		mockRepo := &MockTodoRepository{}
		svc := service.NewAppService(mockRepo)

		// Verify update info is initialized to empty but not nil
		updateInfo := svc.GetUpdateInfo()
		if updateInfo == nil {
			t.Fatal("Expected updateInfo to be initialized as non-nil")
		}

		// Verify default values
		if updateInfo.HasUpdate || updateInfo.ForceUpdate {
			t.Errorf("Expected new updateInfo to have HasUpdate and ForceUpdate set to false")
		}
	})
}

func TestGetTotalTimeSpent(t *testing.T) {
	tests := []struct {
		name         string
		todoID       int64
		mockTodo     *models.Todo
		mockError    error
		wantError    bool
		wantDuration time.Duration
	}{
		{
			name:   "successful retrieval with time spent",
			todoID: 1,
			mockTodo: &models.Todo{
				ID:          1,
				TimeSpent:   3600, // 1 hour in seconds
				TimeStarted: nil,  // Not currently tracking
			},
			mockError:    nil,
			wantError:    false,
			wantDuration: 3600 * time.Second,
		},
		{
			name:   "active tracking adds elapsed time",
			todoID: 2,
			mockTodo: func() *models.Todo {
				startTime := time.Now().Add(-10 * time.Minute)
				return &models.Todo{
					ID:          2,
					TimeSpent:   1800, // 30 minutes in seconds
					Status:      models.Doing,
					TimeStarted: &startTime,
				}
			}(),
			mockError:    nil,
			wantError:    false,
			wantDuration: 1800*time.Second + 10*time.Minute, // Base + ~10 minutes tracking
		},
		{
			name:         "repository error",
			todoID:       3,
			mockTodo:     nil,
			mockError:    errors.New("database error"),
			wantError:    true,
			wantDuration: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For the active tracking test, we need to calculate the expected time dynamically
			var expectedDuration time.Duration = tt.wantDuration
			if tt.mockTodo != nil && tt.mockTodo.TimeStarted != nil && tt.mockTodo.Status == models.Doing {
				baseTime := time.Duration(tt.mockTodo.TimeSpent) * time.Second
				elapsed := time.Since(*tt.mockTodo.TimeStarted)
				expectedDuration = baseTime + elapsed
			}

			mockRepo := &MockTodoRepository{
				MockTodo:  tt.mockTodo,
				MockError: tt.mockError,
			}

			service := service.NewAppService(mockRepo)
			duration, err := service.GetTotalTimeSpent(tt.todoID)

			// Check error
			if (err != nil) != tt.wantError {
				t.Errorf("GetTotalTimeSpent() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// For active tracking, we can't predict the exact elapsed time
			if tt.mockTodo != nil && tt.mockTodo.TimeStarted != nil && tt.mockTodo.Status == models.Doing {
				// Allow a small tolerance of 2 seconds for test execution time
				tolerance := 2 * time.Second
				diff := duration - expectedDuration
				if diff < -tolerance || diff > tolerance {
					t.Errorf("GetTotalTimeSpent() = %v, want approximately %v (diff: %v)",
						duration, expectedDuration, diff)
				}
			} else if duration != tt.wantDuration {
				t.Errorf("GetTotalTimeSpent() = %v, want %v", duration, tt.wantDuration)
			}
		})
	}
}

func TestTimeTrackingMethods(t *testing.T) {
	testCases := []struct {
		name              string
		methodName        string
		todoID            int64
		initialTodo       *models.Todo
		mockError         error
		wantError         bool
		checkUpdatedState func(*testing.T, *models.Todo)
	}{
		{
			name:       "pause tracking - success",
			methodName: "PauseTimeTracking",
			todoID:     1,
			initialTodo: func() *models.Todo {
				startTime := time.Now().Add(-20 * time.Minute)
				return &models.Todo{
					ID:          1,
					Title:       "Test Todo",
					TimeSpent:   1800, // 30 minutes
					Status:      models.Doing,
					TimeStarted: &startTime,
				}
			}(),
			mockError: nil,
			wantError: false,
			checkUpdatedState: func(t *testing.T, todo *models.Todo) {
				if todo.Status != models.Doing {
					t.Errorf("Todo status should remain %v, got %v", models.Doing, todo.Status)
				}
				if todo.TimeStarted != nil {
					t.Errorf("TimeStarted should be nil after pausing, got %v", todo.TimeStarted)
				}
				// Time spent should have increased by approximately 20 minutes
				expectedMin := 1800 + 20*60 - 5 // Allow 5 seconds tolerance
				expectedMax := 1800 + 20*60 + 5
				if todo.TimeSpent < int64(expectedMin) || todo.TimeSpent > int64(expectedMax) {
					t.Errorf("TimeSpent = %v, want approximately %v", todo.TimeSpent, 1800+20*60)
				}
			},
		},
		{
			name:       "pause tracking - not tracking",
			methodName: "PauseTimeTracking",
			todoID:     2,
			initialTodo: &models.Todo{
				ID:          2,
				Title:       "Not Tracking Todo",
				TimeSpent:   600,
				Status:      models.Doing,
				TimeStarted: nil,
			},
			mockError: nil,
			wantError: false,
			checkUpdatedState: func(t *testing.T, todo *models.Todo) {
				if todo.TimeStarted != nil {
					t.Errorf("TimeStarted should remain nil")
				}
				if todo.TimeSpent != 600 {
					t.Errorf("TimeSpent should not change, got %v, want %v", todo.TimeSpent, 600)
				}
			},
		},
		{
			name:        "pause tracking - error",
			methodName:  "PauseTimeTracking",
			todoID:      3,
			initialTodo: nil,
			mockError:   errors.New("database error"),
			wantError:   true,
			checkUpdatedState: func(t *testing.T, todo *models.Todo) {
				// No state to check with error
			},
		},
		{
			name:       "resume tracking - success",
			methodName: "ResumeTimeTracking",
			todoID:     4,
			initialTodo: &models.Todo{
				ID:          4,
				Title:       "Paused Todo",
				TimeSpent:   900,
				Status:      models.Doing,
				TimeStarted: nil,
			},
			mockError: nil,
			wantError: false,
			checkUpdatedState: func(t *testing.T, todo *models.Todo) {
				if todo.TimeStarted == nil {
					t.Errorf("TimeStarted should be set after resuming")
				} else {
					// Start time should be set to approximately now
					now := time.Now()
					startTime := *todo.TimeStarted
					diff := now.Sub(startTime).Seconds()
					if diff < -2 || diff > 2 { // Allow 2 seconds tolerance
						t.Errorf("TimeStarted = %v, want approximately %v (diff: %.2f seconds)",
							startTime, now, diff)
					}
				}
				if todo.TimeSpent != 900 {
					t.Errorf("TimeSpent should remain unchanged, got %v, want %v", todo.TimeSpent, 900)
				}
			},
		},
		{
			name:       "resume tracking - already tracking",
			methodName: "ResumeTimeTracking",
			todoID:     5,
			initialTodo: func() *models.Todo {
				startTime := time.Now().Add(-5 * time.Minute)
				return &models.Todo{
					ID:          5,
					Title:       "Already Tracking Todo",
					TimeSpent:   300,
					Status:      models.Doing,
					TimeStarted: &startTime,
				}
			}(),
			mockError: nil,
			wantError: false,
			checkUpdatedState: func(t *testing.T, todo *models.Todo) {
				if todo.TimeStarted == nil {
					t.Errorf("TimeStarted should remain set")
				}
				// State should remain unchanged
				if todo.TimeSpent != 300 {
					t.Errorf("TimeSpent should not change, got %v, want %v", todo.TimeSpent, 300)
				}
			},
		},
		{
			name:       "reset tracking - success",
			methodName: "ResetTimeTracking",
			todoID:     6,
			initialTodo: func() *models.Todo {
				startTime := time.Now().Add(-30 * time.Minute)
				return &models.Todo{
					ID:          6,
					Title:       "Track Reset Todo",
					TimeSpent:   7200,
					Status:      models.Doing,
					TimeStarted: &startTime,
				}
			}(),
			mockError: nil,
			wantError: false,
			checkUpdatedState: func(t *testing.T, todo *models.Todo) {
				if todo.TimeStarted != nil {
					t.Errorf("TimeStarted should be nil after reset, got %v", todo.TimeStarted)
				}
				if todo.TimeSpent != 0 {
					t.Errorf("TimeSpent should be reset to 0, got %v", todo.TimeSpent)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a copy of the initial todo to avoid modifying the test data
			var initialTodoCopy *models.Todo
			if tc.initialTodo != nil {
				todoCopy := *tc.initialTodo
				initialTodoCopy = &todoCopy
			}

			mockRepo := &MockTodoRepository{
				MockTodo:  initialTodoCopy,
				MockError: tc.mockError,
			}

			service := service.NewAppService(mockRepo)

			var err error
			switch tc.methodName {
			case "PauseTimeTracking":
				err = service.PauseTimeTracking(tc.todoID)
			case "ResumeTimeTracking":
				err = service.ResumeTimeTracking(tc.todoID)
			case "ResetTimeTracking":
				err = service.ResetTimeTracking(tc.todoID)
			default:
				t.Fatalf("Unknown method name: %s", tc.methodName)
			}

			// Check error
			if (err != nil) != tc.wantError {
				t.Errorf("%s() error = %v, wantError %v", tc.methodName, err, tc.wantError)
				return
			}

			// If no error and we updated a todo, check the state
			if err == nil && tc.initialTodo != nil {
				var updatedTodo *models.Todo
				if len(mockRepo.UpdatedTodos) > 0 {
					updatedTodo = mockRepo.UpdatedTodos[0]
				} else {
					updatedTodo = tc.initialTodo
				}
				tc.checkUpdatedState(t, updatedTodo)
			}
		})
	}
}

func TestTimeTrackingProperties(t *testing.T) {
	t.Run("pause_tracking_properties", rapid.MakeCheck(func(t *rapid.T) {
		// Generate a todo with random time spent
		initialTimeSpent := rapid.Int64Range(0, 36000).Draw(t, "initialTimeSpent")

		// Random start time between 1 minute and 2 hours ago
		minMinutesAgo := rapid.IntRange(1, 120).Draw(t, "minutesAgo")
		startTime := time.Now().Add(-time.Duration(minMinutesAgo) * time.Minute)

		initialTodo := &models.Todo{
			ID:          1,
			TimeSpent:   initialTimeSpent,
			Status:      models.Doing,
			TimeStarted: &startTime,
		}

		// Create a copy for testing
		todoCopy := *initialTodo
		mockRepo := &MockTodoRepository{
			MockTodo: &todoCopy,
		}

		service := service.NewAppService(mockRepo)
		err := service.PauseTimeTracking(1)

		// Properties to verify:
		// 1. No error should occur
		if err != nil {
			t.Fatalf("PauseTimeTracking() unexpected error: %v", err)
		}

		// 2. Todo should have updated state
		if len(mockRepo.UpdatedTodos) == 0 {
			t.Fatalf("Expected todo to be updated")
		}

		updatedTodo := mockRepo.UpdatedTodos[0]

		// 3. TimeStarted should be nil
		if updatedTodo.TimeStarted != nil {
			t.Errorf("TimeStarted should be nil after pausing, got %v", updatedTodo.TimeStarted)
		}

		// 4. Time spent should increase by approximately elapsed time
		expectedElapsed := int64(time.Since(startTime).Seconds())
		expectedTimeSpent := initialTimeSpent + expectedElapsed

		// Allow 2 seconds tolerance for test execution time
		if updatedTodo.TimeSpent < expectedTimeSpent-2 || updatedTodo.TimeSpent > expectedTimeSpent+2 {
			t.Errorf("TimeSpent = %v, expected approximately %v (initial: %v + elapsed: %v)",
				updatedTodo.TimeSpent, expectedTimeSpent, initialTimeSpent, expectedElapsed)
		}
	}))
}
