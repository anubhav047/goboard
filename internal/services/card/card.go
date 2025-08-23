package card

import (
	"context"
	"fmt"

	"github.com/anubhav047/goboard/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service handles card-related business logic
type Service struct {
	queries *db.Queries
}

// New creates a new card service
func New(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// CreateCard creates a new card in a list
func (s *Service) CreateCard(ctx context.Context, title, description string, listID int32, position int32) (*db.Card, error) {
	// Validate input
	if title == "" {
		return nil, fmt.Errorf("card title cannot be empty")
	}

	// Create the card
	card, err := s.queries.CreateCard(ctx, db.CreateCardParams{
		Title:       title,
		Description: pgtype.Text{String: description, Valid: true},
		ListID:      listID,
		Position:    position,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	return &card, nil
}

// GetListCards gets all cards for a specific list
func (s *Service) GetListCards(ctx context.Context, listID int32) ([]db.Card, error) {
	cards, err := s.queries.GetCardsByList(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list cards: %w", err)
	}

	return cards, nil
}

// GetCardByID gets a single card by ID
func (s *Service) GetCardByID(ctx context.Context, cardID int32) (*db.Card, error) {
	card, err := s.queries.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	return &card, nil
}

// UpdateCard updates a card's title and description
func (s *Service) UpdateCard(ctx context.Context, cardID int32, title, description string) (*db.Card, error) {
	// Validate input
	if title == "" {
		return nil, fmt.Errorf("card title cannot be empty")
	}

	card, err := s.queries.UpdateCard(ctx, db.UpdateCardParams{
		Title:       title,
		Description: pgtype.Text{String: description, Valid: true},
		ID:          cardID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	return &card, nil
}

// MoveCard moves a card to a different list and/or position
func (s *Service) MoveCard(ctx context.Context, cardID, listID, position int32) (*db.Card, error) {
	card, err := s.queries.MoveCard(ctx, db.MoveCardParams{
		ListID:   listID,
		Position: position,
		ID:       cardID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to move card: %w", err)
	}

	return &card, nil
}

// DeleteCard deletes a card
func (s *Service) DeleteCard(ctx context.Context, cardID int32) error {
	err := s.queries.DeleteCard(ctx, cardID)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return nil
}
