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

	query := "SELECT `page` FROM `users` WHERE `login` = ?"
	res, err := db.Query(query, cookie.Value)
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

	query = "INSERT INTO `basket`(`basketID`, `productID`) VALUES (?, ?)"

	insert, err := db.Query(query, basket.UserID, basket.ProductID)
	if err != nil {
		log.Println(err)
	}

	defer insert.Close()
	http.Redirect(w, r, "/product?id="+productID, http.StatusFound)
}

func basket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userPage := r.URL.Query().Get("id")
	query := "SELECT * FROM `basket` WHERE `basketID` = ?"

	res, err := db.Query(query, userPage)
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
		query := "SELECT * FROM `katalog` WHERE `id` = ?"

		res, err := db.Query(query, basket_prod.ProductID)
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

	query = "SELECT `login` FROM `users` WHERE `page` = ?"

	res, err = db.Query(query, userPage)
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
	var balance float64
	var page string

	balanceQuery := "SELECT `page`, `balance` FROM `users` WHERE `login` = ?"
	balanceRes, err := db.Query(balanceQuery, login)
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

	ProductQuery := "SELECT * FROM `basket` WHERE `BasketID` = ?"

	res, err := db.Query(ProductQuery, page)
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

		for _, prod := range baskets {

			path := fmt.Sprintf("html/static/products/product_%d.bmp", prod.ProductID)
			file, err := os.Open(path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				zipWriter.Close()
				return
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				zipWriter.Close()
				return
			}

			header, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				zipWriter.Close()
				return
			}

			// Устанавливаем имя файла в архиве
			header.Name = filepath.Join("products", fmt.Sprintf("file_%d.bmp", prod.ProductID))

			// Создаем запись в архиве
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				zipWriter.Close()
				return
			}

			_, _ = io.Copy(writer, file)
		}
		err = zipWriter.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		ClsQuery := "DELETE FROM `basket` WHERE `BasketID` = ?"

		_, err = db.Exec(ClsQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else if totalCost == 0 {
		w.Write([]byte("<h1>Надо чё-то выбрать</p>"))
	} else {
		w.Write([]byte("<h1><a href='/'>Каталог</a></h1><p>Недостаточно средств!</p>"))
	}

	http.Redirect(w, r, "/basket?id="+page, http.StatusFound)
}

func clear(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	login := r.FormValue("login")
	var page string

	query := "SELECT page FROM `users` WHERE `login` = ?"

	res, err := db.Query(query, login)
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

	query = "DELETE FROM `basket` WHERE `BasketID` = ?"

	_, err = db.Exec(query, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/basket?id="+page, http.StatusFound)

}
