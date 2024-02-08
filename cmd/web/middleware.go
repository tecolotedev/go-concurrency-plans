package main

import (
	"net/http"
)

func (appConfig *Config) SessionLoad(next http.Handler) http.Handler {
	return appConfig.Session.LoadAndSave(next)
}

func (appConfig *Config) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if appConfig.Session.Exists(r.Context(), "userid") {
			appConfig.Session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}

		next.ServeHTTP(w, r)
	})
}
