package usecase

import (
	"chat-api/domain"
	"context"
	"time"
)

type (
	// Input port
	GetChannelByIdUseCase interface {
		Execute(context.Context, GetChannelByIdInput) (GetChannelByIdOutput, error)
	}

	GetChannelByIdPresenter interface {
		Output(domain.Channel) GetChannelByIdOutput
	}

	GetChannelByIdInput struct {
		Id string `json:"id" validate:"required"`
	}

	Message struct {
		Message     string    `json:"message"`
		MessageFrom string    `json:"messageFrom"`
		Timestamp   time.Time `json:"timeStamp"`
	}

	GetChannelByIdOutput struct {
		Id            string    `json:"id"`
		UserEmail     string    `json:"userEmail"`
		RepEmail      string    `json:"repEmail"`
		CurrentStatus string    `json:"currentStatus"`
		CreatedAt     time.Time `json:"createdAt"`
		Messages      []Message `json:"messages"`
	}

	getChannelByIdInteractor struct {
		repo       domain.ChannelRepository
		presenter  GetChannelByIdPresenter
		ctxTimeout time.Duration
	}
)

// NewFindUserByIdInteractor creates new finduserByIdInteractor with its dependencies
func NewGetChannelByIdInteractor(
	repo domain.ChannelRepository,
	presenter GetChannelByIdPresenter,
	t time.Duration,
) GetChannelByIdUseCase {
	return getChannelByIdInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (a getChannelByIdInteractor) Execute(ctx context.Context, input GetChannelByIdInput) (GetChannelByIdOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	channel, err := a.repo.GetChannelById(ctx, input.Id)
	if err != nil {
		return a.presenter.Output(domain.Channel{}), err
	}

	return a.presenter.Output(channel), nil
}
