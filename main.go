package main

import (
	"database/sql"
	"fmt"
	"gocommerce/handlers"
	"gocommerce/middleware"
	"gocommerce/routes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Assuming you have set the connection string as an environment variable or in your code
const (
	host     = "postgres"
	port     = 5432 // Default port for PostgreSQL
	user     = "mydbuser"
	password = "supersecret"
	dbname   = "gocommerce"
)

var db *sql.DB

func connectToDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initializeDatabase(maxRetries int, delay time.Duration, dsn string) *sql.DB {
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = connectToDatabase(dsn)
		if err == nil {
			log.Println("Successfully connected to the database.")
			return db
		}

		log.Printf("Failed to connect to database: %v. Retrying in %v...\n", err, delay)
		time.Sleep(delay)
	}

	log.Fatalf("Could not connect to the database after %d attempts: %v", maxRetries, err)
	return nil
}

func recoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %s", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db = initializeDatabase(5, 10*time.Second, dsn)
	defer db.Close()

	r := mux.NewRouter()
	r.Use(recoverHandler)
	r.Use(middleware.AuthenticationMiddleware)

	allHandlers := handlers.NewHandlers(db)
	routes.RegisterAll(r, allHandlers)

	// Start server
	http.ListenAndServe(":8088", r)
	log.Println("Server starting on :8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
