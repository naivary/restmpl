package users

import (
	"encoding/json"
	"net/http"

	"github.com/naivary/apitmpl/internal/pkg/models"
)

func (u Users) create(w http.ResponseWriter, r *http.Request) {
	user := models.NewUser()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	u.DB.Model(&user).Insert()
}
