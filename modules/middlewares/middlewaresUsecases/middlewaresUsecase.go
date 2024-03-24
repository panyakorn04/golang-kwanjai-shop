package middlewaresUsecases

import (
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/middlewares"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecases interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecases struct {
	middlewaresRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecases(middlewaresRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecases {
	return &middlewaresUsecases{
		middlewaresRepository: middlewaresRepository,
	}
}

func (u *middlewaresUsecases) FindAccessToken(userId, accessToken string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewaresUsecases) FindRole() ([]*middlewares.Role, error) {
	roles, err := u.middlewaresRepository.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
