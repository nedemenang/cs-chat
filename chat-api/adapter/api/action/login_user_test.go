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

type mockLoginUser struct {
	result usecase.LoginUserOutput
	err    error
}

func (m mockLoginUser) Execute(_ context.Context, _ usecase.LoginUserInput) (usecase.LoginUserOutput, error) {
	return m.result, m.err
}

func TestLoginUserAction_Execute(t *testing.T) {
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
		ucMock             usecase.LoginUserUseCase
		expectedBody       string
		expectedStatusCode int
	}{
		{
			name: "Login user success",
			args: args{
				rawPayload: []byte(`{
					"email": "user_email@gmail.com",
					"password": "supersecurepassword"
				}`),
			},
			ucMock: mockLoginUser{
				result: usecase.LoginUserOutput{
					FirstName: "John",
					LastName:  "James",
					Email:     "user_email@gmail.com",
					Role:      "USER",
					Token:     "03nf0394jf0394rfj0394f0394f0ghy094gh039240954jf093",
				},
				err: nil,
			},
			expectedBody:       `{"firstName":"John","lastName":"James","email":"user_email@gmail.com","role":"USER","token":"03nf0394jf0394rfj0394f0394f0ghy094gh039240954jf093"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "error invalid json",
			args: args{
				rawPayload: []byte(`{
					"email":
					"password": "supersecurepassword",
							}`),
			},
			ucMock: mockLoginUser{
				result: usecase.LoginUserOutput{},
				err:    fmt.Errorf("invalid json"),
			},
			expectedBody:       `{"errors":[{"code":400,"message":"invalid character ':' after object key:value pair","type":"input_error"}]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "error invalid fields",
			args: args{
				rawPayload: []byte(`{
					"emails": "user_email@gmail.com",
					"password": "supersecurepassword"
					}`),
			},
			ucMock: mockLoginUser{
				result: usecase.LoginUserOutput{},
				err:    fmt.Errorf("invalid json"),
			},
			expectedBody:       `{"errors":[{"code":400,"message":"Email is a required field","type":"input_error"}]}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodPost,
				"/user/login",
				bytes.NewReader(tt.args.rawPayload),
			)

			var (
				w      = httptest.NewRecorder()
				action = NewLoginUserAction(tt.ucMock, log.LoggerMock{}, validator)
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
