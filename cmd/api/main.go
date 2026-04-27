package main

import (
	"fmt"
	"log"
	"net/http"
	"uploadyPack/handlers"
)

func main() {
	http.HandleFunc("/upload", handlers.FileUploadHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
