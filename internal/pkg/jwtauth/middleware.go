package jwtauth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Verify the provided Bearer JWT token and set the claims
// of the token in the request context with the key `claimsKey`
func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", 1))
		token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid alg")
			}
			return secret, nil
		})
		if err := checkTokenValidationErr(err); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("token is invalid"))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("claims in wrong format"))
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), claimsKey, claims))
		next.ServeHTTP(w, r)
	})
}
