package Controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"project/Model"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, file_name, description, price, stock, created_at FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Model.Product
		for rows.Next() {
			var p Model.Product
			err := rows.Scan(&p.ID, &p.Name, &p.FileName, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}
		json.NewEncoder(w).Encode(products)
	}
}

func CreateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		name := r.FormValue("name")
		description := r.FormValue("description")
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		stock, _ := strconv.Atoi(r.FormValue("stock"))

		file, handler, err := r.FormFile("file")
		var fileName string
		if err == nil {
			defer file.Close()
			timestamp := time.Now().Unix()
			fileName = fmt.Sprintf("%d_%s", timestamp, handler.Filename)

			dst, err := os.Create(filepath.Join("uploads", fileName))
			if err != nil {
				http.Error(w, "Failed to save file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			io.Copy(dst, file)
		}
		query := "INSERT INTO products (name, file_name, description, price, stock) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(query, name, fileName, description, price, stock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Response
		product := Model.Product{
			Name:        name,
			FileName:    fileName,
			Description: description,
			Price:       price,
			Stock:       stock,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

func GetProductByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		row := db.QueryRow("SELECT id, name, file_name, description, price, stock, created_at FROM products WHERE id = ?", id)

		var p Model.Product
		err := row.Scan(&p.ID, &p.Name, &p.FileName, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(p)
	}
}

func UpdateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var p Model.Product
		json.NewDecoder(r.Body).Decode(&p)

		_, err := db.Exec("UPDATE products SET name=?, description=?, price=? WHERE id=?", p.Name, p.Description, p.Price, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.ID, _ = strconv.Atoi(id)
		json.NewEncoder(w).Encode(p)
	}
}

func DeleteProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		_, err := db.Exec("DELETE FROM products WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Product deleted"))
	}
}
