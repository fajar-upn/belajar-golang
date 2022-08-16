package auth

import (
	"github.com/dgrijalva/jwt-go"
)

/**
this code for JWT
*/

type Service interface {
	GenerateToken(userID int) (string, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}

var SECRET_KEY = []byte("SECRET_KEY") //in the production level, don't write secret key here. recomended for save secret key in the save place

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) //we can choose SigningMethod, JWT have lot of signing method

	signetToken, err := token.SignedString(SECRET_KEY) //we must signature the secret key

	if err != nil {
		return signetToken, err
	}

	return signetToken, nil
}
