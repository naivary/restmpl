package jwtauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/knadh/koanf/v2"
)

type ctxKey string

const (
	claimsKey ctxKey = "claims"
)

var (
	secret []byte
)

type JWTAuth struct {
	signingMethod jwt.SigningMethod
	issuer        string
	exp           time.Duration
	defaultClaims jwt.MapClaims
}

func SetSecret(s []byte) {
	secret = s
}

func New(k *koanf.Koanf) *JWTAuth {
	return &JWTAuth{
		signingMethod: jwt.SigningMethodHS256,
		issuer:        k.String("jwt.issuer"),
		exp:           k.Duration("jwt.expiration"),
		defaultClaims: jwt.MapClaims{
			"iss": k.String("jwt.issuer"),
		},
	}
}

func (j JWTAuth) GenClaims(args ...any) (jwt.Claims, error) {
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
	for key, value := range j.defaultClaims {
		m[key] = value
	}
	m["sub"] = sub
	m["exp"] = time.Now().Add(5 * time.Minute)
	m["iat"] = time.Now().Unix()
	return jwt.NewWithClaims(j.signingMethod, m).SignedString(secret)
}

// Verify the provided Bearer JWT token and set the claims
// of the token in the request context with the key `claims`
func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", -1))
		token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid alg")
			}
			return secret, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("error while parsing the token"))
			fmt.Println(err)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("token is invalid"))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("claims in wrong format"))
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), claimsKey, claims))
		next.ServeHTTP(w, r)
	})
}
