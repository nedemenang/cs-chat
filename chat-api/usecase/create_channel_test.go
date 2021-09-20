package usecase

import (
	"chat-api/domain"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

type mockCreateChannelRepoStore struct {
	domain.ChannelRepository

	result domain.Channel
	err    error
}

func (m *mockCreateChannelRepoStore) Create(ctx context.Context, channel domain.Channel) (domain.Channel, error) {
	return m.result, m.err
}

type mockCreateChannelRepo struct {
	domain.ChannelRepository

	createChannelFake func() (domain.Channel, error)
	invokedCreate     *invoked

	findByIDFake func() (domain.Channel, error)
	invokedFind  *invoked
}

func (m mockCreateChannelRepo) CreateChannel(_ context.Context, _ domain.Channel) (domain.Channel, error) {
	if m.invokedCreate != nil {
		m.invokedCreate.call = true
	}
	return m.createChannelFake()
}

func (m mockCreateChannelRepo) GetChannelById(_ context.Context, _ string) (domain.Channel, error) {

	if m.invokedFind != nil {
		m.invokedFind.call = true
	}
	return m.findByIDFake()
}

type mockCreateChannelPresenter struct {
	result CreateChannelOutput
}

func (m mockCreateChannelPresenter) Output(_ domain.Channel) CreateChannelOutput {
	return m.result
}

func TestCreateChannelInteractor_Execute(t *testing.T) {
	t.Parallel()
	createdTime := time.Now()
	type args struct {
		input CreateChannelInput
	}

	tests := []struct {
		name          string
		args          args
		channelRepo   domain.ChannelRepository
		presenter     CreateChannelPresenter
		expected      CreateChannelOutput
		expectedError string
	}{
		{
			name: "create channel successfully",
			args: args{input: CreateChannelInput{
				UserEmail: "validemail@gmail.com",
			}},
			channelRepo: mockCreateChannelRepo{
				createChannelFake: func() (domain.Channel, error) {
					return domain.NewChannel(
						newChannelId,
						"testemail@email.com",
						domain.ACTIVE,
						createdTime,
						createdTime,
					), nil
				},
				findByIDFake: func() (domain.Channel, error) {
					return domain.NewChannel(
						newChannelId,
						"testemail@email.com",
						domain.ACTIVE,
						createdTime,
						createdTime,
					), nil
				},
			},
			presenter: mockCreateChannelPresenter{result: CreateChannelOutput{
				Id:            newChannelId.Hex(),
				UserEmail:     "testemail@email.com",
				CurrentStatus: domain.ACTIVE,
				// CreatedAt:     createdTime,
			}},
			expected: CreateChannelOutput{
				Id:            newChannelId.Hex(),
				UserEmail:     "testemail@email.com",
				CurrentStatus: domain.ACTIVE,
				// CreatedAt:     createdTime,
			},
		},
		{
			name: "create channels returning error",
			args: args{input: CreateChannelInput{
				UserEmail: "",
			}},
			channelRepo: mockCreateChannelRepo{
				createChannelFake: func() (domain.Channel, error) {
					return domain.Channel{}, errors.New("error")
				},
				findByIDFake: func() (domain.Channel, error) {
					return domain.NewChannel(
						newChannelId,
						"testemail@email.com",
						domain.ACTIVE,
						createdTime,
						createdTime,
					), nil
				},
			},
			presenter:     mockCreateChannelPresenter{result: CreateChannelOutput{}},
			expected:      CreateChannelOutput{},
			expectedError: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewCreateChannelInteractor(tt.channelRepo, tt.presenter, time.Second)

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
