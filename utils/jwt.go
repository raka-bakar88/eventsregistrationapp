package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const secretKey = "supersecretkey"

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	// convert token to a single string
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) (int64, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// check if the token is using the type/method used in this app
		_, isTheCorrectType := token.Method.(*jwt.SigningMethodHMAC)
		if !isTheCorrectType {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, errors.New("Invalid token")
	}
	//check if the token is valid by extracting the valid value from the parsedtoken
	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return 0, errors.New("tokenIsValid")
	}
	// get claims/2nd structure of the jwt
	claims, isValid := parsedToken.Claims.(jwt.MapClaims)
	if !isValid {
		return 0, errors.New("Invalid token")
	}

	//email := claims["email"].(string)
	userId := int64(claims["userId"].(float64))
	return userId, nil
}
