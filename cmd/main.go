package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	logger := slog.Logger{}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	fmt.Println("Server starting!")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

	logger.Info("Server closed!")
}
