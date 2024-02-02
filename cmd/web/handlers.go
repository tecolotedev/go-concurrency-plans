package main

import (
	"net/http"
)

func (appConfig *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	appConfig.render(w, r, "home.page.gohtml", nil)
}

func (appConfig *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	appConfig.render(w, r, "login.page.gohtml", nil)
}

func (appConfig *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	appConfig.render(w, r, ".page.gohtml", nil)
}

func (appConfig *Config) LogoutPage(w http.ResponseWriter, r *http.Request) {
	appConfig.render(w, r, "login.page.gohtml", nil)
}

func (appConfig *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	appConfig.render(w, r, "register.page.gohtml", nil)
}

func (appConfig *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	appConfig.render(w, r, "login.page.gohtml", nil)
}

func (appConfig *Config) ActivateAccountPage(w http.ResponseWriter, r *http.Request) {
	appConfig.render(w, r, "login.page.gohtml", nil)
}
