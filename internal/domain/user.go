package domain

import "errors"

var (
	// ErrEmailBusy is returned when the email is already in use.
	ErrEmailBusy = errors.New("email is busy")

	// ErrUserNotFound is returned when the user is not found.
	ErrUserNotFound = errors.New("user not found")
)

// NewUser creates a new User instance with the provided id, email, and hashed password.
func NewUser(
	id string,
	email string,
	hashedPassword string,
) User {
	return &user{
		id:             id,
		email:          email,
		hashedPassword: hashedPassword,
	}
}

type (
	// User represents a user, providing methods to access ID, email, and hashed password.
	User interface {
		ID() string
		Email() string
		HashedPassword() string
	}

	user struct {
		id             string
		email          string
		hashedPassword string
	}
)

// ID returns the unique identifier of the user as a string.
func (u *user) ID() string {
	return u.id
}

// Email returns the email address associated with the user.
func (u *user) Email() string {
	return u.email
}

// HashedPassword returns the hashed password of the user as a string.
func (u *user) HashedPassword() string {
	return u.hashedPassword
}
