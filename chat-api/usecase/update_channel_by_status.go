package usecase

import (
	"context"
	"time"

	"chat-api/domain"
)

type (
	// Input port
	UpdateChannelStatusUseCase interface {
		Execute(context.Context, UpdateChannelStatusInput) (UpdateChannelStatusOutput, error)
	}

	// Input data
	UpdateChannelStatusInput struct {
		ID        string `json:"id" validate:"required"`
		UpdatedBy string `json:"updatedBy" validate:"required"`
		Status    string `json:"status" validate:"required"`
	}

	// Output port
	UpdateChannelStatusPresenter interface {
		Output(domain.Channel) UpdateChannelStatusOutput
	}

	// Output data
	UpdateChannelStatusOutput struct {
		RepEmail      string    `json:"repEmail"`
		UserEmail     string    `json:"userEmail"`
		CurrentStatus string    `json:"currentStatus"`
		CreatedAt     time.Time `json:"createdAt"`
	}

	updateChannelStatusInteractor struct {
		repo       domain.ChannelRepository
		presenter  UpdateChannelStatusPresenter
		ctxTimeout time.Duration
	}
)

// NewCreateUserInteractor creates new createUserInteractor with its dependencies
func NewUpdateChannelStatusInteractor(
	repo domain.ChannelRepository,
	presenter UpdateChannelStatusPresenter,
	t time.Duration,
) UpdateChannelStatusUseCase {
	return updateChannelStatusInteractor{
		repo:       repo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (a updateChannelStatusInteractor) Execute(ctx context.Context, input UpdateChannelStatusInput) (UpdateChannelStatusOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	channel, err := a.repo.GetChannelById(ctx, input.ID)
	if err != nil {
		switch err {
		case domain.ChannelNotFound:
			return UpdateChannelStatusOutput{}, domain.ErrUserNotFound
		default:
			return UpdateChannelStatusOutput{}, err
		}
	}

	if input.Status == domain.IN_PROGRESS {
		channel.UpdateRepEmail(input.UpdatedBy)
	}
	channel.UpdateStatus(input.Status, input.UpdatedBy, time.Now().Unix())

	err = a.repo.UpdateChannelStatus(ctx, channel)
	if err != nil {
		return a.presenter.Output(domain.Channel{}), err
	}

	return a.presenter.Output(channel), nil
}
