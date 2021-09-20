package presenter

import (
	"chat-api/domain"
	"chat-api/usecase"
)

type getUserByEmailPresenter struct{}

func NewGetUserByEmailPresenter() usecase.GetUserByEmailPresenter {
	return getUserByEmailPresenter{}
}

func (a getUserByEmailPresenter) Output(user domain.User) usecase.GetUserByEmailOutput {
	return usecase.GetUserByEmailOutput{
		Id:        user.Id().Hex(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Email:     user.Email(),
		Role:      user.Role(),
	}
}
