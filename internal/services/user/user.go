package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/anubhav047/goboard/internal/db"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

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

// Logim verifies a user's credentials and return the user if they are valid
func (s *Service) Login(ctx context.Context, email, password string) (db.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		// If the user is not found, pgx returns a special 'ErrNoRows'.
		// We check for this and return our custom, generic error.
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrInvalidCredentials
		}
		// For any other database error, we return a generic failure.
		return db.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	// compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return db.User{}, ErrInvalidCredentials
	}

	return user, nil
}
