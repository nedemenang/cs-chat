package action

import (
	"fmt"
	"net/http"

	"chat-api/adapter/api/logging"
	"chat-api/adapter/api/response"
	"chat-api/adapter/logger"
	"chat-api/adapter/validator"
	"chat-api/domain"
	"chat-api/usecase"

	"go.mongodb.org/mongo-driver/bson"
)

type GetChannelsByQueryAction struct {
	uc        usecase.GetChannelByQueryUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewGetChannelsByQueryAction(uc usecase.GetChannelByQueryUseCase, log logger.Logger) GetChannelsByQueryAction {
	return GetChannelsByQueryAction{
		uc:  uc,
		log: log,
	}
}

func (a GetChannelsByQueryAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "get_channels_by_query"

	query := bson.M{}

	var (
		limit string
		page  string
	)

	if len(r.URL.Query().Get("userEmail")) > 0 {
		query["userEmail"] = fmt.Sprintf("%s", r.URL.Query().Get("userEmail"))
	}

	if len(r.URL.Query().Get("repEmail")) > 0 {
		query["repEmail"] = fmt.Sprintf("%s", r.URL.Query().Get("repEmail"))
	}

	if len(r.URL.Query().Get("currentStatus")) > 0 {
		query["currentStatus"] = fmt.Sprintf("%s", r.URL.Query().Get("currentStatus"))
	}

	page = r.URL.Query().Get("page")
	limit = r.URL.Query().Get("limit")

	output, err := a.uc.Execute(r.Context(), query, limit, page)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			logging.NewError(
				a.log,
				err,
				logKey,
				http.StatusBadRequest,
			).Log("error fetching channels")

			response.NewError("condition_error", http.StatusPreconditionFailed, err, "").Send(w)
			return
		default:
			logging.NewError(
				a.log,
				err,
				logKey,
				http.StatusInternalServerError,
			).Log("error when returning channels")

			response.NewError("internal_server_error", http.StatusInternalServerError, err, "").Send(w)
			return
		}
	}
	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success returning channels")

	response.NewSuccess(output, http.StatusOK).Send(w)
}

// func (a FindAllUserAction) validateInput(input usecase.FindRelationshipChildrenInput) []string {
// 	var msgs []string

// 	err := a.validator.Validate(input)
// 	if err != nil {
// 		for _, msg := range a.validator.Messages() {
// 			msgs = append(msgs, msg)
// 		}
// 	}

// 	return msgs
// }
