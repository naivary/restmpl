package jwtauth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/internal/pkg/random"
)

type ctxKey string

const (
	claimsKey    ctxKey = "claims"
	defaultIDLen int    = 5
)

var (
	secret           []byte
	ErrMissingClaims = errors.New("no claims found in context")
)

type JWTAuth struct {
	signingMethod jwt.SigningMethod
	issuer        string
	expDuration   time.Duration
}

func SetSecret(s string) {
	secret = []byte(s)
}

func New(k *koanf.Koanf) *JWTAuth {
	return &JWTAuth{
		signingMethod: jwt.SigningMethodHS256,
		issuer:        k.String("jwt.issuer"),
		expDuration:   k.Duration("jwt.expiration"),
	}
}

// CustomClaims generates public or private claims in the syntax
// key, value, key, value and so on.
func (j JWTAuth) CustomClaims(args ...any) (jwt.Claims, error) {
	m := jwt.MapClaims{}
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			return nil, fmt.Errorf("%v is not a string", key)
		}
		m[key] = args[i+1]
	}
	return m, nil
}

func (j JWTAuth) NewSignedToken(sub any, claims jwt.MapClaims) (string, error) {
	m := jwt.MapClaims{}
	for key, value := range claims {
		m[key] = value
	}
	j.addRegisteredClaims(sub, m)
	return jwt.NewWithClaims(j.signingMethod, m).SignedString(secret)
}

func (j JWTAuth) addRegisteredClaims(sub any, m jwt.MapClaims) {
	m["sub"] = sub
	m["iss"] = j.issuer
	m["exp"] = time.Now().Add(j.expDuration).Unix()
	m["iat"] = time.Now().Unix()
	m["jti"] = random.ID(defaultIDLen)
}
