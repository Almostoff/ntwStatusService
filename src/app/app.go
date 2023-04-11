package app

import (
	"fProject/src/cases"
	"fProject/src/controller"
	"fProject/src/repositories"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run() {
	repository := repositories.NewRepository()
	useCase := cases.NewUseCase(repository)
	r := mux.NewRouter()
	controller.Build(r, useCase)
	log.Fatal(http.ListenAndServe("localhost:8282", r))
}
