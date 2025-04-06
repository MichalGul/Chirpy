package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateBadWordChirp(chirp string) string {

	log.Printf("Input chirp: %s \n", chirp)

	notValidStrings := []string{"kerfuffle", "sharbert", "fornax"}

	splitedChirp := strings.Split(chirp, " ")

	for _, badString := range notValidStrings {

		for sIndex, ChirpWord := range splitedChirp {
			if badString == strings.ToLower(ChirpWord) {
				splitedChirp[sIndex] = "****"
			}
		}

	}

	filteredChirp := strings.Join(splitedChirp, " ")

	log.Printf("Output chirp: %s \n", filteredChirp)

	return filteredChirp

}

func handleValidation(response http.ResponseWriter, request *http.Request) {
	// Expected request post parameters
	type request_parameters struct {
		Body string `json:body`
	}

	type response_parameters struct {
		Valid        bool   `json:"valid"`
		Error        string `json:"error"`
		Cleaned_body string `json:"cleaned_body"`
	}

	const maxChirpLength = 140

	decoder := json.NewDecoder(request.Body)
	reqParams := request_parameters{}
	reqErr := decoder.Decode(&reqParams)
	if reqErr != nil {
		respondWithError(response, 500, "Error decoding parameters", reqErr)
		return
	}

	log.Printf("req: %s\n", reqParams)

	respBody := response_parameters{}
	response.Header().Set("Content-Type", "application/json")

	if len(reqParams.Body) > maxChirpLength {
		respondWithError(response, 400, "Chirp is too long", nil)
		return
	}

	// todo add profane check
	cleanedChirp := validateBadWordChirp(reqParams.Body)

	respBody.Valid = true
	respBody.Error = ""
	respBody.Cleaned_body = cleanedChirp
	respondWithJSON(response, 200, respBody)

}
