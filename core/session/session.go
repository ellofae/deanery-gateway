package session

import (
	"github.com/ellofae/deanery-gateway/config"
	"github.com/gorilla/sessions"
)

var sessionStorage *sessions.CookieStore

func InitSessionStorage(cfg *config.Config) {
	sessionStorage = sessions.NewCookieStore([]byte(cfg.Session.SessionKey))

	sessionStorage.Options.HttpOnly = true
}

func SessionStorage() *sessions.CookieStore {
	return sessionStorage
}
