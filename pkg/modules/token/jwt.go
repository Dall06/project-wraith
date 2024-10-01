package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateJwtToken(secret string, exp time.Duration, data interface{}) (string, error) {
	claims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ExpireJwtToken(secret string, exp time.Duration, data interface{}) (string, error) {
	claims := jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ValidateJwtToken(tokenStr string, secret string, extraValidation func(jwt.MapClaims) error) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, errors.New("invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return false, errors.New("token has expired")
	}

	if extraValidation != nil {
		if err := extraValidation(claims); err != nil {
			return false, err
		}
	}

	return true, nil
}
