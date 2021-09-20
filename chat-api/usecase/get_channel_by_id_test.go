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

type mockGetChannelByIdRepo struct {
	domain.ChannelRepository

	result domain.Channel
	err    error
}

func (m mockGetChannelByIdRepo) GetChannelById(_ context.Context, _ string) (domain.Channel, error) {
	return m.result, m.err
}

type mockGetChannelByIdPresenter struct {
	result GetChannelByIdOutput
}

func (m mockGetChannelByIdPresenter) Output(_ domain.Channel) GetChannelByIdOutput {
	return m.result
}

func TestGetChannelByIdInteractor_Execute(t *testing.T) {
	t.Parallel()

	type args struct {
		Id string
	}

	var newChannel = domain.NewChannel(primitive.NilObjectID, "user_email@gmail.com", "ACTIVE", time.Time{}, time.Time{})
	var message = Message{
		Message:     "message",
		MessageFrom: "message_from@gmail.com",
		Timestamp:   time.Now(),
	}
	var messages = []Message{message}

	tests := []struct {
		name          string
		args          args
		repository    domain.ChannelRepository
		presenter     GetChannelByIdPresenter
		expected      GetChannelByIdOutput
		expectedError interface{}
	}{
		{
			name:          "Success when returning the channel",
			args:          args{Id: "192039202912"},
			repository:    mockGetChannelByIdRepo{result: newChannel, err: nil},
			presenter:     mockGetChannelByIdPresenter{result: GetChannelByIdOutput{Id: newChannel.Id().Hex(), UserEmail: newChannel.UserEmail(), RepEmail: newChannel.RepEmail(), CurrentStatus: newChannel.CurrentStatus(), CreatedAt: newChannel.CreatedAt(), Messages: messages}},
			expected:      GetChannelByIdOutput{Id: newChannel.Id().Hex(), UserEmail: newChannel.UserEmail(), RepEmail: newChannel.RepEmail(), CurrentStatus: newChannel.CurrentStatus(), CreatedAt: newChannel.CreatedAt(), Messages: messages},
			expectedError: nil,
		},
		{
			name: "Error returning user",
			args: args{
				Id: "192039202912",
			},
			repository: mockGetChannelByIdRepo{
				result: domain.Channel{},
				err:    errors.New("error"),
			},
			presenter: mockGetChannelByIdPresenter{
				result: GetChannelByIdOutput{},
			},
			expectedError: "error",
			expected:      GetChannelByIdOutput{},
		},
	}

	for _, tt := range tests {
		var uc = NewGetChannelByIdInteractor(tt.repository, tt.presenter, time.Second)

		result, err := uc.Execute(context.Background(), GetChannelByIdInput{tt.args.Id})
		if (err != nil) && (err.Error() != tt.expectedError) {
			t.Errorf("[TestCase '%s'] Result: '%v' | ExpectedError: '%v'", tt.name, err, tt.expectedError)
			return
		}

		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'", tt.name, result, tt.expected)
		}
	}
}
