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
	ID        string `json:"id,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
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

// ClearSensitive clears all the data of
// the sensitive fileds like password
func (u *User) ClearSensitive() {
	u.Password = ""
}

// IsValid ensures that the encoded
// data of the struct are valid e.g.
// password strength, email validation...
func (u User) IsValid() error {
	checks := []func() error{
		u.isValidEmail,
		u.isValidPassword,
		u.isValidUsername,
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

func (u User) isValidUsername() error {
	pattern := `^[\w\.][^0-9]+$`
	if ok, err := regexp.Match(pattern, []byte(u.Username)); !ok || err != nil {
		e := fmt.Errorf("username not matching pattern: %s", pattern)
		return errors.Join(e, err)
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
