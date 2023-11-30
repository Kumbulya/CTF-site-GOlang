package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Login    string
	Password string
	Page     string
}

func home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"html/tmpl/home.page.html",
		"html/tmpl/base.layout.html",
		"html/tmpl/footer.partial.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	//http.ServeFile(w, r, "html/index.html")
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func account(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	//http.ServeFile(w, r, "../html/account.html")
	fmt.Fprintf(w, "Отображение выбранного аккаунта с ID %d...", id)
}

func sign_up(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		login_up := r.FormValue("login")
		pass_up := r.FormValue("pass")
		query := fmt.Sprintf("SELECT `login`,`password`,`page` FROM `users` WHERE `login` = '%s' AND `password` = '%s'", login_up, pass_up)

		res, err := db.Query(query)
		if err != nil {
			http.Redirect(w, r, "/sign_in", 301)
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
		http.Redirect(w, r, "/account"+user.Page, 301)

	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "html/tmpl/sign_up.html")
	} else {

		http.Error(w, "Метод запрещен!", 405)
		return
	}

}

func sign_in(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		login_in := r.FormValue("login_si")
		pass_in := r.FormValue("pass_si")
		query := fmt.Sprintf("SELECT COUNT(*) FROM `users` WHERE `login` = '%s'", login_in)
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

		query = fmt.Sprintf("INSERT INTO `users` (`login`, `password`) VALUES ('%s','%s')", login_in, pass_in)
		insert, err := db.Query(query)
		if err != nil {
			log.Println(err)
		}

		defer insert.Close()
		http.Redirect(w, r, "/sign_up", 301)

	} else if r.Method == http.MethodGet {
		http.ServeFile(w, r, "html/tmpl/sign_in.html")
	} else {

		http.Error(w, "Метод запрещен!", 405)
		return
	}

}
