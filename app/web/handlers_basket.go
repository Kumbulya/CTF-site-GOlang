package main

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func add_to_basket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var basket Basket

	productID := r.FormValue("product_id")
	basket.ProductID, _ = strconv.Atoi(productID)
	cookie, _ := getCookie(r, "session")

	query := fmt.Sprintf("SELECT `page` FROM `users` WHERE `login` = '%s'", cookie.Value)
	res, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	for res.Next() {
		err := res.Scan(&basket.UserID)
		if err != nil {
			http.Error(w, "Error scanning result", http.StatusInternalServerError)
			return
		}
	}

	query = fmt.Sprintf("INSERT INTO `basket`(`basketID`, `productID`) VALUES ('%d','%d')", basket.UserID, basket.ProductID)
	insert, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	defer insert.Close()
	http.Redirect(w, r, "/product?id="+productID, http.StatusFound)
}

func basket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userPage := r.URL.Query().Get("id")
	query := fmt.Sprintf("SELECT * FROM `basket` WHERE `basketID` = %s", userPage)
	res, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	var baskets []Basket
	for res.Next() {
		var b Basket
		if err := res.Scan(&b.Id, &b.UserID, &b.ProductID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		baskets = append(baskets, b)
	}

	var products []Product
	var totalCost float32
	totalCost = 0
	for _, basket_prod := range baskets {
		query = fmt.Sprintf("SELECT * FROM `katalog` WHERE `id` = '%d'", basket_prod.ProductID)

		res, err = db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Close()

		for res.Next() {
			var p Product
			if err := res.Scan(&p.Id, &p.Product_name, &p.Category, &p.Seller, &p.Description, &p.Cost); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			totalCost += p.Cost
			products = append(products, p)
		}

	}

	data := struct {
		Login     string
		Products  []Product
		TotalCost float32
	}{
		Login:     "",
		Products:  products,
		TotalCost: totalCost,
	}

	query = fmt.Sprintf("SELECT `login` FROM `users` WHERE `page` = '%s'", userPage)
	res, err = db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	for res.Next() {
		err := res.Scan(&data.Login)
		if err != nil {
			http.Error(w, "Error scanning result", http.StatusInternalServerError)
			return
		}
	}

	files := []string{
		"html/tmpl/basket.html",
		"html/tmpl/base.layout.html",
		"html/tmpl/footer.partial.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func buy(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	login := r.FormValue("login")
	totalCost, _ := strconv.ParseFloat(r.FormValue("cost"), 32)
	balanceQuery := fmt.Sprintf("SELECT `page`,`balance` FROM `users` WHERE `login` = '%s'", login)
	var balance float64
	var page string
	balanceRes, err := db.Query(balanceQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer balanceRes.Close()

	if balanceRes.Next() {
		if err := balanceRes.Scan(&page, &balance); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	ProductQuery := fmt.Sprintf("SELECT  * FROM `basket` WHERE `BasketID` = '%s'", page)
	res, err := db.Query(ProductQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	var baskets []Basket
	for res.Next() {
		var b Basket
		if err := res.Scan(&b.Id, &b.UserID, &b.ProductID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		baskets = append(baskets, b)
	}

	if totalCost <= balance {
		zipFile, err := os.Create("transaction.zip")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		for _, prod := range baskets {

			path := fmt.Sprintf("html/static/products/product_%d.txt", prod.ProductID)
			file, err := os.Open(path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			header, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Устанавливаем имя файла в архиве
			header.Name = filepath.Join("products", fmt.Sprintf("file_%d.txt", prod.ProductID))

			// Создаем запись в архиве
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, _ = io.Copy(writer, file)
		}

		w.Header().Set("Content-Disposition", "attachment; filename=transaction.zip")
		w.Header().Set("Content-Type", "application/zip")

		http.ServeContent(w, r, "transaction.zip", time.Now(), zipFile)

		newBalance := balance - totalCost
		updateBalanceQuery := fmt.Sprintf("UPDATE `users` SET `balance` = %f WHERE `login` = '%s'", newBalance, login)
		_, err = db.Exec(updateBalanceQuery)
		if err != nil {
			log.Println("Failed to update balance:", err)
		}
		clear(w, r, db)
	} else {
		w.Write([]byte("<h1><a href='/'>Каталог</a></h1><p>Недостаточно средств!</p>"))
	}
}

func clear(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	login := r.FormValue("login")

	query := fmt.Sprintf("SELECT page FROM `users` WHERE `login` = '%s'", login)

	var page string

	res, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	for res.Next() {
		err := res.Scan(&page)
		if err != nil {
			http.Error(w, "Error scanning result", http.StatusInternalServerError)
			return
		}
	}

	query = fmt.Sprintf("DELETE FROM `basket` WHERE `BasketID` = '%s'", page)
	_, err = db.Exec(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/basket?id="+page, http.StatusFound)

}
