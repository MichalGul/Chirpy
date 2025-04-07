package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		log.Printf("error hashing password %v", hashErr)
		return "", hashErr
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer: "chirpy",
			Subject: userID.String(),
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		})

	tokenString, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return tokenString, nil
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {


	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func (token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil

	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%v", err)
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		userId, parseErr := uuid.Parse(claims.Subject)
		if parseErr != nil {
			return uuid.UUID{}, fmt.Errorf("error parsing userId")
		}
		// Happy path
		if token.Valid {
			return userId, nil
		}
		return uuid.UUID{}, fmt.Errorf("token was not valid or expired")
	}

	return uuid.UUID{}, fmt.Errorf("unknown claims type")
}