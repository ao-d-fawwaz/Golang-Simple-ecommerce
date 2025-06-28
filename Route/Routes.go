package Route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes(db *sql.DB) http.Handler {
	r := mux.NewRouter()

	// Modular routing
	ProductRoutes(r, db)
	UserRoutes(r, db)
	RegisterCartRoutes(r, db)

	return r
}
