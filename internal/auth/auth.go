package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   stringID})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	id, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}
	return parsedId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authorisationHeader := headers.Get("Authorization")
	if authorisationHeader == "" {
		return "", errors.New("no token avalible")
	}
	sliceHeader := strings.Split(authorisationHeader, " ")
	if len(sliceHeader) != 2 {
		return "", errors.New("incorrect format of authorisation header")
	}
	if strings.ToLower(sliceHeader[0]) != "bearer" {
		return "", errors.New("bearer prefix not found")
	}
	return sliceHeader[1], nil
}

func MakeRefreshToken() (string, error) {
	randomData := make([]byte, 32)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", err
	}
	hexstring := hex.EncodeToString(randomData)
	return hexstring, nil
}
