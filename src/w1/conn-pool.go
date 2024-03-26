package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	chn := make(chan int, 1)
	chn <- 12
	close(chn)
	for i := range chn {
		fmt.Println(i)
	}
	db, err := sql.Open("mysql", "airflow_user:Airflow123!@/airflow_db")
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("select 1")
	if err != nil {
		panic(err)
	}
	for {
		if !rows.Next() {
			break
		}
		cols, _ := rows.Columns()
		for i, v := range cols {
			fmt.Printf("value @ %v: %v", i, v)
		}
	}
}
