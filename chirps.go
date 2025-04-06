package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	"github.com/MichalGul/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

const maxChirpLength = 140

func (cfg *apiConfig) handleChirps(response http.ResponseWriter, request *http.Request) {

	type request_parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := request_parameters{}
	reqErr := decoder.Decode(&reqParams)
	if reqErr != nil {
		respondWithError(response, 500, "Error decoding parameters", reqErr)
		return
	}

	log.Printf("chirp request params: %s\n", reqParams)

	cleanedChirp, validateError := validateChirp(reqParams.Body)
	if validateError != nil {
		respondWithError(response, http.StatusBadRequest, validateError.Error(), validateError)
		return
	}

	createdChirp, createError := cfg.db.CreateChirp(request.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: reqParams.UserID,
	})

	if createError != nil {
		respondWithError(response, 400, "unknown error creating Chirp", createError)
		return
	}

	respBody := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      cleanedChirp,
		UserID:    reqParams.UserID,
	}

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 201, respBody)

}

func validateChirp(body string) (string, error) {
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	cleanedChirp := validateBadWordChirp(body)
	return cleanedChirp, nil
}
