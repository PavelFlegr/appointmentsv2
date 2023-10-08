package global

import (
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
)

type EmailConfig struct {
	Host     string
	From     string
	User     string
	Password string
	Port     int
}

type Config struct {
	Port    int
	ConnStr string
	HashKey string
	Email   EmailConfig
	Host    string
}

var Db *sqlx.DB
var Sc *securecookie.SecureCookie
var Conf *Config
