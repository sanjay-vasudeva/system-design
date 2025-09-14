package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/sanjay-vasudeva/ioutil"
)

func main() {
	// conn := CreateZooKeeperConn()
	// WatchInRepeat(conn, "/test")
	db := ioutil.NewConn("3309", "root", "password", "sharding")
	InsertEntriesWithRetry(db)
}

func InsertEntriesWithRetry(db *sql.DB) {
	defer func() {
		if rec := recover(); rec != nil {
			sleepDuration := 250 * time.Millisecond
			fmt.Printf("Error: %s while inserting data. Retrying after %v", rec, sleepDuration)
			time.Sleep(sleepDuration)
			return
		}
	}()
	for {
		time.Sleep(50 * time.Millisecond)
		db.Exec(
			`
		INSERT INTO test_tbl (value) VALUES 
		('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		,('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test'),('test')
		`)
		fmt.Printf("Inserted data at %s\n", time.Now().Format(time.RFC3339))
	}
}

func WatchInRepeat(conn *zk.Conn, key string) {
	for {
		_, _, testCh, err := conn.GetW(key)
		if err != nil {
			fmt.Printf("Error getting value for %s: %v\n", key, err)
			return
		}
		e := <-testCh
		fmt.Printf("New value for %s: %s\n", key, e.State.String())
	}
}

func CreateZooKeeperConn() *zk.Conn {
	conn, events, err := zk.Connect([]string{"localhost"}, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	go func() {
		for e := range events {
			switch e.Type {
			case zk.EventNodeCreated:
				fmt.Printf("Node created: %s\n", e.Path)
			case zk.EventNodeDeleted:
				fmt.Printf("Node deleted: %s\n", e.Path)
			case zk.EventNodeDataChanged:
				fmt.Printf("Node data changed: %s\n", e.Path)
			case zk.EventSession:
				fmt.Printf("Session event: %s\n", e.State.String())
			default:
				fmt.Printf("Other event: %s\n", e.Type.String())
			}
		}
	}()
	return conn
}
