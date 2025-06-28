package Route

import (
	"database/sql"
	"project/Controller"
	"project/Middleware"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/users", Middleware.MiddlewareAuth(Controller.GetUser(db))).Methods("GET") // dilindungi
	r.HandleFunc("/users/{id}", Middleware.MiddlewareAuth(Controller.GetUserByID(db))).Methods("GET")
	r.HandleFunc("/users", Controller.CreateUser(db)).Methods("POST") // public
	r.HandleFunc("/login", Controller.LoginUser(db)).Methods("POST")  // public
}
