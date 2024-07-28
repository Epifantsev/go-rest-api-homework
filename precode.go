package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
func handleTasksGet(res http.ResponseWriter, req *http.Request) {
	resp, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}
func handleTasksPost(res http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

}

func handleGetTaskFromId(res http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")

	if task, ok := tasks[id]; ok {
		resp, err := json.MarshalIndent(tasks[task.ID], "", "    ")
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(resp)
		return
	}
	res.WriteHeader(http.StatusBadRequest)

}
func handleDeleteTask(res http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "id")

	if val, ok := tasks[id]; ok {
		delete(tasks, val.ID)
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusBadRequest)

}

func main() {
	r := chi.NewRouter()
	r.Get("/tasks", handleTasksGet)
	r.Post("/tasks", handleTasksPost)
	r.Get("/tasks/{id}", handleGetTaskFromId)
	r.Delete("/tasks/{id}", handleDeleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
