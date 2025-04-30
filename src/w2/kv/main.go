package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	io "github.com/sanjay-vasudeva/ioutil"
)

func main() {
	r := gin.Default()

	db := NewConn("3306")
	dbRead := NewConn("3307")
	r.GET("/", func(ctx *gin.Context) {
		key := ctx.Query("key")
		consistent := ctx.Query("consistent")

		cons, err := strconv.ParseBool(consistent)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid consistent parameter"})
			return
		}
		var row *sql.Row
		if cons {
			row = db.QueryRow("SELECT value FROM kv.store WHERE k = ? AND expired_at > UNIX_TIMESTAMP()", key)
		} else {
			row = dbRead.QueryRow("SELECT value FROM kv.store WHERE k = ? AND expired_at > UNIX_TIMESTAMP()", key)
		}
		if row.Err() != nil {
			ctx.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		var value string
		row.Scan(&value)
		ctx.JSON(200, value)
	})

	r.PUT("", func(c *gin.Context) {
		key := c.Query("key")
		value := c.Query("value")
		ttl, err := strconv.Atoi(c.Query("ttl"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ttl"})
			return
		}
		expiredAt := time.Now().Unix() + int64(ttl)

		defer func() {
			if rec := recover(); rec != nil {
				fmt.Println("Recovered in f", rec)
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
		}()
		// putKey1(key, value, expiredAt, db)
		putKey2(key, value, expiredAt, db)
	})

	r.DELETE("/", func(c *gin.Context) {
		key := c.Query("key")
		deleteKey3(key, db)
	})

	go backgroundCleanUp(60, db)
	r.Run(":8080")
}

func backgroundCleanUp(cadence int, db *sql.DB) {
	for {
		res, err := db.Exec("DELETE FROM kv.store WHERE expired_at < UNIX_TIMESTAMP()")
		if err != nil {
			fmt.Println("Error deleting expired keys:", err)
			continue
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			fmt.Println("Error getting rows affected:", err)
			continue
		}
		if rowsAffected > 0 {
			fmt.Printf("Deleted %d expired keys\n", rowsAffected)
		}
		time.Sleep(time.Second * time.Duration(cadence))
	}
}

// approach 1: Check if key exists and decide whether to insert or update
func putKey1(key string, value string, expiredAt int64, db *sql.DB) {
	row := db.QueryRow("SELECT COUNT(1) FROM kv.store WHERE k = ?", key)
	if row.Err() != nil {
		panic(row.Err())
	}
	var count int
	_ = row.Scan(&count)

	var res sql.Result
	var err error
	if count == 0 {
		res, err = db.Exec("INSERT INTO kv.store (k, value, expired_at) VALUES (?, ?, ?)", key, value, expiredAt)
	} else {
		res, err = db.Exec("UPDATE kv.store SET value = ?, expired_at = ? WHERE k = ?", value, expiredAt, key)
	}

	if err != nil {
		panic(err)
	}
	fmt.Println("Rows affected:", res)
}

// approach 2: Insert or update the key in a single query
func putKey2(key string, value string, expiredAt int64, db *sql.DB) {
	db.Exec("INSERT INTO kv.store (k, value, expired_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE value = ?, expired_at = ?", key, value, expiredAt, value, expiredAt)
}

// approach 1: delete the key from the database
func deleteKey1(key string, db *sql.DB) {
	db.Exec("DELETE FROM kv.store WHERE k = ?", key)
}

// approach 2: set the expired at column to special value so we can
// avoid index rebalancing
func deleteKey2(key string, db *sql.DB) {
	db.Exec("UPDATE kv.store set expired_at = -1 where k = ?", key)
}

// approach 3: Include where clause to filter out already expired keys
// so that we can save 2 disk IOs. One in clustered index and other in secondary index
func deleteKey3(key string, db *sql.DB) {
	db.Exec("UPDATE kv.store set expired_at = -1 where k = ? and expired_at > UNIX_TIMESTAMP()", key)
}

func NewConn(port string) *sql.DB {
	return io.NewConn(port, "root", "password", "kv")
}
