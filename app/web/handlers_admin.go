package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func admin_panel(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, _ := getCookie(r, "session")
	if cookie.Value == "admin" {
		rows, err := db.Query("SELECT `login`, `balance` from `users`")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.Login, &u.Balance); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		files := []string{
			"html/tmpl/admin_panel.html",
			"html/tmpl/base.layout.html",
			"html/tmpl/footer.partial.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = ts.Execute(w, users)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	} else {
		w.Write([]byte("Недостаточно прав!"))
	}

}

func balance_change(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	login := r.URL.Query().Get("login")
	balance, _ := strconv.ParseFloat(r.FormValue("account_balance"), 32)
	query := "UPDATE `users` SET `balance` = ? WHERE `login` = ?"
	update, err := db.Exec(query, balance, login)
	if err != nil {
		log.Println(err)
	}

	rowsAffected, err := update.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	if rowsAffected > 0 {
		http.Redirect(w, r, "/admin_panel", http.StatusFound)
	}

}
