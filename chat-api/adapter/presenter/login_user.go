package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
)

type loginUserPresenter struct{}

func NewLoginPresenter() usecase.LoginUserPresenter {
	return loginUserPresenter{}
}

func (a loginUserPresenter) Output(user domain.User, token string) usecase.LoginUserOutput {
	return usecase.LoginUserOutput{
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Token:     token,
		Email:     user.Email(),
		Role:      user.Role(),
	}
}
