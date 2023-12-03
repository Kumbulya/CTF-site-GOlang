package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3307)/katalog") //для докера "test:test@tcp(127.0.0.1:3307)/magazin"
	if err != nil {
		panic(err)
	}

	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		home(w, r, db)
	})
	mux.HandleFunc("/sign_up", func(w http.ResponseWriter, r *http.Request) {
		sign_up(w, r, db)
	})
	mux.HandleFunc("/sign_in", func(w http.ResponseWriter, r *http.Request) {
		sign_in(w, r, db)
	})
	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		account(w, r, db)
	})

	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		product(w, r, db)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		upload(w, r)
	})

	mux.HandleFunc("/upload_product", func(w http.ResponseWriter, r *http.Request) {
		upload_product(w, r, db)
	})

	fileServer := http.FileServer(http.Dir("html/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.ListenAndServe(":1337", mux)

}
