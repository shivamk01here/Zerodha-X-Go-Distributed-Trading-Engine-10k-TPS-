package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API Gateway: Order Received\n")
	})

	log.Println("API Gateway starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
