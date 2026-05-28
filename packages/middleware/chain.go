package middleware

import "net/http"

type Middlerware func(http.Handler) http.Handler

func Chain(middlerwares ...Middlerware) Middlerware {
	return func(next http.Handler) http.Handler {
		for _, MV := range middlerwares {
			next = MV(next)
		}
		return next
	}
}
