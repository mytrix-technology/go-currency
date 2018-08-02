package config

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Parser must implement ParseJSON
type Parser interface {
	ParseJSON([]byte) error
}

// Load the JSON config file
func Load(configFile string, p Parser) {
	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		log.Fatalln(err)
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// Parse the config
	if err := p.ParseJSON(jsonBytes); err != nil {
		log.Fatalln("Could not parse %q: %v", configFile, err)
	}
}

//Function Insert Log Activity
func Logs(action, table, change, users string) {
	var db = Connect()
	defer db.Close()
	var insertlog, _ = db.Prepare("INSERT INTO logs(logs_action,logs_table,logs_change,logs_cdate,logs_cuser) VALUES($1,$2,$3,NOW(),$4)")
	insertlog.Exec(action, table, change, users)
}

//Function string to error validation
func Strings(s string) string {
	var r string
	switch s {
	case "limit":
		r = "LIMIT 1000"
	case "errmethod":
		r = "Error Method"
	case "wrongtoken":
		r = "Wrong Token Key and Secret"
	case "expiredtoken":
		r = "Expired Token"
	case "inactivetoken":
		r = "Inactive Token"
	case "failedauth":
		r = "Authentication Failed"
	case "failedupdate":
		r = "Failed To Update : "
	case "failedinsert":
		r = "Failed To Insert : "
	case "faileddelete":
		r = "Failed To Delete : "
	case "failedurl":
		r = "Failed Url Type(add/edit)"
	case "NA":
		r = "Not Available"
	case "restdelforeign":
		r = "This ID has been foreign key"
	default:
		r = "Unidentified"
	}
	return r
}
