package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Address struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	IsMain  bool   `json:"is_main"`
}

type Contact struct {
	Name        string `json:"name"`
	JobPosition string `json:"job_position"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Mobile      string `json:"mobile"`
	IsMain      bool   `json:"is_main"`
}

type Group struct {
	GroupName string `json:"group_name"`
	Value     string `json:"value"`
	IsActive  bool   `json:"is_active"`
}

type CreateSupplierRequest struct {
	SupplierName string    `json:"supplier_name"`
	NickName     string    `json:"nick_name"`
	Addresses    []Address `json:"addresses"`
	Contacts     []Contact `json:"contacts"`
	Groups       []Group   `json:"groups"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var db *sql.DB

func sendResponse(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Data: data, Message: message})
}

func runMigrations(dbURL string) {
	migrationsPath := "file://db/migrations"
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error during migration 'up': %v", err)
	}
	log.Println("✅ Database migration completed.")
}

func createSupplierHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, nil, "Failed to start transaction")
		return
	}

	defer tx.Rollback()

	var supplierID int
	supplierStatus := "In Progress"
	err = tx.QueryRow(
		"INSERT INTO suppliers (name, nick_name, status) VALUES ($1, $2, $3) RETURNING id",
		req.SupplierName, req.NickName, supplierStatus,
	).Scan(&supplierID)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, nil, "Failed to create supplier")
		return
	}

	for _, addr := range req.Addresses {
		_, err := tx.Exec("INSERT INTO supplier_addresses (supplier_id, name, address, is_main) VALUES ($1, $2, $3, $4)",
			supplierID, addr.Name, addr.Address, addr.IsMain)
		if err != nil {
			sendResponse(w, http.StatusInternalServerError, nil, "Failed to save address")
			return
		}
	}

	for _, contact := range req.Contacts {
		_, err := tx.Exec("INSERT INTO supplier_contacts (supplier_id, name, job_position, email, phone, mobile, is_main) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			supplierID, contact.Name, contact.JobPosition, contact.Email, contact.Phone, contact.Mobile, contact.IsMain)
		if err != nil {
			sendResponse(w, http.StatusInternalServerError, nil, "Failed to save contact")
			return
		}
	}

	for _, group := range req.Groups {
		_, err := tx.Exec("INSERT INTO supplier_groups (supplier_id, group_name, value, is_active) VALUES ($1, $2, $3, $4)",
			supplierID, group.GroupName, group.Value, group.IsActive)
		if err != nil {
			sendResponse(w, http.StatusInternalServerError, nil, "Failed to save group")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		sendResponse(w, http.StatusInternalServerError, nil, "Failed to commit transaction")
		return
	}

	sendResponse(w, http.StatusCreated, map[string]int{"supplier_id": supplierID}, "Supplier created successfully")
}

func getSuppliersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, nick_name, status, created_at FROM suppliers ORDER BY id")
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, nil, "Failed to retrieve suppliers")
		return
	}
	defer rows.Close()

	type SupplierInfo struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		NickName  string    `json:"nick_name"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	suppliers := []SupplierInfo{}
	for rows.Next() {
		var s SupplierInfo
		if err := rows.Scan(&s.ID, &s.Name, &s.NickName, &s.Status, &s.CreatedAt); err != nil {
			sendResponse(w, http.StatusInternalServerError, nil, "Failed to process data")
			return
		}
		suppliers = append(suppliers, s)
	}
	
	sendResponse(w, http.StatusOK, suppliers, "Suppliers retrieved successfully")
}

func main() {
	dbURL := "postgres://postgres:postgres@localhost:5432/suppliers_db?sslmode=disable"
	runMigrations(dbURL)

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}
	log.Println("✅ Successfully connected to database for application.")

	http.HandleFunc("POST /suppliers", createSupplierHandler)
	http.HandleFunc("GET /suppliers", getSuppliersHandler)

	port := ":8080"
	log.Printf("🚀 Server running on http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
