package Controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"project/Model"
	"project/config"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var jwtKey = []byte(config.GetEnv("JWT_SECRET"))

func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, username, password, email, created_at FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []Model.User
		for rows.Next() {
			var user Model.User
			var createdAt sql.NullString

			err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &createdAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if createdAt.Valid {
				user.CreatedAt = &createdAt.String
			} else {
				user.CreatedAt = nil
			}

			users = append(users, user)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func GetUserByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var user Model.User
		var createdAt sql.NullString

		query := "SELECT id, username, password, email, created_at FROM users WHERE id = ?"
		err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &createdAt)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if createdAt.Valid {
			user.CreatedAt = &createdAt.String
		} else {
			user.CreatedAt = nil
		}

		json.NewEncoder(w).Encode(user)
	}
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user Model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"
		_, err := db.Exec(query, user.Username, user.Password, user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
	}
}

func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user Model.User
		var createdAt sql.NullString

		query := "SELECT id, username, password, email, created_at FROM users WHERE username = ?"
		err := db.QueryRow(query, loginData.Username).Scan(
			&user.ID, &user.Username, &user.Password, &user.Email, &createdAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if createdAt.Valid {
			user.CreatedAt = &createdAt.String
		} else {
			user.CreatedAt = nil
		}

		if loginData.Password != user.Password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token":   tokenString,
			"message": "Login successful",
		})
	}
}
