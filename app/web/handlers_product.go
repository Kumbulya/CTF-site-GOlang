package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Basket struct {
	Id        int
	ProductID int
	UserID    int
}

func product(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	query := fmt.Sprintf("SELECT katalog.id, katalog.product_name, katalog.category, katalog.seller, katalog.description, katalog.cost, users.login FROM katalog INNER JOIN users ON katalog.seller = users.page WHERE katalog.id = '%s'", id)
	res, err := db.Query(query)
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusMovedPermanently)
	}
	var product Product

	for res.Next() {
		err := res.Scan(&product.Id, &product.Product_name, &product.Category, &product.Seller, &product.Description, &product.Cost, &product.Seller_name)
		if err != nil {
			http.Error(w, "Error scanning result", http.StatusInternalServerError)
			return
		}
	}

	files := []string{
		"html/tmpl/product.page.html",
		"html/tmpl/base.layout.html",
		"html/tmpl/footer.partial.html",
	}

	templ, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = templ.Execute(w, product)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func search(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	q := r.URL.Query().Get("q")

	query := fmt.Sprintf("SELECT katalog.id, katalog.product_name, katalog.category, katalog.seller, katalog.description, katalog.cost, users.login FROM katalog INNER JOIN users ON katalog.seller = users.page WHERE katalog.product_name LIKE '%%%s%%'", q)
	rows, err := db.Query(query)
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

	err = ts.Execute(w, products)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}
