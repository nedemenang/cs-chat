package action

import (
	"errors"
	"net/http"
	"strings"

	"chat-api/adapter/api/logging"
	"chat-api/adapter/api/response"
	"chat-api/adapter/logger"
	"chat-api/adapter/validator"
	"chat-api/usecase"
)

type GetChannelByIdAction struct {
	uc        usecase.GetChannelByIdUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewGetChannelByIdAction(uc usecase.GetChannelByIdUseCase, log logger.Logger, v validator.Validator) GetChannelByIdAction {
	return GetChannelByIdAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (a GetChannelByIdAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "get_channel_by_id"

	var channelID = r.URL.Query().Get("id")

	input := usecase.GetChannelByIdInput{
		Id: channelID,
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when fetching channel")

		// response.NewError(err, http.StatusInternalServerError).Send(w)
		response.NewError("internal_server_error", http.StatusInternalServerError, err, "").Send(w)

		return
	}
	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success fetching channel")

	response.NewSuccess(output, http.StatusOK).Send(w)
}

func (a GetChannelByIdAction) validateInput(input usecase.GetChannelByIdInput) error {
	err := a.validator.Validate(input)
	if err != nil {
		return errors.New(strings.Join(a.validator.Messages(), ","))
	}
	return nil

}
