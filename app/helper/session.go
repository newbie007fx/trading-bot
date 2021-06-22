package helper

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var sesHelper *sessionHelper

func GetSession(c echo.Context) *sessionHelper {
	if sesHelper == nil {
		sesHelper = new(sessionHelper)
	}

	sesHelper.Setup(c)

	return sesHelper
}

type sessionHelper struct {
	ctx  echo.Context
	sess *sessions.Session
}

func (helper *sessionHelper) Setup(ctx echo.Context) {
	helper.ctx = ctx

	sess, err := session.Get("session", helper.ctx)
	if err != nil {
		log.Println(err.Error())
	}

	helper.sess = sess
}

func (helper *sessionHelper) AddFlashMessage(key, msg string) {
	helper.sess.AddFlash(msg, key)

	helper.sess.Save(helper.ctx.Request(), helper.ctx.Response())
}

func (helper *sessionHelper) GetFlashMessage(key string) *string {
	defer helper.sess.Save(helper.ctx.Request(), helper.ctx.Response())

	msgs := helper.sess.Flashes(key)
	if len(msgs) == 0 {
		return nil
	}

	msg := msgs[0].(string)

	return &msg
}

func (helper *sessionHelper) Set(data map[string]interface{}) {
	for key, val := range data {
		helper.sess.Values[key] = val
	}
	helper.sess.Save(helper.ctx.Request(), helper.ctx.Response())
}

func (helper *sessionHelper) Get(key string) interface{} {
	if val, ok := helper.sess.Values[key]; ok {
		return val
	}

	return nil
}

func (helper *sessionHelper) Destroy() {
	helper.sess.Values = map[interface{}]interface{}{}
	helper.sess.Save(helper.ctx.Request(), helper.ctx.Response())
}
