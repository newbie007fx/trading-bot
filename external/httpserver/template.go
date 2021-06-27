package httpserver

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"telebot-trading/app/helper"
	"telebot-trading/utils"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Template struct{}

func (t *Template) Render(w io.Writer, filename string, data interface{}, c echo.Context) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	tmpData := data.(map[string]interface{})
	tmpData["csrf_token"] = c.Get(echoMiddleware.DefaultCSRFConfig.ContextKey).(string)
	tmpData["param"] = templateParam{}

	resourcePath := utils.Env("RESOURCES_PATH", "resources")
	filename = fmt.Sprintf("%s/views/%s", resourcePath, filename)
	baseLayout := fmt.Sprintf("%s/views/layouts/app.gohtml", resourcePath)

	tmpl := template.New("")
	_, err := tmpl.Funcs(template.FuncMap{
		"session": getFlashMessage(c),
	}).ParseFiles(filename, baseLayout)
	if err != nil {
		panic(err)
	}

	return tmpl.ExecuteTemplate(w, filepath.Base(filename), tmpData)
}

type templateParam map[string]string

func (tmpParam templateParam) Add(key string, val string) interface{} {
	tmpParam[key] = val
	return nil
}

func (tmpParam templateParam) Get(key string) string {
	return tmpParam[key]
}

func getFlashMessage(c echo.Context) func(string) *string {
	return func(key string) *string {
		sess := helper.GetSession(c)
		return sess.GetFlashMessage(key)
	}
}
