package main

import (
	"fmt"
	"log"
	"net/http"
)

func cheersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Cheers!")
}

func main() {
	http.HandleFunc("/cheers", cheersHandler)
	log.Printf("Running server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
