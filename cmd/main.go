package main

import (
	"fmt"
	"log"
	"net/http"

	"workmate/internal/handler"
	"workmate/internal/storage"

	"github.com/gorilla/mux"
)

var (
	addr = ":8080"
)

func main() {
	st := storage.New()
	hn := handler.New(st)

	r := mux.NewRouter()
	r.HandleFunc("/add", hn.AddTask).Methods("POST")
	r.HandleFunc("/tasks", hn.GetTasks).Methods("GET")
	r.HandleFunc("/info", hn.GetTask).Methods("GET")
	r.HandleFunc("/delete", hn.DeleteTask).Methods("DELETE")

	fmt.Println("Listening on", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
