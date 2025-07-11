package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authTestStruct struct {
	test        string
	id          uuid.UUID
	tokenSecred string
	expiresin   time.Duration
	expectedErr error
}

func TestJWTtokens(t *testing.T) {
	tokenSecred := "jvaprenv;jjna'gkaoshyrh;kasv'kafhasld;fnaso;df"

	tests := []authTestStruct{
		{
			test:        "Succesfull test",
			id:          uuid.New(),
			tokenSecred: tokenSecred,
			expiresin:   30 * time.Second,
			expectedErr: nil,
		},
		{
			test:        "wrong token secrit",
			id:          uuid.New(),
			tokenSecred: "hello",
			expiresin:   30 * time.Second,
			expectedErr: jwt.ErrTokenSignatureInvalid,
		},
		{
			test:        "token expired",
			id:          uuid.New(),
			tokenSecred: tokenSecred,
			expiresin:   1 * time.Microsecond,
			expectedErr: jwt.ErrTokenExpired,
		},
	}

	for _, test := range tests {
		fmt.Printf("Test started. Test: %s\n", test.test)
		token, err := MakeJWT(test.id, test.tokenSecred, test.expiresin)
		if err != nil {
			fmt.Printf("Test failed on generating token: %v\n", err)
			t.Fail()
			return
		}
		time.Sleep(1 * time.Second)
		_, err = ValidateJWT(token, tokenSecred)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("test %q: expected %v but recieved %v", test.test, test.expectedErr, err)
		}
	}

}
