package middleware

import (
	"context"
	"fmt"
	"golang/configs"
	"golang/packages/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmail key = "ContextEmail"
)

func WriteUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuth(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		phrase := r.Header.Get("Authorization")
		if !strings.HasPrefix(phrase, "Bearer ") {
			WriteUnauthed(w)
			return
		}
		token := strings.TrimPrefix(phrase, "Bearer ")
		fmt.Println(token)
		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)
		if !isValid {
			WriteUnauthed(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmail, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
