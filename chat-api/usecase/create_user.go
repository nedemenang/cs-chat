package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"chat-api/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// Input port
	CreateUserUseCase interface {
		Execute(context.Context, CreateUserInput) (CreateUserOutput, error)
	}

	CreateUserInput struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Email     string `json:"email" validate:"required"`
		Password  string `json:"password" validate:"required"`
		Role      string `json:"role" validate:"required"`
	}

	// Output port
	CreateUserPresenter interface {
		Output(domain.User) CreateUserOutput
	}

	// Output data
	CreateUserOutput struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Email     string `json:"email" validate:"required"`
	}

	createUserInteractor struct {
		repo       domain.UserRepository
		service    domain.AuthenticationUtilityService
		presenter  CreateUserPresenter
		ctxTimeout time.Duration
	}
)

func NewCreateUserInteractor(
	repo domain.UserRepository,
	service domain.AuthenticationUtilityService,
	presenter CreateUserPresenter,
	t time.Duration,
) CreateUserUseCase {
	return createUserInteractor{
		repo:       repo,
		service:    service,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (c createUserInteractor) Execute(ctx context.Context, input CreateUserInput) (CreateUserOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	existingUser, _ := c.repo.GetUserByEmail(ctx, input.Email)
	// if err != nil {
	// 	return c.presenter.Output(domain.User{}), err
	// }

	if existingUser.Email() != "" {
		return c.presenter.Output(domain.User{}), errors.New(fmt.Sprintf("User already exists with email"))
	}

	hashedPassword, err := c.service.HashPassword(ctx, input.Password)
	if err != nil {
		return c.presenter.Output(domain.User{}), err
	}

	user := domain.NewUser(
		primitive.NewObjectID(),
		input.FirstName,
		input.LastName,
		input.Email,
		hashedPassword,
		time.Now(),
		time.Now(),
	)
	user.UpdateRole(input.Role)

	createdUser, err := c.repo.CreateUser(ctx, user)
	if err != nil {
		return c.presenter.Output(domain.User{}), err
	}

	return c.presenter.Output(createdUser), nil
}
