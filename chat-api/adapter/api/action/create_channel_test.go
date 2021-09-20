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

type mockCreateChannel struct {
	result usecase.CreateChannelOutput
	err    error
}

func (m mockCreateChannel) Execute(_ context.Context, _ usecase.CreateChannelInput) (usecase.CreateChannelOutput, error) {
	return m.result, m.err
}

func TestCreateChannelAction_Execute(t *testing.T) {
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
		ucMock             usecase.CreateChannelUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "AddUserIdentityAction success",
			args: args{
				rawPayload: []byte(`{
					"userEmail": "user_email@gmail.com"
				}`),
			},
			ucMock: mockCreateChannel{
				result: usecase.CreateChannelOutput{
					Id:      "3c096a40-ccba-4b58-93ed-57379ab04679",
					UserEmail:   "user_email@gmail.com",
					CurrentStatus:  "ACTIVE",
				},
				err: nil,
			},
			expectedBody:       `{"id":"3c096a40-ccba-4b58-93ed-57379ab04679","userEmail":"user_email@gmail.com","currentStatus":"ACTIVE"}`,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "error invalid json",
			args: args{
				rawPayload: []byte(`{
							"userEmails":"user_email@gmail.com"
							}`),
			},
			ucMock: mockCreateChannel{
				result: usecase.CreateChannelOutput{},
				err:    fmt.Errorf("invalid json"),
			},
			expectedBody:       `{"errors":[{"code":400,"message":"UserEmail is a required field","type":"input_error"}]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "error invalid fields",
			args: args{
				rawPayload: []byte(`{
					"userEmail": ""
					}`),
			},
			ucMock: mockCreateChannel{
				result: usecase.CreateChannelOutput{},
				err:    fmt.Errorf("invalid json"),
			},
			expectedBody:       `{"errors":[{"code":400,"message":"userEmail is a required field","type":"input_error"}]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/channel",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewCreateChannelAction(tt.ucMock, log.LoggerMock{}, validator)
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
