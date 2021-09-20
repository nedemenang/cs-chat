package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_addMessagePresenter_Output(t *testing.T) {
	type args struct {
		channel domain.Channel
	}
	createdAt := time.Now()
	channeId := primitive.NewObjectID()
	channel := domain.NewChannel(
		channeId,
		"anthony.jones@gmail.com",
		"ACTIVE",
		createdAt,
		createdAt,
	)
	channel.AddMessage("anthony.jones@gmail.com", "Hello World", createdAt)
	tests := []struct {
		name string
		args args
		want usecase.CreateMessageOutput
	}{
		{
			name: "Add Message",
			args: args{
				channel: channel,
			},
			want: usecase.CreateMessageOutput{
				Id:            channeId.Hex(),
				UserEmail:     "anthony.jones@gmail.com",
				CurrentStatus: "ACTIVE",
				CreatedAt:     createdAt,
				Messages: []usecase.MessageOutput{{
					MessageFrom: "anthony.jones@gmail.com",
					Message:     "Hello World",
					Timestamp:   createdAt,
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewCreateMessagePresenter()
			if got := pre.Output(tt.args.channel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}
