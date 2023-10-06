package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/radu0v/1quoteEveryDay/internal/config"
	"github.com/radu0v/1quoteEveryDay/internal/models"
)

var app *config.AppConfig

var functions = template.FuncMap{
	"inc": inc,
}

func NewRenderer(a *config.AppConfig) {
	app = a
}

func addDefaultData(data *models.Data, r *http.Request) *models.Data {
	data.CSRFToken = nosurf.Token(r)
	return data
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.Data) error {
	var tc map[string]*template.Template

	if app.UseCache {
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// for testing purposes , it rebuilds the cache on
		// every request
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}
	buf := new(bytes.Buffer)
	data = addDefaultData(data, r)
	err := t.Execute(buf, data)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser: ", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

func inc(i int) int {
	return i + 1
}
