package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MichalGul/Chirpy/internal/auth"
)

type Token struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handleRefresh(response http.ResponseWriter, request *http.Request) {

	tokenString, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, 401, "Unauthorized", err)
		return
	}

	refreshTokenDb, err := cfg.db.GetRefreshToken(request.Context(), tokenString)
	if err != nil {
		if err.Error() == NoRowError.Error() {
			respondWithError(response, 404, "Refresh token not found", err)
			return
		}
	}

	fmt.Printf("expires at: %v", refreshTokenDb.ExpiresAt)

	if refreshTokenDb.ExpiresAt.Before(time.Now()) {
		respondWithError(response, 404, "Refresh token expired", err)
		return
	}

	if refreshTokenDb.RevokedAt.Valid {
		respondWithError(response, 401, "Unauthorized", err)
		return
	}

	userId, err := cfg.db.GetUserFromRefreshToken(request.Context(), tokenString)
	if err != nil {
		respondWithError(response, 401, "error extracting user by refresh token", err)
		return
	}

	token, err := auth.MakeJWT(userId, cfg.secret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(response, 401, "error creating token", err)
		return
	}

	respBody := Token{
		Token: token,
	}

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 200, respBody)

}

func (cfg *apiConfig) handleRevoke(response http.ResponseWriter, request *http.Request) {

	tokenString, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, 401, "Unauthorized", err)
		return
	}

	refreshTokenDb, err := cfg.db.GetRefreshToken(request.Context(), tokenString)
	if err != nil {
		if err.Error() == NoRowError.Error() {
			respondWithError(response, 404, "Refresh token not found", err)
			return
		}
	}

	updatedToken, err := cfg.db.SetRevokeOnToken(request.Context(), refreshTokenDb.Token)
	if err != nil {
		respondWithError(response, 401, "error revoking token", err)
		return
	}
	if !updatedToken.RevokedAt.Valid {
		respondWithError(response, 401, "error revoking token", err)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	respondWithJSON(response, 204, nil)

}
