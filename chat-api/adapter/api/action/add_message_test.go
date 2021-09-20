package action

import (
	"bytes"
	"chat-api/infrastructure/log"
	"chat-api/infrastructure/validation"
	"chat-api/usecase"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockAddMessage struct {
	result usecase.CreateMessageOutput
	err    error
}

func (m mockAddMessage) Execute(_ context.Context, _ usecase.CreateMessageInput) (usecase.CreateMessageOutput, error) {
	return m.result, m.err
}

func TestAddUserIdentityAction_Execute(t *testing.T) {
	t.Parallel()

	validator, _ := validation.NewValidatorFactory(validation.InstanceGoPlayground)
	// createdAt := time.Now().UTC()
	// createdAtString := createdAt.String()

	type args struct {
		rawPayload []byte
	}

	tests := []struct {
		name               string
		args               args
		ucMock             usecase.CreateMessageUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "error invalid json",
			args: args{
				rawPayload: []byte(`{
							"channelId": "30fj3094rjf0t934059rjf3094",
							"messageFrom": ,
							"message":
							}`),
			},
			ucMock: mockAddMessage{
				result: usecase.CreateMessageOutput{},
				err:    fmt.Errorf("invalid json"),
			},
			expectedBody:       `{"errors":[{"code":400,"message":"invalid character ',' looking for beginning of value","type":"input_error"}]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "error invalid fields",
			args: args{
				rawPayload: []byte(`{
					"channelId": "30fj3094rjf0t934059rjf3094",
					"messageFrom": "message_from@gmail.com",
					"message123": "message"
					}`),
			},
			ucMock: mockAddMessage{
				result: usecase.CreateMessageOutput{},
				err:    fmt.Errorf("invalid json"),
			},
			expectedBody:       `{"errors":[{"code":400,"message":"Message is a required field","type":"input_error"}]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/message",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewCreateMessageAction(tt.ucMock, log.LoggerMock{}, validator)
			)

			action.Execute(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf(
					"[TestCase '%s'] HTTP handler returned wrong statusCode: recieved '%v' expected '%v'",
					tt.name,
					w.Code,
					tt.expectedStatusCode,
				)
			}

			var result = strings.TrimSpace(w.Body.String())
			if !strings.EqualFold(result, tt.expectedBody) {
				t.Errorf(
					"[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					result,
					tt.expectedBody,
				)
			}
		})
	}
}
