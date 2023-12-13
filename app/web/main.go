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
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		home(w, r, db)
	})
	mux.HandleFunc("/sign_up", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		sign_up(w, r, db)
	})
	mux.HandleFunc("/sign_in", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		sign_in(w, r, db)
	})
	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		account(w, r, db)
	})

	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		product(w, r, db)
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		search(w, r, db)
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		upload(w, r)
	})

	mux.HandleFunc("/upload_product", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		upload_product(w, r, db)
	})

	mux.HandleFunc("/admin_panel", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		admin_panel(w, r, db)
	})

	mux.HandleFunc("/balance_change", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		balance_change(w, r, db)
	})

	mux.HandleFunc("/add_to_basket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		add_to_basket(w, r, db)
	})

	mux.HandleFunc("/basket", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		basket(w, r, db)
	})

	mux.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		buy(w, r, db)
	})

	mux.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		clear(w, r, db)
	})
	fileServer := http.FileServer(http.Dir("html/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.ListenAndServe(":1337", mux)

}
