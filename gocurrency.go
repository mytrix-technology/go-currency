package main

import (
	"flag"
	"go-currency/modul"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

func init() {
	//Logging file name and line number
	log.SetFlags(log.Lshortfile)
}

func main() {
	route()

	var port = flag.String("port", "3000", "isi port")
	flag.Parse()
	var timenow = time.Now().Format(time.RFC822)
	c := color.New(color.FgGreen).Add(color.Underline)
	c.Printf("%s %s %s", timenow, "Berjalan di port", *port)
	http.ListenAndServe(":3000"+*port, nil)
}

func route() {
	//Insert Data
	http.HandleFunc("/apiInsertData", modul.ApiInsertData)
	//List Exchange and Average 7 Days
	http.HandleFunc("/apiListData", modul.ApiListData)
	//Search Data Exchange from most recent 7 data points
	http.HandleFunc("/apiListDataPoints", modul.ApiListDataPoints)
	//Insert Data Sysmbols
	http.HandleFunc("/apiInsertDataSymbols", modul.ApiInsertDataSymbols)
	//Delete Data
	http.HandleFunc("/apiDeleteData", modul.ApiDeleteData)
}
