package main

import "net/http"

func (appConfig *Config) SessionLoad(next http.Handler) http.Handler {
	return appConfig.Session.LoadAndSave(next)
}
