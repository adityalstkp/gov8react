package handler

import (
	"encoding/json"
	"net/http"

	"github.com/adityalstkp/gov8react/internal/usecase"
)

type introHandler struct {
	introUsecase usecase.IntroUsecase
}

type IntroHandlerOpts struct {
	IntroUsecase usecase.IntroUsecase
}

func NewIntroHandler(opts IntroHandlerOpts) introHandler {
	return introHandler{introUsecase: opts.IntroUsecase}
}

type IntroHandlerRouter interface {
	Greet(w http.ResponseWriter, r *http.Request)
}

func (iH introHandler) Greet(w http.ResponseWriter, r *http.Request) {
	d, err := iH.introUsecase.Greet()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}
