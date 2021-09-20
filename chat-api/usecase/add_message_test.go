package usecase

import (
	"chat-api/domain"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockAddMessageRepoStore struct {
	domain.ChannelRepository

	result domain.Channel
	err    error
}

var (
	newChannelId = primitive.NewObjectID()
)

type invoked struct {
	call bool
}

func (m mockAddMessageRepoStore) Create(_ context.Context, _ domain.Channel) (domain.Channel, error) {
	return m.result, m.err
}

type mockAddMessageRepo struct {
	domain.ChannelRepository

	addMessageFake func() error
	invokedCreate  *invoked

	findByIDFake func() (domain.Channel, error)
	invokedFind  *invoked
}

func (m mockAddMessageRepo) AddMessage(_ context.Context, _ domain.Channel) error {

	if m.invokedCreate != nil {
		m.invokedCreate.call = true
	}
	return m.addMessageFake()
}

func (m mockAddMessageRepo) GetChannelById(_ context.Context, _ string) (domain.Channel, error) {

	if m.invokedFind != nil {
		m.invokedFind.call = true
	}
	return m.findByIDFake()
}

type mockAddMessagePresenter struct {
	result CreateMessageOutput
}

func (m mockAddMessagePresenter) Output(_ domain.Channel) CreateMessageOutput {
	return m.result
}

func TestAddMessageInteractor_Execute(t *testing.T) {
	t.Parallel()

	type args struct {
		input CreateMessageInput
	}

	var messages = make([]MessageOutput, 0)
	for i := 0; i < 1; i++ {
		messages = append(messages, MessageOutput{
			MessageFrom: "validemail@gmail.com",
			Message:     "This is a test message",
			Timestamp:   time.Now(),
		})
	}

	tests := []struct {
		name          string
		args          args
		channelRepo   domain.ChannelRepository
		presenter     CreateMessagePresenter
		expected      CreateMessageOutput
		expectedError string
	}{
		{
			name: "create message successfully",
			args: args{input: CreateMessageInput{
				ChannelId:   newChannelId.Hex(),
				MessageFrom: "validemail@gmail.com",
				Message:     "this is a test email",
			}},
			channelRepo: mockAddMessageRepo{
				addMessageFake: func() error {
					return nil
				},
				findByIDFake: func() (domain.Channel, error) {
					return domain.NewChannel(
						newChannelId,
						"testemail@email.com",
						domain.ACTIVE,
						time.Now(),
						time.Now(),
					), nil
				},
			},
			presenter: mockAddMessagePresenter{result: CreateMessageOutput{
				Id:            newChannelId.Hex(),
				UserEmail:     "testemail@email.com",
				RepEmail:      "repTestemail@email.com",
				CurrentStatus: domain.ACTIVE,
				// CreatedAt:     time.Now(),
				Messages: messages,
			}},
			expected: CreateMessageOutput{
				Id:            newChannelId.Hex(),
				UserEmail:     "testemail@email.com",
				RepEmail:      "repTestemail@email.com",
				CurrentStatus: domain.ACTIVE,
				Messages:      messages,
				// CreatedAt:     time.Now(),
			},
		},
		{
			name: "create messages returning error",
			args: args{input: CreateMessageInput{
				ChannelId:   "",
				MessageFrom: "",
				Message:     "",
			}},
			channelRepo: mockAddMessageRepo{
				addMessageFake: func() error {
					return errors.New("error")
				},
				findByIDFake: func() (domain.Channel, error) {
					return domain.NewChannel(
						newChannelId,
						"testemail@email.com",
						domain.ACTIVE,
						time.Now(),
						time.Now(),
					), nil
				},
			},
			presenter:     mockAddMessagePresenter	{result: CreateMessageOutput{}},
			expected:      CreateMessageOutput{},
			expectedError: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uc = NewCreateMessageInteractor(tt.channelRepo, tt.presenter, time.Second)

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
