package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type todos struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var todo []todos

func main() {

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			getAllTasks(w)uyyyyyyyy8
		case http.MethodPost:
			postTask(w, r)
		case http.MethodPut:
			putTask(w, r)
		case http.MethodDelete:
			deleteTask(w, r)
		default:
			http.Error(w, "Metod not found", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
	})
	http.ListenAndServe(":8080", nil)
}

func getAllTasks(w http.ResponseWriter) {
	sendingJson, _ := json.Marshal(todo)               // Преобразуем массив структур в формат json
	w.Header().Set("Content-type", "application/json") // Появляется красивый сайт, а не пустой сайт, можно и без этого
	w.WriteHeader(http.StatusOK)                       // Пишем в заголовок, что всё норм
	w.Write(sendingJson)                               // Пишем в ответ
}

func postTask(w http.ResponseWriter, r *http.Request) {
	newTodo := decodeRequest(w, r)
	todo = append(todo, newTodo)
	w.WriteHeader(http.StatusOK) // Пишем в заголовок, что всё норм
}

func putTask(w http.ResponseWriter, r *http.Request) {
	taskIndex := getIndexTask(w, r)
	newTodo := decodeRequest(w, r)
	todo[taskIndex] = newTodo
	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	taskIndex := getIndexTask(w, r)
	todo = append(todo[:taskIndex], todo[taskIndex+1:]...) // уменьшаем слайс у учётом сдвижения на один индекс
	w.WriteHeader(http.StatusNoContent)
}

func getIndexTask(w http.ResponseWriter, r *http.Request) int {
	indexId := r.URL.Path[len("/todos/"):]
	taskIndex, err := strconv.Atoi(indexId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
	}

	if taskIndex < 0 || taskIndex >= len(todo) {
		http.Error(w, "Task not found", http.StatusNotFound)
	}
	return taskIndex
}

func decodeRequest(w http.ResponseWriter, r *http.Request) todos {
	decodedRequest := json.NewDecoder(r.Body) // распаршиваем тело запроса
	var newTodo todos
	err := decodedRequest.Decode(&newTodo) // из распаршенного тела запроса кодируем в указатель структуру newTodo
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest) // Пишем в заголовок, что всё плохо
	}
	return newTodo
}