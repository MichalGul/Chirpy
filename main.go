package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

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

func (cfg *apiConfig) resetServerHitsHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)

	body := fmt.Sprintf("Hits reset!")
	_, err := response.Write([]byte(body))
	if err != nil {
		http.Error(response, "Unable to write response", http.StatusInternalServerError)
		return
	}
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

func main() {
	fmt.Printf("Hallo Chirpy\n")
	const port = "8080"
	const filepathRoot = "."

	apiConfig := &apiConfig{}

	httpMultiplexer := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: httpMultiplexer,
	}

	httpMultiplexer.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))) // This will handle all paths with app prefix. Eg. app/assets
	httpMultiplexer.HandleFunc("GET /api/healthz", healthHandler)
	httpMultiplexer.HandleFunc("GET /admin/metrics", apiConfig.returnServerHitsHandler)
	httpMultiplexer.HandleFunc("POST /admin/reset", apiConfig.resetServerHitsHandler)

	log.Printf("Serving on port 8080\n")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
