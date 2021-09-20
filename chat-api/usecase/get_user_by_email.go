package usecase

import (
	"chat-api/domain"
	"context"
	"time"
)

type (
	// Input port
	GetUserByEmailUseCase interface {
		Execute(context.Context, GetUserByEmailInput) (GetUserByEmailOutput, error)
	}

	GetUserByEmailPresenter interface {
		Output(domain.User) GetUserByEmailOutput
	}

	GetUserByEmailInput struct {
		Email string `json:"email" validate:"required"`
	}

	GetUserByEmailOutput struct {
		Id        string    `json:"_id"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"createdAt"`
		Email     string    `json:"email"`
	}

	getUserByEmailInteractor struct {
		repo       domain.UserRepository
		presenter  GetUserByEmailPresenter
		ctxTimeout time.Duration
	}
)

func NewUserByEmailInteractor(
	repo domain.UserRepository,
	presenter GetUserByEmailPresenter,
	t time.Duration,
) GetUserByEmailUseCase {
	return getUserByEmailInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (a getUserByEmailInteractor) Execute(ctx context.Context, input GetUserByEmailInput) (GetUserByEmailOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()
	user, err := a.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return a.presenter.Output(domain.User{}), err
	}

	return a.presenter.Output(user), nil
}
