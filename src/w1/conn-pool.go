package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Database struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
}

func explaination() {
	f, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	LockAndGoroutineTest()
	var cfg Config
	json.NewDecoder(f).Decode(&cfg)
	fmt.Println(cfg)

	chn := make(chan int, 1)
	chn <- 12
	close(chn)
	for i := range chn {
		fmt.Println(i)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"root",
		"password",
		"localhost",
		"3306",
		"sakila",
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from actor limit 10")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	for i, v := range cols {
		fmt.Printf("value @ %v: %v\n", i, v)
	}
	for rows.Next() {
		var actorId int
		var firstname string
		var lastname string
		var lastUpdate time.Time

		err = rows.Scan(&actorId, &firstname, &lastname, &lastUpdate)
		if err != nil {
			panic(err)
		}
		lastUpdateStr := lastUpdate.Format(time.RFC1123)
		fmt.Printf("value: %d %s %s %s\n", actorId, firstname, lastname, lastUpdateStr)
	}
}

func LockAndGoroutineTest() {
	var mu sync.Mutex = sync.Mutex{}
	func() {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("Lock acquired")
	}()
	mu.Lock()
}
