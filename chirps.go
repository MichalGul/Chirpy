package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/MichalGul/Chirpy/internal/auth"
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
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := request_parameters{}
	reqErr := decoder.Decode(&reqParams)
	if reqErr != nil {
		respondWithError(response, 500, "Error decoding parameters", reqErr)
		return
	}

	log.Printf("chirp request params: %s\n", reqParams)

	tokenString, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, 401, "Unauthorized", err)
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(response, 401, "Unauthorized", err)
		return
	}

	cleanedChirp, validateError := validateChirp(reqParams.Body)
	if validateError != nil {
		respondWithError(response, http.StatusBadRequest, validateError.Error(), validateError)
		return
	}

	createdChirp, createError := cfg.db.CreateChirp(request.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: userId,
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
		UserID:    userId,
	}

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 201, respBody)

}

func (cfg *apiConfig) handleGetChirps(response http.ResponseWriter, request *http.Request) {

	authorId := request.URL.Query().Get("author_id") // Extract value of authorId from querry parameters

	var chirps []database.Chirp
	var err error

	if authorId != "" {
		userId, _ := uuid.Parse(authorId)
		chirps, err = cfg.db.GetChirpsByAuthorId(request.Context(), userId)
	} else {
		chirps, err = cfg.db.GetChirps(request.Context())

	}

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

	sortType := request.URL.Query().Get("sort")
	if sortType == "desc" {
		sort.Slice(chirpResponse, func(i, j int) bool { return chirpResponse[i].CreatedAt.After(chirpResponse[j].CreatedAt) })
	} else {
		sort.Slice(chirpResponse, func(i, j int) bool { return chirpResponse[i].CreatedAt.Before(chirpResponse[j].CreatedAt) })
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

	chirp := Chirp{
		ID:        chirpDb.ID,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body:      chirpDb.Body,
		UserID:    chirpDb.UserID,
	}
	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 200, chirp)

}

func (cfg *apiConfig) handleDeleteChirpByID(response http.ResponseWriter, request *http.Request) {

	chirpID := request.PathValue("chirpID") // Extract value of chirpID

	chirpUUID, parseErr := uuid.Parse(chirpID)
	if parseErr != nil {
		respondWithError(response, 500, "error parsing id to uuid", parseErr)
		return
	}

	tokenString, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, 401, "Unauthorized", err)
		return
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(response, 401, "Unauthorized", err)
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

	if userId != chirpDb.UserID {
		respondWithError(response, 403, "User is not owner of that chirp", getErr)
		return
	}

	delErr := cfg.db.DeleteChirpById(request.Context(), chirpUUID)
	if delErr != nil {
		if err.Error() == NoRowError.Error() {
			respondWithError(response, 404, "Chirp not found", delErr)
			return
		}

		respondWithError(response, 500, "unknown error deleting chirp by id", delErr)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 204, nil)

}

func validateChirp(body string) (string, error) {
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	cleanedChirp := validateBadWordChirp(body)
	return cleanedChirp, nil
}
