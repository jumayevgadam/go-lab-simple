package main

import (
	"log"
	"net/http"
)

func main() {
	// implement handlers in this place
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/logout", Logout)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
