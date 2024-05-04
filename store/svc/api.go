package svc

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func CreateDBConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "stnduser:stnduser@tcp(127.0.0.1:3309)/store")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func ReserveFood(db *sql.DB) error {
	// Query to reserve a food packet
	// Need to first start a transaction to ensure atomicity - Query non-reserved food packet then update the record

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Query to get a non-reserved food packet - With a exclusive lock to ensure no other transaction can access the same record
	query := "SELECT id FROM food_packet WHERE is_reserved = 0 and order_id is NULL LIMIT 1 FOR UPDATE SKIP LOCKED"
	row := tx.QueryRow(query)

	var foodPacketId int
	if err := row.Scan(&foodPacketId); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No results found.")
		}
		return err
	}

	// Update the record to mark it as reserved
	updateQuery := "UPDATE food_packet SET is_reserved = 1 WHERE id = ?"
	_, err = tx.Exec(updateQuery, foodPacketId)
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

func BookFood(db *sql.DB, orderId int) error {
	// Query to book a food packet
	// Need to first start a transaction to ensure atomicity - Query reserved food packet then update the record

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Query to get a reserved food packet - With a exclusive lock to ensure no other transaction can access the same record
	query := "SELECT id FROM food_packet WHERE is_reserved = 1 and order_id is NULL LIMIT 1 FOR UPDATE SKIP LOCKED"

	row := tx.QueryRow(query)

	var foodPacketId int
	if err := row.Scan(&foodPacketId); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No results found.")
		}
		return err
	}

	// Update the record to mark it as booked and assign the orderId
	updateQuery := "UPDATE food_packet SET order_id = ?, is_reserved = 0 WHERE id = ?"
	_, err = tx.Exec(updateQuery, orderId, foodPacketId)
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
