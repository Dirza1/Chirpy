package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	bitePassword := []byte(password)
	cost, err := bcrypt.Cost(bitePassword)
	if err != nil {
		return "", err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(bitePassword, cost)
	if err != nil {
		return "", err
	}
	returnPassword := string(hashedPassword)
	return returnPassword, nil
}

func CheckPasswordHash(password, hash string) error {
	bytepassword := []byte(password)
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, bytepassword)
	if err != nil {
		return err
	}
	return nil
}
