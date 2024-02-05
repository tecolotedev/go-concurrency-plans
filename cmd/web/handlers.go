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
	appConfig.Session.RenewToken(r.Context())

	// parse form post
	err := r.ParseForm()
	if err != nil {
		appConfig.ErrorLog.Println(err)
	}

	// get email  and password from post form
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := appConfig.Models.User.GetByEmail(email)
	if err != nil {
		appConfig.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// check password
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		appConfig.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !validPassword {
		appConfig.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// ok, so user log in
	appConfig.Session.Put(r.Context(), "userID", user.ID)
	appConfig.Session.Put(r.Context(), "user", user)
	appConfig.Session.Put(r.Context(), "flash", "Successsful login")

	//redirect the user
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (appConfig *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean all session
	appConfig.Session.Destroy(r.Context())
	appConfig.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (appConfig *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	appConfig.render(w, r, "register.page.gohtml", nil)
}

func (appConfig *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	//create user

	// send an activation email

	// subscribe the user to an account

}

func (appConfig *Config) ActivateAccountPage(w http.ResponseWriter, r *http.Request) {
	// validate url

	// generate an invoice

	// send the email with attatchments

	// send an email with the invoice attatched
}
