package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
)

type createUserPresenter struct{}

func NewCreateUserPresenter() usecase.CreateUserPresenter {
	return createUserPresenter{}
}

func (a createUserPresenter) Output(user domain.User) usecase.CreateUserOutput {
	return usecase.CreateUserOutput{
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Email:     user.Email(),
	}
}
