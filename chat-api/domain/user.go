package domain

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ADMIN = "ADMIN"
	USER  = "USER"
)

var (
	ErrUserNotFound                = errors.New("user not found")
	ErrUsernameOrPasswordIncorrect = errors.New("username or password incorrect")
)

type (
	AuthenticationUtilityService interface {
		HashPassword(context.Context, string) (string, error)
		CheckPasswordHash(context.Context, string, string) bool
		GenerateToken(context.Context, User) (string, error)
	}

	UserRepository interface {
		CreateUser(context.Context, User) (User, error)
		GetUserByEmail(context.Context, string) (User, error)
	}

	User struct {
		id        primitive.ObjectID
		firstName string
		lastName  string
		email     string
		password  string
		role      string
		createdAt time.Time
		updatedAt time.Time
	}
)

func NewUser(id primitive.ObjectID, firstName, lastName, email, password string, createdAt, updatedAt time.Time) User {
	return User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		password:  password,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) UpdateRole(role string) {
	u.role = role
}

func (u User) FirstName() string {
	return u.firstName
}

func (u User) Password() string {
	return u.password
}

func (u User) LastName() string {
	return u.lastName
}

func (u User) Role() string {
	return u.role
}

func (u User) Email() string {
	return u.email
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u User) Id() primitive.ObjectID {
	return u.id
}
