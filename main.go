package main

import (
	"database/sql"
	"log"
	"net/http"
	"project/Route"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Koneksi ke DB
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/toko")
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	// Buat router utama
	r := mux.NewRouter()

	// Register routes untuk API
	Route.UserRoutes(r, db)
	Route.ProductRoutes(r, db)
	Route.RegisterCartRoutes(r, db)

	// Serve static files dari folder "public"
	// Semua request yang tidak cocok dengan API -> coba diambil dari "public"
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
