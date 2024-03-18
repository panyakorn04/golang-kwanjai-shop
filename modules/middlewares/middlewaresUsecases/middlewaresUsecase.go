package middlewaresUsecases

import (
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecases interface {
	FindAccessToken(userId, accessToken string) bool
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
