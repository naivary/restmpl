package auth

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/naivary/apitmpl/internal/pkg/random"
)

func (a Auth) signin(w http.ResponseWriter, r *http.Request) {
	t, err := a.jwt.NewSignedToken(random.ID(6), jwt.MapClaims{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(t))
}
func (a Auth) signup(w http.ResponseWriter, r *http.Request) {}
