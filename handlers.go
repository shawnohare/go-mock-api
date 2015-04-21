package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	m := make(map[string]interface{})
	m["data"] = todos

	if err := json.NewEncoder(w).Encode(m); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["todoID"]
	fmt.Fprintln(w, "Todo show:", todoID)
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var req = make(map[string]Todo)
	// b := io.LimitReader(r.Body, 1048567)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	newTodo := req["data"]
	t := RepoCreateTodo(&newTodo)
	res := map[string]Todo{"data": t}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func TodoCreateMany(w http.ResponseWriter, r *http.Request) {
	// Request body should conform to a JSON API format.
	var newTodos = make(map[string]TodosToAdd)
	// var newTodos TodosToAdd
	var response []Todo
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&newTodos); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	for _, todo := range newTodos["data"].Todos {
		t := RepoCreateTodo(&todo)
		response = append(response, t)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
