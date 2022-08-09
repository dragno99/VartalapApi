package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dragno99/vartalapAPI/router"
	"github.com/joho/godotenv"
)

func main() {

	// start the router
	myRouter := router.InitRouter()

	godotenv.Load()
	PORT := string(os.Getenv("PORT"))

	// start listening
	log.Fatal(http.ListenAndServe(PORT, myRouter))

}
