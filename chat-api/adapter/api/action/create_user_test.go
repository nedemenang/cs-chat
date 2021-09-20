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

type mockCreateUser struct {
	result usecase.CreateUserOutput
	err    error
}

func (m mockCreateUser) Execute(_ context.Context, _ usecase.CreateUserInput) (usecase.CreateUserOutput, error) {
	return m.result, m.err
}

func TestCreateUserAction_Execute(t *testing.T) {
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
		ucMock             usecase.CreateUserUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "create user success",
			args: args{
				rawPayload: []byte(`{
					"firstName": "John",
					"lastName": "James",
					"email": "user_email@gmail.com",
					"password": "supersecurepassword",
					"role": "USER"
				}`),
			},
			ucMock: mockCreateUser{
				result: usecase.CreateUserOutput{
					FirstName: "John",
					LastName:  "James",
					Email:     "user_email@gmail.com",
				},
				err: nil,
			},
			expectedBody:       `{"firstName":"John","lastName":"James","email":"user_email@gmail.com"}`,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "error invalid json",
			args: args{
				rawPayload: []byte(`{
					"firstName": "John",
					"last": "James",
					"email": "user_email@gmail.com",
					"password": "supersecurepassword",
					"role": "USER"
							}`),
			},
			ucMock: mockCreateUser{
				result: usecase.CreateUserOutput{},
				err:    fmt.Errorf("invalid json"),
			},
			expectedBody:       `{"errors":[{"code":400,"message":"lastName is a required field","type":"input_error"}]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "error invalid fields",
			args: args{
				rawPayload: []byte(`{
					"firstName": "",
					"last": "James",
					"email": "user_email@gmail.com",
					"password": "supersecurepassword",
					"role": "USER"
					}`),
			},
			ucMock: mockCreateUser{
				result: usecase.CreateUserOutput{},
				err:    fmt.Errorf("invalid json"),
			},
			expectedBody:       `{"errors":[{"code":400,"message":"FirstName is a required field,LastName is a required field","type":"input_error"}]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/user",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewCreateUserAction(tt.ucMock, log.LoggerMock{}, validator)
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
