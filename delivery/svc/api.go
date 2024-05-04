package svc

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func CreateDBConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "stnduser:stnduser@tcp(127.0.0.1:3309)/delivery")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func ReserveAgent(db *sql.DB) error {
	// Query to reserve a delivery
	// Need to first start a transaction to ensure atomicity - Query non-reserved delivery agent then update the record

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Query to get a non-reserved delivery agent - With a exclusive lock to ensure no other transaction can access the same record
	query := "SELECT id FROM delivery_agent WHERE is_reserved = 0 and order_id is NULL LIMIT 1 FOR UPDATE SKIP LOCKED"
	row := tx.QueryRow(query)

	//fmt.Println("Waiting for lock")
	var deliveryAgentID int
	if err := row.Scan(&deliveryAgentID); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No results found.")
		}
		return err
	}

	// Update the record to mark it as reserved
	updateQuery := "UPDATE delivery_agent SET is_reserved = 1 WHERE id = ?"
	_, err = tx.Exec(updateQuery, deliveryAgentID)
	if err != nil {
		tx.Rollback()
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func BookAgent(db *sql.DB, orderId int) error {
	// Query to book a delivery agent
	// Need to first start a transaction to ensure atomicity - Query reserved delivery agent then update the record

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Query to get a reserved delivery agent - With a exclusive lock to ensure no other transaction can access the same record
	query := "SELECT id FROM delivery_agent WHERE is_reserved = 1 and order_id is NULL LIMIT 1 FOR UPDATE SKIP LOCKED"
	row := tx.QueryRow(query)

	var deliveryAgentID int
	if err := row.Scan(&deliveryAgentID); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No results found.")
		}
		return err
	}

	// Update the record to mark it as booked and assign the orderId
	updateQuery := "UPDATE delivery_agent SET order_id = ?, is_reserved = 0 WHERE id = ?"
	_, err = tx.Exec(updateQuery, orderId, deliveryAgentID)
	if err != nil {
		tx.Rollback()
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
