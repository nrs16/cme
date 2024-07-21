package middleware

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username string
}

func ParseClaims(claims jwt.MapClaims) Claims {
	var c Claims
	for k, v := range claims {
		switch k {
		case "username":
			c.Username = v.(string)
		}

	}
	return c
}
func validateToken(token *jwt.Token) (interface{}, error) {
	keyString := "A8S/VfWYAEZA7lxs3MqOVpi/GvZrVdARGkNtODgJy1Y="
	key := []byte(keyString)
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("invalid_signing_method")
	}
	return key, nil
}

func GenerateJWT(claims jwt.MapClaims) (string, error) {
	keyString := "A8S/VfWYAEZA7lxs3MqOVpi/GvZrVdARGkNtODgJy1Y="
	key := []byte(keyString)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(token)

	tokenString, err := token.SignedString(key)
	return tokenString, err

}

func ValidateJWT(token string) (jwt.MapClaims, error) {
	var mapClaims jwt.MapClaims
	t, err := jwt.Parse(token, validateToken)
	if err != nil {
		return mapClaims, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return mapClaims, errors.New("invalid_claims")
	}
	return claims, nil
}
