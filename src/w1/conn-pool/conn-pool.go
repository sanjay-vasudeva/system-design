package connPool

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	q "w1/queue"
)

type ConnectionPool struct {
	maxConnections int
	pool           *q.BlockingQueue[*sql.DB]
}

func (cp *ConnectionPool) Take() *sql.DB {
	db := cp.pool.Take()
	return db
}

func (cp *ConnectionPool) Put(db *sql.DB) {
	cp.pool.Put(db)
}

func NewConnectionPool(maxConn int) *ConnectionPool {
	cp := ConnectionPool{
		maxConnections: maxConn,
		pool:           q.NewBlockingQueue[*sql.DB](maxConn),
	}
	for range maxConn {
		cp.pool.Put(newConn())
	}
	return &cp
}

func benchmarkNonPool(count int) {
	now := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(count)
	for range count {
		go func() {
			defer wg.Done()
			db := newConn()
			defer db.Close()

			_, err := db.Exec("SELECT SLEEP(0.01);")
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Time taken: %v\n", time.Since(now))
}

func benchmarkPool(count int) {
	now := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(count)
	cp := NewConnectionPool(10)
	for range count {
		go func() {
			defer wg.Done()
			db := cp.Take()
			defer cp.Put(db)

			_, err := db.Exec("SELECT SLEEP(0.01);")
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Time taken: %v\n", time.Since(now))
}
