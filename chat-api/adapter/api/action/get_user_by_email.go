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

type GetUserByEmailAction struct {
	uc        usecase.GetUserByEmailUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewGetUserByEmailAction(uc usecase.GetUserByEmailUseCase, log logger.Logger, v validator.Validator) GetUserByEmailAction {
	return GetUserByEmailAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (a GetUserByEmailAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "get_user_by_email"

	var email = r.URL.Query().Get("email")

	input := usecase.GetUserByEmailInput{
		Email: email,
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("error when fetching user")

		// response.NewError(err, http.StatusInternalServerError).Send(w)
		response.NewError("internal_server_error", http.StatusInternalServerError, err, "").Send(w)

		return
	}
	logging.NewInfo(a.log, logKey, http.StatusOK).Log("successfully fetched user")

	response.NewSuccess(output, http.StatusOK).Send(w)
}

func (a GetUserByEmailAction) validateInput(input usecase.GetChannelByIdInput) error {
	err := a.validator.Validate(input)
	if err != nil {
		return errors.New(strings.Join(a.validator.Messages(), ","))
	}
	return nil

}
