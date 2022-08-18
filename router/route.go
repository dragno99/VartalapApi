package router

import (
	"net/http"

	"github.com/dragno99/VartalapApi/controller"
	"github.com/dragno99/VartalapApi/middleware"

	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello", controller.SayHello).Methods("GET")
	router.HandleFunc("/login", controller.UserLogIn).Methods("POST")
	router.HandleFunc("/signup", controller.UserSignUp).Methods("POST")
	router.Handle("/userchats/", middleware.IsAuthorized(controller.GetUserChats)).Methods("GET")
	router.Handle("/appusers/", middleware.IsAuthorized(controller.GetAppUser)).Methods("GET")
	router.Handle("/chatmessages/{chatId}", middleware.IsAuthorized(controller.GetChatMessages)).Methods("GET")
	router.Handle("/startchat/", middleware.IsAuthorized(controller.StartChat)).Methods("POST")
	router.Handle("/joinChatRoom/{chatId}", http.HandlerFunc(controller.JoinChatRoom))
	router.Handle("/message/", middleware.IsAuthorized(controller.AddMessage)).Methods("POST")
	router.Handle("/updatefullname/", middleware.IsAuthorized(controller.UpdateFullName)).Methods("POST")
	return router
}
