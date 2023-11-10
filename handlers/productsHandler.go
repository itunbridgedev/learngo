// productsHandler.go

package handlers

import (
	"database/sql"
	"encoding/json"
	"gocommerce/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	DB *sql.DB
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("GetProducts: Starting to retrieve products")

	// Prepare SQL query
	query := "SELECT id, name, price FROM products"
	rows, err := h.DB.Query(query)
	if err != nil {
		log.Printf("GetProducts: Error executing query: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			log.Printf("GetProducts: Error scanning row: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetProducts: Error iterating over rows: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("GetProducts: Successfully retrieved products")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product models.Product
	row := h.DB.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id)

	err = row.Scan(&product.ID, &product.Name, &product.Price)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create an SQL INSERT statement
	sqlStatement := `INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id`

	// Execute the SQL statement
	err := h.DB.QueryRow(sqlStatement, product.Name, product.Price).Scan(&product.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Authentication and admin check should be implemented here

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updated models.Product
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE products SET name = $1, price = $2 WHERE id = $3`
	_, err = h.DB.Exec(sqlStatement, updated.Name, updated.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Authentication and admin check should be implemented here

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	sqlStatement := `DELETE FROM products WHERE id = $1`
	_, err = h.DB.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
