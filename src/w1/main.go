package main

import (
	"net/http"
)

func main() {
	//basic_single()
	//basic_multi(10)
	// fair_multi(10)

	//cp := NewConnectionPool(5)
	// benchmarkNonPool(150)
	// benchmarkPool(10000)

	http.HandleFunc("/events", eventHandler)
	http.ListenAndServe(":8080", nil)
}
