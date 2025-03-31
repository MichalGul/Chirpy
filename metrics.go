package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) returnServerHitsHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(`
	<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body
	
	</html>`, cfg.fileserverHits.Load())

	_, err := response.Write([]byte(body))
	if err != nil {
		http.Error(response, "Unable to write response", http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}


func healthHandler(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	body := "OK"

	_, err := response.Write([]byte(body))
	if err != nil {
		http.Error(response, "Unable to write response", http.StatusInternalServerError)
		return
	}

}