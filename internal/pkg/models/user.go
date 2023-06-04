package models

import (
	"errors"
	"fmt"
	"net/mail"
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

// IsValid ensures that the encoded
// data of the struct are valid
func (u User) IsValid() error {
	checks := []func() error{
		u.isValidEmail,
		u.isValidPassword,
	}
	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}
	return nil
}

func (u User) isValidEmail() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return err
	}
	return nil
}

// isValidPassword ensures strength and character constraints
func (u User) isValidPassword() error {
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
