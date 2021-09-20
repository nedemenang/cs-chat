package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
)

type createMessagePresenter struct{}

func NewCreateMessagePresenter() usecase.CreateMessagePresenter {
	return createMessagePresenter{}
}

func (a createMessagePresenter) Output(channel domain.Channel) usecase.CreateMessageOutput {

	messages := make([]usecase.MessageOutput, 0)
	for _, message := range channel.Messages() {
		messages = append(messages, usecase.MessageOutput{
			MessageFrom: message.MessageFrom,
			Message:     message.Message,
			Timestamp:   message.Timestamp,
		})
	}

	return usecase.CreateMessageOutput{

		Id:            channel.Id().Hex(),
		UserEmail:     channel.UserEmail(),
		CurrentStatus: channel.CurrentStatus(),
		CreatedAt:     channel.CreatedAt(),
		Messages:      messages,
	}
}
