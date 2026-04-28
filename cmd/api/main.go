package main

import (
	"fmt"
	"log"
	"net/http"
	handlersv1 "uploadyPack/handlers/v1"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", handlersv1.FileUploadHandler)
	mux.HandleFunc("/uploadchunk", handlersv1.ChunkUploadHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow the browser to read the response
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle the browser's "preflight" request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
