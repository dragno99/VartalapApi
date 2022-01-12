package main

import (
	"log"
	"net/http"
	"vartalap/router"
)

func main() {

	// fmt.Println("My simple token")

	// tokenString, err := middleware.GenerateJWT()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(tokenString)

	// start the router
	myRouter := router.InitRouter()

	// start listening
	log.Fatal(http.ListenAndServe(":8000", myRouter))

}
