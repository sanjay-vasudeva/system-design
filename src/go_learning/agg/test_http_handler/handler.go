package main

import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
	// Handle the HTTP request here
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("FOO"))
}
