package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
)

type getChannelByIdPresenter struct{}

func NewGetChannelByIdPresenter() usecase.GetChannelByIdPresenter {
	return getChannelByIdPresenter{}
}

func (a getChannelByIdPresenter) Output(channel domain.Channel) usecase.GetChannelByIdOutput {

	messages := make([]usecase.Message, 0)
	for _, message := range channel.Messages() {
		messages = append(messages, usecase.Message{
			MessageFrom: message.MessageFrom,
			Message:     message.Message,
			Timestamp:   message.Timestamp,
		})
	}

	return usecase.GetChannelByIdOutput{

		Id:            channel.Id().Hex(),
		UserEmail:     channel.UserEmail(),
		RepEmail:      channel.RepEmail(),
		CurrentStatus: channel.CurrentStatus(),
		CreatedAt:     channel.CreatedAt(),
		Messages:      messages,
	}
}
