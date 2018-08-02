package modul

import (
	"encoding/json"
	"go-currency/config"
	"log"
	"net/http"
	"strings"
)

func ApiInsertData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application-json")

	if r.Method == "POST" {
		var db = config.Connect()
		defer db.Close()

		var date = strings.Replace(r.FormValue("date"), "'", "\\'", -1)
		var iso_code_from = strings.Replace(r.FormValue("iso_code_from"), "'", "\\'", -1)
		var iso_code_to = strings.Replace(r.FormValue("iso_code_to"), "'", "\\'", -1)
		var rate = strings.Replace(r.FormValue("rate"), "'", "\\'", -1)

		var sql, err = db.Prepare(`INSERT INTO list_currencies (date, iso_code_from, iso_code_to, rate)
									 VALUES($1,$2,$3,$4)`)
		_, errs := sql.Exec(date, iso_code_from, iso_code_to, rate)
		config.Logs("INSERT", "list_currencies", "date="+date+",iso_code_from="+iso_code_from+",iso_code_to="+iso_code_to+",rate="+rate, "")

		if err != nil {
			log.Print(err.Error())
		}

		if errs != nil {
			http.Error(w, errs.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ApiListData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application-json")

	if r.Method == "GET" {
		type DataList struct {
			Date        string `json:"date"`
			IsoCodeFrom string `json:"iso_code_from"`
			IsoCodeTo   string `json:"iso_code_to"`
			Rate        string `json:"rate"`
		}
		var res = []DataList{}
		var each = DataList{}

		var db = config.Connect()
		defer db.Close()
		var date = strings.Replace(r.FormValue("date"), "'", "\\'", -1)
		var iso_code_from = strings.Replace(r.FormValue("iso_code_from"), "'", "\\'", -1)
		var iso_code_to = strings.Replace(r.FormValue("iso_code_to"), "'", "\\'", -1)
		var rate = strings.Replace(r.FormValue("rate"), "'", "\\'", -1)
		var order_by = strings.Replace(r.FormValue("order_by"), "'", "\\'", -1)
		var order_type = strings.Replace(r.FormValue("order_type"), "'", "\\'", -1)
		var limit = strings.Replace(r.FormValue("limit"), "'", "\\'", -1)
		var partWhere, where, limits, orderby string
		limits = config.Strings("limit")
		orderby = ""
		partWhere = ""

		if date != "" {
			if partWhere == "" {
				partWhere = "date ='" + date + "'"
			}
		}
		if iso_code_from != "" {
			if partWhere == "" {
				partWhere = "iso_code_from ='" + iso_code_from + "'"
			}
		}
		if iso_code_to != "" {
			if partWhere == "" {
				partWhere = "iso_code_to ='" + iso_code_to + "'"
			}
		}
		if rate != "" {
			if partWhere == "" {
				partWhere = "rate ='" + rate + "'"
			}
		}
		if partWhere == "" {
			where = ""
		} else {
			where = "WHERE " + partWhere
		}
		if order_by != "" && order_type != "" {
			orderby = "ORDER BY " + order_by + " " + order_type
		}

		if limit != "" {
			limits = "LIMIT " + limit
		}

		var sql = `select p.iso_code_from, p.iso_code_to, d.date, sum(usagecount),
					sum(sum(usagecount)) over (partition by p.iso_code_from order by d.date 
						rows between 6 preceding and current row) as Sum7day
 					from (select distinct iso_code_from from list_currencies) p cross join
					  (select distinct date from list_currencies) d
					  left join list_currencies h on h.iso_code_from = p.iso_code_from
					and h.date = p.date ` + where + " group by p.iso_code_from, d.date " + orderby + " " + limits
		col, errsql := db.Query(sql)

		if errsql != nil {
			res = append(res, each)
		} else {
			for col.Next() {
				var errs = col.Scan(&each.Date, &each.IsoCodeFrom, &each.IsoCodeTo, &each.Rate)
				if errs != nil {
					log.Print(errs.Error())
				}
				res = append(res, each)
			}
		}

		var result, error = json.Marshal(res)
		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
		return
	}
}

func ApiListDataPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application-json")

	if r.Method == "POST" {
		type DataPointList struct {
			Date        string `json:"date"`
			IsoCodeFrom string `json:"iso_code_from"`
			IsoCodeTo   string `json:"iso_code_to"`
			Rate        string `json:"rate"`
		}
		var res = []DataPointList{}
		var each = DataPointList{}

		var db = config.Connect()
		defer db.Close()
		var date = strings.Replace(r.FormValue("date"), "'", "\\'", -1)
		var iso_code_from = strings.Replace(r.FormValue("iso_code_from"), "'", "\\'", -1)
		var iso_code_to = strings.Replace(r.FormValue("iso_code_to"), "'", "\\'", -1)
		var rate = strings.Replace(r.FormValue("rate"), "'", "\\'", -1)
		var order_by = strings.Replace(r.FormValue("order_by"), "'", "\\'", -1)
		var order_type = strings.Replace(r.FormValue("order_type"), "'", "\\'", -1)
		var limit = strings.Replace(r.FormValue("limit"), "'", "\\'", -1)
		var partWhere, where, limits, orderby string
		limits = config.Strings("limit")
		orderby = ""
		partWhere = ""

		if date != "" {
			if partWhere == "" {
				partWhere = "date ='" + date + "'"
			}
		}
		if iso_code_from != "" {
			if partWhere == "" {
				partWhere = "iso_code_from ='" + iso_code_from + "'"
			}
		}
		if iso_code_to != "" {
			if partWhere == "" {
				partWhere = "iso_code_to ='" + iso_code_to + "'"
			}
		}
		if rate != "" {
			if partWhere == "" {
				partWhere = "rate ='" + rate + "'"
			}
		}
		if partWhere == "" {
			where = ""
		} else {
			where = "WHERE " + partWhere
		}
		if order_by != "" && order_type != "" {
			orderby = "ORDER BY " + order_by + " " + order_type
		}

		if limit != "" {
			limits = "LIMIT " + limit
		}

		var sql = `SELECT max(rate) as maxrate, min(rate) as minrate, max(rate)-min(rate) as variance, avg(rate) as average, date, rate
					FROM list_currencies
					WHERE iso_code_from = $2 AND iso_code_to = $3 AND date BETWEEN $1 AND $1 ` + where + " GROUP BY 1,2 " + orderby + " " + limits
		col, errsql := db.Query(sql)

		if errsql != nil {
			res = append(res, each)
		} else {
			for col.Next() {
				var errs = col.Scan(&each.Date, &each.IsoCodeFrom, &each.IsoCodeTo, &each.Rate)
				if errs != nil {
					log.Print(errs.Error())
				}
				res = append(res, each)
			}
		}

		var result, error = json.Marshal(res)
		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
		return
	}
}

func ApiInsertDataSymbols(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application-json")

	if r.Method == "POST" {
		var db = config.Connect()
		defer db.Close()

		var iso_code_from = strings.Replace(r.FormValue("iso_code_from"), "'", "\\'", -1)
		var iso_code_to = strings.Replace(r.FormValue("iso_code_to"), "'", "\\'", -1)

		var sql, err = db.Prepare(`INSERT INTO list_currencies (iso_code_from, iso_code_to)
									 VALUES($2,$3)`)
		_, errs := sql.Exec(iso_code_from, iso_code_to)
		config.Logs("INSERT", "list_currencies", "iso_code_from="+iso_code_from+",iso_code_to="+iso_code_to, "")

		if err != nil {
			log.Print(err.Error())
		}

		if errs != nil {
			http.Error(w, errs.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ApiDeleteData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application-json")

	if r.Method == "POST" {
		var iso_code_from = strings.Replace(r.FormValue("iso_code_from"), "'", "\\'", -1)
		var iso_code_to = strings.Replace(r.FormValue("iso_code_to"), "'", "\\'", -1)
		var db = config.Connect()
		db.Close()

		var sql = "DELETE FROM list_currencies WHERE iso_code_from = $2 AND iso_code_to = $3"
		var prepare, err = db.Prepare(sql)
		_, err = prepare.Exec(iso_code_from, iso_code_to)
		config.Logs("DELETE", "list_currencies", "iso_code_from="+iso_code_from+"iso_code_to="+iso_code_to, "")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
