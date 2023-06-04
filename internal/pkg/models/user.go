package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func NewUser() User {
	return User{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().Unix(),
	}
}

// implementing dbx.TableModel
func (u User) TableName() string {
	return "users"
}

// CheckPassword ensures strength and character constraints
func (u User) CheckPassword() error {
	pattern := `^[[:ascii:]]\S+$`
	if !(len(u.Password) >= 10) {
		return errors.New("password has to be longer than 10 characters")
	}
	if ok, err := regexp.Match(pattern, []byte(u.Password)); !ok || err != nil {
		e := fmt.Errorf("password not matching pattern: %s", pattern)
		return errors.Join(e, err)
	}
	return nil
}
