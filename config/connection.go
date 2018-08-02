package config

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	UrlFiles = "config" + string(os.PathSeparator) + "config.json"
)

var configs = &configuration{}

// configuration contains the application settings
type configuration struct {
	DB_username string `json:"DB_username"`
	DB_password string `json:"DB_password"`
	DB_name     string `json:"DB_name"`
	DB_host     string `json:"DB_host"`
	DB_SSL      string `json:"DB_SSL"`
	FILE_URL    string `json:"FILE_URL"`
}

func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

func UrlSSO() string {
	Load(UrlFiles, configs)
	var Url string
	Url = configs.FILE_URL
	return Url
}

func ConfigURL() string {
	Load(UrlFiles, configs)
	var Url string
	Url = configs.FILE_URL
	return Url
}

func Connect() *sql.DB {
	Load(UrlFiles, configs)
	dbDriver := configs.DB_username
	dbUser := configs.DB_username
	dbURL := configs.DB_host
	// dbPort   := "3306"
	dbPass := configs.DB_password
	dbName := configs.DB_name
	dbSSL := configs.DB_SSL
	db, err := sql.Open(dbDriver, "postgres://"+dbUser+":"+dbPass+"@"+dbURL+"/"+dbName+"?sslmode="+dbSSL)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
