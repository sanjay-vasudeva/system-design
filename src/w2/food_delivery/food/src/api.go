package src

import (
	"database/sql"
)

func Reserve(db *sql.DB) (*stock, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(`
		SELECT id, is_reserved, order_id FROM stock
		WHERE
			food_id = 1 AND is_reserved = 0 AND order_id is NULL
		LIMIT 1
		FOR UPDATE`)
	if row.Err() != nil {
		tx.Rollback()
		return nil, err
	}
	var stock stock
	if err := row.Scan(&stock.id, &stock.isReserved, &stock.orderID); err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec("UPDATE stock SET is_reserved = 1 WHERE id = ?", stock.id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func Book(orderID string, db *sql.DB) (*stock, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(`
		SELECT 
			id, is_reserved, order_id 
		FROM 
			stock 
		WHERE
			food_id = 1 AND is_reserved = 1 
		LIMIT 1 
		FOR UPDATE`)

	if row.Err() != nil {
		tx.Rollback()
		return nil, err
	}
	var stock stock
	if err := row.Scan(&stock.id, &stock.isReserved, &stock.orderID); err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.Exec("UPDATE stock SET is_reserved = 0, order_id = ? WHERE id = ?", orderID, stock.id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &stock, nil
}
