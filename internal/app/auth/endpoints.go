package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (a Auth) signin(w http.ResponseWriter, r *http.Request) {
	t, err := a.jwt.NewSignedToken("something", jwt.MapClaims{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(t))
}

func (a Auth) signup(w http.ResponseWriter, r *http.Request) {}
