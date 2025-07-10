package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	bitePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bitePassword, 10)
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
