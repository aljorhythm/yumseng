package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func cheersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Cheers!")
}

func getPort() string {
	value := os.Getenv("PORT")
	if value == "" {
		log.Printf("Resorting to default port 80")
		return "80"
	}
	return value
}

func main() {
	http.HandleFunc("/cheers", cheersHandler)
	port := getPort()
	log.Printf("Running server PORT=%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
