package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_loginUserPresenter_Output(t *testing.T) {
	type args struct {
		user  domain.User
		token string
	}
	createdAt := time.Now()
	userId := primitive.NewObjectID()
	user := domain.NewUser(
		userId,
		"FirstName",
		"LastName",
		"firstname.lastname@gmail.com",
		"superSecurePassword",
		createdAt,
		createdAt,
	)
	user.UpdateRole("USER")
	tests := []struct {
		name string
		args args
		want usecase.LoginUserOutput
	}{
		{
			name: "Login User",
			args: args{
				user:  user,
				token: "04jf03945r0394w;th3490594j03",
			},
			want: usecase.LoginUserOutput{
				FirstName: "FirstName",
				LastName:  "LastName",
				Email:     "firstname.lastname@gmail.com",
				Role:      "USER",
				Token:     "04jf03945r0394w;th3490594j03",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewLoginPresenter()
			if got := pre.Output(tt.args.user, tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}
