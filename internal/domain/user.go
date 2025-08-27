package domain

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

func (u *user) ID() string {
	return u.id
}

func (u *user) Email() string {
	return u.email
}

func (u *user) HashedPassword() string {
	return u.hashedPassword
}
