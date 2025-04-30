package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/events", EventHandler)
	http.ListenAndServe(":8080", nil)
}
func EventHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "data: %d\n", i)
		time.Sleep(1 * time.Second)
		w.(http.Flusher).Flush() // Flush the response to the client
	}
	// closeNotifier := w.(http.CloseNotifier).CloseNotify()
	closeNotifier := r.Context().Done()
	<-closeNotifier
	fmt.Println("Client disconnected")
}
