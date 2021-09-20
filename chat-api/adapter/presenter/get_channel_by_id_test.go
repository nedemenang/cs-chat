package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_getChannelByIdPresenter_Output(t *testing.T) {
	type args struct {
		channel domain.Channel
	}
	channelId := primitive.NewObjectID()
	createdAt := time.Now()
	channel := domain.NewChannel(
		channelId,
		"anthony.jones@gmail.com",
		"ACTIVE",
		createdAt,
		createdAt,
	)
	channel.UpdateRepEmail("jones.anthony@gmail.com")
	channel.AddMessage("jones.anthony@gmail.com", "Hello world", createdAt)

	tests := []struct {
		name string
		args args
		want usecase.GetChannelByIdOutput
	}{
		{
			name: "Create channel",
			args: args{
				channel: channel,
			},
			want: usecase.GetChannelByIdOutput{
				Id:            channelId.Hex(),
				UserEmail:     "anthony.jones@gmail.com",
				RepEmail:      "jones.anthony@gmail.com",
				CurrentStatus: "ACTIVE",
				CreatedAt:     createdAt,
				Messages: []usecase.Message{{
					Message:     "Hello world",
					MessageFrom: "jones.anthony@gmail.com",
					Timestamp:   createdAt,
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewGetChannelByIdPresenter()
			if got := pre.Output(tt.args.channel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}
