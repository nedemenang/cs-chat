package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
)

type createChannelPresenter struct{}

func NewCreateChannelPresenter() usecase.CreateChannelPresenter {
	return createChannelPresenter{}
}

func (a createChannelPresenter) Output(channel domain.Channel) usecase.CreateChannelOutput {
	return usecase.CreateChannelOutput{
		Id:            channel.Id().Hex(),
		UserEmail:     channel.UserEmail(),
		CurrentStatus: channel.CurrentStatus(),
		// CreatedAt:     channel.CreatedAt(),
	}
}
