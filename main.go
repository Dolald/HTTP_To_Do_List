package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type todos struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var todo []todos
var a any

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/todos", getAllTasks).Methods("GET")
	router.HandleFunc("/todos", createTask).Methods("POST")
	router.HandleFunc("/todos/{id}", editTask).Methods("PUT")
	router.HandleFunc("/todos/{id}", deleteTask).Methods("DELETE")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	sendingJson, err := json.Marshal(todo)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(sendingJson)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	newTodo, err := decodeRequest(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
	}
	todo = append(todo, newTodo)
	w.WriteHeader(http.StatusOK)
}

func editTask(w http.ResponseWriter, r *http.Request) {
	taskIndex, err := getIndexTask(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
	}
	newTodo, err := decodeRequest(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
	}
	todo[taskIndex] = newTodo
	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	taskIndex, err := getIndexTask(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
	}

	todo = append(todo[:taskIndex], todo[taskIndex+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func getIndexTask(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	taskIndex, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
	}
	return taskIndex, nil
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (todos, error) {
	decodedRequest := json.NewDecoder(r.Body)
	var newTodo todos
	err := decodedRequest.Decode(&newTodo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return newTodo, err
	}
	return newTodo, nil
}
