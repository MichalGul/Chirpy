package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MichalGul/Chirpy/internal/auth"
	"github.com/MichalGul/Chirpy/internal/database"
)

func parseExpirationTime(expirationTime *int) int {
	var expiresIn int = 3600
	if expirationTime != nil {
		expiresIn = *expirationTime

		if expiresIn > 3600 {
			expiresIn = 3600
		}
	}
	return expiresIn

}

func (cfg *apiConfig) handleLogin(response http.ResponseWriter, request *http.Request) {

	type request_parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := request_parameters{}
	reqErr := decoder.Decode(&reqParams)
	if reqErr != nil {
		respondWithError(response, 500, "Error decoding parameters", reqErr)
		return
	}

	expiresIn := parseExpirationTime(nil)
	fmt.Printf("Expires in %d", expiresIn)

	// get user by email
	user, getErr := cfg.db.GetUserByEmail(request.Context(), reqParams.Email)
	if getErr != nil {
		respondWithError(response, 401, "Incorrect email or password", getErr)
		return
	}

	authErr := auth.CheckPasswordHash(user.HashedPassword, reqParams.Password)
	if authErr != nil {
		respondWithError(response, 401, "Incorrect email or password", authErr)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		respondWithError(response, 401, "error creating token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(response, 401, "error creating refresh token", err)
		return
	}

	createdRefreshToken, err := cfg.db.CreateRefreshToken(request.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Duration(60*24)*time.Hour),
	})
	if err != nil {
		respondWithError(response, 401, "error adding refresh token to database", err)
		return
	}

	returnUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     &token,
		RefreshToken: &createdRefreshToken.Token,
	}

	respondWithJSON(response, 200, returnUser)

}

func (cfg *apiConfig) handleUsers(response http.ResponseWriter, request *http.Request) {

	//expected request post parameters

	type request_parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := request_parameters{}
	reqErr := decoder.Decode(&reqParams)
	if reqErr != nil {
		respondWithError(response, 500, "Error decoding parameters", reqErr)
		return
	}

	log.Printf("req: %s\n", reqParams)

	passwordHash, hashErr := auth.HashPassword(reqParams.Password)
	if hashErr != nil {
		respondWithError(response, 500, "Error hashing user password", hashErr)
		return
	}

	createdUser, createErr := cfg.db.CreateUser(request.Context(), database.CreateUserParams{
		Email:          reqParams.Email,
		HashedPassword: passwordHash,
	})
	if createErr != nil {
		respondWithError(response, 500, "Error creating user parameters", createErr)
		return
	}

	respBody := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 201, respBody)
}


func (cfg *apiConfig) handleEditUser(response http.ResponseWriter, request *http.Request) {
	
	//expected request put parameters
	type request_parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	decoder := json.NewDecoder(request.Body)
	reqParams := request_parameters{}
	reqErr := decoder.Decode(&reqParams)
	if reqErr != nil {
		respondWithError(response, 500, "Error decoding parameters", reqErr)
		return
	}

	log.Printf("req: %s\n", reqParams)

	passwordHash, hashErr := auth.HashPassword(reqParams.Password)
	if hashErr != nil {
		respondWithError(response, 500, "Error hashing user password", hashErr)
		return
	}

	editedUser, err := cfg.db.UpdateUser(request.Context(), database.UpdateUserParams{
		ID: userId,
		Email: reqParams.Email,
		HashedPassword: passwordHash,
	})

	if err != nil {
		respondWithError(response, 404, "error uupdating user to edit", err)
		return
	}

	respBody := User{
		ID:        editedUser.ID,
		CreatedAt: editedUser.CreatedAt,
		UpdatedAt: editedUser.UpdatedAt,
		Email:     editedUser.Email,
	}

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 200, respBody)

}