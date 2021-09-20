package usecase

import (
	"context"
	"time"

	"chat-api/domain"
)

type (
	// Input port
	CreateMessageUseCase interface {
		Execute(context.Context, CreateMessageInput) (CreateMessageOutput, error)
	}

	CreateMessageInput struct {
		ChannelId   string `json:"channelId" validate:"required"`
		MessageFrom string `json:"messageFrom" validate:"required"`
		Message     string `json:"message" validate:"required"`
	}

	// Output port
	CreateMessagePresenter interface {
		Output(domain.Channel) CreateMessageOutput
	}

	MessageOutput struct {
		MessageFrom string    `json:"messageFrom"`
		Message     string    `json:"message"`
		Timestamp   time.Time `json:"timestamp"`
	}

	// Output data
	CreateMessageOutput struct {
		Id            string          `json:"id"`
		UserEmail     string          `json:"userEmail"`
		RepEmail      string          `json:"repEmail"`
		CurrentStatus string          `json:"currentStatus"`
		Messages      []MessageOutput `json:"messages"`
		CreatedAt     time.Time       `json:"createdAt"`
	}

	createMessageInteractor struct {
		repo       domain.ChannelRepository
		presenter  CreateMessagePresenter
		ctxTimeout time.Duration
	}
)

func NewCreateMessageInteractor(
	repo domain.ChannelRepository,
	presenter CreateMessagePresenter,
	t time.Duration,
) CreateMessageUseCase {
	return createMessageInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (c createMessageInteractor) Execute(ctx context.Context, input CreateMessageInput) (CreateMessageOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	channel, err := c.repo.GetChannelById(ctx, input.ChannelId)
	if err != nil {
		return c.presenter.Output(domain.Channel{}), err
	}

	channel.AddMessage(input.MessageFrom, input.Message, time.Now())

	err = c.repo.AddMessage(ctx, channel)
	if err != nil {
		return c.presenter.Output(domain.Channel{}), err
	}

	return c.presenter.Output(channel), nil
}
