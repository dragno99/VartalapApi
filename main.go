package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dragno99/VartalapApi/router"
	"github.com/joho/godotenv"
)

func main() {

	// start the router
	myRouter := router.InitRouter()

	godotenv.Load(".env")

	PORT := string(os.Getenv("PORT"))

	if PORT == "" {
		PORT = ":8080"
	}
	// start listening
	log.Fatal(http.ListenAndServe(PORT, myRouter))
	// log.Fatal(http.ListenAndServe(PORT, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(myRouter)))

}
