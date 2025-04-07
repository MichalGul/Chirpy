package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	userId, _ := uuid.Parse("05ad1d63-5b1f-4ff7-acdf-5e628e505366")

	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Correct Make",
			userID:      userId,
			tokenSecret: "secret",
			expiresIn:   time.Minute * 5,
			wantErr:     false,
		},
		{
			name:        "Correct Make 1",
			userID:      userId,
			tokenSecret: "secretwwww",
			expiresIn:   time.Second * 5,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if !tt.wantErr && len(tokenString) == 0 {
				t.Error("Expected a valid token string, but got an empty string")
			}
			fmt.Printf("Created token string %s\n", tokenString)

			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.tokenSecret), nil
			})

			if err != nil || !token.Valid {
				t.Errorf("Failed to parse or validate created token: %v", err)
			}

			claims, ok := token.Claims.(*jwt.RegisteredClaims)
			if !ok {
				t.Error("Invalid claims type")
			}
			if claims.Subject != tt.userID.String() {
				t.Errorf("Expected %s", tt.userID)
			}

		})
	}
}


func TestValidateJWT(t *testing.T) {
	userId, _ := uuid.Parse("05ad1d63-5b1f-4ff7-acdf-5e628e505366")

	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecretIN string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Correct Make and Valid",
			userID:      userId,
			tokenSecretIN: "secret",
			expiresIn:   time.Minute * 5,
			wantErr:     false,
		},
		{
			name:        "Correct Make and expired",
			userID:      userId,
			tokenSecretIN: "secretwwww",
			expiresIn:   time.Second * 1,
			wantErr:     true,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tokenString, err := MakeJWT(tt.userID, tt.tokenSecretIN, tt.expiresIn)
			if !tt.wantErr && len(tokenString) == 0 {
				t.Error("Expected a valid token string, but got an empty string")
			}
			fmt.Printf("Created token string %s\n", tokenString)

			if (err != nil)  {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			time.Sleep(5 * time.Second) 
			userId, err := ValidateJWT(tokenString, tt.tokenSecretIN)

			if !tt.wantErr && userId != tt.userID {
				t.Errorf("Bad user returned from token validation")
			}

			if tt.wantErr && err == nil {
				t.Errorf("Token should be expired but is valid")
			}


		})
	}



}