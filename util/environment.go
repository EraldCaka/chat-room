package util

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DB_CONN_STR string
var SECRET_KEY string
var COOKIE_NAME string
var COOKIE_VALUE string

func InitEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	DB_CONN_STR = os.Getenv("DB_CONN_STR")
	SECRET_KEY = os.Getenv("SECRET_KEY")
}

func StoreCookieToEnv(cookieName, cookieValue string) error {
	err := godotenv.Write(map[string]string{
		"COOKIE_NAME":  cookieName,
		"COOKIE_VALUE": cookieValue,
	}, "cookies.env")
	if err != nil {
		return err
	}
	return nil
}
func EmptyCookieEnv() error {
	err := godotenv.Write(map[string]string{
		"COOKIE_NAME":  "",
		"COOKIE_VALUE": "",
	}, "cookies.env")
	if err != nil {
		return err
	}
	return nil
}
