package src

import "database/sql"

type agent struct {
	id         int
	isReserved bool
	orderID    sql.NullString
}
