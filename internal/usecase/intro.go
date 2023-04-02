package usecase

import "github.com/adityalstkp/gov8react/internal/model"

type introUsecase struct{}

func NewIntroUsecase() introUsecase {
	return introUsecase{}
}

type IntroUsecase interface {
	Greet() (model.GreetResponse, error)
}

func (iU introUsecase) Greet() (model.GreetResponse, error) {
	return model.GreetResponse{Message: "Hi, we are still working progress!"}, nil
}
