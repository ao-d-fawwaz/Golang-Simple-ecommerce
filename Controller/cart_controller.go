package Controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"project/Model"
	"strconv"
	"time"
)

// GET /cart?user_id=1
func GetCart(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, user_id, product_id, quantity, added_at FROM cart WHERE user_id = ?", userID)
		if err != nil {
			http.Error(w, "failed to query cart: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var carts []Model.Cart
		for rows.Next() {
			var c Model.Cart
			var addedAtStr string

			err := rows.Scan(&c.ID, &c.UserID, &c.ProductID, &c.Quantity, &addedAtStr)
			if err != nil {
				http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			c.AddedAt, _ = time.Parse("2006-01-02 15:04:05", addedAtStr) // Sesuaikan format sesuai DB
			carts = append(carts, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(carts)
	}
}

// POST /cart
func CreateCartItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cart Model.Cart
		err := json.NewDecoder(r.Body).Decode(&cart)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validasi
		if cart.UserID == 0 || cart.ProductID == 0 || cart.Quantity <= 0 {
			http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
			return
		}

		// Tambahkan ke database
		result, err := db.Exec("INSERT INTO cart (user_id, product_id, quantity, added_at) VALUES (?, ?, ?, NOW())",
			cart.UserID, cart.ProductID, cart.Quantity)
		if err != nil {
			http.Error(w, "Insert error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		cart.ID = int(id)
		cart.AddedAt = time.Now()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(cart)
	}
}

// DELETE /cart?id=123
func DeleteCartItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM cart WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Delete error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
