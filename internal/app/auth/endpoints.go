package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (a Auth) signin(w http.ResponseWriter, r *http.Request) {
	_, err := a.jwt.NewSignedToken("must", jwt.MapClaims{})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
func (a Auth) signup(w http.ResponseWriter, r *http.Request) {}
