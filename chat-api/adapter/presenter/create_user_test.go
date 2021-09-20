package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_createUserPresenter_Output(t *testing.T) {
	type args struct {
		user domain.User
	}
	userId := primitive.NewObjectID()
	user := domain.NewUser(
		userId,
		"FirstName",
		"LastName",
		"firstname.lastname@gmail.com",
		"superSecurePassword",
		time.Now(),
		time.Now(),
	)
	tests := []struct {
		name string
		args args
		want usecase.CreateUserOutput
	}{
		{
			name: "Create user",
			args: args{
				user: user,
			},
			want: usecase.CreateUserOutput{
				FirstName: "FirstName",
				LastName:  "LastName",
				Email:     "firstname.lastname@gmail.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewCreateUserPresenter()
			if got := pre.Output(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}
