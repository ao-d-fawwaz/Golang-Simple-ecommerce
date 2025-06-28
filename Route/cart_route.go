package Route

import (
	"database/sql"
	"net/http"
	"project/Controller"

	"github.com/gorilla/mux"
)

// RegisterCartRoutes mendaftarkan semua endpoint terkait cart
func RegisterCartRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/cart", Controller.GetCart(db)).Methods(http.MethodGet)
	r.HandleFunc("/cart", Controller.CreateCartItem(db)).Methods(http.MethodPost)
	r.HandleFunc("/cart", Controller.DeleteCartItem(db)).Methods(http.MethodDelete)
}
