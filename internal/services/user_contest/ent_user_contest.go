package user_contest

import "go-rest-api-boilerplate/types"

func NewResults(numCorrect, numIncorrect int, timeSubmit float64) *types.Results {
	return &types.Results{
		NumIncorrect: numIncorrect,
		NumCorrect:   numCorrect,
		TimeSubmit:   timeSubmit,
	}
}
