package src

import (
	"context"
	"database/sql"
	"time"
)

func Reserve(db *sql.DB) (*agent, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}
	timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	row := tx.QueryRowContext(timeout, `
		SELECT 
			id, is_reserved, order_id 
		FROM 
			agent
		WHERE
			is_reserved = 0 AND order_id IS NULL
		LIMIT 1
		FOR UPDATE`)
	if row.Err() != nil {
		tx.Rollback()
		return nil, row.Err()
	}
	var agent agent
	if err := row.Scan(&agent.id, &agent.isReserved, &agent.orderID); err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec("UPDATE agent SET is_reserved = 1 WHERE id = ?", agent.id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func Book(orderID string, db *sql.DB) (*agent, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(`
		SELECT	
			id, is_reserved, order_id
		FROM
			agent
		WHERE
			is_reserved = 1 
		LIMIT 1 
		FOR UPDATE
	`)
	if row.Err() != nil {
		tx.Rollback()
		return nil, err
	}
	var agent agent
	if err := row.Scan(&agent.id, &agent.isReserved, &agent.orderID); err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec("UPDATE agent SET is_reserved = 0, order_id = ? WHERE id = ?", orderID, agent.id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &agent, nil
}
