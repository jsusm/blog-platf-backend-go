package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/JesusJMM/blog-plat-go/postgres"
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	UID int `json:"uid"`
	UserName string `json:"userName"`
	UserImg *string `json:"userImg"`
	jwt.RegisteredClaims
}

var SECRET_KEY []byte

const TOKEN_DURATION = 24 * time.Hour

func init() {
	SECRET_KEY = []byte(os.Getenv("TOKEN_SECRET_KEY"))
}

func SignToken(user postgres.User) (string, error) {
	claims := TokenClaims{
    user.ID,
    user.Name,
    user.Img,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_DURATION)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(SECRET_KEY)
	return ss, err
}

func ParseToken(tokenString string) (*jwt.Token, *TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unspected signing method: %v", t.Header["alg"])
		}
		return SECRET_KEY, nil
	})
  if err != nil {
		return nil, nil, err
  }
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return token, claims, nil
	} else {
		return nil, nil, err
	}
}
