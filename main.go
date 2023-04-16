package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var schema = `
create extension if not exists "uuid-ossp";


create table if not exists Todo (
	id uuid default uuid_generate_v4() primary key,
	content varchar not null,
	checked bool not null
);`

type Todo struct {
	Id      string `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
	Checked bool   `json:"checked" db:"checked"`
}

type ICreateTodo struct {
	Content string `json:"content" db:"content"`
	Checked bool   `json:"checked" db:"checked"`
}

func connectDb() sqlx.DB {
	db, dbErr := sqlx.Connect("postgres",
		fmt.Sprintf(
			"user=%s dbname=%s password=%s sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_PASSWORD"),
		),
	)

	if dbErr != nil {
		log.Fatalln(dbErr)
	}

	return *db
}

func createTodo(w http.ResponseWriter, req *http.Request) {
	db := connectDb()

	body := &ICreateTodo{}

	decodeError := json.NewDecoder(req.Body).Decode(body)

	if decodeError != nil {
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		`INSERT INTO public.todo ("content", checked) 
			VALUES ($1, $2)`,
		body.Content,
		body.Checked,
	)

	if err != nil {
		log.Fatal(err)
	}

	bodyJson, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bodyJson))
}

func getManyTodo(w http.ResponseWriter, req *http.Request) {
	db := connectDb()

	todos := []Todo{}

	dbErr := db.Select(&todos, `SELECT id, "content", checked FROM public."todo"`)

	if dbErr != nil {
		log.Fatal(dbErr)
		return
	}

	todosJson, _ := json.Marshal(todos)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(todosJson))
}

func getOneTodo(w http.ResponseWriter, req *http.Request) {
	db := connectDb()
	vars := mux.Vars(req)
	todo := Todo{}
	dbErr := db.Get(&todo, `SELECT id, "content", checked FROM public."todo" WHERE id=$1`, vars["todoId"])

	if dbErr != nil {
		log.Fatal(dbErr)
		return
	}

	todoJson, _ := json.Marshal(todo)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(todoJson))
}

func updateOneTodo(w http.ResponseWriter, req *http.Request) {
	db := connectDb()
	vars := mux.Vars(req)
	body := &ICreateTodo{}

	decodeError := json.NewDecoder(req.Body).Decode(body)

	if decodeError != nil {
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		`UPDATE public.todo
			SET 
				"content"=$2,
				checked=$3
			WHERE id=$1`,
		vars["todoId"],
		body.Content,
		body.Checked,
	)

	if err != nil {
		log.Fatal(err)
	}

	bodyJson, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bodyJson))
}

func deleteOneTodo(w http.ResponseWriter, req *http.Request) {
	db := connectDb()
	vars := mux.Vars(req)

	_, err := db.Exec(
		`DELETE FROM public.todo
			WHERE id=$1`,
		vars["todoId"],
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Delete todo")
}

func main() {
	envLoadError := godotenv.Load()

	if envLoadError != nil {
		log.Fatal("Error loading .env file")
	}

	db := connectDb()
	db.MustExec(schema)

	r := mux.NewRouter()
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos", getManyTodo).Methods("GET")
	r.HandleFunc("/todos/{todoId}", getOneTodo).Methods("GET")
	r.HandleFunc("/todos/{todoId}", updateOneTodo).Methods("PATCH")
	r.HandleFunc("/todos/{todoId}", deleteOneTodo).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
