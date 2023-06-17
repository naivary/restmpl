package auth

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/naivary/restmpl/internal/pkg/models"
	"github.com/pocketbase/dbx"
)

func (a Auth) signin(w http.ResponseWriter, r *http.Request) {
	si := models.NewSignin()
	if err := json.NewDecoder(r.Body).Decode(&si); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.NewUser()
	q := a.DB.Select("id", "created_at", "email", "username", "password").From("users").Where(dbx.HashExp{
		"email": si.Email,
	})

	if err := q.One(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !si.ComparePasswordHash(user.Password) {
		http.Error(w, "wrong password or email", http.StatusBadRequest)
		return
	}

	acToken, err := a.jwt.NewSignedToken(user.ID, jwt.MapClaims{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	si.Tokens.AccessToken = acToken
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&si.Tokens); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
