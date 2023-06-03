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

// Verify the provided Bearer JWT token and set the claims
// of the token in the request context with the key `claimsKey`
func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", 1))
		token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid alg")
			}
			return secret, nil
		})
		if err := checkTokenValidationErr(err); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
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

func GetClaims(ctx context.Context) (jwt.MapClaims, error) {
	val := ctx.Value(claimsKey)
	c, ok := val.(jwt.MapClaims)
	if !ok {
		return nil, ErrMissingClaims
	}
	return c, nil
}

func (j JWTAuth) addRegisteredClaims(sub any, m jwt.MapClaims) {
	m["sub"] = sub
	m["iss"] = j.issuer
	m["exp"] = time.Now().Add(j.expDuration).Unix()
	m["iat"] = time.Now().Unix()
	m["jti"] = random.ID(defaultIDLen)
}

func checkTokenValidationErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, jwt.ErrTokenExpired) {
		return errors.New("token is expired and should be renewed using the refresh token")
	}
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return jwt.ErrSignatureInvalid
	}
	return err
}
