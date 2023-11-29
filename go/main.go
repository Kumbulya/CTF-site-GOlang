package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Login    string
	Password string
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/index.html")

}

func account(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/account.html")
}

func sign_up(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		login_up := r.FormValue("login")
		pass_up := r.FormValue("pass")
		query := fmt.Sprintf("SELECT `login`,`password` FROM `users` WHERE `login` = '%s' AND `password` = '%s'", login_up, pass_up)

		res, err := db.Query(query)
		if err != nil {
			http.Redirect(w, r, "/sign_in", 301)
		}
		defer res.Close()

		user := User{}
		for res.Next() {
			err := res.Scan(&user.Login, &user.Password)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		http.Redirect(w, r, "/account", 301)

	} else {
		http.ServeFile(w, r, "html/sign_up.html")
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

	} else {
		http.ServeFile(w, r, "html/sign_in.html")
	}

}

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3307)/katalog") //для докера "test:test@tcp(127.0.0.1:3307)/magazin"
	if err != nil {
		panic(err)
	}

	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index(w, r)
	})
	http.HandleFunc("/sign_up", func(w http.ResponseWriter, r *http.Request) {
		sign_up(w, r, db)
	})
	http.HandleFunc("/sign_in", func(w http.ResponseWriter, r *http.Request) {
		sign_in(w, r, db)
	})
	http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		account(w, r)
	})

	http.ListenAndServe(":1337", nil)

}
