package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
)

type updateChannelStatusPresenter struct{}

func NewUpdateChannelStatusPresenter() usecase.UpdateChannelStatusPresenter {
	return updateChannelStatusPresenter{}
}

func (a updateChannelStatusPresenter) Output(channel domain.Channel) usecase.UpdateChannelStatusOutput {
	return usecase.UpdateChannelStatusOutput{
		UserEmail:     channel.UserEmail(),
		CurrentStatus: channel.CurrentStatus(),
		RepEmail:      channel.RepEmail(),
		CreatedAt:     channel.CreatedAt(),
	}
}
