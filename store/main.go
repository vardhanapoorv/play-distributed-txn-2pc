package main

import (
	"database/sql"
	store "dt2pc/store/svc"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	//"strconv"
)

type App struct {
	DB *sql.DB
}

func main() {

	db, err := store.CreateDBConn()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		panic(err)
	}
	defer db.Close()
	// Pass db connection to reserveAgentHandler
	app := &App{DB: db}

	r := mux.NewRouter()
	api := r.PathPrefix("/store").Subrouter()
	api.HandleFunc("/reserve", app.ReserveFoodHandler)
	api.HandleFunc("/book", app.BookFoodHandler)
	fmt.Println("Server listening on port 8081...")
	http.ListenAndServe(":8081", r)
}

func (a *App) ReserveFoodHandler(w http.ResponseWriter, r *http.Request) {
	// Reserve a food packet
	// Reserve a food packet by calling the ReserveFoodpacket function
	// If successful, return a success message
	// If error, return an error message
	err := store.ReserveFood(a.DB)
	if err != nil {
		fmt.Println("Error reserving food:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Food reserved successfully")
}

func (a *App) BookFoodHandler(w http.ResponseWriter, r *http.Request) {
	// Book a food packet
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
	//orderID, _ := strconv.Atoi(orderId)
	err := store.BookFood(a.DB, orderId)
	if err != nil {
		fmt.Println("Error booking food packet:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Agent booked successfully")
}
