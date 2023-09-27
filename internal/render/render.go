package render

import (
	"html/template"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/radu0v/1quoteEveryDay/internal/models"
)

func addDefaultData(data *models.Data, r *http.Request) *models.Data {
	data.CSRFToken = nosurf.Token(r)
	return data
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.Data) {
	page := "./templates/" + tmpl
	var funcMap = template.FuncMap{
		"inc": inc,
	}
	tpl, err := template.New(page).Funcs(funcMap).ParseFiles(page, "./templates/base.layout.tmpl", "./templates/admin.layout.tmpl")
	if err != nil {
		log.Println(err)
	}

	data = addDefaultData(data, r)
	err = tpl.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatal(err)
	}
}

func inc(i int) int {
	return i + 1
}
