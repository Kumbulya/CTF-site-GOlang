package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

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
		log.Println("Ошибочка вышла(((")
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
	product_cost, _ := strconv.ParseFloat(r.FormValue("product_cost"), 32)

	query := "INSERT INTO `katalog` (`product_name`, `category`, `seller`, `description`, `cost`) VALUES (?, ?, ?, ?, ?)"
	insert, err := db.Exec(query, product_name, product_category, user, product_description, product_cost)
	if err != nil {
		log.Println(err)
	}

	rowsAffected, err := insert.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	log.Printf("Inserted %d rows", rowsAffected)

	query = "SELECT `id` FROM `katalog` WHERE `product_name` = ?"
	res, err := db.Query(query, product_name)
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusMovedPermanently)
		return
	}

	for res.Next() {
		err := res.Scan(&new_product.Id)
		if err != nil {
			http.Error(w, "Error scanning result", http.StatusInternalServerError)
			return
		}
	}

	avatar, _, err := r.FormFile("product_avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer avatar.Close()

	dst, err := os.Create("html/static/img/product_" + strconv.Itoa(new_product.Id) + ".jpg")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, avatar); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prod, _, err := r.FormFile("product_self")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer prod.Close()

	dst2, err := os.Create("html/static/products/product_" + strconv.Itoa(new_product.Id) + ".bmp")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst2.Close()

	if _, err := io.Copy(dst2, prod); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/account?id="+user, http.StatusMovedPermanently)
}
