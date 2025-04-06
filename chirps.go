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

type Chirps struct {
	Chirp
}

const maxChirpLength = 140

var NoRowError = errors.New("sql: no rows in result set")

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

func (cfg *apiConfig) handleGetChirps(response http.ResponseWriter, request *http.Request) {

	chirps, err := cfg.db.GetChirps(request.Context())
	if err != nil {
		respondWithError(response, 400, "error getting chirps", err)
	}

	chirpResponse := make([]Chirp, len(chirps))

	for i, c := range chirps {
		chirpResponse[i] = Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}

	}

	log.Printf("chirps: %v \n", chirpResponse)

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 200, chirpResponse)

}

func (cfg *apiConfig) handleGetChirpByID(response http.ResponseWriter, request *http.Request) {

	chirpID := request.PathValue("chirpID") // Extract value of chirpID
	chirpUUID, parseErr := uuid.Parse(chirpID)
	if parseErr != nil {
		respondWithError(response, 500, "error parsing id to uuid", parseErr)
		return
	}

	chirpDb, getErr := cfg.db.GetChirpById(request.Context(), chirpUUID)
	if getErr != nil {
		if getErr.Error() == NoRowError.Error() {
			respondWithError(response, 404, "Chirp not found", getErr)
			return
		}

		respondWithError(response, 500, "unknown error getting chirp by id", getErr)
		return
	}

	chirp := Chirp {
		ID: chirpDb.ID,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body: chirpDb.Body,
		UserID: chirpDb.UserID,
	}
	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 200, chirp)

}

func validateChirp(body string) (string, error) {
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	cleanedChirp := validateBadWordChirp(body)
	return cleanedChirp, nil
}
