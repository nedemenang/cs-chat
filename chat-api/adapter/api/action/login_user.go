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

type LoginUserAction struct {
	uc        usecase.LoginUserUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewLoginUserAction(uc usecase.LoginUserUseCase, log logger.Logger, v validator.Validator) LoginUserAction {
	return LoginUserAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (a LoginUserAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "login_user"

	var input usecase.LoginUserInput
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
			http.StatusBadRequest,
		).Log("error when logging in")

		// response.NewError(err, http.StatusInternalServerError).Send(w)
		response.NewError("internal_server_error", http.StatusBadRequest, err, "").Send(w)

		return
	}
	logging.NewInfo(a.log, logKey, http.StatusOK).Log("login successful")

	response.NewSuccess(output, http.StatusOK).Send(w)
}

func (a LoginUserAction) validateInput(input usecase.LoginUserInput) error {
	err := a.validator.Validate(input)
	if err != nil {
		return errors.New(strings.Join(a.validator.Messages(), ","))
	}
	return nil

}
