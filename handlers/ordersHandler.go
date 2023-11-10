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

type OrderHandler struct {
	DB *sql.DB
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	log.Println("GetOrders: Starting to retrieve orders")

	query := "SELECT id, customer_id, total_price, status FROM orders"
	rows, err := h.DB.Query(query)
	if err != nil {
		log.Printf("GetOrders: Error executing query: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.TotalPrice, &o.Status); err != nil {
			log.Printf("GetOrders: Error scanning row: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetOrders: Error iterating over rows: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("GetOrders: Successfully retrieved orders")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var order models.Order
	row := h.DB.QueryRow("SELECT id, customer_id, total_price, status FROM orders WHERE id = $1", id)

	err = row.Scan(&order.ID, &order.CustomerID, &order.TotalPrice, &order.Status)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create an SQL INSERT statement
	sqlStatement := `INSERT INTO orders (customer_id, total_price, status) VALUES ($1, $2, $3) RETURNING id`

	// Execute the SQL statement
	err := h.DB.QueryRow(sqlStatement, order.CustomerID, order.TotalPrice, order.Status).Scan(&order.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var updated models.Order
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE orders SET customer_id = $1, total_price = $2, status = $3 WHERE id = $4`
	_, err = h.DB.Exec(sqlStatement, updated.CustomerID, updated.TotalPrice, updated.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	sqlStatement := `DELETE FROM orders WHERE id = $1`
	_, err = h.DB.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
