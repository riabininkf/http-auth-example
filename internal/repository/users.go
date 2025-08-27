package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/riabininkf/http-auth-example/internal/domain"
)

func NewUsers(conn Conn) *Users {
	return &Users{
		conn: conn,
	}
}

type Users struct {
	conn Conn
}

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

func (u *Users) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	query := `UPDATE public.users SET password = $1 WHERE id = $2`
	if _, err := u.conn.Exec(ctx, query, hashedPassword, userID); err != nil {
		return err
	}

	return nil
}
