package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...
func getTasks(res http.ResponseWriter, req *http.Request) {
	resp, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		res.Header().Set("Content-Type", "application/json")
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		log.Println(err, "ошибка записи ответа")
		return
	}
}
func addTask(res http.ResponseWriter, req *http.Request) {
	var task Task
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		res.Header().Set("Content-Type", "application/json")
		return
	}
	if _, ok := tasks[task.ID]; ok {
		res.WriteHeader(http.StatusBadRequest)
		res.Header().Set("Content-Type", "application/json")
		return
	}

	tasks[task.ID] = task

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

}

func getTask(res http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
	}
	resp, err := json.MarshalIndent(tasks[task.ID], "", "    ")
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		log.Println(err, "ошибка записи ответа")
		return
	}
}
func deleteTask(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	val, ok := tasks[id]
	if !ok {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	delete(tasks, val.ID)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	r.Get("/tasks", getTasks)
	r.Post("/tasks", addTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
