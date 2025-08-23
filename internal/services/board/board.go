package board

import (
	"context"
	"fmt"

	"github.com/anubhav047/goboard/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service handles board-related business logic
type Service struct {
	queries *db.Queries
}

// New creates a new board service
func New(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// CreateBoard creates a new board for a user
func (s *Service) CreateBoard(ctx context.Context, name, description string, userID int32) (*db.Board, error) {
	// Validate Input
	if name == "" {
		return nil, fmt.Errorf("board name cannot be empty")
	}

	// Create a board
	board, err := s.queries.CreateBoard(ctx, db.CreateBoardParams{
		Name:        name,
		Description: pgtype.Text{String: description, Valid: true},
		CreatedBy:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	return &board, nil
}

// GetUserBoards gets all boards for a specific user
func (s *Service) GetUserBoards(ctx context.Context, userID int32) ([]db.Board, error) {
	boards, err := s.queries.GetBoardsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fet user boards; %w", err)
	}

	return boards, nil
}

// GetBoardById gets a sinle board by ID
func (s *Service) GetBoardByID(ctx context.Context, boardID int32) (*db.Board, error) {
	board, err := s.queries.GetBoardByID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	return &board, nil
}

// UpdateBoard updates a board's name and description
func (s *Service) UpdateBoard(ctx context.Context, boardID int32, name, description string) (*db.Board, error) {
	// Validate input
	if name == "" {
		return nil, fmt.Errorf("board name cannot be empty")
	}

	board, err := s.queries.UpdateBoard(ctx, db.UpdateBoardParams{
		Name:        name,
		Description: pgtype.Text{String: description, Valid: true},
		ID:          boardID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update board: %w", err)
	}

	return &board, nil
}

// DeleteBoard deletes a board
func (s *Service) DeleteBoard(ctx context.Context, boardID int32) error {
	err := s.queries.DeleteBoard(ctx, boardID)
	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	return nil
}
