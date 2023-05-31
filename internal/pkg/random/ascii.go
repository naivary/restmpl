package random

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	Whitespace    = 32
	QuotationMark = 34
	Backslash     = 92
	LessThanSign  = 60
	Delete        = 127
)

// https://gist.github.com/denisbrodbeck/635a644089868a51eccd6ae22b2eb800 (source)
func ascii(length int) (string, error) {
	result := make([]byte, length)
	_, err := rand.Read(result)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		result[i] &= 0x7F
		for result[i] <= Whitespace || result[i] == Delete || result[i] == QuotationMark || result[i] == Backslash || result[i] == LessThanSign {
			_, err = rand.Read(result[i : i+1])
			if err != nil {
				return "", err
			}
			result[i] &= 0x7F
		}
	}

	return string(result), nil
}

func ID(len int) string {
	id, _ := ascii(len)
	return hex.EncodeToString([]byte(id))
}
