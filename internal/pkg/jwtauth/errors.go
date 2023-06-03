package jwtauth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func checkTokenValidationErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, jwt.ErrTokenExpired) {
		return errors.New("token is expired and should be renewed using the refresh token")
	}
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return jwt.ErrSignatureInvalid
	}
	return err
}
