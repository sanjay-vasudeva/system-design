package main

import (
	utils "airline_checkin/internal/utils"
	"fmt"
	"sync"
	"time"
)

type user struct {
	id   int
	name string
}

func main() {
	book()
	utils.PrintSeats()
	utils.Clean()
}

func book() {
	start := time.Now()
	db := utils.NewConn()
	rows, err := db.Query("SELECT id,name from airline_checkin.users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	users := []user{}
	for rows.Next() {
		var user user
		rows.Scan(&user.id, &user.name)
		users = append(users, user)
	}
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		// time.Sleep(10 * time.Millisecond)
		go func(userId int) {
			tx, err := db.Begin()
			if err != nil {
				panic(err)
			}
			row := tx.QueryRow("SELECT id,name FROM airline_checkin.seats WHERE user_id is NULL ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED")
			if row.Err() != nil {
				fmt.Printf("Unable to pick a seat. Message: %v", err)
			}
			var id int
			var name string
			row.Scan(&id, &name)

			tx.Exec("UPDATE airline_checkin.seats set user_id = ? where id = ?", userId, id)
			tx.Commit()
			fmt.Printf("User %d got seat %s\n", userId, name)
			wg.Done()
		}(user.id)
	}
	wg.Wait()

	fmt.Printf("\nTook %d milliseconds\n", time.Since(start).Milliseconds())
}
