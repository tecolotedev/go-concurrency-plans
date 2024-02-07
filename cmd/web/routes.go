package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (appConfig *Config) routes() http.Handler {
	//create router
	mux := chi.NewRouter()

	// setup middleware
	mux.Use(middleware.Recoverer)
	mux.Use(appConfig.SessionLoad)

	// define application routes
	mux.Get("/", appConfig.HomePage)

	mux.Get("/login", appConfig.LoginPage)
	mux.Post("/login", appConfig.PostLoginPage)
	mux.Get("/logout", appConfig.Logout)
	mux.Get("/register", appConfig.RegisterPage)
	mux.Post("/register", appConfig.PostRegisterPage)
	mux.Get("/activate", appConfig.ActivateAccountPage)

	mux.Get("/plans", appConfig.ChooseSubscription)
	mux.Get("/subscribe", appConfig.SubscribeToPlan)

	return mux

}
