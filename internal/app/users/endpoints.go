package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/naivary/apitmpl/internal/pkg/hash"
	"github.com/naivary/apitmpl/internal/pkg/jwtauth"
	"github.com/naivary/apitmpl/internal/pkg/models"
	"github.com/pocketbase/dbx"
)

func (u Users) create(w http.ResponseWriter, r *http.Request) {
	user := models.NewUser()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := user.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPass, err := hash.Password([]byte(user.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = hashedPass
	if err := u.DB.Model(&user).Insert(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// clear any sensitive data from the struct to
	// not expose them to the public.
	user.ClearSensitive()
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u Users) single(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userID")
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := models.NewUser()
	q := u.DB.Select("id", "username", "email", "created_at").
		From("users").
		Where(dbx.HashExp{
			"id": id,
		})
	if err := q.One(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u Users) list(w http.ResponseWriter, r *http.Request) {
	var limit int64 = -1
	var users []models.User
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if params.Has("limit") && params.Get("limit") != "" {
		limit, err = strconv.ParseInt(params.Get("limit"), 10, 64)
	}
	if errors.Is(err, strconv.ErrRange) {
		http.Error(w, "limit is too big", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if limit > 0 {
		users = make([]models.User, 0, limit)
	}
	q := u.DB.Select("id", "username", "email", "created_at").From("users").Limit(int64(limit))
	if err := q.All(&users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u Users) delete(w http.ResponseWriter, r *http.Request) {
	claims, err := jwtauth.GetClaims(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sub, err := claims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := u.DB.Delete("users", dbx.HashExp{"id": sub}).Execute(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
