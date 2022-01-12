package main

import (
	"log"
	"net/http"
	"vartalap/router"
)

func main() {

	// start the router
	myRouter := router.InitRouter()

	// start listening
	log.Fatal(http.ListenAndServe(":8000", myRouter))

}
