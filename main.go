package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

// web ui static assets
//go:embed webui/*
var webuiFs embed.FS

func generateUiHandler() (http.Handler, error) {
	uiFileSystem, err := fs.Sub(webuiFs, "webui")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(uiFileSystem)), nil
}

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
	webuiHandler, err := generateUiHandler()

	if err != nil {
		log.Panicf("[main.go#main] Error generateUiHandler %s", err)
	}

	http.Handle("/", webuiHandler)
	http.HandleFunc("/cheers", cheersHandler)

	port := getPort()
	log.Printf("Running server PORT=%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
