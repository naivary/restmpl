package jwtauth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

func GetClaims(ctx context.Context) (jwt.MapClaims, error) {
	val := ctx.Value(claimsKey)
	c, ok := val.(jwt.MapClaims)
	if !ok {
		return nil, ErrMissingClaims
	}
	return c, nil
}
