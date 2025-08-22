package user

import (
	"context"
	"fmt"

	"github.com/anubhav047/goboard/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// Service handles the business logic for users
type Service struct {
	queries *db.Queries
}

// New creates a new user Service.
func New(queries *db.Queries) *Service {
	return &Service{
		queries: queries,
	}
}

// Register creates a new user, hashes their password, and saves it to the database.
func (s *Service) Register(ctx context.Context, name, email, password string) (db.User, error) {
	//Hash the pw using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user in the database
	params := db.CreateUserParams{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	user, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		// TODO : check for db errors like duplicate email.
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
