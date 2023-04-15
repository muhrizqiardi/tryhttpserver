package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func createTodo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Create todo")
}

func getManyTodo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Get many todo")
}

func getOneTodo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Get todo")
}

func updateOneTodo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Update todo")
}

func deleteOneTodo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Delete todo")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos", getManyTodo).Methods("GET")
	r.HandleFunc("/todos/{todoId}", getOneTodo).Methods("GET")
	r.HandleFunc("/todos/{todoId}", updateOneTodo).Methods("PATCH")
	r.HandleFunc("/todos/{todoId}", deleteOneTodo).Methods("DELETE")

	http.Handle("/", r)
}
