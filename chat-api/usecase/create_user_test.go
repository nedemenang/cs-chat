package usecase

import (
	"chat-api/domain"
	"context"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	newUserId = primitive.NewObjectID()
)

type mockAuthenticationService struct {
	domain.AuthenticationUtilityService

	result string
	err    error
}

func (m mockAuthenticationService) HashPassword(_ context.Context, _ string) (string, error) {
	return m.result, m.err
}

type mockCreateUserRepoStore struct {
	domain.UserRepository

	result domain.User
	err    error
}

func (m mockCreateUserRepoStore) Create(_ context.Context, _ domain.User) (domain.User, error) {
	return m.result, m.err
}

type mockCreateUserRepo struct {
	domain.UserRepository

	createUserFake func() (domain.User, error)
	invokedCreate  *invoked

	getUserByEmailFake func() (domain.User, error)
	invokedFind        *invoked
}

func (m mockCreateUserRepo) CreateUser(_ context.Context, _ domain.User) (domain.User, error) {

	if m.invokedCreate != nil {
		m.invokedCreate.call = true
	}
	return m.createUserFake()
}

func (m mockCreateUserRepo) GetUserByEmail(_ context.Context, _ string) (domain.User, error) {

	if m.invokedFind != nil {
		m.invokedFind.call = true
	}
	return m.getUserByEmailFake()
}

type mockCreateUserPresenter struct {
	result CreateUserOutput
}

func (m mockCreateUserPresenter) Output(_ domain.User) CreateUserOutput {
	return m.result
}

func TestCreateUserInteractor_Execute(t *testing.T) {
	t.Parallel()
	createdTime := time.Now()
	type args struct {
		input CreateUserInput
	}

	tests := []struct {
		name          string
		args          args
		userRepo      domain.UserRepository
		service       domain.AuthenticationUtilityService
		presenter     CreateUserPresenter
		expected      CreateUserOutput
		expectedError string
	}{
		{
			name: "create a user successful",
			args: args{input: CreateUserInput{Email: "newEmail@email.com", Password: "password", FirstName: "firstName", LastName: "lastName", Role: "user"}},
			userRepo: mockCreateUserRepo{createUserFake: func() (domain.User, error) {
				return domain.NewUser(newUserId, "firstName", "lastName", "newEmail@email.com", "password", createdTime, createdTime), nil
			}, getUserByEmailFake: func() (domain.User, error) {
				return domain.User{}, nil
			}},
			service:   mockAuthenticationService{result: "03jr04jf03jlkjfeo3nflp23049tfj30"},
			presenter: mockCreateUserPresenter{result: CreateUserOutput{FirstName: "firstName", LastName: "lastName", Email: "newEmail@email.com"}},
			expected:  CreateUserOutput{FirstName: "firstName", LastName: "lastName", Email: "newEmail@email.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewCreateUserInteractor(tt.userRepo, tt.service, tt.presenter, time.Second)

			got, err := uc.Execute(context.Background(), tt.args.input)
			if (err != nil) && (err.Error() != tt.expectedError) {
				t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, got, tt.expected)
			}
		})
	}
}
