package action

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"chat-api/adapter/api/logging"
	"chat-api/adapter/api/response"
	"chat-api/adapter/logger"
	"chat-api/adapter/validator"
	"chat-api/usecase"
)

type CreateChannelAction struct {
	uc        usecase.CreateChannelUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewCreateChannelAction(uc usecase.CreateChannelUseCase, log logger.Logger, v validator.Validator) CreateChannelAction {
	return CreateChannelAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (a CreateChannelAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "add_channel"

	var input usecase.CreateChannelInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log("error when decoding json")

		response.NewError("input_error", http.StatusBadRequest, err, "").Send(w)
		return
	}
	defer r.Body.Close()

	if err := a.validateInput(input); err != nil {
		logging.NewError(
			a.log,
			response.ErrInvalidInput,
			logKey,
			http.StatusBadRequest,
		).Log("invalid input")

		response.NewError("input_error", http.StatusBadRequest, err, "").Send(w)
		return
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when creating channel")

		// response.NewError(err, http.StatusInternalServerError).Send(w)
		response.NewError("internal_server_error", http.StatusInternalServerError, err, "").Send(w)

		return
	}
	logging.NewInfo(a.log, logKey, http.StatusCreated).Log("success creating channel")

	response.NewSuccess(output, http.StatusCreated).Send(w)
}

func (a CreateChannelAction) validateInput(input usecase.CreateChannelInput) error {
	err := a.validator.Validate(input)
	if err != nil {
		return errors.New(strings.Join(a.validator.Messages(), ","))
	}
	return nil

}
