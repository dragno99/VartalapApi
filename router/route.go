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
	router.HandleFunc("/appusers/{userId}", controller.GetAppUser).Methods("GET")
	router.HandleFunc("/chatmessages/{chatId}", controller.GetChatMessages).Methods("GET")
	router.HandleFunc("/startchat/{userId}", controller.StartChat).Methods("POST")
	router.HandleFunc("/message/{userId}", controller.AddMessage).Methods("POST")
	router.HandleFunc("/updatefullname/{userId}", controller.UpdateFullName).Methods("POST")

	return router
}
