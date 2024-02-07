package main

import (
	"final-project/cmd/web/data"
	"fmt"
	"html/template"
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
		msg := Message{
			To:      email,
			Subject: "Failed login attempt",
			Data:    "Invalid login attempt",
		}
		fmt.Println("here1")
		appConfig.sendEmail(msg)
		fmt.Println("here2")

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
	err := r.ParseForm()
	if err != nil {
		appConfig.ErrorLog.Println(err)
	}
	// TODO validate data

	// create user
	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
		IsAdmin:   0,
	}

	_, err = u.Insert(u)

	if err != nil {
		appConfig.Session.Put(r.Context(), "error", "Unable to create user.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send an activation email
	url := fmt.Sprintf("http://localhost/activate?email=%s", u.Email)
	signedUrl := GenerateTokenFromString(url)
	appConfig.InfoLog.Println(signedUrl)

	msg := Message{
		To:       u.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedUrl),
	}

	appConfig.sendEmail(msg)

	appConfig.Session.Put(r.Context(), "flash", "Confirmation email sent. Check your email")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (appConfig *Config) ActivateAccountPage(w http.ResponseWriter, r *http.Request) {
	// validate url
	url := r.RequestURI
	testUrl := fmt.Sprintf("http://localhost%s", url)
	ok := VerifyToken(testUrl)

	if !ok {
		appConfig.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// activate account
	u, err := appConfig.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		appConfig.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u.Active = 1
	err = u.Update()
	if err != nil {
		appConfig.Session.Put(r.Context(), "error", "Unable to update user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	appConfig.Session.Put(r.Context(), "flash", "Account activated")
	http.Redirect(w, r, "/", http.StatusSeeOther)

	// generate an invoice

	// send the email with attatchments

	// send an email with the invoice attatched
}
