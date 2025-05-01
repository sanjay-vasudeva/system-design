package src

import "database/sql"

type stock struct {
	id         int
	orderID    sql.NullString
	isReserved bool
}
