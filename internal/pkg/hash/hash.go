package hash

import "golang.org/x/crypto/bcrypt"

func Password(raw []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(raw, bcrypt.DefaultCost)
	return string(bytes), err
}
