package usecase

import (
	"chat-api/domain"
	"context"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type (
	// Input port
	GetChannelByQueryUseCase interface {
		Execute(context.Context, interface{}, string, string) (GetChannelByQueryOutput, error)
	}

	GetChannelByQueryPresenter interface {
		Output([]domain.Channel, int, int, int, int) GetChannelByQueryOutput
	}

	ChannelByQueryOutput struct {
		Id            string    `json:"id"`
		UserEmail     string    `json:"userEmail"`
		RepEmail      string    `json:"repEmail"`
		CurrentStatus string    `json:"currentStatus"`
		UserFullName  string    `json:"userFullName"`
		CreatedAt     time.Time `json:"createdAt"`
	}

	GetChannelByQueryOutput struct {
		Page       int                    `json:"page"`
		Count      int                    `json:"count"`
		Limit      int                    `json:"limit"`
		TotalCount int                    `json:"totalCount"`
		Data       []ChannelByQueryOutput `json:"data"`
	}

	getChannelByQueryInteractor struct {
		repo       domain.ChannelRepository
		userRepo   domain.UserRepository
		presenter  GetChannelByQueryPresenter
		ctxTimeout time.Duration
	}
)

func NewGetChannelByQueryInteractor(
	repo domain.ChannelRepository,
	userRepo domain.UserRepository,
	presenter GetChannelByQueryPresenter,
	t time.Duration,
) GetChannelByQueryUseCase {
	return getChannelByQueryInteractor{
		repo:       repo,
		userRepo:   userRepo,
		presenter:  presenter,
		ctxTimeout: t,
	}
}

// Execute orchestrates the use case
func (a getChannelByQueryInteractor) Execute(ctx context.Context, query interface{}, limit, page string) (GetChannelByQueryOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	var intLimit int
	var intPage int

	if page != "" {
		parsedPage, err := strconv.Atoi(page)
		if err != nil {
			intPage = 1
		}
		intPage = parsedPage
	} else {
		intPage = 1
	}

	if limit != "" {
		parsedLim, err := strconv.Atoi(limit)
		if err != nil {
			parsedLim = 10
		}

		if parsedLim > 50 {
			intLimit = 50
		} else {
			intLimit = parsedLim
		}
	} else {
		intLimit = 10
	}

	start := (intPage - 1) * intLimit

	channels, err := a.repo.GetChannelsByQuery(ctx, query, start, intLimit)
	if err != nil {
		return a.presenter.Output([]domain.Channel{}, 0, 0, 0, 0), err
	}

	for i, channel := range channels {
		u, _ := a.userRepo.GetUserByEmail(ctx, channel.UserEmail())
		channels[i].UpdateUserFullName(fmt.Sprintf("%s %s", u.FirstName(), u.LastName()))
	}

	i := bson.M{}
	channelsCount, err := a.repo.GetChannelsByQueryCount(ctx, i)
	if err != nil {
		return a.presenter.Output([]domain.Channel{}, 0, 0, 0, 0), err
	}

	return a.presenter.Output(channels, intPage, intLimit, len(channels), int(channelsCount)), nil
}
