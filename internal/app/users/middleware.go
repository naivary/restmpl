package users

import "net/http"

func (u Users) middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		u.incomingMetricCounter,
	}
}

func (u Users) incomingMetricCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := u.m.GetCounter("incomingReq")
		c.Inc()
		next.ServeHTTP(w, r)
	})
}
