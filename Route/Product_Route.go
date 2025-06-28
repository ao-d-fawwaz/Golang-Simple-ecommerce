package Route

import (
	"net/http"

	"database/sql"
	"project/Controller"

	"github.com/gorilla/mux"
)

func ProductRoutes(r *mux.Router, db *sql.DB) http.Handler {

	r.HandleFunc("/products", Controller.GetProducts(db)).Methods("GET")
	r.HandleFunc("/products", Controller.CreateProduct(db)).Methods("POST")
	r.HandleFunc("/products/{id}", Controller.GetProductByID(db)).Methods("GET")
	r.HandleFunc("/products/{id}", Controller.UpdateProduct(db)).Methods("PUT")
	r.HandleFunc("/products/{id}", Controller.DeleteProduct(db)).Methods("DELETE")

	return r
}
