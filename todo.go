package main

import "time"

type Todo struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Due         time.Time `json:"due"`
}

type TodosToAdd struct {
	Todos []Todo `json:"todos"`
}

type Todos []Todo
