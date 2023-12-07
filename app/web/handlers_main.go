package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	Id       int
	Login    string
	Password string
	Page     string
	Avatar   bool
	Own      bool
}

type Product struct {
	Id           int
	Product_name string
	Category     string
	Seller       string
	Seller_name  string
	Description  string
}

func getCookie(r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func home(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	rows, err := db.Query("SELECT katalog.id, katalog.product_name, katalog.category, katalog.seller, katalog.description, users.login FROM katalog INNER JOIN users ON katalog.seller = users.page")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Id, &p.Product_name, &p.Category, &p.Seller, &p.Description, &p.Seller_name); err != nil {
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
		err := res.Scan(&user.Id, &user.Login, &user.Password, &user.Page)
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

func upload(w http.ResponseWriter, r *http.Request) {
	page := r.FormValue("page")
	file, _, err := r.FormFile("account_avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	dst, err := os.Create("html/static/img/account_" + page + ".jpg")

	if err != nil {
		log.Println("Pizda rulu")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/account?id="+page, http.StatusMovedPermanently)

}

func upload_product(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	new_product := Product{}

	user := r.FormValue("user")
	product_name := r.FormValue("product_name")
	product_category := r.FormValue("product_category")
	product_description := r.FormValue("product_description")

	query := fmt.Sprintf("INSERT INTO `katalog` (`product_name`, `category`, `seller`, `description`) VALUES ('%s','%s','%s','%s')",
		product_name, product_category, user, product_description)
	insert, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	defer insert.Close()

	query = fmt.Sprintf("SELECT `id` FROM `katalog` WHERE `product_name` = '%s'", product_name)
	res, err := db.Query(query)
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusMovedPermanently)
	}

	for res.Next() {
		err := res.Scan(&new_product.Id)
		if err != nil {
			http.Error(w, "Error scanning result", http.StatusInternalServerError)
			return
		}
	}

	file, _, err := r.FormFile("product_avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	dst, err := os.Create("html/static/img/product_" + strconv.Itoa(new_product.Id) + ".jpg")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/account?id="+user, http.StatusMovedPermanently)
}

func product(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	query := fmt.Sprintf("SELECT katalog.id, katalog.product_name, katalog.category, katalog.seller, katalog.description, users.login FROM katalog INNER JOIN users ON katalog.seller = users.page WHERE katalog.id = '%s'", id)
	res, err := db.Query(query)
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusMovedPermanently)
	}
	var product Product

	for res.Next() {
		err := res.Scan(&product.Id, &product.Product_name, &product.Category, &product.Seller, &product.Description, &product.Seller_name)
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

	query := fmt.Sprintf("SELECT katalog.id, katalog.product_name, katalog.category, katalog.seller, katalog.description, users.login FROM katalog INNER JOIN users ON katalog.seller = users.page WHERE katalog.product_name LIKE '%%%s%%'", q)
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Id, &p.Product_name, &p.Category, &p.Seller, &p.Description, &p.Seller_name); err != nil {
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
