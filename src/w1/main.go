package main

func main() {
	//basic_single()
	//basic_multi(10)
	// fair_multi(10)

	//cp := NewConnectionPool(5)
	// benchmarkNonPool(150)
	benchmarkPool(10000)
}
