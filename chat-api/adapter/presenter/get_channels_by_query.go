package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
)

type getChannelsByQueryPresenter struct{}

func NewGetChannelsByQueryPresenter() usecase.GetChannelByQueryPresenter {
	return getChannelsByQueryPresenter{}
}

func (a getChannelsByQueryPresenter) Output(channels []domain.Channel, page, limit, queryChannelCount, totalChannelCount int) usecase.GetChannelByQueryOutput {

	var channelList = make([]usecase.ChannelByQueryOutput, 0)

	for _, channel := range channels {

		channelList = append(channelList, usecase.ChannelByQueryOutput{
			CurrentStatus: channel.CurrentStatus(),
			Id:            channel.Id().Hex(),
			CreatedAt:     channel.CreatedAt(),
			UserEmail:     channel.UserEmail(),
			UserFullName:  channel.UserFullName(),
			RepEmail:      channel.RepEmail(),
		})

	}

	channelOutput := usecase.GetChannelByQueryOutput{
		Page:       page,
		Count:      queryChannelCount,
		Limit:      limit,
		TotalCount: totalChannelCount,
		Data:       channelList,
	}
	return channelOutput
}
