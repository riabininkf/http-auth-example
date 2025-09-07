package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/riabininkf/http-auth-example/internal/domain"
)

// NewUsers creates a new instance of Users using the provided Conn interface for database operations.
func NewUsers(conn Conn) *Users {
	return &Users{
		conn: conn,
	}
}

// Users provides methods to interact with the users table in the database. It uses Conn for database operations.
type Users struct {
	conn Conn
}

// Save inserts a new user into the database. Returns ErrEmailBusy if the email is already in use or any other errors.
func (u *Users) Save(ctx context.Context, user domain.User) error {
	query := `INSERT INTO public.users (id, email, password) VALUES ($1, $2, $3)`

	if _, err := u.conn.Exec(ctx, query, user.ID(), user.Email(), user.HashedPassword()); err != nil {
		if isUniqueConstraintViolation(err) {
			return domain.ErrEmailBusy
		}

		return err
	}

	return nil
}

// GetByEmail retrieves a user by their email address from the database. Returns a User and error if applicable.
func (u *Users) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, password FROM public.users WHERE email = $1`

	var (
		userID         string
		hashedPassword string
	)
	if err := u.conn.QueryRow(ctx, query, email).Scan(&userID, &hashedPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return domain.NewUser(
		userID,
		email,
		hashedPassword,
	), nil
}

// GetByID retrieves a user by their unique identifier from the database, returning a domain.User or an error.
func (u *Users) GetByID(ctx context.Context, userID string) (domain.User, error) {
	query := `SELECT email, password FROM public.users WHERE id = $1`

	var (
		email          string
		hashedPassword string
	)
	if err := u.conn.QueryRow(ctx, query, userID).Scan(&email, &hashedPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return domain.NewUser(
		userID,
		email,
		hashedPassword,
	), nil
}

// UpdatePassword updates the password of a user identified by userID with the provided hashed password in the database.
func (u *Users) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	query := `UPDATE public.users SET password = $1 WHERE id = $2`
	if _, err := u.conn.Exec(ctx, query, hashedPassword, userID); err != nil {
		return err
	}

	return nil
}
