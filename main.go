package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("Hallo Chirpy\n")
	const port = ":8080"
	const filepathRoot = "."

	mutex := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    port,
		Handler: mutex,
	}

	mutex.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	log.Printf("Serving on port 8080\n")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
