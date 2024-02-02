package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

type TemplateData = struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	Data          map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	// User *data.User
}

func (appConfig *Config) render(
	w http.ResponseWriter,
	r *http.Request,
	t string, //tempalte
	td *TemplateData,
) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}

	templateSlice := []string{}
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplates, t))

	templateSlice = append(templateSlice, partials...)
	// for _, x := range partials {
	// 	templateSlice = append(templateSlice, x)
	// }

	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		appConfig.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, appConfig.AddDefaultData(td, r)); err != nil {
		appConfig.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (appConfig *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = appConfig.Session.PopString(r.Context(), "flash")
	td.Warning = appConfig.Session.PopString(r.Context(), "warning")
	td.Error = appConfig.Session.PopString(r.Context(), "error")
	if appConfig.isAuthenticated(r) {
		td.Authenticated = true
		// TODO: get more user information
	}
	td.Now = time.Now()

	return td

}

func (appConfig *Config) isAuthenticated(r *http.Request) bool {
	return appConfig.Session.Exists(r.Context(), "userID")
}
