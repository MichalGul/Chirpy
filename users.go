package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleUsers(response http.ResponseWriter, request *http.Request) {

	//expected request post parameters

	type request_parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := request_parameters{}
	reqErr := decoder.Decode(&reqParams)
	if reqErr != nil {
		respondWithError(response, 500, "Error decoding parameters", reqErr)
		return
	}

	log.Printf("req: %s\n", reqParams)

	createdUser, createErr := cfg.db.CreateUser(request.Context(), reqParams.Email)
	if createErr != nil {
		respondWithError(response, 500, "Error creating user parameters", createErr)
		return
	}

	respBody := User{
		ID:         createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:      createdUser.Email,
	}

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 201, respBody)
}
