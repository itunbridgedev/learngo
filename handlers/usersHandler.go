package handlers

import (
	"database/sql"
	"encoding/json"
	"gocommerce/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the new user into the database.
	// Here, you should also handle password hashing
	sqlStatement := `INSERT INTO users (username, email, passwordhash) VALUES ($1, $2, $3) RETURNING id`
	err := h.DB.QueryRow(sqlStatement, user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	row := h.DB.QueryRow("SELECT id, username, email FROM users WHERE id = $1", id)

	err = row.Scan(&user.ID, &user.Username, &user.Email)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updated models.User
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Here, add logic to handle password updates if necessary (e.g., hashing)

	sqlStatement := `UPDATE users SET username = $1, email = $2, passwordhash = $3 WHERE id = $4`
	_, err = h.DB.Exec(sqlStatement, updated.Username, updated.Email, updated.PasswordHash, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	sqlStatement := `DELETE FROM users WHERE id = $1`
	_, err = h.DB.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
