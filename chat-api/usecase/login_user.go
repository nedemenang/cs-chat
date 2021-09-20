package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"chat-api/domain"
)

type (
	// Input port
	LoginUserUseCase interface {
		Execute(context.Context, LoginUserInput) (LoginUserOutput, error)
	}

	LoginUserInput struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// Output port
	LoginUserPresenter interface {
		Output(domain.User, string) LoginUserOutput
	}

	// Output data
	LoginUserOutput struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		Token     string `json:"token"`
	}

	loginUserInteractor struct {
		repo       domain.UserRepository
		service    domain.AuthenticationUtilityService
		presenter  LoginUserPresenter
		ctxTimeout time.Duration
	}
)

func NewLoginUserInteractor(
	repo domain.UserRepository,
	service domain.AuthenticationUtilityService,
	presenter LoginUserPresenter,
	t time.Duration,
) LoginUserUseCase {
	return loginUserInteractor{
		repo:       repo,
		service:    service,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (l loginUserInteractor) Execute(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, l.ctxTimeout)
	defer cancel()

	existingUser, err := l.repo.GetUserByEmail(ctx, input.Email)

	if existingUser.Email() == "" {
		return l.presenter.Output(domain.User{}, ""), errors.New(fmt.Sprintf("Username or password incorrect"))
	}

	isPasswordCorrect := l.service.CheckPasswordHash(ctx, input.Password, existingUser.Password())
	if err != nil {
		return l.presenter.Output(domain.User{}, ""), errors.New(fmt.Sprintf("Username or password incorrect"))
	}

	if !isPasswordCorrect {
		return l.presenter.Output(domain.User{}, ""), errors.New(fmt.Sprintf("Username or password incorrect"))
	}

	token, err := l.service.GenerateToken(ctx, existingUser)
	if err != nil {
		return l.presenter.Output(domain.User{}, ""), err
	}

	return l.presenter.Output(existingUser, token), nil
}
