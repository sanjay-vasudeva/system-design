package io

import "fmt"

func PrintSeats() {
	fmt.Println()
	conn := NewConn()
	rows, err := conn.Query("SELECT id,name,user_id FROM airline_checkin.seats")
	if err != nil {
		panic(err)
	}
	i := 1
	for rows.Next() {
		var id, user_id int
		var name string

		rows.Scan(&id, &name, &user_id)
		if user_id == 0 {
			fmt.Print(" - ")
		} else {
			fmt.Print(" * ")
		}
		if i%5 == 0 {
			fmt.Println()
		}
		i++
	}
}
func Clean() {
	tx, err := NewConn().Begin()
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec("UPDATE airline_checkin.seats set user_id = NULL")
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
