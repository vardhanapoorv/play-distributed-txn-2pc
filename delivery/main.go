package main

import (
	"database/sql"
	del "dt2pc/delivery/svc"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	DB *sql.DB
}

func main() {
	// Create basic server
	// Create a basic server that listens on port 8080
	// Service is delivery_agent service should have the following endpoints
	// Reserve - POST - /delivery_agent/Reserve
	// Book - POST - /delivery_agent/Book

	// Call database connection function
	// Call the createDBConn function to create a connection to the database
	// If successful, return a success message
	// If error, return an error message
	db, err := del.CreateDBConn()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		panic(err)
	}
	defer db.Close()
	// Pass db connection to reserveAgentHandler
	app := &App{DB: db}
	r := mux.NewRouter()

	api := r.PathPrefix("/delivery-agent").Subrouter()
	api.HandleFunc("/reserve", app.ReserveAgentHandler)
	api.HandleFunc("/book", app.BookAgentHandler)
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

func (a *App) ReserveAgentHandler(w http.ResponseWriter, r *http.Request) {
	// Reserve a delivery agent
	// Reserve a delivery agent by calling the ReserveAgent function
	// If successful, return a success message
	// If error, return an error message
	err := del.ReserveAgent(a.DB)
	if err != nil {
		fmt.Println("Error reserving agent:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Agent reserved successfully")
}

func (a *App) BookAgentHandler(w http.ResponseWriter, r *http.Request) {
	// Book a delivery agent
	// Book a delivery agent by calling the BookAgent function
	// If successful, return a success message
	// If error, return an error message
	// Get orderId from Request
	// Parse JSON request body
	var requestData map[string]int
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get value of the "orderId" key
	orderId, ok := requestData["orderId"]
	if !ok {
		http.Error(w, "key not found or not a string", http.StatusBadRequest)
		return
	}
	err := del.BookAgent(a.DB, orderId)
	if err != nil {
		fmt.Println("Error booking agent:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Agent booked successfully")
}
