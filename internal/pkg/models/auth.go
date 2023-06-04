package models

import "golang.org/x/crypto/bcrypt"

type Signin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Tokens   Tokens
}

type Tokens struct {
	AccessToken string `json:"accessToken"`
}

func NewSignin() Signin {
	return Signin{}
}

func (s Signin) ComparePasswordHash(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s.Password))
	return err == nil
}
