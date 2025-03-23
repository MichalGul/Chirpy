package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("Hallo Chirpy\n")

	mutex := http.NewServeMux()

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: mutex,
	}

	err := httpServer.ListenAndServe()
	log.Printf("Serving on port 8080\n")
	if err != nil {
		log.Fatal(err)
	}

}
