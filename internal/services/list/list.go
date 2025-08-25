package list

import (
	"context"
	"fmt"

	"github.com/anubhav047/goboard/internal/db"
)

// Service handles list-related business logic
type Service struct {
	queries *db.Queries
}

// New creates a new list service
func New(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// CreateList creates a new list in a board
func (s *Service) CreateList(ctx context.Context, name string, boardID int32, position int32) (*db.List, error) {
	// Validate input
	if name == "" {
		return nil, fmt.Errorf("fist name cannot be empty")
	}

	// Create the list
	list, err := s.queries.CreateList(ctx, db.CreateListParams{
		Name:     name,
		BoardID:  boardID,
		Position: position,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	return &list, nil
}

// GetBoardLists gets all lists for a specified board
func (s *Service) GetBoardLists(ctx context.Context, boardID int32) ([]db.List, error) {
	lists, err := s.queries.GetListsByBoard(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed tp get board lists: %w", err)
	}

	// Ensure we return an empty slice instead of nil
	if lists == nil {
		return []db.List{}, nil
	}

	return lists, nil
}

// GetListByID gets a single list by ID
func (s *Service) GetListByID(ctx context.Context, listID int32) (*db.List, error) {
	list, err := s.queries.GetListByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	return &list, nil
}

// UpdateList updates a list's name and position
func (s *Service) UpdateList(ctx context.Context, listID int32, name string, position int32) (*db.List, error) {
	// Validate input
	if name == "" {
		return nil, fmt.Errorf("list name cannot be empty")
	}

	list, err := s.queries.UpdateList(ctx, db.UpdateListParams{
		Name:     name,
		Position: position,
		ID:       listID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	return &list, nil
}

// DeleteList deletes a list
func (s *Service) DeleteList(ctx context.Context, listID int32) error {
	err := s.queries.DeleteList(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to delete list: %w", err)
	}

	return nil
}
