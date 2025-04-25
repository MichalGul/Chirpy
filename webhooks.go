package main

import (
	"encoding/json"
	"net/http"

	"github.com/MichalGul/Chirpy/internal/auth"
	"github.com/MichalGul/Chirpy/internal/database"
	"github.com/google/uuid"
)

type WebhookParameters struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}
type Data struct {
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlePolkaWebhook(response http.ResponseWriter, request *http.Request) {

	apiKey, errGetApi := auth.GetAPIKey(request.Header)
	if errGetApi != nil {
		respondWithError(response, 401, "Unauthorized", errGetApi)
		return
	}

	if apiKey != cfg.polka_key {
		respondWithError(response, 401, "Unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := WebhookParameters{}
	reqErr := decoder.Decode(&reqParams)
	if reqErr != nil {
		respondWithError(response, 500, "Error decoding parameters", reqErr)
		return
	}

	if reqParams.Event != "user.upgraded" {
		// We don't do anything
		respondWithJSON(response, 204, nil)
		return
	}

	_, err := cfg.db.SetChirpyRed(request.Context(), database.SetChirpyRedParams{
		ID:          reqParams.Data.UserID,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(response, 404, "error updating user", reqErr)
		return
	}

	respondWithJSON(response, 204, nil)
	return

}
