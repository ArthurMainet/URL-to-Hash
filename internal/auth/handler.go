package auth

import (
	"fmt"
	"golang/configs"
	"golang/packages/jwt"
	"golang/packages/req"
	"golang/packages/res"
	"net/http"
)

type AuthHandlerDeps struct {
	Config      *configs.Config
	AuthService *AuthService
}

type AuthHandler struct {
	Config      *configs.Config
	AuthService *AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	auth := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", auth.Login())
	router.HandleFunc("POST /auth/register", auth.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{Email: email})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		res.Json(w, token, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		check, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		fmt.Println(check)

	}
}
