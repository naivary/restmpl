package users

import (
	"encoding/json"
	"net/http"

	"github.com/naivary/apitmpl/internal/pkg/models"
)

func (u Users) create(w http.ResponseWriter, r *http.Request) {
	user := models.NewUser()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := user.CheckPassword(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := u.DB.Model(&user).Insert(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
