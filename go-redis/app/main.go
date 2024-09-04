package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/products", httpHandler())
	http.ListenAndServe(":8080", nil)
}

// httpHandler is
func httpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := getProducts()
		if err != nil {
			fmt.Fprintf(w, err.Error(), "\r\n")
		} else {
			enc := json.NewEncoder(w)
			enc.SetIndent("", " ")

			if err := enc.Encode(response); err != nil {
				fmt.Println(err.Error())
			}
		}

	}
}
