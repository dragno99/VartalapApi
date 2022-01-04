package main

import (
	"fmt"
	"log"
	"net/http"
	"vartalap/router"
)

func main() {
	myRouter := router.InitRouter()
	err := http.ListenAndServe(":8000", myRouter)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Server started")
	}
}
