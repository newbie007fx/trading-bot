package httpserver

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"telebot-trading/app/helper"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	dataMap := data.(map[string]interface{})
	dataMap["csrf_token"] = c.Get(echoMiddleware.DefaultCSRFConfig.ContextKey).(string)

	tmpl := t.templates.Funcs(template.FuncMap{
		"session": getFlashMessage(c),
	})

	return tmpl.ExecuteTemplate(w, name, dataMap)
}

func getFlashMessage(c echo.Context) func(string) *string {
	return func(key string) *string {
		sess := helper.GetSession(c)
		return sess.GetFlashMessage(key)
	}
}

func getTemplate() *Template {
	templ := template.New("").Funcs(template.FuncMap{
		"session": func(key string) *string {
			return nil
		},
	})

	err := filepath.Walk("./resources/views", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gohtml") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		panic(err)
	}

	return &Template{
		templates: templ,
	}
}
