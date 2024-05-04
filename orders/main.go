package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Order service will call the four endpoints - two of delivery service and two of store service
// Since this would be a distributed transaction, we need to implement a two-phase commit protocol
// The two-phase commit protocol will have the following steps:
// 1. Prepare phase - Call the reserve endpoints of delivery and store service
// 2. Commit phase - Call the book endpoints of delivery and store service
// If any of the endpoints fail, we need to rollback the transaction

// Write a function that makes the above calls
// To simulate that want fire 10 go routines that will call the above function

func main() {

	// Now we will call the endpoints in a loop with 10 go routines
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			callEndpoints()
		}()
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	duration := time.Since(start)
	fmt.Println("Time spent", duration)
}

func callEndpoints() {

	// Reserve food packet

	response, err := http.Post("http://localhost:8081/store/reserve", "application/json", nil)
	if err != nil {
		fmt.Println("Error reserving food packet:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error reserving food packet:", response.Status)
		return
	}

	fmt.Println("Food packet reserved successfully")

	// Reserve delivery agent

	response, err = http.Post("http://localhost:8080/delivery-agent/reserve", "application/json", nil)
	if err != nil {
		fmt.Println("Error reserving delivery agent:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {

		fmt.Println("Error reserving delivery agent:", response.Status)
		return

	}

	fmt.Println("Delivery agent reserved successfully")

	// Book food packet
	orderID := rand.Intn(1000)                        // Generate a random order ID
	requestData := map[string]int{"orderId": orderID} // Replace 123 with the actual order ID

	// Convert the request data to JSON
	payload, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	response, err = http.Post("http://localhost:8081/store/book", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error booking food packet:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error booking food packet:", response.Status)
		return
	}

	fmt.Println("Food packet booked successfully - Order ID:", orderID)
	// Book delivery agent
	requestData = map[string]int{"orderId": orderID}

	// Convert the request data to JSON
	payload, err = json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	response, err = http.Post("http://localhost:8080/delivery-agent/book", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error booking delivery agent:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error booking delivery agent:", response.Status)
		return
	}

	fmt.Println("Delivery agent booked successfully - Order ID:", orderID)
}
