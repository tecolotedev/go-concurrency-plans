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
	mux.Get("/activate-account", appConfig.ActivateAccountPage)

	mux.Get("/test-email", func(w http.ResponseWriter, r *http.Request) {
		m := Mail{
			Domain:      "localhost",
			Host:        "localhost",
			Port:        1025,
			Encryption:  "none",
			FromAddress: "info@mycompany.com",
			FromName:    "info",

			ErrorChan: make(chan error),
		}

		msg := Message{
			To:      "me@here.com",
			Subject: "Test email",
			Data:    "hello, world",
		}
		m.sendEmail(msg, make(chan error))

	})

	return mux

}
