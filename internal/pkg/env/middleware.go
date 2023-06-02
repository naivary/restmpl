package env

import "net/http"

func (a API) middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{}
}
