package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetServerHitsHandler(response http.ResponseWriter, request *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(response, 403, "Forbiden action for this platform setting", nil)
	}
	
	
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)
	body := fmt.Sprintf("Hits and users reset!")
	cfg.db.DeleteUsers(request.Context())	


	_, err := response.Write([]byte(body))
	if err != nil {
		http.Error(response, "Unable to write response", http.StatusInternalServerError)
		return
	}
}
