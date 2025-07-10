package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	stringID := userID.String()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{Issuer: "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   stringID})
	signedToken, err := token.SignedString(token)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	return uuid.Nil, nil
}
