package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func setCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func sign_up(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		login_in := r.FormValue("login_si")
		pass_in := r.FormValue("pass_si")
		query := fmt.Sprintf("SELECT COUNT(*) FROM `users` WHERE `login` = '%s';", login_in)
		res, err := db.Query(query)
		if err != nil {
			http.Error(w, "Error executing query", http.StatusInternalServerError)
			return
		}
		defer res.Close()

		var count int
		for res.Next() {
			err := res.Scan(&count)
			if err != nil {
				http.Error(w, "Error scanning result", http.StatusInternalServerError)
				return
			}
		}

		if count > 0 {

			return
		}

		page := rand.Intn(100000)
		query = fmt.Sprintf("INSERT INTO `users` (`login`, `password`, `page`) VALUES ('%s','%s','%d');", login_in, pass_in, page)
		insert, err := db.Query(query)
		if err != nil {
			log.Println(err)
		}

		defer insert.Close()
		setCookie(w, "session", login_in)
		http.Redirect(w, r, "/sign_in", http.StatusMovedPermanently)

	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "html/tmpl/sign_up.html")
	} else {

		http.Error(w, "Метод запрещен!", http.StatusMethodNotAllowed)
		return
	}

}

func sign_in(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		login_up := r.FormValue("login")
		pass_up := r.FormValue("pass")
		query := fmt.Sprintf("SELECT COUNT(*) FROM `users` WHERE `login` = '%s' ;", login_up)
		res, err := db.Query(query)
		if err != nil {
			http.Error(w, "Error executing query", http.StatusInternalServerError)
			return
		}
		defer res.Close()

		var count int
		for res.Next() {
			err := res.Scan(&count)
			if err != nil {
				http.Error(w, "Error scanning result", http.StatusInternalServerError)
				return
			}
		}

		if count != 1 {
			http.Redirect(w, r, "/sign_up", http.StatusMovedPermanently)
			return
		}
		query = fmt.Sprintf("SELECT `login`,`password`,`page` FROM `users` WHERE `login` = '%s' AND `password` = '%s';", login_up, pass_up)

		res, err = db.Query(query)
		if err != nil {
			http.Redirect(w, r, "/sign_up", http.StatusMovedPermanently)
		}
		defer res.Close()

		user := User{}
		for res.Next() {
			err := res.Scan(&user.Login, &user.Password, &user.Page)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		setCookie(w, "session", login_up)
		http.Redirect(w, r, "/account?id="+user.Page, http.StatusMovedPermanently)

	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "html/tmpl/sign_in.html")
	} else {

		http.Error(w, "Метод запрещен!", http.StatusMethodNotAllowed)
		return
	}

}
