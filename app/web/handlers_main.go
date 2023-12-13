package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type User struct {
	Id       int
	Login    string
	Password string
	Page     string
	Avatar   bool
	IsAdmin  int
	Balance  float32
	Own      bool
}

type Product struct {
	Id           int
	Product_name string
	Category     string
	Seller       string
	Seller_name  string
	Description  string
	Cost         float32
}

func getCookie(r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func home(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	rows, err := db.Query("SELECT katalog.id, katalog.product_name, katalog.category, katalog.seller, katalog.description, katalog.cost, users.login FROM katalog INNER JOIN users ON katalog.seller = users.page")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Id, &p.Product_name, &p.Category, &p.Seller, &p.Description, &p.Cost, &p.Seller_name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		products = append(products, p)
	}

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

	err = ts.Execute(w, products)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func account(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")

	query := fmt.Sprintf("SELECT COUNT(*) FROM `users` WHERE `page` = '%s'", id)

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
		w.Write([]byte("Такого аккаунта нет"))
		return
	}

	query = fmt.Sprintf("SELECT * FROM `users` WHERE `page` = '%s'", id)

	res, err = db.Query(query)
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusMovedPermanently)
	}

	defer res.Close()

	user := User{}
	for res.Next() {
		err := res.Scan(&user.Id, &user.Login, &user.Password, &user.Page, &user.IsAdmin, &user.Balance)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	cookie, _ := getCookie(r, "session")

	if cookie.Value != user.Login {
		user.Own = false
	} else {
		user.Own = true
	}

	if err != nil {

		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}
	user.Avatar = true
	if _, err := os.Stat("html/static/img/account_" + user.Page + ".jpg"); os.IsNotExist(err) {
		user.Avatar = false
	}

	files := []string{
		"html/tmpl/account.page.html",
		"html/tmpl/base.layout.html",
		"html/tmpl/footer.partial.html",
	}

	templ, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = templ.Execute(w, user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
