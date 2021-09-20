package usecase

import (
	"context"
	"time"

	"chat-api/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// Input port
	CreateChannelUseCase interface {
		Execute(context.Context, CreateChannelInput) (CreateChannelOutput, error)
	}

	CreateChannelInput struct {
		UserEmail string `json:"userEmail" validate:"required"`
	}

	// Output port
	CreateChannelPresenter interface {
		Output(domain.Channel) CreateChannelOutput
	}

	// Output data
	CreateChannelOutput struct {
		Id            string `json:"id"`
		UserEmail     string `json:"userEmail"`
		CurrentStatus string `json:"currentStatus"`
		// CreatedAt     time.Time `json:"createdAt"`
	}

	createChannelInteractor struct {
		repo       domain.ChannelRepository
		presenter  CreateChannelPresenter
		ctxTimeout time.Duration
	}
)

func NewCreateChannelInteractor(
	repo domain.ChannelRepository,
	presenter CreateChannelPresenter,
	t time.Duration,
) CreateChannelUseCase {
	return createChannelInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (c createChannelInteractor) Execute(ctx context.Context, input CreateChannelInput) (CreateChannelOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	channel := domain.NewChannel(
		primitive.NewObjectID(),
		input.UserEmail,
		domain.ACTIVE,
		time.Now(),
		time.Now(),
	)

	channel.UpdateStatus(domain.ACTIVE, input.UserEmail, time.Now().Unix())

	createdChannel, err := c.repo.CreateChannel(ctx, channel)
	if err != nil {
		return c.presenter.Output(domain.Channel{}), err
	}

	return c.presenter.Output(createdChannel), nil
}
