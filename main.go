package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dragno99/VartalapApi/router"
)

func main() {

	// start the router
	myRouter := router.InitRouter()

	PORT := ":" + os.Getenv("PORT")

	// start listening
	log.Fatal(http.ListenAndServe(PORT, myRouter))

}
