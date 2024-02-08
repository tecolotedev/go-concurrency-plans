package main

import (
	"errors"
	"final-project/cmd/web/data"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
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
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (appConfig *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	// Get the id of the plan that is choosen
	id := r.URL.Query().Get("id")
	planID, _ := strconv.Atoi(id)

	// Get plan for the database
	plan, err := appConfig.Models.Plan.GetOne(planID)
	if err != nil {
		appConfig.Session.Put(r.Context(), "error", "Unable to find plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	// Get the user from the session
	user, ok := appConfig.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		appConfig.Session.Put(r.Context(), "error", "Login first")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Generate an invoice and email it
	appConfig.Wait.Add(1)

	go func() {
		defer appConfig.Wait.Done()
		invoice, err := appConfig.getInvoice(user, plan)
		if err != nil {
			appConfig.ErrorChan <- err
		}
		msg := Message{
			To:       user.Email,
			Subject:  "Your Invoice",
			Data:     invoice,
			Template: "invoice",
		}
		appConfig.sendEmail(msg)
	}()

	// Generate a manual
	appConfig.Wait.Add(1)
	go func() {
		defer appConfig.Wait.Done()

		pdf := appConfig.generateManual(user, plan)
		err := pdf.OutputFileAndClose(fmt.Sprintf("./tmp/%d_manual.pdf", user.ID))
		if err != nil {
			appConfig.ErrorChan <- err
		}

		msg := Message{
			To:      user.Email,
			Subject: "Your Manual",
			Data:    "Your user manual is attached",
			AttatchmentsMap: map[string]string{
				"Manual.pdf": fmt.Sprintf("./tmp/%d_manual.pdf", user.ID),
			},
		}

		appConfig.sendEmail(msg)

		//test a error chan

		appConfig.ErrorChan <- errors.New("some custom errror")

	}()

	// Subscribe the user to an account

	// Redirect
	appConfig.Session.Put(r.Context(), "flash", "Subscribed!")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)

}
func (appConfig *Config) generateManual(user data.User, plan *data.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)
	importer := gofpdi.NewImporter()
	time.Sleep(time.Second * 5)

	t := importer.ImportPage(pdf, "./pdf/manual.pdf", 1, "/MediaBox")

	pdf.AddPage()

	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)
	pdf.SetX(75)
	pdf.SetY(150)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", user.FirstName, user.LastName), "", "C", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", user.FirstName), "", "C", false)
	return pdf
}
func (appConfig *Config) getInvoice(user data.User, plan *data.Plan) (string, error) {
	return plan.PlanAmountFormatted, nil
}

func (appConfig *Config) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	plans, err := appConfig.Models.Plan.GetAll()
	if err != nil {
		appConfig.ErrorLog.Panicln(err)
		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	appConfig.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}
