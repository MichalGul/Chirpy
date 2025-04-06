package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)


func HashPassword(password string) (string, error) {
	cost := 12
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), cost)
	if hashErr != nil {
		log.Printf("error hashing password %v", hashErr)
		return "", hashErr
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("Passwords does not match")
		return err
	}
	return nil

}
