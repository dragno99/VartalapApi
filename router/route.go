package router

import (
	"vartalap/controller"

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello", controller.SayHello).Methods("GET")
	router.HandleFunc("/login", controller.UserLogIn).Methods("POST")
	router.HandleFunc("/signup", controller.UserSignUp).Methods("POST")
	return router
}
