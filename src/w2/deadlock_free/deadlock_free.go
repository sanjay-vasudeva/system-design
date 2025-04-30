package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NUM_RECORDS = 3
	NUM_CONN    = 6
)

type RecordData struct {
	attrs [3]int
}

type Record struct {
	data RecordData
	mu   sync.Mutex
}

type Database struct {
	records [NUM_RECORDS]Record
}

var db Database

func acquireLock(txn string, recIdx int) {
	// Acquire the lock on the record at recIdx
	fmt.Printf("%s: Trying to acquire lock on record %d\n", txn, recIdx)
	db.records[recIdx].mu.Lock()
	fmt.Printf("%s: Acquired lock on record %d\n", txn, recIdx)
}

func releaseLock(txn string, recIdx int) {
	// Release the lock on the record at recIdx
	db.records[recIdx].mu.Unlock()
	fmt.Printf("%s: Released lock on record %d\n", txn, recIdx)
}

func InitDatabase() {
	for i := 0; i < NUM_RECORDS; i++ {
		db.records[i] = Record{data: RecordData{attrs: [3]int{i, rand.Intn(100), 0}}}
	}
}

func MimicLoad(txn string) {
	for {
		rec1 := rand.Intn(1000) % NUM_RECORDS
		rec2 := rand.Intn(1000) % NUM_RECORDS

		if rec1 == rec2 {
			continue
		}
		//Total Order imposed
		if rec1 > rec2 {
			rec1, rec2 = rec2, rec1
		}
		acquireLock(txn, rec1)
		acquireLock(txn, rec2)

		time.Sleep(2 * time.Second)

		releaseLock(txn, rec1)
		releaseLock(txn, rec2)

		time.Sleep(1 * time.Second)
	}
}

func main() {
	wg := sync.WaitGroup{}
	/*

		for range 1000 {
			wg.Add(1)
			go func() {
				ol.IncrementCounter()
				wg.Done()
			}()
		}
		wg.Wait()
	*/
	InitDatabase()
	for i := range NUM_CONN {
		wg.Add(1)
		go func(i int) {
			MimicLoad("Transaction" + fmt.Sprint(i))
		}(i)
	}
	wg.Wait()
}
